package videoEncoding

// px
var Resolutions = []string{
	"1920x1080", "1280x720", "640x360",
}

func Resolution1920x1080() string {
	return Resolutions[0]
}

func Resolution1280x720() string {
	return Resolutions[1]
}

func Resolution640x360() string {
	return Resolutions[2]
}