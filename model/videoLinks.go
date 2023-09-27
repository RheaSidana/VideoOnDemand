package model

import "gorm.io/gorm"

type VideoLinks struct {
	gorm.Model
	VideoMetaDataID uint        `gorm:"not null" json:"videoID" binding:"required"`
	VideoMetaData   VideoMetaData `gorm:"foreignKey:VideoMetaDataID; references:ID" json:"-" binding:"-"`
	EncodedLink     string        `gorm:"not null; unique" json:"encodedLink" binding:"required"`
}
