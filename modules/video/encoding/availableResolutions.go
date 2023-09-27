package videoEncoding

// px
var Resolutions = []string{
	"1920x1080", "1280x720", 
	"320x240", "640x480", "720x480",
	"720x576", "1440x1080", "3840x2160",
}

func Resolution1920x1080() string {
	return Resolutions[0]
}

func Resolution1280x720() string {
	return Resolutions[1]
}