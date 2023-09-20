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
)

func TestCreateUserHandlerWhenEmptyUser(t *testing.T) {
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

func TestCreateUserHandlerWhenUnableToCreateUser(t *testing.T) {
	repo := new(mocks.Repository)
	handler := Handler{repository: repo}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/user", handler.SignUpHandler)
	newUser := model.User{
		Name:     "test",
		Email:    "test@example.com",
		Password: "Test@3_5",
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

	assert.Equal(t, respR.Code, http.StatusInternalServerError)
	assert.NotEqual(t, expectedUser, actualUser)
}

func TestCreateUserHandler(t *testing.T) {
	repo := new(mocks.Repository)
	handler := Handler{repository: repo}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/user", handler.SignUpHandler)
	newUser := model.User{
		Name:     "test",
		Email:    "test@example.com",
		Password: "Test@3_5",
	}
	b, _ := json.Marshal(newUser)
	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "/user", body)
	respR := httptest.NewRecorder()
	expectedUser := newUser
	expectedUser.ID = 1
	repo.On("Create", newUser).Return(expectedUser, nil)

	actualUser, _ := repo.Create(newUser)
	r.ServeHTTP(respR, req)

	assert.Equal(t, respR.Code, http.StatusOK)
	assert.Equal(t, expectedUser, actualUser)
}

// func (h *Handler) FindUser(c *gin.Context, user int) (model.User, error){
// 	return h.repository.Find(user)
// }
