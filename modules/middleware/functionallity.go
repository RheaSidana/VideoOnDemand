package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func isEmpty(token string) bool {
	return token == ""
}

func bearerToken(token string) []string {
	return strings.Split(token, " ")
}

func isNotBearerToken(token string) bool {
	splits := bearerToken(token)

	return len(splits) != 2 || splits[0] != "Bearer"
}

func tokenParse(token string) (*jwt.Token, error) {
	splits := bearerToken(token)
	return jwt.Parse(splits[1], func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}

		var secretKey string = (os.Getenv("SECRET"))
		return []byte(secretKey), nil
	})
}

func isExpiredToken(exp float64) bool {
	return int64(exp) < time.Now().Unix()
}
