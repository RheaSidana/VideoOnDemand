package initializer

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var RedisDB *redis.Client
var Ctx context.Context

func NewRedisClient(ctx context.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password by default
		DB:       0,  // Default DB
	})

	// Check if Redis is reachable
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	RedisDB = rdb
	Ctx = ctx
}
