package users

import (
	"vod/model"
)

type ErrorResponse struct {
	Message string
}

type UserResponse struct {
	Message string
}

type LoginResponse struct {
	Token   string
	Message string
	User    model.User
}

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

