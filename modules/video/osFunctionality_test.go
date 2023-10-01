package videos

import (
	// "fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFolderWhenError(t *testing.T) {
	tmpFolder := ""

	defer func() {
		err := os.RemoveAll(tmpFolder)
		assert.Equal(t, nil, err)
	}()

	err := createFolder(tmpFolder)
	assert.NotEqual(t, nil, err)

	_, err = os.Stat(tmpFolder)
	assert.NotEqual(t, nil, err)
}

func TestCreateFolder(t *testing.T) {
	tmpFolder := "./test_tmp_folder"

	defer func() {
		err := os.RemoveAll(tmpFolder)
		assert.Equal(t, nil, err)
	}()

	err := createFolder(tmpFolder)
	assert.Equal(t, nil, err)

	_, err = os.Stat(tmpFolder)
	assert.Equal(t, nil, err)
}

func TestDeleteFolder(t *testing.T) {
	folderPath := "path/to/folder"

	err := deleteFolder(folderPath)

	assert.NoError(t, err)
}

func TestCreateVideosTxtFile_CreatesFileWithGivenFilename(t *testing.T) {
	videoFiles := []string{"video1.mp4", "video2.mp4"}
	filenameToCreate := "videos.txt"
	videoTitle := "video_title"

	err := createVideosTxtFile(videoFiles, filenameToCreate, videoTitle)

	assert.NoError(t, err)
	_, err = os.Stat(filenameToCreate)
	assert.False(t, os.IsNotExist(err))
	deleteFile(filenameToCreate)
}

func TestRemoveRedundantPath(t *testing.T) {
	videoPath := BASE_DIR_PROCESSING + BASE_DIR_PROCESSING + "path/to/video/originalVideoTitle/video.mp4"
	originalVideoTitle := BASE_DIR_PROCESSING + "path/to/video/originalVideoTitle"
	expectedFileName := "video.mp4"

	actualFileName := removeRedundantPath(videoPath, originalVideoTitle)

	assert.Equal(t, expectedFileName, actualFileName)
}

func TestDeleteFile(t *testing.T) {
	t.Run("DeleteExistingFile", func(t *testing.T) {
		// Define a temporary file for testing.
		tmpFile := "./test_tmp_file_existing.txt"

		// Create the temporary file for testing.
		file, err := os.Create(tmpFile)
		if err != nil {
			t.Fatalf("Failed to create temporary file for testing: %v", err)
		}
		file.Close() // Close the file here.

		// Test deleting an existing file.
		err = deleteFile(tmpFile)
		assert.NoError(t, err)

		// Check if the file no longer exists.
		_, err = os.Stat(tmpFile)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("DeleteNonExistingFile", func(t *testing.T) {
		// Test deleting a non-existing file.
		nonExistentFile := "./non_existent_file.txt"
		err := deleteFile(nonExistentFile)
		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
	})
}

