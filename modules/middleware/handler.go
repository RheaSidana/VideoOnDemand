package middleware

import (
	// "fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	Auth AuhorisationUtils
}

func (h *Handler) AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if h.Auth.IsEmpty(tokenString) {
		c.JSON(401, ErrorResponse{
			Message: "No token provided, Unauthorized Access",
		})
		return
	}

	if h.Auth.IsNotBearerToken(tokenString) {
		c.JSON(401, ErrorResponse{
			Message: "Invalid token format, Unauthorized Access",
		})
		return
	}

	token, err := h.Auth.TokenParse(tokenString)
	if err != nil || !token.Valid {
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
	if !ok || h.Auth.IsExpiredToken(exp) {
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

	// fmt.Println("userEmail", userEmail)

	c.Set("userEmail", userEmail)
	// val, _ := c.Get("userEmail")
	// fmt.Println("userEmail", val)

	c.Next()
}
