package videoEncoding

// kbps
var Bitrates = []string{
	"1500k", "800k", 
	"512k", "768k", "900k",
	"1200k", "2000k",
}

func Bitrates1500k() string {
	return Bitrates[0]
}

func Bitrates800k() string {
	return Bitrates[1]
}
