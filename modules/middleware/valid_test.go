package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.New()
	return r
}

func TestIsAuthorisedAccess_Authorized(t *testing.T) {
	r := setupRouter()

	r.GET("/protected", func(c *gin.Context) {
		c.Set("userEmail", "test@example.com")
		isAuthorized := IsAuthorisedAccess(c)
		assert.True(t, isAuthorized, "Expected IsAuthorisedAccess to return true")
		c.String(http.StatusOK, "Authorized")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status OK")
	assert.Equal(t, "Authorized", w.Body.String(), "Expected 'Authorized' user")
}

func TestIsAuthorisedAccess_Unauthorized(t *testing.T) {
	r := setupRouter()

	r.GET("/protected", func(c *gin.Context) {
		isAuthorized := IsAuthorisedAccess(c)
		assert.False(t, isAuthorized, "Expected IsAuthorisedAccess to return false")
		c.String(http.StatusOK, "Unauthorized")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status OK")
	assert.Equal(t, "Unauthorized", w.Body.String(), "Expected 'Unauthorized' user")
}

