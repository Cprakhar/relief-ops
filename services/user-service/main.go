package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cprakhar/relief-ops/services/user-service/event"
	"github.com/cprakhar/relief-ops/services/user-service/mail"
	"github.com/cprakhar/relief-ops/services/user-service/repo"
	"github.com/cprakhar/relief-ops/services/user-service/service"
	"github.com/cprakhar/relief-ops/shared/db"
	"github.com/cprakhar/relief-ops/shared/env"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
	"github.com/cprakhar/relief-ops/shared/observe/logs"
	"github.com/cprakhar/relief-ops/shared/observe/traces"
)

var (
	addr           = env.GetString("USER_GRPC_ADDR", ":9001")
	environment    = env.GetString("ENVIRONMENT", "development")
	webURL         = env.GetString("WEB_URL", "http://localhost:3000")
	fromEmail      = env.GetString("FROM_EMAIL", "developerluffy23@gmail.com")
	sendGridAPIKey = env.GetString("SENDGRID_API_KEY", "")
	brokers        = env.GetString("KAFKA_BROKERS", "apache-kafka:9092")

	// JWT configuration
	jwtSecret = env.GetString("JWT_SECRET", "")
	jwtExpiry = env.GetTimeDuration("JWT_EXPIRY", time.Hour*24*7) // 7 days

	// Redis configuration
	redisAddr     = env.GetString("REDIS_ADDR", "redis-db:6379")
	redisUsername = env.GetString("REDIS_USERNAME", "")
	redisMaxConn  = env.GetInt("REDIS_MAX_CONN", 10)
	redisMinIdle  = env.GetInt("REDIS_MIN_IDLE", 2)
	redisMaxIdle  = env.GetInt("REDIS_MAX_IDLE", 5)
	redisPassword = env.GetString("REDIS_PASSWORD", "")
	redisDB       = env.GetInt("REDIS_DB", 0)

	// MongoDB configuration
	mongoURI     = env.GetString("MONGODB_URI", "")
	mongoDB      = env.GetString("MONGODB_DB", "relief_ops")
	mongoMaxIdle = env.GetTimeDuration("MONGODB_MAX_IDLE", 5*time.Second)
	mongoMaxPool = uint64(env.GetInt("MONGODB_MAX_POOL", 10))
	mongoMinPool = uint64(env.GetInt("MONGODB_MIN_POOL", 2))
	mongoTimeout = env.GetTimeDuration("MONGODB_TIMEOUT", 30*time.Second)

	// OTLP configuration
	otlpEndpoint = env.GetString("OTLP_ENDPOINT", "otel-collector:4317")
	otlpInsecure = env.GetBool("OTLP_INSECURE", true)
)

func main() {
	// Initialize logger
	logger, err := logs.Init("user-service")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Errorw("Failed to sync logger", "error", err)
		}
	}()
	logger.Info("Logger initialized")

	// Set up context with signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tracerCfg := &traces.TracerConfig{
		ServiceName:      "user-service",
		Environment:      environment,
		ExporterEndpoint: otlpEndpoint,
		Secure:           !otlpInsecure,
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

	redisCfg := &db.RedisConfig{
		Addr:           redisAddr,
		Username:       redisUsername,
		Password:       redisPassword,
		DB:             int(redisDB),
		MaxActiveConns: int(redisMaxConn),
		MaxIdleConns:   int(redisMaxIdle),
		MinIdleConns:   int(redisMinIdle),
	}
	if err := db.InitRedis(redisCfg); err != nil {
		logger.Fatalw("Failed to connect to Redis", "error", err)
	}
	logger.Info("Connected to Redis")

	// Initialize MongoDB client
	mongoCfg := &db.MongoDBConfig{
		URI:        mongoURI,
		Database:   mongoDB,
		Collection: "users",
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
		GroupID: "user-service-group",
	}

	kafkaClient, err := messaging.NewKafkaClient("user-service", kafkaCfg)
	if err != nil {
		logger.Fatalw("Failed to create Kafka client", "error", err)
	}
	defer kafkaClient.Close()
	logger.Info("Kafka client initialized")

	mailer := mail.NewSendGrid(fromEmail, sendGridAPIKey)

	// Initialize repository and service
	userRepo, err := repo.NewUserRepo(ctx, mongoClient)
	if err != nil {
		logger.Fatalw("Failed to create user repository", "error", err)
	}
	userService := service.NewUserService(userRepo, jwtSecret, jwtExpiry)

	// Initialize and start the disaster consumer
	topics := []string{events.UserNotifyAdminReview}
	disasterConsumer := event.NewDisasterConsumer(kafkaClient, userService, mailer, webURL)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		if err := disasterConsumer.Consumer(ctx, topics); err != nil {
			logger.Errorw("Error in disaster consumer", "error", err)
		}
	}()

	// Initialize and run the gRPC server
	gRPCServer := newgRPCServer(addr, userService, jwtSecret)
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Infow("User service running", "addr", addr)
		if err := gRPCServer.run(ctx); err != nil {
			logger.Errorw("gRPC server error", "error", err)
		}
	}()
	<-ctx.Done()
	wg.Wait()
	logger.Info("User service stopped")
}
