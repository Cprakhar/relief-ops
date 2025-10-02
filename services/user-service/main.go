package main

import (
	"context"
	"log"
	"os"
	"os/signal"
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
)

var (
	addr           = env.GetString("USER_GRPC_ADDR", ":9001")
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
)

func main() {
	// Set up context with signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

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
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

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
	userRepo := repo.NewUserRepo(mongoClient)
	userService := service.NewUserService(userRepo, jwtSecret, jwtExpiry)

	// Initialize and start the disaster consumer
	topics := []string{events.UserNotifyAdminReview}
	disasterConsumer := event.NewDisasterConsumer(kafkaClient, userService, mailer)

	go func() {
		if err := disasterConsumer.Consumer(ctx, topics); err != nil {
			log.Printf("Error in disaster consumer: %v", err)
		}
	}()

	// Initialize and run the gRPC server
	gRPCServer := newgRPCServer(addr, userService, jwtSecret)

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
