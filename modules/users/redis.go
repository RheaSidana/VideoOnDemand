package users

import (
	"context"
	"errors"
	"vod/model"

	"github.com/go-redis/redis/v8"
)

var SetName = "Users"

type RedisRepository interface {
	SetInRedis(user model.User, token string) error
	GetFromRedis(user model.User) (model.User, error)
}

type redisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(client *redis.Client, ctx context.Context) RedisRepository {
	return &redisRepository{
		client: client,
		ctx:    ctx,
	}
}

func (r *redisRepository) SetInRedis(user model.User, token string) error {
	key, val := createRedisKeyValuePair(user, token)

	exists, _ := r.client.SIsMember(r.ctx, key, val).Result()
	if exists {
		return nil
	}

	multipleUsers, err := r.client.SMembers(r.ctx, SetName).Result()
	foundRedisData, _ := findFromSet(user.Email, multipleUsers, err)
	if foundRedisData != (RedisData{}) {
		foundRedisData, _ := foundRedisData.MarshalBinary()
		r.client.SRem(r.ctx, key, foundRedisData)
	}

	return r.client.SAdd(r.ctx, key, val).Err()
}

func (r *redisRepository) GetFromRedis(user model.User) (model.User, error) {
	multipleUsers, err := r.client.SMembers(r.ctx, SetName).Result()

	data, err := findFromSet(user.Email, multipleUsers, err)

	return data.User, err
}

func findFromSet(searchForEmail string, redisDatas []string, err error) (
	RedisData, error,
) {
	if err != nil {
		return RedisData{}, err
	}

	for _, rdata := range redisDatas {

		var data RedisData
		data.UnmarshalBinary([]byte(rdata))

		if data.User.Email == searchForEmail {
			return data, nil
		}
	}

	return RedisData{}, errors.New("unable to find in redis")
}

func createRedisKeyValuePair(user model.User, token string) (
	string, []byte,
) {
	key := SetName
	val, _ := RedisData{
		User:  user,
		Token: token,
	}.MarshalBinary()

	return key, val
}
