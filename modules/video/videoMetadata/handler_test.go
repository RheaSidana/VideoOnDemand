package videoMetadata

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"vod/model"
	"vod/modules/middleware"
	authMocks "vod/modules/middleware/mocks"
	"vod/modules/video/videoMetadata/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPlayVideoHandlerWhenUnauthorizedAccess(t *testing.T) {
	redisRepo := new(mocks.RedisRepository)
	repo := new(mocks.Repository)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/play", handler.PlayVideoHandler)
	req, _ := http.NewRequest(
		http.MethodPost, "/protected/video/play", nil)
	respR := httptest.NewRecorder()
	auth.On(
		"IsEmpty", mock.AnythingOfType("string"),
	).Return(true)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"No token provided, Unauthorized Access\"}"
	assert.Equal(t, http.StatusUnauthorized, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func mockAuth(
	auth **authMocks.AuhorisationUtils,
) {
	email := "test@example.com"
	claims := jwt.MapClaims{
		"userEmail": email,
		"exp":       float64(time.Now().Add(12 * time.Second).Unix()),
	}
	token := &jwt.Token{
		Claims: claims,
		Valid:  true,
	}
	(*auth).On(
		"IsEmpty", mock.AnythingOfType("string"),
	).Return(false)
	(*auth).On(
		"IsNotBearerToken", "Bearer your-access-token",
	).Return(false)
	(*auth).On(
		"TokenParse", "Bearer your-access-token",
	).Return(token, nil)
	(*auth).On(
		"IsExpiredToken", mock.AnythingOfType("float64"),
	).Return(false)
}

func TestPlayVideoHandlerWhenBadRequest(t *testing.T) {
	redisRepo := new(mocks.RedisRepository)
	repo := new(mocks.Repository)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.GET("/video/play", handler.PlayVideoHandler)
	req, _ := http.NewRequest(
		http.MethodGet, "/protected/video/play", nil)
	respR := httptest.NewRecorder()
	mockAuth(&auth)

	req.Header.Set("Authorization", "Bearer your-access-token")
	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Bad Request: Unable to play video.\"}"
	assert.Equal(t, http.StatusBadRequest, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestPlayVideoHandlerWhenSuccessfulRedisLookup(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
	}
	r := gin.Default()
	rGroup := r.Group(
		"/protected", authHandler.AuthMiddleware)
	rGroup.POST(
		"/video/play", handler.PlayVideoHandler)
	videoMD := model.VideoMetaData{
		Title:      "Sample Video",
		Format:     "mp4",
		UploadedBy: 1,
		Length:     3.0,
		Size:       3.4,
		Resolution: "23x23",
		Topic:      "Math",
	}
	b, _ := json.Marshal(videoMD)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(
		http.MethodPost, "/protected/video/play", body)
	req.Header.Set("Authorization", "Bearer your-access-token")

	mockAuth(&auth)
	expectedVideoMD := videoMD
	redisRepo.On(
		"GetFromRedis", mock.Anything,
	).Return(expectedVideoMD, nil)

	respR := httptest.NewRecorder()
	r.ServeHTTP(respR, req)

	expectedResponse := "{\"" +
		"Message\":\"Found in Redis, Play the video \"," +
		"\"VideoMD\":{\"" +
		"ID\":0," +
		"\"CreatedAt\":\"0001-01-01T00:00:00Z\"," +
		"\"UpdatedAt\":\"0001-01-01T00:00:00Z\"," +
		"\"DeletedAt\":null," +
		"\"title\":\"Sample Video\"," +
		"\"format\":\"mp4\"," +
		"\"uploadedBy\":1," +
		"\"length\":3," +
		"\"size\":3.4," +
		"\"resolution\":\"23x23\"," +
		"\"topic\":\"Math\"" +
		"}" +
		"}"
	assert.Equal(t, expectedResponse, respR.Body.String())
	assert.Equal(t, http.StatusOK, respR.Code)
}

func TestPlayVideoHandlerWhenRepoThrowsError(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
	}
	r := gin.Default()
	rGroup := r.Group(
		"/protected", authHandler.AuthMiddleware)
	rGroup.POST(
		"/video/play", handler.PlayVideoHandler)
	videoMD := model.VideoMetaData{
		Title:      "Sample Video",
		Format:     "mp4",
		UploadedBy: 1,
		Length:     3.0,
		Size:       3.4,
		Resolution: "23x23",
		Topic:      "Math",
	}
	b, _ := json.Marshal(videoMD)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(
		http.MethodPost, "/protected/video/play", body)
	req.Header.Set("Authorization", "Bearer your-access-token")

	mockAuth(&auth)
	expectedVideoMD := videoMD
	redisRepo.On(
		"GetFromRedis", mock.Anything,
	).Return(expectedVideoMD, errors.New("error occurred"))
	repo.On(
		"Find", mock.Anything,
	).Return(expectedVideoMD, errors.New("Database error"))

	respR := httptest.NewRecorder()
	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Unable to play video.\"}"
	assert.Equal(t, expectedResponse, respR.Body.String())
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
}

func TestPlayVideoHandler(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
	}
	r := gin.Default()
	rGroup := r.Group(
		"/protected", authHandler.AuthMiddleware)
	rGroup.POST(
		"/video/play", handler.PlayVideoHandler)
	videoMD := model.VideoMetaData{
		Title:      "Sample Video",
		Format:     "mp4",
		UploadedBy: 1,
		Length:     3.0,
		Size:       3.4,
		Resolution: "23x23",
		Topic:      "Math",
	}
	b, _ := json.Marshal(videoMD)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(
		http.MethodPost, "/protected/video/play", body)
	req.Header.Set("Authorization", "Bearer your-access-token")

	mockAuth(&auth)
	expectedVideoMD := videoMD
	redisRepo.On(
		"GetFromRedis", mock.Anything,
	).Return(expectedVideoMD, errors.New("error occurred"))
	repo.On(
		"Find", mock.Anything,
	).Return(expectedVideoMD, nil)

	respR := httptest.NewRecorder()
	r.ServeHTTP(respR, req)

	expectedResponse := "{\""+
	"Message\":\"Found, Play the video \","+
	"\"VideoMD\":{\"ID\":0,"+
	"\"CreatedAt\":\"0001-01-01T00:00:00Z\","+
	"\"UpdatedAt\":\"0001-01-01T00:00:00Z\","+
	"\"DeletedAt\":null,"+
	"\"title\":\"Sample Video\","+
	"\"format\":\"mp4\","+
	"\"uploadedBy\":1,"+
	"\"length\":3,"+
	"\"size\":3.4,"+
	"\"resolution\":\"23x23\","+
	"\"topic\":\"Math\"}}"
	assert.Equal(t, expectedResponse, respR.Body.String())
	assert.Equal(t, http.StatusOK, respR.Code)
}
