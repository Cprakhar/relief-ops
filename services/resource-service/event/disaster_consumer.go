package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cprakhar/relief-ops/services/resource-service/service"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
)

type disasterConsumer struct {
	kafkaClient *messaging.KafkaClient
	svc         service.ResourceService
}

// NewDisasterConsumer creates a new instance of disasterConsumer.
func NewDisasterConsumer(kc *messaging.KafkaClient, svc service.ResourceService) *disasterConsumer {
	return &disasterConsumer{kafkaClient: kc, svc: svc}
}

// DisasterConsumer starts consuming messages from the specified topics.
func (dc *disasterConsumer) DisasterConsumer(ctx context.Context, topics []string) error {
	return dc.kafkaClient.Consume(ctx, topics, func(eventType, key string, value []byte) error {
		// Handle the event based on the topic
		switch eventType {
		case events.ResourceCommandFind:
			if err := dc.handleFindResources(ctx, value); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown event type: %s", eventType)
		}
		return nil
	})
}

// handleFindResources processes the find resources command.
func (dc *disasterConsumer) handleFindResources(ctx context.Context, value []byte) error {
	var payload events.DisasterEventCreatedPayload
	if err := json.Unmarshal(value, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	// For simplicity, just log the received payload
	log.Printf("Finding resources for disaster: %s...", payload.DisasterID)
	if err := dc.svc.SaveResources(ctx, payload.Range, payload.Location.Latitude, payload.Location.Longitude); err != nil {
		return fmt.Errorf("failed to save resources: %w", err)
	}
	log.Printf("Resources saved for disaster: %s", payload.DisasterID)

	// If Successful, Notify User Service to notify admin to review
	if err := dc.kafkaClient.Produce(ctx, events.UserNotifyAdminReview, payload.DisasterID, value); err != nil {
		return fmt.Errorf("failed to notify user service for admin review: %w", err)
	}
	return nil
}
