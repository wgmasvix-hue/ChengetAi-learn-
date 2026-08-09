package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func New(redisURL string) (*goredis.Client, error) {
	opts, err := goredis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := goredis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return client, nil
}

func Ping(ctx context.Context, client *goredis.Client) error {
	if client == nil {
		return errors.New("redis client is nil")
	}
	return client.Ping(ctx).Err()
}
