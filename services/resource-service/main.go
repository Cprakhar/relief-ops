package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cprakhar/relief-ops/services/resource-service/event"
	"github.com/cprakhar/relief-ops/services/resource-service/repo"
	"github.com/cprakhar/relief-ops/services/resource-service/service"
	"github.com/cprakhar/relief-ops/shared/db"
	"github.com/cprakhar/relief-ops/shared/env"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
)

var (
	// Kafka configuration
	brokers = env.GetString("KAFKA_BROKERS", "apache-kafka:9092")

	// MongoDB configuration
	mongoURI     = env.GetString("MONGODB_URI", "")
	mongoDB      = env.GetString("MONGODB_DB", "relief_ops")
	mongoTimeout = env.GetTimeDuration("MONGODB_TIMEOUT", 30*time.Second)
	mongoMaxIdle = env.GetTimeDuration("MONGODB_MAX_IDLE", 5*time.Second)
	mongoMaxPool = uint64(env.GetInt("MONGODB_MAX_POOL", 10))
	mongoMinPool = uint64(env.GetInt("MONGODB_MIN_POOL", 2))
)

func main() {
	// Set up context with signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize MongoDB client
	mongoCfg := &db.MongoDBConfig{
		URI:        mongoURI,
		Database:   mongoDB,
		Collection: "resources",
		Timeout:    &mongoTimeout,
		MaxIdle:    &mongoMaxIdle,
		MaxPool:    &mongoMaxPool,
		MinPool:    &mongoMinPool,
	}

	mongoClient, err := db.NewMongoDBClient(mongoCfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

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

	resourceRepo, err := repo.NewResourceRepo(ctx, mongoClient)
	if err != nil {
		log.Fatalf("Failed to create resource repository: %v", err)
	}
	resourceService := service.NewResourceService(resourceRepo)

	// Initialize and start the disaster consumer
	topics := []string{events.ResourceCommandFind}
	disasterConsumer := event.NewDisasterConsumer(kafkaClient, resourceService)

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
