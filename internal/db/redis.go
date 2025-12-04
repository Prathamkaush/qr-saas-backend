package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedis(url string) *redis.Client {
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opts)

	// simple ping
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return client
}
