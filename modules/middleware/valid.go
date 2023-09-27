package middleware

import "github.com/gin-gonic/gin"

func IsAuthorisedAccess(c *gin.Context) bool {
	_, exists := c.Get("userEmail")
	return exists
}
