package videoMetadata

import "vod/model"

type ErrorResponse struct {
	Message string
}

type VideoMDResponse struct {
	Message string
	VideoMD model.VideoMetaData
}

type VideoMDRequest struct {
	VideoMD model.VideoMetaData
	Category string `json:"category" binding:"required"`
}
