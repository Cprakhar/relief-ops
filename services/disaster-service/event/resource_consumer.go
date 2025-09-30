package event

import (
	"context"
	"fmt"

	"github.com/cprakhar/relief-ops/services/disaster-service/service"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
)

type resourceConsumer struct {
	kafkaClient *messaging.KafkaClient
	svc         service.DisasterService
}

// NewResourceConsumer creates a new instance of resourceConsumer.
func NewResourceConsumer(kc *messaging.KafkaClient, svc service.DisasterService) *resourceConsumer {
	return &resourceConsumer{kafkaClient: kc, svc: svc}
}

// Consume starts consuming messages from the specified topics.
func (rc *resourceConsumer) Consume(ctx context.Context, topics []string) error {
	return rc.kafkaClient.Consume(ctx, topics, func(eventType, key string, value []byte) error {
		// Handle the event based on the topic
		switch eventType {
		case events.DisasterCommandDelete:
			if err := rc.svc.DeleteDisaster(ctx, key); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown event type: %s", eventType)
		}
		return nil
	})
}
