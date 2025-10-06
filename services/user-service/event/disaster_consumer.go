package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cprakhar/relief-ops/services/user-service/mail"
	"github.com/cprakhar/relief-ops/services/user-service/service"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
)

type disasterConsumer struct {
	kafkaClient *messaging.KafkaClient
	svc         service.UserService
	mailer      mail.Client
	webURL      string
}

// NewDisasterConsumer creates a new instance of disasterConsumer.
func NewDisasterConsumer(kc *messaging.KafkaClient, svc service.UserService, mailer mail.Client, w string) *disasterConsumer {
	return &disasterConsumer{kafkaClient: kc, svc: svc, mailer: mailer, webURL: w}
}

// Consumer starts consuming messages from the specified topics.
func (dc *disasterConsumer) Consumer(ctx context.Context, topics []string) error {
	return dc.kafkaClient.Consume(ctx, topics, func(ctx context.Context, eventType, key string, value []byte) error {
		switch eventType {
		case events.UserNotifyAdminReview:
			if err := dc.handleAdminNotify(ctx, value); err != nil {
				return err
			}
		}
		return nil
	})
}

// handleAdminNotify processes the admin notification event.
func (dc *disasterConsumer) handleAdminNotify(ctx context.Context, value []byte) error {
	users, err := dc.svc.GetAdmins(ctx)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	var data events.DisasterEventCreatedPayload
	if err := json.Unmarshal(value, &data); err != nil {
		return err
	}

	adminData := struct {
		DisasterID  string
		VolunteerID string
		ReviewURL   string
	}{
		DisasterID:  data.DisasterID,
		VolunteerID: data.VolunteerID,
		ReviewURL:   fmt.Sprintf("%s/admin/review/%s", dc.webURL, data.DisasterID),
	}

	return dc.mailer.NotifyMultiple(users, adminData, false)
}
