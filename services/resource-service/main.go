package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/relief-ops/services/resource-service/event"
	"github.com/cprakhar/relief-ops/shared/env"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
)

var (
	brokers = env.GetString("KAFKA_BROKERS", "apache-kafka:9092")
)

func main() {
	// Set up context with signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize Kafka client
	kafkaCfg := &messaging.KafkaConfig{
		Brokers: brokers,
		GroupID: "resource-service-group",
	}

	kafkaClient, err := messaging.NewKafkaClient("resource-service", kafkaCfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer kafkaClient.Close()
	log.Println("Kafka client initialized")

	// Initialize and start the disaster consumer
	topics := []string{events.ResourceCommandFind}
	disasterConsumer := event.NewDisasterConsumer(kafkaClient)

	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := disasterConsumer.DisasterConsumer(ctx, topics); err != nil {
			log.Printf("Error in disaster consumer: %v", err)
		}
	}()
	<-ctx.Done()
	<-done
	log.Println("Resource service stopped")
}
