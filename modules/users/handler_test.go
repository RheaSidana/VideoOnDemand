package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"vod/model"
	"vod/modules/users/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignupHandlerWhenEmptyUser(t *testing.T) {
	repo := new(mocks.Repository)
	handler := Handler{repository: repo}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/user", handler.SignUpHandler)
	newUser := model.User{}
	b, _ := json.Marshal(newUser)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/user", body)
	respR := httptest.NewRecorder()
	expectedUser := newUser
	repo.On("Create", newUser).Return(model.User{}, errors.New("Empty User JSON."))

	actualUser, _ := repo.Create(newUser)
	r.ServeHTTP(respR, req)

	assert.Equal(t, respR.Code, http.StatusBadRequest)
	assert.Equal(t, expectedUser, actualUser)
}

func TestSignupHandlerWhenUnableToCreateUserInDB(t *testing.T) {
	repo := new(mocks.Repository)
	handler := Handler{repository: repo}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/user", handler.SignUpHandler)
	newUser := model.User{
		Name:     "test",
		Email:    "test@example.com",
		Password: "Test@3_5",
		Role:     "admin",
	}
	b, _ := json.Marshal(newUser)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/user", body)
	respR := httptest.NewRecorder()
	expectedUser := newUser
	expectedUser.ID = 1
	repo.On("Create", newUser).Return(model.User{}, errors.New("Error while creating user"))

	actualUser, _ := repo.Create(newUser)
	r.ServeHTTP(respR, req)

	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.NotEqual(t, expectedUser, actualUser)
}

func TestSignUpHandlerWhenUnableToCreateUserInRedisDB(t *testing.T) {
	repo := new(mocks.Repository)
	redis := new(mocks.RedisRepository)
	handler := Handler{
		repository: repo,
		redisRepository: redis,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/user", handler.SignUpHandler)
	newUser := model.User{
		Name:     "test",
		Email:    "test@example.com",
		Password: "Test@3_5",
		Role:     "admin",
	}
	b, _ := json.Marshal(newUser)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/user", body)
	respR := httptest.NewRecorder()
	expectedUser := newUser
	expectedUser.ID = 1
	repo.On("Create", newUser).Return(expectedUser, nil)
	redis.On(
		"SetInRedis", 
		expectedUser,
		mock.AnythingOfType("string"),
	).Return(errors.New("error occurred"))

	actualUser, _ := repo.Create(newUser)
	r.ServeHTTP(respR, req)

	assert.Equal(t, http.StatusInternalServerError, respR.Code)
	assert.Equal(t, expectedUser, actualUser)
}

func TestSignUpHandler(t *testing.T) {
	repo := new(mocks.Repository)
	redis := new(mocks.RedisRepository)
	handler := Handler{
		repository: repo,
		redisRepository: redis,
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/user", handler.SignUpHandler)
	newUser := model.User{
		Name:     "test",
		Email:    "test@example.com",
		Password: "Test@3_5",
		Role:     "admin",
	}
	b, _ := json.Marshal(newUser)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/user", body)
	respR := httptest.NewRecorder()
	expectedUser := newUser
	expectedUser.ID = 1
	repo.On("Create", newUser).Return(expectedUser, nil)
	redis.On(
		"SetInRedis", 
		expectedUser,
		mock.AnythingOfType("string"),
	).Return(nil)

	actualUser, _ := repo.Create(newUser)
	r.ServeHTTP(respR, req)

	assert.Equal(t, http.StatusOK, respR.Code)
	assert.Equal(t, expectedUser, actualUser)
}
