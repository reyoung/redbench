package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/reyoung/piperedis"
)

type redisClient interface {
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
	io.Closer
}

var (
	gRedisAddrs               = flag.String("redis_addrs", "localhost:6379", "comma separated redis addresses")
	gRedisDb                  = flag.Int("redis_db", 0, "redis db")
	gRedisPwd                 = flag.String("redis_pwd", "", "redis password")
	gRedisUseAutoPipe         = flag.Bool("redis_use_auto_pipe", false, "use auto pipe")
	gRedisAutoPipeNWorkers    = flag.Int("redis_auto_pipe_n_workers", 4, "number of workers for auto pipe")
	gRedisAutoPipeBufSize     = flag.Int("redis_auto_pipe_buf_size", 32, "min size for auto pipe")
	gRedisAutoPipeMaxInterval = flag.Duration("redis_auto_pipe_worker_timeout", time.Millisecond,
		"timeout for auto pipe worker")
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
	if !*gRedisUseAutoPipe {
		return client, nil
	}

	return piperedis.New(client, piperedis.Option{
		NumBackgroundWorker: *gRedisAutoPipeNWorkers,
		ChannelBufferSize:   *gRedisAutoPipeBufSize,
		MinCollectInterval:  *gRedisAutoPipeMaxInterval,
	})
}
