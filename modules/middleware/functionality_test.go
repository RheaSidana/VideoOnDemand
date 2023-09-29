package middleware

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateValidJWTToken() string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "1234567890"
	claims["name"] = "John Doe"
	claims["iat"] = 1516239022

	secretKey := os.Getenv("SECRET")
	tokenString, _ := token.SignedString(
		[]byte(secretKey))

	return "Bearer " + tokenString
}

func TestAuthorizationUtils_IsEmpty(t *testing.T) {
	utils := NewAuthorisation()
	assert.True(
		t, utils.IsEmpty(""), 
		"Expected true for isEmpty(\"\")",
	)
	assert.False(
		t, utils.IsEmpty("some_token"), 
		"Expected false for isEmpty(\"some_token\")",
	)
}

func TestAuthorizationUtils_BearerToken(t *testing.T) {
	utils := NewAuthorisation()
	token := "Bearer some_token"
	splits := utils.BearerToken(token)

	assert.Equal(t, 2, len(splits), "Expected 2 splits")
	assert.Equal(t, "Bearer", splits[0], "Expected 'Bearer'")
	assert.Equal(t, "some_token", splits[1], "Expected 'some_token'")
}

func TestAuthorizationUtils_IsNotBearerToken(t *testing.T) {
	utils := NewAuthorisation()
	assert.True(
		t, utils.IsNotBearerToken("some_token"), 
		"Expected true for IsNotBearerToken(\"some_token\")",
	)
	assert.False(
		t, utils.IsNotBearerToken("Bearer some_token"), 
		"Expected false for IsNotBearerToken(\"Bearer some_token\")",
	)
}

func TestAuthorizationUtils_TokenParse(t *testing.T) {
	os.Setenv("SECRET", "NewToken12390334edefv")
	utils := NewAuthorisation()

	token := generateValidJWTToken()

	parsedToken, err := utils.TokenParse(token)

	assert.NoError(t, err, "Token parsing failed with error")
	assert.NotNil(t, parsedToken, "Parsed token should not be nil")
	assert.True(t, parsedToken.Valid, "Token should be valid")
}

func TestAuthorizationUtils_IsExpiredToken(t *testing.T) {
	utils := NewAuthorisation()
	expiration := float64(time.Now().Add(1 * time.Hour).Unix())
	assert.False(
		t, utils.IsExpiredToken(expiration), 
		"Expected false for IsExpiredToken with future expiration",
	)

	expiration = float64(time.Now().Add(-1 * time.Second).Unix())
	assert.True(
		t, utils.IsExpiredToken(expiration), 
		"Expected true for IsExpiredToken with past expiration",
	)
}
