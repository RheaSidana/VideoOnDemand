package videoEncoding

import (
	"vod/initializer"

	"github.com/gin-gonic/gin"
)

func Apis(r *gin.RouterGroup) {
	repository := InitRepository(initializer.Db)
	redisRepository := InitRedisRepository(initializer.RedisDB, initializer.Ctx)
	// videoEncodingHandler
	handler := InitHandler(repository, redisRepository)

	// r.GET("/video/play", videoMDHandler.PlayVideoHandler)
	r.POST("/video/encode", handler.VideoEncodeHandler)
	r.POST("video/temp", handler.tempForConcat)
}
