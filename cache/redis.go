package cache

import (
	"EMTestTask/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{client: client}
}

func (c *RedisClient) SetData(key string, data *model.User, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshal: %v", err)
	}
	err = c.client.Set(context.Background(), key, jsonData, expiration).Err()
	if err != nil {
		return fmt.Errorf("error set: %v", err)
	}

	return nil
}
