package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"

	"github.com/go-redis/redis/v8"
)

type redisClient interface {
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
	io.Closer
}

var (
	gRedisAddrs = flag.String("redis_addrs", "localhost:6379", "comma separated redis addresses")
	gRedisDb    = flag.Int("redis_db", 0, "redis db")
	gRedisPwd   = flag.String("redis_pwd", "", "redis password")
)

func newRedisClientFromArgs() (redisClient, error) {
	addrs := strings.Split(*gRedisAddrs, ",")
	if len(addrs) == 0 {
		return nil, errors.New("no redis address")
	}
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    addrs,
		DB:       *gRedisDb,
		Password: *gRedisPwd,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
