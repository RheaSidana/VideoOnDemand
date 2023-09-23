package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if isEmpty(tokenString) {
		c.JSON(401, ErrorResponse{
			Message: "No token provided, Unauthorized Access",
		})
		return
	}

	if isNotBearerToken(tokenString) {
		c.JSON(401, ErrorResponse{
			Message: "Invalid token format, Unauthorized Access",
		})
		return
	}

	token, err := tokenParse(tokenString)
	if err != nil || !token.Valid{
		c.JSON(401, ErrorResponse{
			Message: "Invalid token, Unauthorized Access",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(401, ErrorResponse{
			Message: "Failed to extract user information",
		})
		return
	}

	exp, ok := claims["exp"].(float64)
	if !ok || isExpiredToken(exp) {
		c.JSON(401, ErrorResponse{
			Message: "Token has expired",
		})
		return
	}

	userEmail, ok := claims["userEmail"].(string)
	if !ok {
		c.JSON(401, ErrorResponse{
			Message: "Invalid user information in token",
		})
		return
	}

	c.Set("userEmail", userEmail)
	// val, _ := c.Get("userEmail")
	// fmt.Println("userEmail", val)
	
	c.Next()
}
