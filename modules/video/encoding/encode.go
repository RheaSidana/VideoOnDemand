package videoEncoding

import (
	"fmt"
	"strings"
	"sync"
	"vod/model"
	videoProcessing "vod/modules/video/encoding/pocessVideo"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// current function is encoding the video in
// Bitrates1500k and Resolution1280x720
func encode(
	inputVideoFile string, inputFileLoc string,
	videoData model.VideoData,	
) ([]string, error) {
	// bitRate := Bitrates800k()
	// resolution := Resolution1280x720()
	bitRate, resolution := Bitrates1500k(), Resolution1280x720()

	inputVideoTitle := strings.Split(inputVideoFile, ".")[0]
	outputVideoFileLoc, err := getOutputFilePath(
		inputVideoTitle, bitRate, resolution)
	if err != (nil) {
		return []string{}, err
	}
	
	var wg sync.WaitGroup
	wg.Add(1)
	
	assignConstantsToProcessing()
	splitedVideos, err := videoProcessing.SplitVideoIntoFiveParts(
		inputVideoFile, inputVideoTitle, createFolder, 
		videoData, &wg)
	if err != (nil) {
		return []string{}, err
	}

	wg.Wait()

	processedVideos := videoProcessing.ProcessVideo(
		inputVideoFile, resolution, bitRate, createFolder, splitedVideos)

	// create txt file saving the video in it
	filenameTXT := BASE_DIR_PROCESSING + inputVideoTitle +
		"/" + "videos.txt"
	err = createVideosTxtFile(processedVideos, filenameTXT, inputVideoTitle)
	if err != nil {
		fmt.Println("Error: (create File) :", err.Error())
		return []string{}, err
	}

	// run command to concat the vid files to create one file
	// cmd := exec.Command(
	// 	"ffmpeg",
	// 	"-f", "concat",
	// 	"-safe", "0",
	// 	"-i", filenameTXT,
	// 	"-c:v", "copy",
	// 	"-c:a", "copy",
	// 	outputVideoFileLoc,
	// 	"-y",
	// 	// "-v", "info", "-report", "log.txt",
	// )
	input := ffmpeg.Input(filenameTXT, ffmpeg.KwArgs{
		"f":    "concat",
		"safe": 0,
	})
	cmd := input.Output(
		outputVideoFileLoc,
		ffmpeg.KwArgs{
			"c:v": "copy",
			"c:a": "copy",
			// "v" : "info",
			// "report": "log.txt",
		},
	)

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error: (concat) :", err.Error())
		return []string{}, err
	}

	// delete processing folder
	err = deleteFolder(inputVideoTitle)
	if err != nil {
		fmt.Println("Error: (delete folder) :", err.Error())
		return []string{}, err
	}

	//delete original file
	err = deleteFile(inputFileLoc)
	if err != nil {
		fmt.Println("Error: (delete file) :", err.Error())
		return []string{}, err
	}

	return []string{outputVideoFileLoc}, nil
}

func getOutputFilePath(
	inputVideoTitle, bitrate, resolution string,
) (string, error) {

	videoTitleBitrateLocation, err := createOutputFolderForEncoding(
		inputVideoTitle, bitrate,
	)
	if err != nil {
		return "", err
	}

	outputVido := "R" + resolution + ".mp4"
	outputVideoLocation := videoTitleBitrateLocation +
		outputVido

	return outputVideoLocation, nil
}
