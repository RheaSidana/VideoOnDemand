package videoEncoding

import (
	"context"
	"errors"

	// "fmt"
	"vod/model"

	"github.com/go-redis/redis/v8"
)

type RedisRepository interface {
	SetInRedis(videoMetadata model.VideoLinks) error
	GetFromRedis(videoMD model.VideoLinks) ([]model.VideoLinks, error)
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

func (r *redisRepository) SetInRedis(videoMD model.VideoLinks) error {
	key, val := createRedisKeyValuePair(videoMD)

	exists, _ := r.client.SIsMember(r.ctx, key, val).Result()
	if exists {
		return nil
	}

	multipleVideoMD, err := r.client.SMembers(r.ctx, key).Result()
	foundRedisData, _ := findFromSet(
		int(videoMD.VideoMetaDataID),
		videoMD.Link,
		multipleVideoMD,
		err,
	)
	if foundRedisData != (RedisData{}) {
		foundRedisData, _ := foundRedisData.MarshalBinary()
		r.client.SRem(r.ctx, key, foundRedisData)
	}

	return r.client.SAdd(r.ctx, key, val).Err()
}

func (r *redisRepository) GetFromRedis(videoMD model.VideoLinks) (
	[]model.VideoLinks, error,
) {
	key := SET_NAME
	redisDataString, err := r.client.SMembers(r.ctx, key).Result()

	rdatas, err := findAllFromSet(
		int(videoMD.VideoMetaDataID),
		redisDataString,
		err,
	)
	return extractLinksFrom(rdatas), err
}

func findAllFromSet(searchForID int, redisDataString []string, err error) (
	[]RedisData, error,
) {
	if err != nil {
		return []RedisData{}, err
	}

	var redisDatas []RedisData
	for _, rdata := range redisDataString {

		var data RedisData
		data.UnmarshalBinary([]byte(rdata))

		if data.VideoLinks.VideoMetaDataID == uint(searchForID) {
			redisDatas = append(redisDatas, data)
		}
	}

	return redisDatas, nil

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

func createRedisKeyValuePair(videoLinks model.VideoLinks) (string, []byte) {
	key := SET_NAME
	val, _ := RedisData{
		VideoLinks: videoLinks,
	}.MarshalBinary()

	return key, val
}

func extractLinksFrom(redisDatas []RedisData) []model.VideoLinks {
	var links []model.VideoLinks
	for _, data := range redisDatas {
		links = append(links, data.VideoLinks)
	}

	return links
}
