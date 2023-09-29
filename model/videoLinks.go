package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type VideoLinks struct {
	gorm.Model
	VideoMetaDataID uint          `gorm:"not null" json:"videoID" binding:"required"`
	VideoMetaData   VideoMetaData `gorm:"foreignKey:VideoMetaDataID; references:ID" json:"-" binding:"-"`
	Link            string        `gorm:"not null; unique" json:"encodedLink" binding:"required"`
}

func (v VideoLinks) ToString() string {
	return fmt.Sprintf(
		"VideoMetaDataID: %d, EncodedLink: %s",
		v.VideoMetaDataID,
		v.Link,
	)
}

func StringToVideoLinks(input string) (VideoLinks, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return VideoLinks{}, errors.New("invalid input format")
	}

	// Parse VideoMetaDataID from the string
	videoMetaDataID, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return VideoLinks{}, errors.New("invalid video_metadata_id")
	}

	// Create a VideoLinks instance
	videoLinks := VideoLinks{
		VideoMetaDataID: uint(videoMetaDataID),
		Link:            parts[1],
	}

	return videoLinks, nil
}
