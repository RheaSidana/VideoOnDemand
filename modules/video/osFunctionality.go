package videos

import (
	"fmt"
	"os"
	"strings"
	videoEncoding "vod/modules/video/encoding"
	videoEncryption "vod/modules/video/encrypt"
)

func createFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err == nil || os.IsExist(err) {
		fmt.Printf(
			"Folder '%s' already exists or was created successfully.\n",
			folderPath,
		)
		return nil
	}
	return err
}

func deleteFolder(foldername string) error {
	folderPath := foldername

	err := os.RemoveAll(folderPath)
	if err != nil {
		return err
	}

	return nil
}

func createVideosTxtFile(
	videoFiles []string,
	filenameToCreate, videoTitle string,
) error {

	file, err := os.Create(filenameToCreate)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, fileName := range videoFiles {
		fileName = removeRedundantPath(fileName, videoTitle)
		_, err := file.WriteString(fmt.Sprintf("file '%s'\n", fileName))
		if err != nil {
			return err
		}
	}

	return nil
}

func removeRedundantPath(videoPath, originalVideoTitle string) string {
	substr := BASE_DIR_PROCESSING + originalVideoTitle +"/"
	fileName := strings.Replace(videoPath, substr, "", -1)
	
	return fileName
}

func deleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func assignOsFunctions() {
	videoEncoding.AssignOsFunctions(
		createFolder,
		createVideosTxtFile,
		deleteFolder,
		deleteFile,
	)

	videoEncryption.AssignOsFunctions(
		createFolder,
		createVideosTxtFile,
		deleteFolder,
		deleteFile,
	)
}