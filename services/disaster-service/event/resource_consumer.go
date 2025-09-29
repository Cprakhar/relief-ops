package event

import (
	"context"

	"github.com/cprakhar/relief-ops/services/disaster-service/service"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
)

type resourceConsumer struct {
	kafkaClient *messaging.KafkaClient
	svc         service.DisasterService
}

func NewResourceConsumer(kc *messaging.KafkaClient, svc service.DisasterService) *resourceConsumer {
	return &resourceConsumer{kafkaClient: kc, svc: svc}
}

func (rc *resourceConsumer) Consume(ctx context.Context, topics []string) error {
	return rc.kafkaClient.Consume(ctx, topics, func(eventType, key string, value []byte) error {
		// Handle the event based on the topic
		switch eventType {
		case events.DisasterCommandDelete:
			if err := rc.svc.DeleteDisaster(ctx, key); err != nil {
				return err
			}
		}
		return nil
	})
}
