package videoMetadata

import (
	"vod/initializer"

	"github.com/gin-gonic/gin"
)

func Apis(r *gin.RouterGroup) {
	repository := InitRepository(initializer.Db)
	redisRepository := InitRedisRepository(initializer.RedisDB, initializer.Ctx)
	// videoMDHandler
	videoMDHandler := InitHandler(repository, redisRepository)

	r.GET("/video/play", videoMDHandler.PlayVideoHandler)
}
