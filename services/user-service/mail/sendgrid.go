package mail

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/cprakhar/relief-ops/shared/observe/logs"
	"github.com/cprakhar/relief-ops/shared/tools"
	"github.com/cprakhar/relief-ops/shared/types"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

// NewSendGrid creates a new SendGridMailer instance.
func NewSendGrid(fromEmail, apiKey string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

// Send sends an email using the specified template and data.
func (s *SendGridMailer) Send(templateFile, name, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, s.fromEmail)
	to := mail.NewEmail(name, email)

	// template parsing and dynamic data handling
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	retryCfg := &tools.RetryConfig{
		MaxAttempts:   MaxRetries,
		InitialDelay:  2 * time.Second,
		BackoffFactor: 2.0,
		MaxDelay:      30 * time.Second,
		Jitter:        true,
	}

	var statusCode int
	err = tools.RetryWithBackoff(context.Background(), retryCfg, func() error {
		response, err := s.client.Send(message)
		if err != nil {
			return err
		}
		statusCode = response.StatusCode
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			return nil
		}
		return fmt.Errorf("failed to send email, status code: %d, body: %s", response.StatusCode, response.Body)
	})
	if err != nil {
		return statusCode, err
	}
	return statusCode, nil
}

// NotifyMultiple sends the same email to multiple users with proper error handling.
// Returns a slice of errors (one per user) and an aggregate error if any sends failed.
func (s *SendGridMailer) NotifyMultiple(users []*types.User, data any, isSandbox bool) error {
	logger := logs.L()

	if len(users) == 0 {
		return nil
	}

	type result struct {
		email      string
		err        error
		statusCode int
	}

	results := make(chan result, len(users))
	semaphore := make(chan struct{}, 5) // Limit to 5 concurrent sends

	// Send emails concurrently
	for _, user := range users {
		semaphore <- struct{}{} // Acquire semaphore
		go func(u *types.User) {
			defer func() { <-semaphore }() // Release semaphore
			statusCode, err := s.Send(AdminNotifyTemplate, u.Name, u.Email, data, isSandbox)
			results <- result{email: u.Email, err: err, statusCode: statusCode}
		}(user)
	}

	// Collect results
	var failedEmails []string
	successCount := 0
	for range users {
		res := <-results
		if res.err != nil {
			failedEmails = append(failedEmails, res.email)
			logger.Errorw("Failed to send email", "email", res.email, "error", res.err, "statusCode", res.statusCode)
		} else {
			successCount++
		}
	}

	if len(failedEmails) > 0 {
		return fmt.Errorf("failed to send %d/%d emails to: %v", len(failedEmails), len(users), failedEmails)
	}

	logger.Infow("Successfully sent emails", "count", successCount)
	return nil
}
