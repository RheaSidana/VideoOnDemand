package videoEncoding

import (
	"vod/initializer"

	"github.com/gin-gonic/gin"
)

func Apis(r *gin.RouterGroup) {
	repository := InitRepository(initializer.Db)
	redisRepository := InitRedisRepository(initializer.RedisDB, initializer.Ctx)
	handler := InitHandler(repository, redisRepository)

	r.POST("/video/encode", handler.VideoEncodeHandler)
}
