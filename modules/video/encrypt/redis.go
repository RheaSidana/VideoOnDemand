package videoEncryption

import (
	"context"
	"vod/model"
	"errors"

	"github.com/go-redis/redis/v8"
)

type RedisRepository interface {
	UpdateInRedis(encodeToEncrypt map[model.VideoLinks]model.VideoLinks) (bool, error)
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

func (r *redisRepository) UpdateInRedis(
	encodeToEncrypt map[model.VideoLinks]model.VideoLinks,
) (
	bool, error,
) {
	KEY := SET_NAME
	multipleVideoMD, err := r.client.SMembers(r.ctx, KEY).Result()

	for key, val := range encodeToEncrypt {
		foundRedisData, _ := findFromSet(
			int(key.VideoMetaDataID),
			key.Link,
			multipleVideoMD,
			err,
		)

		foundRedisDataByte, _ := foundRedisData.MarshalBinary()
		r.client.SRem(r.ctx, KEY, foundRedisDataByte)

		foundRedisData.VideoLinks.Link = val.Link
		err := r.client.SAdd(r.ctx, KEY, foundRedisData).Err()
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func findFromSet(searchForID int, link string, redisDatas []string, err error) (
	RedisData, error,
) {
	if err != nil {
		return RedisData{}, err
	}

	for _, rdata := range redisDatas {

		var data RedisData
		data.UnmarshalBinary([]byte(rdata))

		if data.VideoLinks.VideoMetaDataID == uint(searchForID) &&
			data.VideoLinks.Link == link {
			return data, nil
		}
	}

	return RedisData{}, errors.New("unable to find in redis")
}
