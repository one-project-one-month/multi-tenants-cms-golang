package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisClient *redis.Client

func InitRedis(addr, password string, db int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func SetRedis(key string, value string, expiration time.Duration) error {
	return redisClient.Set(ctx, key, value, expiration).Err()
}

func GetRedis(key string) (string, error) {
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func GetClient() *redis.Client {
	return redisClient
}

func DeleteRedis(email string) error {
	return redisClient.Del(ctx, email).Err()
}

func RevokeJWTToken(tokenID string, expiration time.Duration) error {
	return redisClient.Set(ctx, "revoked:"+tokenID, "true", expiration).Err()
}

func IsTokenRevoked(tokenID string) (bool, error) {
	val, err := redisClient.Get(ctx, "revoked:"+tokenID).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "true", nil
}

func QueueRevokeToken(tokenID string, expiration time.Duration) error {
	revokeData := map[string]interface{}{
		"token_id":   tokenID,
		"revoked_at": time.Now().Unix(),
		"expires_at": time.Now().Add(expiration).Unix(),
	}

	jsonData, err := json.Marshal(revokeData)
	if err != nil {
		return err
	}

	return redisClient.LPush(ctx, "jwt_revoke_queue", jsonData).Err()
}

func ProcessRevokeQueue() (map[string]interface{}, error) {
	result, err := redisClient.BRPop(ctx, 0, "jwt_revoke_queue").Result()
	if err != nil {
		return nil, err
	}

	var revokeData map[string]interface{}
	err = json.Unmarshal([]byte(result[1]), &revokeData)
	if err != nil {
		return nil, err
	}

	return revokeData, nil
}

func BulkRevokeTokens(tokenIDs []string, expiration time.Duration) error {
	pipe := redisClient.Pipeline()

	for _, tokenID := range tokenIDs {
		pipe.Set(ctx, "revoked:"+tokenID, "true", expiration)
	}

	_, err := pipe.Exec(ctx)
	return err
}
