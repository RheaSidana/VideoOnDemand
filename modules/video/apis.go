package videos

import (
	videoEncoding "vod/modules/video/encoding"
	"vod/modules/video/videoMetadata"

	"github.com/gin-gonic/gin"
)

func Apis(r *gin.RouterGroup) {
	videoMetadata.Apis(r)
	videoEncoding.Apis(r)
}