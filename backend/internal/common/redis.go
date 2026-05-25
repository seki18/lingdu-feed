package common

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seki18/lingdu-feed/config"
)

// Redis is the shared Redis client singleton.
var Redis *redis.Client

// InitRedis initializes the Redis client.
func InitRedis(cfg config.Config) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Redis.Ping(ctx).Err(); err != nil {
		log.Printf("[Redis] Connection failed (degraded mode): %v", err)
		Redis = nil
	} else {
		log.Println("[Redis] Connected successfully")
	}
}
