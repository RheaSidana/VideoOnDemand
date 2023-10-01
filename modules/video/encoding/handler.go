package videoEncoding

import (
	// "fmt"
	// "strings"
	"vod/modules/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository      Repository
	redisRepository RedisRepository
	functionality   IFunctionality
	encodingVideo	IEncoding
}

func (h *Handler) VideoEncodeHandler(c *gin.Context) {
	if !middleware.IsAuthorisedAccess(c) {
		return
	}

	videoToEncode, err := c.FormFile("video")
	if err != nil {
		c.JSON(400, ErrorResponse{
			Message: err.Error() + "   Bad Request: Unable to encode video.",
		})
		return
	}

	fileLoc, err := h.functionality.SaveOriginalVideo(c, videoToEncode)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "Unable to save uploaded video."})
		return
	}

	videoData := h.functionality.GetVideoData(fileLoc)
	fileSavedAt, err := h.functionality.SaveVideoToLoc(fileLoc, videoData)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "Unable to save video."})
		return
	}

	videoEncodedLinks, err := h.encodingVideo.Encode(
		videoToEncode.Filename, fileLoc, videoData)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "Unable to save encode video."})
		return
	}

	videoEncodedLinks = append(videoEncodedLinks, fileSavedAt)
	videoId := c.Request.FormValue("videoId")
	videoLinks, videoMD, err := h.repository.Create(videoEncodedLinks, videoId)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "Unable to save encode video."})
		return
	}

	//redis
	err = saveDataInRedis(*h, videoLinks)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Message: "error : Error adding token to redis: " + err.Error()})
		return
	}

	c.JSON(200, VideoEncodeResponse{
		Message:    "Encoded Successfully!",
		VideoLinks: videoLinks,
		VideoMD:    videoMD,
	})

}
