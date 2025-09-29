package event

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
)

type disasterConsumer struct {
	kafkaClient *messaging.KafkaClient
}

func NewDisasterConsumer(kc *messaging.KafkaClient) *disasterConsumer {
	return &disasterConsumer{kafkaClient: kc}
}

func (dc *disasterConsumer) DisasterConsumer(ctx context.Context, topics []string) error {
	return dc.kafkaClient.Consume(ctx, topics, func(eventType, key string, value []byte) error {
		// Handle the event based on the topic
		switch eventType {
		case events.ResourceCommandFind:
			dc.handleFindResources(key, value)
		default:
			// Unknown event type
		}
		return nil
	})
}

func (dc *disasterConsumer) handleFindResources(key string, value []byte) error {
	var payload events.DisasterEventCreatedPayload
	if err := json.Unmarshal(value, &payload); err != nil {
		// Handle error
		if err := dc.kafkaClient.Produce(events.DisasterCommandDelete, key, []byte(payload.DisasterID)); err != nil {
			// Log the failure to produce the delete command
			log.Printf("Failed to produce delete command for disaster %s: %v", payload.DisasterID, err)
		}
		return err
	}

	// Implement logic to find and store resources in the database
	return nil
}
