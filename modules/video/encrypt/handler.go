package videoEncryption

import (
	// "fmt"
	"vod/modules/middleware"
	videoEncoding "vod/modules/video/encoding"
	"vod/modules/video/videoMetadata"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository            Repository
	redisRepository       RedisRepository
	redisRepoVideoMD      videoMetadata.RedisRepository
	redisRepoVideoEncoded videoEncoding.RedisRepository
	videoEncrypt          IEncryption
}

func (h *Handler) VideoEncryptHandler(c *gin.Context) {
	if !middleware.IsAuthorisedAccess(c) {
		return
	}

	var videoToEncrypt VideoEncryptRequest
	c.BindJSON(&videoToEncrypt)
	if videoToEncrypt == (VideoEncryptRequest{}) {
		c.JSON(400, ErrorResponse{
			Message: "Bad Request: Unable to encrypt video.",
		})
		return
	}

	//videoMD
	videoMD, err := videoMDfromRedis(
		videoToEncrypt.VideoID,
		h.redisRepoVideoMD,
	)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "Invalid Video Meta Data."})
		return
	}

	// videoLinks = encoded
	videoLinks, err := videoLinksFromRedis(
		videoMD,
		h.redisRepoVideoEncoded,
	)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "Error Occured! Unable to find encoded videos."})
		return
	}
	if videoLinks == (nil) {
		c.JSON(500, ErrorResponse{
			Message: "Unable to find encoded videos."})
		return
	}

	mapVideoEncodedToVideoEncrypted, err := h.videoEncrypt.Encrypt(videoLinks)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "Unable to encrypt video. " + err.Error()})
		return
	}

	//update psql
	isUpdated, err := h.repository.Update(mapVideoEncodedToVideoEncrypted)
	if err != nil && !isUpdated {
		c.JSON(500, ErrorResponse{
			Message: "Unable to encrypt video, save to db. " + err.Error()})
		return
	}

	//update redis
	isUpdated, err = h.redisRepository.UpdateInRedis(mapVideoEncodedToVideoEncrypted)
	if err != nil && !isUpdated {
		c.JSON(500, ErrorResponse{
			Message: "Unable to encrypt video, save to rdb. " + err.Error()})
		return
	}

	c.JSON(200, VideoEncryptResponse{
		Message:            "Video Encrypted successfully!",
		VideoMD:            videoMD,
		EncodedToEncrypted: mapVideoEncodedToVideoEncrypted,
	})
}
