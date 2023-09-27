package videoMetadata

import (
	// "fmt"
	"vod/model"
	"vod/modules/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository      Repository
	redisRepository RedisRepository
}

func (h *Handler) PlayVideoHandler(c *gin.Context) {
	if !middleware.IsAuthorisedAccess(c) {
		return
	}

	// should have videoMD ID
	var videoToPlay model.VideoMetaData
	c.BindJSON(&videoToPlay)
	if videoToPlay == (model.VideoMetaData{}) {
		c.JSON(400, ErrorResponse{
			Message: "Bad Request: Unable to play video."})
		return
	}

	// check in redis
	videoMD, err := h.redisRepository.GetFromRedis(
		videoToPlay)
	if err == nil {
		// found
		c.JSON(200, VideoMDResponse{
			Message: "Found in Redis, Play the video ",
			VideoMD: videoMD,
		})
		return
	}

	// check in db
	videoMD, err = h.repository.Find(videoToPlay)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "Unable to play video."})
		return
	}

	c.JSON(200, VideoMDResponse{
		Message: "Found, Play the video ",
		VideoMD: videoMD,
	})
}
