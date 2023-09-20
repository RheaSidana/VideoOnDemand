package users

import (
	// "fmt"
	"encoding/json"
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

type RedisData struct {
	User model.User `json:"User"`
	Token string `json:"Token"`
}

// func (RedisData) ToSting() string {
// 	return fmt.Sprintf("{User: %s")
// }

func (rd RedisData) MarshalBinary() ([]byte, error) {
    return json.Marshal(rd) // You can choose any binary serialization format you prefer
}

func (rd *RedisData) UnmarshalBinary(data []byte) error {
    return json.Unmarshal(data, rd) // You should use the same serialization format as MarshalBinary
}