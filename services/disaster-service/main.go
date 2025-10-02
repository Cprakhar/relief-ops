package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cprakhar/relief-ops/services/disaster-service/repo"
	"github.com/cprakhar/relief-ops/services/disaster-service/service"
	"github.com/cprakhar/relief-ops/shared/db"
	"github.com/cprakhar/relief-ops/shared/env"
	"github.com/cprakhar/relief-ops/shared/messaging"
)

var (
	addr         = env.GetString("DISASTER_GRPC_ADDR", ":9002")

	// Kafka configuration
	brokers      = env.GetString("KAFKA_BROKERS", "apache-kafka:9092")

	// MongoDB configuration
	mongoURI     = env.GetString("MONGODB_URI", "")
	mongoTimeout = 30 * time.Second
	mongoMaxIdle = 5 * time.Second
	mongoMaxPool = uint64(10)
	mongoMinPool = uint64(2)
)

func main() {
	// Set up context with signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize MongoDB client
	mongoCfg := &db.MongoDBConfig{
		URI:        mongoURI,
		Database:   "relief_ops",
		Collection: "disasters",
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
		GroupID: "disaster-service-group",
	}

	kafkaClient, err := messaging.NewKafkaClient("disaster-service", kafkaCfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer kafkaClient.Close()
	log.Println("Kafka client initialized")

	// Initialize repository and service
	userRepo := repo.NewMongodbDisasterRepo(mongoClient)
	userService := service.NewDisasterService(userRepo)

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
