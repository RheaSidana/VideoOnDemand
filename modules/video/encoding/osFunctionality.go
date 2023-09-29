package videoEncoding


var createFolder func(string) error
var createVideosTxtFile func([]string, string, string) error
var deleteFolder func(string) error
var deleteFile func(string) error

func AssignOsFunctions(
	CreateFolder func(string) error,
	CreateVideosTxtFile func([]string, string, string) error,
	DeleteFolder func(string) error,
	DeleteFile func(string) error,
) {
	createFolder = CreateFolder
	createVideosTxtFile = CreateVideosTxtFile
	deleteFile = DeleteFile
	deleteFolder = DeleteFolder
}

func createOutputFolderForEncoding(inputVideoTitle, bitrate string) (string, error){
	saveToLocation := BASE_DIR_ENCODED

	videoTitleBitrateLocation := saveToLocation +
		inputVideoTitle + "/" +
		"B" + bitrate + "/"

	err := createFolder(videoTitleBitrateLocation)
	if err != nil {
		return "", err
	}

	return videoTitleBitrateLocation, nil
}