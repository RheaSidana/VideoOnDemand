package videoEncryption

import "os"

var createFolder func(string) error

// var createVideosTxtFile func([]string, string, string) error
var deleteFolder func(string) error

// var deleteFile func(string) error

func AssignOsFunctions(
	CreateFolder func(string) error,
	CreateVideosTxtFile func([]string, string, string) error,
	DeleteFolder func(string) error,
	DeleteFile func(string) error,
) {
	createFolder = CreateFolder
	// createVideosTxtFile = CreateVideosTxtFile
	// deleteFile = DeleteFile
	deleteFolder = DeleteFolder
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false // File does not exist
	}
	return err == nil // File exists (or there was an error, such as permission issues)
}
