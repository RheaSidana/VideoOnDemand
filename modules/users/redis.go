package users

import (
	"context"
	"vod/model"

	"github.com/go-redis/redis/v8"
)

type RedisRepository interface {
	SetInRedis(user model.User, token string) error
	GetFromRedis(user model.User) (string, error)
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
	value, _ := r.GetFromRedis(user)

	isAvailable, err := isAvailableInRedis(val, value)
	if err != nil {
		return err
	}
	if isAvailable {
		return nil
	}

	return r.client.Set(r.ctx, key, val, 0).Err()
}

func (r *redisRepository) GetFromRedis(user model.User) (string, error) {
	return r.client.Get(r.ctx, user.Email).Result()
}

func isAvailableInRedis(val RedisData, dataInRedis string) (bool, error) {
	jsonString, err := val.MarshalBinary()
	if err != nil {
		return false, err
	}
	valString := string(jsonString)

	if valString == dataInRedis {
		return true, nil
	}

	return false, nil
}

func createRedisKeyValuePair(user model.User, token string) (
	string, RedisData,
) {
	key := user.Email
	val := RedisData{
		User:  user,
		Token: token,
	}

	return key, val
}
