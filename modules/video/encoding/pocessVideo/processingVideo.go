package videoProcessing

import (
	"fmt"
	// "os"

	"sync"

	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ProcessVideo(
	fileName, resolution, bitrate string, 
	createFolder func(string) error,
	inputPaths []string,
) []string {

	outputPaths := getOutputFilePathForProccessing_output(
		fileName, bitrate, resolution,
		createFolder,
	)
	
	var wg sync.WaitGroup

	for i, inputPath := range inputPaths {

		wg.Add(1)

		go encodeVideo(
			inputPath, outputPaths[i],
			resolution, bitrate,
			&wg,
		)
	}

	wg.Wait()
	return outputPaths
}

func encodeVideo(
	inputPath, outputPath, resolution, bitrate string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	input := ffmpeg.Input(inputPath)
	cmd := input.Output(
		outputPath,
		ffmpeg.KwArgs{
			"b:v": bitrate,
			"s":   resolution,
		},
	).OverWriteOutput()

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func getOutputFilePathForProccessing_output(
	inputFilename, bitrate, resolution string,
	createFolder func(string) error,
) []string {
	inputVideoTitle := strings.Split(inputFilename, ".")[0]
	saveToLocation := BASE_DIR_PROCESSING

	loc := saveToLocation + inputVideoTitle

	err := createFolder(loc)
	if err != nil {
		fmt.Println("error: ", err.Error())
	}

	bitrateResolution := "_B" + bitrate + "_R" + resolution

	numberOfVideos := 5
	return createVideoChunks(
		loc, "vid_out_00", (bitrateResolution + ".mp4"), 
		numberOfVideos,
	)
}
