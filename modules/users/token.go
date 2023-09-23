package users

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func generateToken(userEmail string) (string, error) {
	var secretKey string = (os.Getenv("SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userEmail"] = userEmail
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}