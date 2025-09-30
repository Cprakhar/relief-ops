package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
)

type disasterConsumer struct {
	kafkaClient *messaging.KafkaClient
}

// NewDisasterConsumer creates a new instance of disasterConsumer.
func NewDisasterConsumer(kc *messaging.KafkaClient) *disasterConsumer {
	return &disasterConsumer{kafkaClient: kc}
}

// DisasterConsumer starts consuming messages from the specified topics.
func (dc *disasterConsumer) DisasterConsumer(ctx context.Context, topics []string) error {
	return dc.kafkaClient.Consume(ctx, topics, func(eventType, key string, value []byte) error {
		// Handle the event based on the topic
		switch eventType {
		case events.ResourceCommandFind:
			if err := dc.handleFindResources(ctx, key, value); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown event type: %s", eventType)
		}
		return nil
	})
}

// handleFindResources processes the find resources command.
func (dc *disasterConsumer) handleFindResources(ctx context.Context, key string, value []byte) error {
	var payload events.DisasterEventCreatedPayload
	if err := json.Unmarshal(value, &payload); err != nil {
		if err := dc.kafkaClient.Produce(ctx, events.DisasterCommandDelete, key, []byte(payload.DisasterID)); err != nil {
			log.Printf("Failed to produce delete command for disaster %s: %v", payload.DisasterID, err)
		}
		return err
	}
	// For simplicity, just log the received payload
	log.Printf("Finding resources for disaster: %s...", payload.DisasterID)
	// Implement logic to find and store resources in the database
	return nil
}
