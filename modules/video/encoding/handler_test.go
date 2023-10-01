package videoEncoding

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"vod/model"
	"vod/modules/middleware"
	authMocks "vod/modules/middleware/mocks"
	"vod/modules/video/encoding/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVideoEncodeHandlerWhenUnauthorised(t *testing.T){
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
	rGroup.POST("/video/encode", handler.VideoEncodeHandler)
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encode", nil)
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

func TestVideoEncodeHandlerInvalidFormFile(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	funcImpl := new(mocks.IFunctionality)
	encImpl := new(mocks.IEncoding)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		functionality:   funcImpl,
		encodingVideo:   encImpl,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encode", handler.VideoEncodeHandler)
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encode", nil)
	req.Header.Set("Authorization", "Bearer your-access-token")
	respR := httptest.NewRecorder()
	mockAuth(&auth)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"request Content-Type isn't multipart/form-data   Bad Request: Unable to encode video.\"}"
	assert.Equal(t, http.StatusBadRequest, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func createTempFile(t *testing.T, content []byte, filename string) string {
	file, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		t.Fatalf("Failed to write temporary file: %v", err)
	}

	return file.Name()
}


func TestVideoEncodeHandlerWhenErrorSavingOriginalVideo(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	funcImpl := new(mocks.IFunctionality)
	encImpl := new(mocks.IEncoding)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		functionality:   funcImpl,
		encodingVideo:   encImpl,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encode", handler.VideoEncodeHandler)
	videoContent := []byte("test video content")
	tempVideoFile := createTempFile(t, videoContent, "test_video.mp4")
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	writer.WriteField("videoId", "123")
	part, _ := writer.CreateFormFile("video", "test_video.mp4")
	_, _ = part.Write(videoContent)
	writer.Close()
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encode", reqBody)
	req.Header.Set("Authorization", "Bearer your-access-token")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	respR := httptest.NewRecorder()
	mockAuth(&auth)
	funcImpl.On(
		"SaveOriginalVideo", 
		mock.Anything, 
		mock.Anything,
	).Return(
		tempVideoFile, 
		errors.New("error occurred"),
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Unable to save uploaded video.\"}"
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncodeHandlerWhenErrorSaveVideoToLoc(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	funcImpl := new(mocks.IFunctionality)
	encImpl := new(mocks.IEncoding)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		functionality:   funcImpl,
		encodingVideo:   encImpl,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encode", handler.VideoEncodeHandler)
	videoContent := []byte("test video content")
	tempVideoFile := createTempFile(t, videoContent, "test_video.mp4")
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	writer.WriteField("videoId", "123")
	part, _ := writer.CreateFormFile("video", "test_video.mp4")
	_, _ = part.Write(videoContent)
	writer.Close()
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encode", reqBody)
	req.Header.Set("Authorization", "Bearer your-access-token")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	respR := httptest.NewRecorder()
	mockAuth(&auth)
	funcImpl.On(
		"SaveOriginalVideo", mock.Anything, mock.Anything,
	).Return(
		tempVideoFile, nil,
	)
	funcImpl.On(
		"GetVideoData", mock.Anything, 
	).Return(
		model.VideoData{},
	)
	funcImpl.On(
		"SaveVideoToLoc", mock.Anything, mock.Anything,
	).Return(
		"savedAt", errors.New("error occurred"),
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Unable to save video.\"}"
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncodeHandlerWhenErrorEncode(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	funcImpl := new(mocks.IFunctionality)
	encImpl := new(mocks.IEncoding)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		functionality:   funcImpl,
		encodingVideo:   encImpl,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encode", handler.VideoEncodeHandler)
	videoContent := []byte("test video content")
	tempVideoFile := createTempFile(t, videoContent, "test_video.mp4")
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	writer.WriteField("videoId", "123")
	part, _ := writer.CreateFormFile("video", "test_video.mp4")
	_, _ = part.Write(videoContent)
	writer.Close()
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encode", reqBody)
	req.Header.Set("Authorization", "Bearer your-access-token")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	respR := httptest.NewRecorder()
	mockAuth(&auth)
	funcImpl.On(
		"SaveOriginalVideo", mock.Anything, mock.Anything,
	).Return(
		tempVideoFile, nil,
	)
	funcImpl.On(
		"GetVideoData", mock.Anything, 
	).Return(
		model.VideoData{},
	)
	funcImpl.On(
		"SaveVideoToLoc", mock.Anything, mock.Anything,
	).Return(
		"savedAt", nil,
	)
	encImpl.On(
		"Encode", mock.Anything, mock.Anything, mock.Anything,
	).Return(
		[]string{},
		errors.New("error occurred"),
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Unable to save encode video.\"}"
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncodeHandlerWhenErrorCreate(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	funcImpl := new(mocks.IFunctionality)
	encImpl := new(mocks.IEncoding)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		functionality:   funcImpl,
		encodingVideo:   encImpl,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encode", handler.VideoEncodeHandler)
	videoContent := []byte("test video content")
	tempVideoFile := createTempFile(t, videoContent, "test_video.mp4")
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	writer.WriteField("videoId", "123")
	part, _ := writer.CreateFormFile("video", "test_video.mp4")
	_, _ = part.Write(videoContent)
	writer.Close()
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encode", reqBody)
	req.Header.Set("Authorization", "Bearer your-access-token")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	respR := httptest.NewRecorder()
	mockAuth(&auth)
	funcImpl.On(
		"SaveOriginalVideo", mock.Anything, mock.Anything,
	).Return(
		tempVideoFile, nil,
	)
	funcImpl.On(
		"GetVideoData", mock.Anything, 
	).Return(
		model.VideoData{},
	)
	funcImpl.On(
		"SaveVideoToLoc", mock.Anything, mock.Anything,
	).Return(
		"savedAt", nil,
	)
	encImpl.On(
		"Encode", mock.Anything, mock.Anything, mock.Anything,
	).Return(
		[]string{"temp1", "temp2"},
		nil,
	)
	repo.On(
		"Create", mock.Anything, mock.Anything,
	).Return(
		[]model.VideoLinks{},
		model.VideoMetaData{},
		errors.New("error occurred"),
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"Unable to save encode video.\"}"
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncodeHandlerWhenErrorSaveDataInRedis(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	funcImpl := new(mocks.IFunctionality)
	encImpl := new(mocks.IEncoding)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		functionality:   funcImpl,
		encodingVideo:   encImpl,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encode", handler.VideoEncodeHandler)
	videoContent := []byte("test video content")
	tempVideoFile := createTempFile(t, videoContent, "test_video.mp4")
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	writer.WriteField("videoId", "123")
	part, _ := writer.CreateFormFile("video", "test_video.mp4")
	_, _ = part.Write(videoContent)
	writer.Close()
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encode", reqBody)
	req.Header.Set("Authorization", "Bearer your-access-token")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	respR := httptest.NewRecorder()
	mockAuth(&auth)
	funcImpl.On(
		"SaveOriginalVideo", mock.Anything, mock.Anything,
	).Return(
		tempVideoFile, nil,
	)
	funcImpl.On(
		"GetVideoData", mock.Anything, 
	).Return(
		model.VideoData{},
	)
	funcImpl.On(
		"SaveVideoToLoc", mock.Anything, mock.Anything,
	).Return(
		"savedAt", nil,
	)
	encImpl.On(
		"Encode", mock.Anything, mock.Anything, mock.Anything,
	).Return(
		[]string{"temp1", "temp2"},
		nil,
	)
	videoLinksSlice := []model.VideoLinks{
		{
			VideoMetaDataID: 1,
			Link: "Link 1",
		},{
			VideoMetaDataID: 1,
			Link: "Link 2",
		},
	}
	videoMD := model.VideoMetaData{
		Title: "testing title",
	}
	videoMD.ID = 1
	repo.On(
		"Create", mock.Anything, mock.Anything,
	).Return(
		videoLinksSlice,
		videoMD,
		nil,
	)
	redisRepo.On(
		"SetInRedis", mock.Anything,
	).Return(
		errors.New("error occured"),
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\"Message\":\"error : Error adding token to redis: error occured\"}"
	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

func TestVideoEncodeHandler(t *testing.T) {
	repo := new(mocks.Repository)
	redisRepo := new(mocks.RedisRepository)
	funcImpl := new(mocks.IFunctionality)
	encImpl := new(mocks.IEncoding)
	handler := &Handler{
		repository:      repo,
		redisRepository: redisRepo,
		functionality:   funcImpl,
		encodingVideo:   encImpl,
	}
	auth := new(authMocks.AuhorisationUtils)
	authHandler := middleware.Handler{
		Auth: auth,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	rGroup := r.Group("/protected", authHandler.AuthMiddleware)
	rGroup.POST("/video/encode", handler.VideoEncodeHandler)
	videoContent := []byte("test video content")
	tempVideoFile := createTempFile(t, videoContent, "test_video.mp4")
	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	writer.WriteField("videoId", "123")
	part, _ := writer.CreateFormFile("video", "test_video.mp4")
	_, _ = part.Write(videoContent)
	writer.Close()
	req, _ := http.NewRequest(http.MethodPost, "/protected/video/encode", reqBody)
	req.Header.Set("Authorization", "Bearer your-access-token")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	respR := httptest.NewRecorder()
	mockAuth(&auth)
	funcImpl.On(
		"SaveOriginalVideo", mock.Anything, mock.Anything,
	).Return(
		tempVideoFile, nil,
	)
	funcImpl.On(
		"GetVideoData", mock.Anything, 
	).Return(
		model.VideoData{},
	)
	funcImpl.On(
		"SaveVideoToLoc", mock.Anything, mock.Anything,
	).Return(
		"savedAt", nil,
	)
	encImpl.On(
		"Encode", mock.Anything, mock.Anything, mock.Anything,
	).Return(
		[]string{"temp1", "temp2"},
		nil,
	)
	videoLinksSlice := []model.VideoLinks{
		{
			VideoMetaDataID: 1,
			Link: "Link 1",
		},{
			VideoMetaDataID: 1,
			Link: "Link 2",
		},
	}
	videoMD := model.VideoMetaData{
		Title: "testing title",
	}
	videoMD.ID = 1
	repo.On(
		"Create", mock.Anything, mock.Anything,
	).Return(
		videoLinksSlice,
		videoMD,
		nil,
	)
	redisRepo.On(
		"SetInRedis", mock.Anything,
	).Return(
		nil,
	)

	r.ServeHTTP(respR, req)

	expectedResponse := "{\""+
	"Message\":\"Encoded Successfully!\","+
	"\"VideoMD\":{\""+
		"ID\":1,"+
		"\"CreatedAt\":\"0001-01-01T00:00:00Z\","+
		"\"UpdatedAt\":\"0001-01-01T00:00:00Z\","+
		"\"DeletedAt\":null,"+
		"\"title\":\"testing title\","+
		"\"format\":\"\","+
		"\"uploadedBy\":0,"+
		"\"length\":0,"+
		"\"size\":0,"+
		"\"resolution\":\"\","+
		"\"topic\":\"\""+
	"},\""+
	"VideoLinks\":["+
		"{\"ID\":0,"+
			"\"CreatedAt\":\"0001-01-01T00:00:00Z\","+
			"\"UpdatedAt\":\"0001-01-01T00:00:00Z\","+
			"\"DeletedAt\":null,"+
			"\"videoID\":1,"+
			"\"encodedLink\":\"Link 1\""+
		"},"+
		"{\"ID\":0,"+
			"\"CreatedAt\":\"0001-01-01T00:00:00Z\","+
			"\"UpdatedAt\":\"0001-01-01T00:00:00Z\","+
			"\"DeletedAt\":null,"+
			"\"videoID\":1,"+
			"\"encodedLink\":\"Link 2\""+
		"}"+
	"]}"
	assert.Equal(t, http.StatusOK, respR.Code)
	assert.Equal(t, expectedResponse, respR.Body.String())
}

