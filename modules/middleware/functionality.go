package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuhorisationUtils interface{
	IsEmpty(token string) bool
	BearerToken(token string) []string
	IsNotBearerToken(token string) bool 
	TokenParse(token string) (*jwt.Token, error) 
	IsExpiredToken(exp float64) bool 	
}

type authorizationUtils struct{}

func NewAuthorisation() AuhorisationUtils {
	return &authorizationUtils{}
}

func (a *authorizationUtils) IsEmpty(token string) bool {
	return token == ""
}

func (a *authorizationUtils) BearerToken(token string) []string {
	return strings.Split(token, " ")
}

func (a *authorizationUtils) IsNotBearerToken(token string) bool {
	splits := a.BearerToken(token)

	return len(splits) != 2 || splits[0] != "Bearer"
}

func (a *authorizationUtils) TokenParse(token string) (*jwt.Token, error) {
	splits := a.BearerToken(token)
	return jwt.Parse(splits[1], func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}

		var secretKey string = (os.Getenv("SECRET"))
		return []byte(secretKey), nil
	})
}

func (a *authorizationUtils) IsExpiredToken(exp float64) bool {
	return int64(exp) < time.Now().Unix()
}
