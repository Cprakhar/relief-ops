package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/relief-ops/services/user-service/event"
	"github.com/cprakhar/relief-ops/services/user-service/mail"
	"github.com/cprakhar/relief-ops/services/user-service/repo"
	"github.com/cprakhar/relief-ops/services/user-service/service"
	"github.com/cprakhar/relief-ops/shared/env"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
)

var (
	addr           = env.GetString("USER_GRPC_ADDR", ":9001")
	fromEmail      = env.GetString("FROM_EMAIL", "developerluffy23@gmail.com")
	sendGridAPIKey = env.GetString("SENDGRID_API_KEY", "")
	brokers        = env.GetString("KAFKA_BROKERS", "apache-kafka:9092")
)

func main() {
	// Set up context with signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize Kafka client
	kafkaCfg := &messaging.KafkaConfig{
		Brokers: brokers,
		GroupID: "user-service-group",
	}

	kafkaClient, err := messaging.NewKafkaClient("user-service", kafkaCfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer kafkaClient.Close()
	log.Println("Kafka client initialized")

	mailer := mail.NewSendGrid(fromEmail, sendGridAPIKey)

	// Initialize repository and service
	userRepo := repo.NewUserRepo()
	userService := service.NewUserService(userRepo)

	// Initialize and start the disaster consumer
	topics := []string{events.UserNotifyAdminReview}
	disasterConsumer := event.NewDisasterConsumer(kafkaClient, userService, mailer)

	go func() {
		if err := disasterConsumer.Consumer(ctx, topics); err != nil {
			log.Printf("Error in disaster consumer: %v", err)
		}
	}()
	
	// Initialize and run the gRPC server
	gRPCServer := newgRPCServer(addr, userService)

	done := make(chan struct{})
	go func() {
		defer close(done)
		log.Printf("User service running on %s", addr)
		if err := gRPCServer.run(ctx); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()
	<-ctx.Done()
	<-done
	log.Println("User service stopped")
}
