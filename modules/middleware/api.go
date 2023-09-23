package middleware

import "github.com/gin-gonic/gin"

func Apis(r *gin.Engine) (*gin.RouterGroup){
	return r.Group("/protected", AuthMiddleware)
}