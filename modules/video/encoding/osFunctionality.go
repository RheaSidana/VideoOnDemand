package videoEncoding

import (
	"fmt"
	"os"
	"strings"
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

func deleteFolder(filename string) error {
	folderPath := BASE_DIR_PROCESSING + filename

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
