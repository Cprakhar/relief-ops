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
	"github.com/cprakhar/relief-ops/shared/observe/logs"
	"github.com/cprakhar/relief-ops/shared/observe/traces"
)

var (
	addr        = env.GetString("DISASTER_GRPC_ADDR", ":9002")
	environment = env.GetString("ENVIRONMENT", "development")

	// Kafka configuration
	brokers = env.GetString("KAFKA_BROKERS", "apache-kafka:9092")

	// MongoDB configuration
	mongoURI     = env.GetString("MONGODB_URI", "")
	mongoTimeout = 30 * time.Second
	mongoMaxIdle = 5 * time.Second
	mongoMaxPool = uint64(10)
	mongoMinPool = uint64(2)

	// OTLP configuration
	otlpEndpoint = env.GetString("OTLP_ENDPOINT", "otel-collector:4317")
	otlpInsecure = env.GetBool("OTLP_INSECURE", true)
)

func main() {
	// Initialize logger
	logger, err := logs.Init("disaster-service")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()
	logger.Info("Logger initialized")

	// Set up context with signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tracerCfg := &traces.TracerConfig{
		ServiceName:      "disaster-service",
		Environment:      environment,
		Secure:           !otlpInsecure,
		ExporterEndpoint: otlpEndpoint,
	}

	shutdown, err := traces.InitTrace(ctx, tracerCfg)
	if err != nil {
		logger.Fatalw("Failed to initialize tracing", "error", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			logger.Fatalw("Failed to shutdown tracer provider", "error", err)
		}
	}()
	logger.Info("Tracing initialized")

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
		logger.Fatalw("Failed to connect to MongoDB", "error", err)
	}
	logger.Info("Connected to MongoDB")

	// Initialize Kafka client
	kafkaCfg := &messaging.KafkaConfig{
		Brokers: brokers,
		GroupID: "disaster-service-group",
	}

	kafkaClient, err := messaging.NewKafkaClient("disaster-service", kafkaCfg)
	if err != nil {
		logger.Fatalw("Failed to create Kafka client", "error", err)
	}
	defer kafkaClient.Close()
	logger.Info("Kafka client initialized")

	// Initialize repository and service
	userRepo, err := repo.NewMongodbDisasterRepo(ctx, mongoClient)
	if err != nil {
		logger.Fatalw("Failed to create disaster repository", "error", err)
	}
	userService := service.NewDisasterService(userRepo)

	// Initialize and run the gRPC server
	gRPCServer := newgRPCServer(addr, userService, kafkaClient)

	done := make(chan struct{})
	go func() {
		defer close(done)
		logger.Infow("Disaster service running", "addr", addr)
		if err := gRPCServer.run(ctx); err != nil {
			logger.Errorw("gRPC server error", "error", err)
		}
	}()
	<-ctx.Done()
	<-done
	logger.Info("Disaster service stopped")
}
