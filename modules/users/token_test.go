package users

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	// Set the secret key in the environment.
	os.Setenv("SECRET", "your_secret_key_here")

	// Initialize the token service with a user email.
	userEmail := "test@example.com"
	tokenService := NewToken(userEmail)

	// Generate a token.
	tokenString, err := tokenService.GenerateToken()
	assert.NoError(t, err, "Expected no error when generating a token")
	assert.NotEmpty(t, tokenString, "Expected a non-empty token string")

	// Parse the token to verify its claims.
	token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method and return the secret key.
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	assert.NoError(t, parseErr, "Expected no error when parsing the token")

	// Verify the token claims.
	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok, "Expected token claims to be of type jwt.MapClaims")
	assert.Equal(t, userEmail, claims["userEmail"], "Expected userEmail claim to match")
	exp, ok := claims["exp"].(float64)
	assert.True(t, ok, "Expected exp claim to be of type float64")
	expTime := time.Unix(int64(exp), 0)
	assert.WithinDuration(t, time.Now().Add(time.Hour*24), expTime, time.Second,
		"Expected exp claim to be roughly 24 hours from now")
}
