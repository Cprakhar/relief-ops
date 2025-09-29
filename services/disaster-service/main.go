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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	userRepo := repo.NewDisasterRepo()
	userService := service.NewDisasterService(userRepo)

	resourceConsumer := event.NewResourceConsumer(kafkaClient, userService)
	topics := []string{events.DisasterCommandDelete}

	go func() {
		if err := resourceConsumer.Consume(ctx, topics); err != nil {
			log.Printf("Error in resource consumer: %v", err)
		}
	}()

	gRPCServer := newgRPCServer(addr, userService, kafkaClient)

	done := make(chan struct{})
	go func() {
		defer close(done)
		log.Printf("Server running on %s", addr)
		if err := gRPCServer.run(ctx); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()
	<-ctx.Done()
	<-done
}
