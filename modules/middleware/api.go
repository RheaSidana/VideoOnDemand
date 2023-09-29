package middleware

import "github.com/gin-gonic/gin"

func Apis(r *gin.Engine) (*gin.RouterGroup){
	handler := Handler{
		auth: NewAuthorisation(),
	}
	return r.Group("/protected", handler.AuthMiddleware)
}