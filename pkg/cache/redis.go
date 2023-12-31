package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Zavr22/EMTestTask/pkg/models"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{client: client}
}

func (c *RedisClient) SetData(key string, data *models.User, expiration time.Duration) error {
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
