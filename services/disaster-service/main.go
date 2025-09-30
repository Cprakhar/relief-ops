package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/relief-ops/services/disaster-service/event"
	"github.com/cprakhar/relief-ops/services/disaster-service/repo"
	"github.com/cprakhar/relief-ops/services/disaster-service/service"
	"github.com/cprakhar/relief-ops/shared/env"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
)

var (
	addr    = env.GetString("DISASTER_GRPC_ADDR", ":9002")
	brokers = env.GetString("KAFKA_BROKERS", "apache-kafka:9092")
)

func main() {
	// Set up context with signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize Kafka client
	kafkaCfg := &messaging.KafkaConfig{
		Brokers: brokers,
		GroupID: "disaster-service-group",
	}

	kafkaClient, err := messaging.NewKafkaClient("disaster-service", kafkaCfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer kafkaClient.Close()
	log.Println("Kafka client initialized")

	// Initialize repository and service
	userRepo := repo.NewDisasterRepo()
	userService := service.NewDisasterService(userRepo)

	// Initialize and start the resource consumer
	resourceConsumer := event.NewResourceConsumer(kafkaClient, userService)
	topics := []string{events.DisasterCommandDelete}

	go func() {
		if err := resourceConsumer.Consume(ctx, topics); err != nil {
			log.Printf("Error in resource consumer: %v", err)
		}
	}()

	// Initialize and run the gRPC server
	gRPCServer := newgRPCServer(addr, userService, kafkaClient)

	done := make(chan struct{})
	go func() {
		defer close(done)
		log.Printf("Disaster service running on %s", addr)
		if err := gRPCServer.run(ctx); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()
	<-ctx.Done()
	<-done
	log.Print("Disaster service stopped")
}
