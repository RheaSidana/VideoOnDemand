package videoEncoding

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"vod/model"

	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func videoData(fileLoc string) model.VideoData {
	// fileLoc := BASE_DIR + "Sample/Title 1.mp4"
	ffprobe, err := ffmpeg.Probe(fileLoc)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return model.VideoData{}
	}

	type VideoInfo struct {
		Streams []struct {
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			BitRate  string `json:"bit_rate"`
			Duration string `json:"duration"`
		} `json:"streams"`
	}

	var info VideoInfo
	if err := json.Unmarshal([]byte(ffprobe), &info); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return model.VideoData{}
	}

	var data model.VideoData
	// fmt.Println("Input command run res")
	// fmt.Println(ffprobe)
	if len(info.Streams) > 0 {
		stream := info.Streams[0]
		// fmt.Printf("Resolution: %dx%d\n", stream.Width, stream.Height)
		data.Resolution = strconv.Itoa(stream.Width) + "x" + strconv.Itoa(stream.Height)

		// fmt.Printf("Bitrate: %s\n", stream.BitRate)
		if len(stream.BitRate) >= 3 {
			data.BitRate = stream.BitRate[:len(stream.BitRate)-3] + "K"
		} else {
			data.BitRate = stream.BitRate
		}

		// fmt.Printf("Duration: %s\n", stream.Duration)
		dur, err := strconv.ParseFloat(stream.Duration, 64)
		if err != nil {
			return model.VideoData{}
		}
		data.Duration, err = strconv.ParseFloat(fmt.Sprintf("%.3f", dur), 64)
		if err != nil {
			return model.VideoData{}
		}
	} else {
		fmt.Println("No video streams found in the file.")
		return model.VideoData{}
	}

	return data
}

func saveOriginalVideo(c *gin.Context, videoFile *multipart.FileHeader) (string, error) {
	saveToLocation := BASE_DIR_ORIGINAL
	fileLoc := saveToLocation + videoFile.Filename
	return fileLoc, c.SaveUploadedFile(videoFile, fileLoc)
}

func saveVideoToLoc(
	inputFileLoc string,
	data model.VideoData,
) (string, error) {
	videoFileTitle := strings.Split((strings.Replace(inputFileLoc, BASE_DIR_ORIGINAL, "", -1)), ".")[0]
	outputFileLoc, err := getOutputFilePath(
		videoFileTitle,
		data.BitRate,
		data.Resolution,
	)
	if err != nil {
		return "",err
	}

	cmd := ffmpeg.Input(
		inputFileLoc,
	).Output(
		outputFileLoc,
		ffmpeg.KwArgs{
			"c:v": "copy",
			"c:a": "copy",
		},
	)

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error: (concat) :", err.Error())
		return "",err
	}

	return outputFileLoc, nil
}

func saveDataInRedis(
	h Handler, 
	videoLinks []model.VideoLinks,
) error {
	for _, videoLink := range videoLinks {

		err :=h.redisRepository.SetInRedis(videoLink)
		if err!=nil {
			return err
		}
	}

	return nil
}