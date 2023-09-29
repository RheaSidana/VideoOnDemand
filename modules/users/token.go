package users

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IToken interface {
	GenerateToken() (string, error)
}

type tokenService struct{
	userEmail string
}

func NewToken(userEmail string) IToken{
	return &tokenService{userEmail: userEmail}
}

func (t *tokenService) GenerateToken() (string, error) {
	var secretKey string = (os.Getenv("SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userEmail"] = t.userEmail
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
