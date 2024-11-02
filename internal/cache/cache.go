package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Get(key string, value interface{}) error {
	data, err := c.client.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return ErrCacheMiss
	} else if err != nil {
		return ErrCacheFailed
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return ErrUnmarshal
	}

	return nil
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return ErrMarshal
	}

	if _, err := c.client.Set(context.Background(), key, data, expiration).Result(); err != nil {
		return ErrCacheFailed
	}

	return nil
}

func (c *Cache) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

func (c *Cache) DeletePattern(pattern string) error {
	keys, err := c.client.Keys(context.Background(), pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(context.Background(), keys...).Err()
	}

	return nil
}
