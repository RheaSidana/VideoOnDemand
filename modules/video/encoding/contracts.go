package videoEncoding

import "vod/model"

type ErrorResponse struct {
	Message string
}

type VideoEncodeResponse struct {
	Message    string
	VideoMD    model.VideoMetaData
	VideoLinks []model.VideoLinks
}
