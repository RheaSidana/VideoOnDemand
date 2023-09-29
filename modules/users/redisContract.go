package users

import (
	"encoding/json"
	"vod/model"
)

type RedisData struct {
	User model.User `json:"User"`
	Token string `json:"Token"`
}

func (rd RedisData) MarshalBinary() ([]byte, error) {
    return json.Marshal(rd) 
}

func (rd *RedisData) UnmarshalBinary(data []byte) error {
    return json.Unmarshal(data, rd) 
}