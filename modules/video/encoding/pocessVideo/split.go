package videoProcessing

import (
	"fmt"
	"strconv"
	"sync"
	"vod/model"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func SplitVideoIntoFiveParts(
	videoFileName, videoTitle string,
	createFolder func(string) error,
	videoData model.VideoData,
	wg *sync.WaitGroup,
) ([]string, error) {
	defer wg.Done()

	fileLoc := BASE_DIR_ORIGINAL + videoFileName

	segmentDuration := findSegmentDuration(videoData)

	loc, err := createOutputFolderForProcessing(
		videoTitle, createFolder)
	if err != nil {
		return []string{}, err
	}

	outputPaths := getOutputFilePathForProccessing_input(loc)

	err = split(fileLoc, outputPaths, segmentDuration)
	if err != nil {
		return []string{}, err
	}

	return outputPaths, nil
}

func split(intputPath string, outputPaths []string,
	segmentDuration float64,
) error {
	seekTo := 0.0
	
	for _, outputFile := range outputPaths {

		input := ffmpeg.Input(intputPath)

		cmd := input.Output(
			outputFile,
			ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{
				{"ss": strconv.FormatFloat(seekTo, 'f', -1, 64)},
				{"t": strconv.FormatFloat(segmentDuration, 'f', -1, 64)},
				{"c": "copy"},
			}),
		).OverWriteOutput()

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return err
			// os.Exit(1)
		}

		seekTo += segmentDuration
	}

	return nil
}

func findSegmentDuration(videoData model.VideoData) float64 {
	totalDuration := videoData.Duration
	return totalDuration / 5
}

func getOutputFilePathForProccessing_input(loc string) []string {
	return createVideoChunks(loc, "vid00", ".mp4", 5)
}

func createOutputFolderForProcessing(
	inputVideoTitle string, createFolder func(string) error) (string, error) {
	saveToLocation := BASE_DIR_PROCESSING
	loc := saveToLocation + inputVideoTitle

	err := createFolder(loc)
	if err != nil {
		return "", err
	}
	return loc, nil
}
