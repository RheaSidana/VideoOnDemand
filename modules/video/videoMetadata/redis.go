package videoMetadata

import (
	"context"
	"errors"

	// "strconv"
	"vod/model"

	"github.com/go-redis/redis/v8"
)

var SET_NAME = "VideoMetaData"

type RedisRepository interface {
	SetInRedis(videoMetadata model.VideoMetaData) error
	GetFromRedis(videoMD model.VideoMetaData) (model.VideoMetaData, error)
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

func (r *redisRepository) SetInRedis(videoMD model.VideoMetaData) error {
	key, val := createRedisKeyValuePair(videoMD)

	exists, _ := r.client.SIsMember(r.ctx, key, val).Result()
	if exists {
		return nil
	}

	multipleVideoMD, err := r.client.SMembers(r.ctx, key).Result()
	foundRedisData, _ := findFromSet(int(videoMD.ID), multipleVideoMD, err)
	if foundRedisData != (RedisData{}) {
		foundRedisData, _ := foundRedisData.MarshalBinary()
		r.client.SRem(r.ctx, key, foundRedisData)
	}

	return r.client.SAdd(r.ctx, key, val).Err()
}

func (r *redisRepository) GetFromRedis(videoMD model.VideoMetaData) (
	model.VideoMetaData, error,
) {
	// key := "vidMD:" + strconv.Itoa(int(videoMD.ID))
	key := SET_NAME
	redisDataString, err := r.client.SMembers(r.ctx, key).Result()

	rdata, err := findFromSet(int(videoMD.ID), redisDataString, err)
	return rdata.VideoMetaData, err
}

func findFromSet(searchForID int, redisDatas []string, err error) (
	RedisData, error,
) {
	if err != nil {
		return RedisData{}, err
	}

	for _, rdata := range redisDatas {

		var data RedisData
		data.UnmarshalBinary([]byte(rdata))

		if data.VideoMetaData.ID == uint(searchForID) {
			return data, nil
		}
	}

	return RedisData{}, errors.New("unable to find in redis")
}

func createRedisKeyValuePair(videoMD model.VideoMetaData) (string, []byte) {
	// key := "vidMD:" + strconv.Itoa(int(videoMD.ID))
	key := SET_NAME
	val, _ := RedisData{
		VideoMetaData: videoMD,
	}.MarshalBinary()

	return key, val
}
