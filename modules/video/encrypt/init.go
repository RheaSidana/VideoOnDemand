package videoEncryption

import (
	"context"
	videoEncoding "vod/modules/video/encoding"
	"vod/modules/video/videoMetadata"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InitRepository(client *gorm.DB) Repository {
	return NewRepository(client)
}

func InitRedisRepository(client *redis.Client, ctx context.Context) RedisRepository {
	return NewRedisRepository(client, ctx)
}

func InitHandler(
	userRepository Repository, 
	userRedisRepository RedisRepository,
	vmdRedisRepo videoMetadata.RedisRepository,
	vencodedRedisRepo videoEncoding.RedisRepository,

) Handler {
	return Handler{
		repository: userRepository,
		redisRepository: userRedisRepository,
		redisRepoVideoMD: vmdRedisRepo,
		redisRepoVideoEncoded: vencodedRedisRepo,
		videoEncrypt: NewEncryption(),
	}

}
