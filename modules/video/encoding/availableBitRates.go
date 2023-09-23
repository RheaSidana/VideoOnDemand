package videoEncoding

var Bitrates = []string{
	"1500k", "800k",
}

func Bitrates1500k() string {
	return Bitrates[0]
}

func Bitrates800k() string {
	return Bitrates[1]
}
