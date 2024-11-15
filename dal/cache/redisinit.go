package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/kackerx/go-mall/config"
)

var redisClient *redis.Client

func Redis() *redis.Client {
	return redisClient
}

func init() {
	redisConf := config.Conf.Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:         redisConf.Addr,
		Password:     redisConf.Password,
		DB:           redisConf.Db,
		PoolSize:     redisConf.PoolSize,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		PoolTimeout:  10 * time.Second,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
