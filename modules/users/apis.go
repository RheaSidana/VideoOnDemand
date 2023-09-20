package users

import (
	"vod/initializer"

	"github.com/gin-gonic/gin"
)

func Apis(r *gin.Engine) {
	repository := InitRepository(initializer.Db)
	redisRepository := InitRedisRepository(initializer.RedisDB, initializer.Ctx)
	userHandler := InitHandler(repository, redisRepository)

	r.POST("/user/signup", userHandler.SignUpHandler)
	r.POST("/user/login", userHandler.LoginHandler)
}
