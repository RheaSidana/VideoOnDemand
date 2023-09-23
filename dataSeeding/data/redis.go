package data

import (
	"context"
	"vod/initializer"
	"vod/model"
	"vod/modules/users"
	"vod/modules/video/videoMetadata"
)

func SetInRedis(dummyUsers []model.User, videoMetaDatas []model.VideoMetaData) {
	ctx := context.Background()
    initializer.NewRedisClient(ctx)
	
	setUsers(dummyUsers)
	setVideoMetaData(videoMetaDatas)
}

func setVideoMetaData(videoMetaDatas []model.VideoMetaData) {
	repo := videoMetadata.NewRedisRepository(initializer.RedisDB, initializer.Ctx)

	for _, videoMD := range videoMetaDatas {
		repo.SetInRedis(videoMD)
	}
}

func setUsers(dummyUsers []model.User) {
	repo := users.NewRedisRepository(initializer.RedisDB, initializer.Ctx)

	for _, user := range dummyUsers {
		repo.SetInRedis(user, "")
	}
}