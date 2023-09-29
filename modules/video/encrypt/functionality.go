package videoEncryption

import (
	"vod/initializer"
	"vod/model"
	videoEncoding "vod/modules/video/encoding"
	"vod/modules/video/videoMetadata"
)

func videoMDfromRedis(videoID int) (model.VideoMetaData, error) {
	videoMDRedisRepo := videoMetadata.InitRedisRepository(
		initializer.RedisDB,
		initializer.Ctx,
	)
	var videoMD model.VideoMetaData
	videoMD.ID = uint(videoID)
	return videoMDRedisRepo.GetFromRedis(videoMD)
}

func videoLinksFromRedis(videoMD model.VideoMetaData) (
	[]model.VideoLinks, error,
) {
	videoLinksRedisRepo := videoEncoding.InitRedisRepository(
		initializer.RedisDB,
		initializer.Ctx,
	)
	
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