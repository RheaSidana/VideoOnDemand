package videoEncoding

import videoProcessing "vod/modules/video/encoding/pocessVideo"

var BASE_DIR string
var BASE_DIR_PROCESSING string
var BASE_DIR_ENCODED string
var BASE_DIR_ORIGINAL string

func AssignConstants(
	baseDIR, 
	baseDirORIGINAL, 
	baseDirENCODED, 
	baseDirPROCESSING string,
) {
	BASE_DIR = baseDIR
	BASE_DIR_ORIGINAL = baseDirORIGINAL
	BASE_DIR_PROCESSING = baseDirPROCESSING
	BASE_DIR_ENCODED = baseDirENCODED
}

func assignConstantsToProcessing() {
	videoProcessing.AssignConstants(
		BASE_DIR,
		BASE_DIR_ORIGINAL,
		BASE_DIR_ENCODED,
		BASE_DIR_PROCESSING,
	)
}