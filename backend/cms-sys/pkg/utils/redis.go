package utils

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	host := GetEnv("REDIS_HOST", "localhost")
	port := GetEnv("REDIS_PORT", "6379")
	password := GetEnv("REDIS_PASSWORD", "")
	dbStr := GetEnv("REDIS_DB", "0")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return fmt.Errorf("invalid REDIS_DB value: %v", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}

func RevokeToken(tokenID string, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("revoked_token:%s", tokenID)
	return RedisClient.Set(ctx, key, "revoked", expiration).Err()
}

func IsTokenRevoked(tokenID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("revoked_token:%s", tokenID)
	result := RedisClient.Get(ctx, key)

	if result.Err() == redis.Nil {
		return false, nil
	}

	if result.Err() != nil {
		return false, result.Err()
	}

	return true, nil
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
