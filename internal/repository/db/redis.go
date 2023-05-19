package db

import (
	"context"
	"fmt"
	"time"

	"github.com/berkantay/firefly-weather-condition-api/config"
	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/redis/go-redis/v9"
)

// RedisDatabase is responsible for interacting with a Redis database.
type RedisDatabase struct {
	client *redis.Client
}

// NewRedisStorage creates a new instance of RedisDatabase.
func NewRedisStorage(config *config.Config) (*RedisDatabase, error) {
	hostUrl := fmt.Sprintf("%s:%s",
		config.Redis.Host,
		config.Redis.Port,
	)

	rdb := redis.NewClient(&redis.Options{
		Addr:     hostUrl,
		Password: config.Redis.Password, // no password set
		DB:       0,                     // use default DB
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisDatabase{
		client: rdb,
	}, nil
}

// Set stores the given cache data in Redis with a specified time-to-live (TTL).
func (rdb *RedisDatabase) Set(ctx context.Context, cacheWeather *domain.Cache, ttl time.Duration) (*domain.Cache, error) {
	status := rdb.client.Set(ctx, cacheWeather.Key, cacheWeather.Value, ttl)
	if status.Err() != nil {
		return nil, status.Err()
	}
	return cacheWeather, nil
}

// Get retrieves the value associated with the specified key from Redis.
func (rdb *RedisDatabase) Get(ctx context.Context, key string) (*domain.Cache, error) {
	v, err := rdb.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return &domain.Cache{
		Key:   key,
		Value: v,
	}, nil
}

// Exists checks if a key exists in Redis.
func (rdb *RedisDatabase) Exists(ctx context.Context, key string) bool {
	exists, err := rdb.client.Exists(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists != 0
}
