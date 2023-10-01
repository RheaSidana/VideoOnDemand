package videoEncryption

import (
	"vod/initializer"
	videoEncoding "vod/modules/video/encoding"
	"vod/modules/video/videoMetadata"

	"github.com/gin-gonic/gin"
)

func Apis(r *gin.RouterGroup) {
	repository := InitRepository(initializer.Db)
	redisRepository := InitRedisRepository(initializer.RedisDB, initializer.Ctx)
	redisRepositoryVMD := videoMetadata.InitRedisRepository(initializer.RedisDB, initializer.Ctx)
	redisRepositoryVEncoded := videoEncoding.InitRedisRepository(initializer.RedisDB, initializer.Ctx)
	// handler
	handler := InitHandler(
		repository, 
		redisRepository,
		redisRepositoryVMD,
		redisRepositoryVEncoded,
	)

	r.POST("/video/encrypt", handler.VideoEncryptHandler)
}
