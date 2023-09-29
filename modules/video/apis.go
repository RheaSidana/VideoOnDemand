package videos

import (
	videoEncoding "vod/modules/video/encoding"
	videoEncryption "vod/modules/video/encrypt"
	"vod/modules/video/videoMetadata"

	"github.com/gin-gonic/gin"
)

func Apis(r *gin.RouterGroup) {
	assignConstantsToEncoding()
	assignOsFunctions()
	
	videoMetadata.Apis(r)
	videoEncoding.Apis(r)
	videoEncryption.Apis(r)
}