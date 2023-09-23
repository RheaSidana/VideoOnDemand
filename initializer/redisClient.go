package initializer

import (
	"context"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var RedisDB *redis.Client
var Ctx context.Context

func NewRedisClient(ctx context.Context) {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_DB_ADDR"),
		Password: os.Getenv("REDIS_DB_PASSWORD"),
		DB:       db,  // Default DB
	})

	// Check if Redis is reachable
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	RedisDB = rdb
	Ctx = ctx
}
