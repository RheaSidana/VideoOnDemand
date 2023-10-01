package videoEncryption

import (
	"vod/model"
	videoEncoding "vod/modules/video/encoding"
	"vod/modules/video/videoMetadata"
)

func videoMDfromRedis(
	videoID int,
	videoMDRedisRepo videoMetadata.RedisRepository,
) (model.VideoMetaData, error) {
	var videoMD model.VideoMetaData
	videoMD.ID = uint(videoID)
	return videoMDRedisRepo.GetFromRedis(videoMD)
}

func videoLinksFromRedis(
	videoMD model.VideoMetaData,
	videoLinksRedisRepo videoEncoding.RedisRepository,
) (
	[]model.VideoLinks, error,
) {
	var videoLink model.VideoLinks
	videoLink.VideoMetaDataID = videoMD.ID
	return videoLinksRedisRepo.GetFromRedis(
		videoLink,
	)
}

// func updateVideoLinksToRedis(
// 	encodeToEncrypt map[model.VideoLinks]model.VideoLinks,
// ) (
// 	bool, error,
// ){
// 	videoLinksRedisRepo := videoEncoding.InitRedisRepository(
// 		initializer.RedisDB,
// 		initializer.Ctx,
// 	)

// 	videoLinksRedisRepo.UpdateInRedis(encodeToEncrypt)

// 	return true, nil
// }