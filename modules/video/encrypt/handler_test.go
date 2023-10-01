package videoEncryption

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"vod/modules/middleware"
	authMocks "vod/modules/middleware/mocks"
	videoEncodingMocks "vod/modules/video/encoding/mocks"
	"vod/modules/video/encrypt/mocks"
	VmdMocks "vod/modules/video/videoMetadata/mocks"

	"vod/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVideoEncryptHandlerUnauthorizedAccess(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
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
	rGroup.POST("/video/encrypt", handler.VideoEncryptHandler)
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encrypt", nil)
	respR := httptest.NewRecorder()
	auth.On("IsEmpty", mock.AnythingOfType("string")).Return(true)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"No token provided, Unauthorized Access\"}"
	assert.Equal(t, http.StatusUnauthorized, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func mockAuth(auth **authMocks.AuhorisationUtils) {
	email := "test@example.com"
	claims := jwt.MapClaims{
		"userEmail": email,
		"exp":       float64(time.Now().Add(12 * time.Second).Unix()),
	}
	token := &jwt.Token{
		Claims: claims,
		Valid:  true,
	}
	(*auth).On("IsEmpty", mock.AnythingOfType("string")).Return(false)
	(*auth).On("IsNotBearerToken", "Bearer your-access-token").Return(false)
	(*auth).On("TokenParse", "Bearer your-access-token").Return(token, nil)
	(*auth).On("IsExpiredToken", mock.AnythingOfType("float64")).Return(false)
}

func TestVideoEncryptHandlerInvalidRequestBody(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
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
	rGroup.POST("/video/encrypt", handler.VideoEncryptHandler)
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encrypt", nil)
	req.Header.Set("Authorization", "Bearer your-access-token")
	respR := httptest.NewRecorder()
	mockAuth(&auth)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Bad Request: Unable to encrypt video.\"}"
	assert.Equal(t, http.StatusBadRequest, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncryptHandlerWhenErrorInVMDRedis(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	redisRepoVMD := new(VmdMocks.RedisRepository)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		redisRepoVideoMD: redisRepoVMD,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encrypt", handler.VideoEncryptHandler)

	videoToEncrypt := VideoEncryptRequest{
		VideoID: 1,
	}

	b, _ := json.Marshal(videoToEncrypt)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encrypt", body)
	req.Header.Set("Authorization", "Bearer your-access-token")
	respR := httptest.NewRecorder()
	mockAuth(&auth)

	// Mock an error in Redis lookup
	redisRepoVMD.On(
		"GetFromRedis", mock.Anything,
	).Return(
		model.VideoMetaData{}, 
		errors.New("Redis error"),
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Invalid Video Meta Data.\"}"
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncryptHandlerWhenErrorInVEncodedRepo(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	redisRepoVMD := new(VmdMocks.RedisRepository)
	redisRepoVE := new(videoEncodingMocks.RedisRepository)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		redisRepoVideoMD: redisRepoVMD,
		redisRepoVideoEncoded: redisRepoVE,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encrypt", handler.VideoEncryptHandler)

	videoToEncrypt := VideoEncryptRequest{
		VideoID: 1,
	}

	b, _ := json.Marshal(videoToEncrypt)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encrypt", body)
	req.Header.Set("Authorization", "Bearer your-access-token")
	respR := httptest.NewRecorder()
	mockAuth(&auth)

	// Mock an error in Redis lookup
	redisRepoVMD.On(
		"GetFromRedis", mock.Anything,
	).Return(
		model.VideoMetaData{}, 
		nil,
	)

	redisRepoVE.On(
		"GetFromRedis", mock.Anything,
	).Return(
		[]model.VideoLinks{},
		errors.New("Redis error"),
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Error Occured! Unable to find encoded videos.\"}"
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncryptHandlerWhenEmptyLinksInVEncodedRepo(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	redisRepoVMD := new(VmdMocks.RedisRepository)
	redisRepoVE := new(videoEncodingMocks.RedisRepository)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		redisRepoVideoMD: redisRepoVMD,
		redisRepoVideoEncoded: redisRepoVE,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encrypt", handler.VideoEncryptHandler)

	videoToEncrypt := VideoEncryptRequest{
		VideoID: 1,
	}

	b, _ := json.Marshal(videoToEncrypt)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encrypt", body)
	req.Header.Set("Authorization", "Bearer your-access-token")
	respR := httptest.NewRecorder()
	mockAuth(&auth)

	// Mock an error in Redis lookup
	redisRepoVMD.On(
		"GetFromRedis", mock.Anything,
	).Return(
		model.VideoMetaData{}, 
		nil,
	)

	redisRepoVE.On(
		"GetFromRedis", mock.Anything,
	).Return(
		nil,
		nil,
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Unable to find encoded videos.\"}"
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncryptHandlerWhenErrorInEncrypt(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	redisRepoVMD := new(VmdMocks.RedisRepository)
	redisRepoVE := new(videoEncodingMocks.RedisRepository)
	vidEncrypt := new(mocks.IEncryption)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		redisRepoVideoMD: redisRepoVMD,
		redisRepoVideoEncoded: redisRepoVE,
		videoEncrypt: vidEncrypt,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encrypt", handler.VideoEncryptHandler)

	videoToEncrypt := VideoEncryptRequest{
		VideoID: 1,
	}

	b, _ := json.Marshal(videoToEncrypt)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encrypt", body)
	req.Header.Set("Authorization", "Bearer your-access-token")
	respR := httptest.NewRecorder()
	mockAuth(&auth)

	// Mock an error in Redis lookup
	redisRepoVMD.On(
		"GetFromRedis", mock.Anything,
	).Return(
		model.VideoMetaData{}, 
		nil,
	)

	redisRepoVE.On(
		"GetFromRedis", mock.Anything,
	).Return(
		[]model.VideoLinks{},
		nil,
	)

	vidEncrypt.On(
		"Encrypt", mock.Anything,
	).Return(
		make(map[model.VideoLinks]model.VideoLinks),
		errors.New("Encrypting error"),
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Unable to encrypt video. Encrypting error\"}"
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncryptHandlerWhenErrorInRepo(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	redisRepoVMD := new(VmdMocks.RedisRepository)
	redisRepoVE := new(videoEncodingMocks.RedisRepository)
	vidEncrypt := new(mocks.IEncryption)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		redisRepoVideoMD: redisRepoVMD,
		redisRepoVideoEncoded: redisRepoVE,
		videoEncrypt: vidEncrypt,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encrypt", handler.VideoEncryptHandler)

	videoToEncrypt := VideoEncryptRequest{
		VideoID: 1,
	}

	b, _ := json.Marshal(videoToEncrypt)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encrypt", body)
	req.Header.Set("Authorization", "Bearer your-access-token")
	respR := httptest.NewRecorder()
	mockAuth(&auth)

	// Mock an error in Redis lookup
	redisRepoVMD.On(
		"GetFromRedis", mock.Anything,
	).Return(
		model.VideoMetaData{}, 
		nil,
	)

	redisRepoVE.On(
		"GetFromRedis", mock.Anything,
	).Return(
		[]model.VideoLinks{},
		nil,
	)

	vidEncrypt.On(
		"Encrypt", mock.Anything,
	).Return(
		make(map[model.VideoLinks]model.VideoLinks),
		nil,
	)

	repo.On(
		"Update", mock.Anything,
	).Return(
		false,
		errors.New("error"),
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Unable to encrypt video, save to db. error\"}"
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncryptHandlerWhenErrorInRedisRepo(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	redisRepoVMD := new(VmdMocks.RedisRepository)
	redisRepoVE := new(videoEncodingMocks.RedisRepository)
	vidEncrypt := new(mocks.IEncryption)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		redisRepoVideoMD: redisRepoVMD,
		redisRepoVideoEncoded: redisRepoVE,
		videoEncrypt: vidEncrypt,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encrypt", handler.VideoEncryptHandler)

	videoToEncrypt := VideoEncryptRequest{
		VideoID: 1,
	}

	b, _ := json.Marshal(videoToEncrypt)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encrypt", body)
	req.Header.Set("Authorization", "Bearer your-access-token")
	respR := httptest.NewRecorder()
	mockAuth(&auth)

	// Mock an error in Redis lookup
	redisRepoVMD.On(
		"GetFromRedis", mock.Anything,
	).Return(
		model.VideoMetaData{}, 
		nil,
	)

	redisRepoVE.On(
		"GetFromRedis", mock.Anything,
	).Return(
		[]model.VideoLinks{},
		nil,
	)

	vidEncrypt.On(
		"Encrypt", mock.Anything,
	).Return(
		make(map[model.VideoLinks]model.VideoLinks),
		nil,
	)

	repo.On(
		"Update", mock.Anything,
	).Return(
		true,
		nil,
	)
	
	redisRepo.On(
		"UpdateInRedis", mock.Anything,
	).Return(
		false,
		errors.New("error"),
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Unable to encrypt video, save to rdb. error\"}"
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncryptHandler(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	redisRepoVMD := new(VmdMocks.RedisRepository)
	redisRepoVE := new(videoEncodingMocks.RedisRepository)
	vidEncrypt := new(mocks.IEncryption)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		redisRepoVideoMD: redisRepoVMD,
		redisRepoVideoEncoded: redisRepoVE,
		videoEncrypt: vidEncrypt,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encrypt", handler.VideoEncryptHandler)

	videoToEncrypt := VideoEncryptRequest{
		VideoID: 1,
	}

	b, _ := json.Marshal(videoToEncrypt)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encrypt", body)
	req.Header.Set("Authorization", "Bearer your-access-token")
	respR := httptest.NewRecorder()
	mockAuth(&auth)

	// Mock an error in Redis lookup
	redisRepoVMD.On(
		"GetFromRedis", mock.Anything,
	).Return(
		model.VideoMetaData{}, 
		nil,
	)

	redisRepoVE.On(
		"GetFromRedis", mock.Anything,
	).Return(
		[]model.VideoLinks{},
		nil,
	)

	vidEncrypt.On(
		"Encrypt", mock.Anything,
	).Return(
		make(map[model.VideoLinks]model.VideoLinks),
		nil,
	)

	repo.On(
		"Update", mock.Anything,
	).Return(
		true,
		nil,
	)
	
	redisRepo.On(
		"UpdateInRedis", mock.Anything,
	).Return(
		true,
		nil,
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"message\":\"Video Encrypted successfully!\",\"video_metadata\":{\"ID\":0,\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null,\"title\":\"\",\"format\":\"\",\"uploadedBy\":0,\"length\":0,\"size\":0,\"resolution\":\"\",\"topic\":\"\"},\"encoded_to_encrypted\":{}}"
	assert.Equal(t, http.StatusOK, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}