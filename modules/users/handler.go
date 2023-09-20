package users

import (
	"vod/model"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository Repository
	redisRepository RedisRepository
}

func (h *Handler) SignUpHandler(c *gin.Context) {
	var newUser model.User
	c.BindJSON(&newUser)
	if newUser == (model.User{}) {
		c.JSON(400, ErrorResponse{
			Message: "Bad Request: Unable to add user."})
		return
	}

	user, err := h.repository.Create(newUser)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "Unable to add user."})
		return
	}

	err = h.redisRepository.SetInRedis(user, "")
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "Unable to add user."})
		return
	}

	c.JSON(200, UserResponse{
		Message: user.Name + " created successfully!!"})
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var loginUser LoginUser
	c.BindJSON(&loginUser)
	if loginUser == (LoginUser{}) {
		c.JSON(400, ErrorResponse{
			Message: "Bad Request: Unable to login user."})
		return
	}

	user := model.User{
		Email: loginUser.Email,
		Password: loginUser.Password,
	}
	user, err := h.repository.Find(user)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "Unable to login user."})
		return
	}

	token, err := generateToken(user.Email)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "error : Error generating token"})
		return
	}

	err = h.redisRepository.SetInRedis(user, token)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "error : Error adding token to redis"+ err.Error()})
		return
	}

	c.JSON(200, LoginResponse{
		Token: token,
		Message: "User Logged in successfully!",
		User: user,
	})
}
