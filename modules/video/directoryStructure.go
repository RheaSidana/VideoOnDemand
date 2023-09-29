package videos

import (
	videoEncoding "vod/modules/video/encoding"
	videoEncryption "vod/modules/video/encrypt"
)

var BASE_DIR string = "./modules/video/videosCollection/"
var BASE_DIR_ORIGINAL string = BASE_DIR + "Original/"
var BASE_DIR_PROCESSING string = BASE_DIR + "Processing/"
var BASE_DIR_ENCODED string = BASE_DIR + "Encoded/"

func assignConstantsToEncoding(){
	videoEncoding.AssignConstants(
		BASE_DIR,
		BASE_DIR_ORIGINAL,
		BASE_DIR_ENCODED,
		BASE_DIR_PROCESSING,
	)

	videoEncryption.AssignConstants(
		BASE_DIR,
		BASE_DIR_ORIGINAL,
		BASE_DIR_ENCODED,
		BASE_DIR_PROCESSING,
	)
}