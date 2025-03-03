// Using Redis as cache backend
package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	rdb *redis.Client
	ctx context.Context
}

// Establish connection to Redis server
func New() *CacheClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	return &CacheClient{
		rdb: rdb,
		ctx: context.Background(),
	}
}

// Get cached value
func (c *CacheClient) Get(key string) (string, error) {
	return c.rdb.Get(c.ctx, key).Result()
}

// Get from cache, if available. If noe available (value = ""), execute given function,
// save function retrun value in cache and return this value.
func (c *CacheClient) CachedGet(key string, f func() (string, error), expiration time.Duration) (string, error) {
	value, err := c.Get(key)

	// Cache not found, set new cache value
	if err == redis.Nil {
		value, err = f()
		if err != nil {
			return "", fmt.Errorf("failed to get cache: %v", err)
		}

		if err = c.Set(key, value, expiration); err != nil {
			return "", err
		}

		return value, err
	}

	if err != nil {
		return "", fmt.Errorf("failed to get cache: %v", err)
	}

	// Cache found
	return value, nil
}

func (c *CacheClient) Set(key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(c.ctx, key, value, expiration).Err() // 0 means no expiry
}
