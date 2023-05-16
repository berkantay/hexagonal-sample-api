package db

import (
	"context"
	"fmt"

	"github.com/berkantay/firefly-weather-condition-api/config"
	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	client *redis.Client
}

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

func (rdb *RedisDatabase) Client() *redis.Client {
	return rdb.client
}
