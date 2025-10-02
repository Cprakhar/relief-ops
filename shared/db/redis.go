package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr           string
	Username       string
	Password       string
	DB             int
	MaxIdleConns   int
	MaxActiveConns int
	MinIdleConns   int
}

var (
	globalRedisClient *redis.Client
)

func InitRedis(cfg *RedisConfig) error {
	client := redis.NewClient(&redis.Options{
		Addr:           cfg.Addr,
		Username:       cfg.Username,
		Password:       cfg.Password,
		DB:             cfg.DB,
		MinIdleConns:   cfg.MinIdleConns,
		MaxActiveConns: cfg.MaxActiveConns,
		MaxIdleConns:   cfg.MaxIdleConns,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := client.Ping(ctx)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	globalRedisClient = client
	return nil
}

func GetRedisClient() *redis.Client {
	if globalRedisClient == nil {
		panic("Redis client is not initialized. Call InitRedis first.")
	}
	return globalRedisClient
}
