package mail

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/cprakhar/relief-ops/services/user-service/repo"
	"github.com/cprakhar/relief-ops/shared/tools"
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

	err = tools.RetryWithBackoff(context.Background(), retryCfg, func() error {
		response, err := s.client.Send(message)
		if err != nil {
			return err
		}
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			return nil
		}
		return fmt.Errorf("failed to send email, status code: %d, body: %s", response.StatusCode, response.Body)
	})
	return -1, err
}

// SpamMail sends the same email to multiple users.
func (s *SendGridMailer) SpamMail(users []*repo.User, data any, isSandbox bool) error {
	var err error
	for _, user := range users {
		_, err = s.Send(AdminNotifyTemplate, user.Name, user.Email, data, isSandbox)
	}
	return err
}
