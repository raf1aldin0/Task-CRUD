package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// InitRedis mengembalikan *redis.Client dan error
func InitRedis(cfg *Config) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr: cfg.RedisHost + ":" + cfg.RedisPort,
        DB:   0,
    })

    _, err := client.Ping(context.Background()).Result()
    if err != nil {
        log.Printf("Failed to connect to Redis: %v", err)
        return nil, err
    }

    log.Println("Connected to Redis successfully")
    return client, nil
}
