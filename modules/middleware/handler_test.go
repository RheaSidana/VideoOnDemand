package middleware

import (
	"errors"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"vod/modules/middleware/mocks"

	// "github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/http/httptest"
)

func TestAuthMiddlewareWhenIsEmptyTrue(t *testing.T) {
	authMock := new(mocks.AuhorisationUtils)
	handler := Handler{Auth: authMock}
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Add("Authorization", "Bearer ")
	respR := httptest.NewRecorder()
	router := gin.New()

	authMock.On("IsEmpty", "Bearer ").Return(true)
	router.Use(func(c *gin.Context) {
		handler.AuthMiddleware(c)
	})
	router.ServeHTTP(respR, req)

	assert.Equal(
		t, http.StatusUnauthorized, respR.Code,
		"Expected HTTP status Unauthorized",
	)
	assert.Contains(
		t, respR.Body.String(),
		"No token provided", "Expected error message in response",
	)
	authMock.AssertExpectations(t)
}

func TestAuthMiddlewareWhenIsNotBearerTokenTrue(t *testing.T) {
	authMock := new(mocks.AuhorisationUtils)
	handler := Handler{Auth: authMock}
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Add("Authorization", "InvalidToken")
	respR := httptest.NewRecorder()
	router := gin.New()

	authMock.On("IsEmpty", "InvalidToken").Return(false)
	authMock.On("IsNotBearerToken", "InvalidToken").Return(true)
	router.Use(func(c *gin.Context) {
		handler.AuthMiddleware(c)
	})
	router.ServeHTTP(respR, req)

	assert.Equal(
		t, http.StatusUnauthorized, respR.Code,
		"Expected HTTP status Unauthorized",
	)
	assert.Contains(
		t, respR.Body.String(),
		"Invalid token format", "Expected error message in response",
	)
	authMock.AssertExpectations(t)
}

func TestAuthMiddlewareWhenTokenParseFails(t *testing.T) {
	authMock := new(mocks.AuhorisationUtils)
	handler := Handler{Auth: authMock}
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Add("Authorization", "Bearer some_invalid_token")
	respR := httptest.NewRecorder()
	router := gin.New()

	authMock.On(
		"IsEmpty", "Bearer some_invalid_token",
	).Return(false)
	authMock.On(
		"IsNotBearerToken", "Bearer some_invalid_token",
	).Return(false)
	authMock.On(
		"TokenParse", "Bearer some_invalid_token",
	).Return(nil, errors.New("error occured"))

	router.Use(func(c *gin.Context) {
		handler.AuthMiddleware(c)
	})
	router.ServeHTTP(respR, req)

	assert.Equal(
		t, http.StatusUnauthorized, respR.Code,
		"Expected HTTP status Unauthorized",
	)
	assert.Contains(
		t, respR.Body.String(), "Invalid token",
		"Expected error message in response",
	)
	authMock.AssertExpectations(t)
}

func TestAuthMiddlewareWhenClaimsTypeExpFails(t *testing.T) {
	authMock := new(mocks.AuhorisationUtils)
	handler := Handler{Auth: authMock}
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Add("Authorization", "Bearer some_invalid_token")
	respR := httptest.NewRecorder()
	router := gin.New()

	claims := jwt.MapClaims{
		"userEmail": "test@example.com",
		"exp":       float64(time.Now().Add(-1 * time.Second).Unix()),
	}
	token := &jwt.Token{
		Claims: claims,
		Valid:  true,
	}
	authMock.On(
		"IsEmpty", "Bearer some_invalid_token",
	).Return(false)
	authMock.On(
		"IsNotBearerToken", "Bearer some_invalid_token",
	).Return(false)
	authMock.On(
		"TokenParse", "Bearer some_invalid_token",
	).Return(token, nil)
	authMock.On(
		"IsExpiredToken", mock.AnythingOfType("float64"),
	).Return(true)

	router.Use(func(c *gin.Context) {
		handler.AuthMiddleware(c)
	})
	router.ServeHTTP(respR, req)

	assert.Equal(
		t, http.StatusUnauthorized, respR.Code,
		"Expected HTTP status Unauthorized",
	)
	assert.Contains(
		t, respR.Body.String(), "Token has expired",
		"Expected error message in response",
	)
	authMock.AssertExpectations(t)
}

func TestAuthMiddlewareWhenClaimsTypeEmailFails(t *testing.T) {
	authMock := new(mocks.AuhorisationUtils)
	handler := Handler{Auth: authMock}
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Add("Authorization", "Bearer some_invalid_token")
	respR := httptest.NewRecorder()
	router := gin.New()

	claims := jwt.MapClaims{
		"exp": float64(time.Now().Add(12 * time.Second).Unix()),
	}
	token := &jwt.Token{
		Claims: claims,
		Valid:  true,
	}
	authMock.On(
		"IsEmpty", "Bearer some_invalid_token",
	).Return(false)
	authMock.On(
		"IsNotBearerToken", "Bearer some_invalid_token",
	).Return(false)
	authMock.On(
		"TokenParse", "Bearer some_invalid_token",
	).Return(token, nil)
	authMock.On(
		"IsExpiredToken", mock.AnythingOfType("float64"),
	).Return(false)

	router.Use(func(c *gin.Context) {
		handler.AuthMiddleware(c)
	})
	router.ServeHTTP(respR, req)

	assert.Equal(
		t, http.StatusUnauthorized, respR.Code,
		"Expected HTTP status Unauthorized",
	)
	assert.Contains(
		t, respR.Body.String(), "Invalid user information in token",
		"Expected error message in response",
	)
	authMock.AssertExpectations(t)
}

func TestAuthMiddleware(t *testing.T) {
	authMock := new(mocks.AuhorisationUtils)
	handler := Handler{Auth: authMock}
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/protected", handler.AuthMiddleware)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Add("Authorization", "Bearer some_invalid_token")
	respR := httptest.NewRecorder()

	email := "test@example.com"
	claims := jwt.MapClaims{
		"userEmail": email,
		"exp":       float64(time.Now().Add(12 * time.Second).Unix()),
	}
	token := &jwt.Token{
		Claims: claims,
		Valid:  true,
	}
	authMock.On(
		"IsEmpty", "Bearer some_invalid_token",
	).Return(false)
	authMock.On(
		"IsNotBearerToken", "Bearer some_invalid_token",
	).Return(false)
	authMock.On(
		"TokenParse", "Bearer some_invalid_token",
	).Return(token, nil)
	authMock.On(
		"IsExpiredToken", mock.AnythingOfType("float64"),
	).Return(false)

	router.Use(func(c *gin.Context) {
		handler.AuthMiddleware(c)

		// Assert that c.Get("userEmail") is equal to "test@example.com".
		// userEmail, exists := req.Context().Value("userEmail").(string)
		userEmail, exists := c.Get("userEmail")
		assert.True(
			t, exists,
			"Expected 'userEmail' key in context",
		)
		assert.Equal(
			t, email, userEmail,
			"Expected userEmail to be 'test@example.com'",
		)
	})
	router.ServeHTTP(respR, req)

	assert.Equal(
		t, http.StatusOK, respR.Code,
		"Expected HTTP status OK",
	)
	assert.Empty(
		t, respR.Body.String(),
		"Expected empty response body",
	)
	authMock.AssertExpectations(t)
}
