package videoEncryption

import (
	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InitRepository(client *gorm.DB) Repository {
	return NewRepository(client)
}

func InitRedisRepository(client *redis.Client, ctx context.Context) RedisRepository {
	return NewRedisRepository(client, ctx)
}

func InitHandler(userRepository Repository, userRedisRepository RedisRepository) Handler {
	return Handler{
		repository: userRepository,
		redisRepository: userRedisRepository,
	}
}
