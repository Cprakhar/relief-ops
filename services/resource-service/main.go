package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cprakhar/relief-ops/services/resource-service/event"
	"github.com/cprakhar/relief-ops/services/resource-service/repo"
	"github.com/cprakhar/relief-ops/services/resource-service/service"
	"github.com/cprakhar/relief-ops/shared/db"
	"github.com/cprakhar/relief-ops/shared/env"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
	"github.com/cprakhar/relief-ops/shared/observe/logs"
	"github.com/cprakhar/relief-ops/shared/observe/traces"
)

var (
	addr        = env.GetString("RESOURCE_GRPC_ADDR", ":9003")
	environment = env.GetString("ENVIRONMENT", "development")

	// Kafka configuration
	brokers = env.GetString("KAFKA_BROKERS", "apache-kafka:9092")

	// MongoDB configuration
	mongoURI     = env.GetString("MONGODB_URI", "")
	mongoDB      = env.GetString("MONGODB_DB", "relief_ops")
	mongoTimeout = env.GetTimeDuration("MONGODB_TIMEOUT", 30*time.Second)
	mongoMaxIdle = env.GetTimeDuration("MONGODB_MAX_IDLE", 5*time.Second)
	mongoMaxPool = uint64(env.GetInt("MONGODB_MAX_POOL", 10))
	mongoMinPool = uint64(env.GetInt("MONGODB_MIN_POOL", 2))

	// OTLP configuration
	otlpEndpoint = env.GetString("OTLP_ENDPOINT", "otel-collector:4317")
	otlpInsecure = env.GetBool("OTLP_INSECURE", true)
)

func main() {
	// Initialize logger
	logger, err := logs.Init("resource-service")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatalw("Failed to sync logger", "error", err)
		}
	}()
	logger.Info("Logger initialized")

	// Set up context with signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tracerCfg := &traces.TracerConfig{
		ServiceName:      "resource-service",
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
		Database:   mongoDB,
		Collection: "resources",
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
		GroupID: "resource-service-group",
	}

	kafkaClient, err := messaging.NewKafkaClient("resource-service", kafkaCfg)
	if err != nil {
		logger.Fatalw("Failed to create Kafka client", "error", err)
	}
	defer kafkaClient.Close()
	logger.Info("Kafka client initialized")

	resourceRepo, err := repo.NewResourceRepo(ctx, mongoClient)
	if err != nil {
		logger.Fatalw("Failed to create resource repository", "error", err)
	}
	resourceService := service.NewResourceService(resourceRepo)

	// Initialize and start the disaster consumer
	topics := []string{events.ResourceCommandFind}
	disasterConsumer := event.NewDisasterConsumer(kafkaClient, resourceService)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := disasterConsumer.DisasterConsumer(ctx, topics); err != nil {
			logger.Errorw("Error in disaster consumer", "error", err)
		}
	}()

	gRPCServer := newgRPCServer(addr, resourceService, kafkaClient)
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Infow("Disaster service running", "addr", addr)
		if err := gRPCServer.run(ctx); err != nil {
			logger.Errorw("gRPC server error", "error", err)
		}
	}()
	<-ctx.Done()
	wg.Wait()
	logger.Info("Resource service stopped")
}
