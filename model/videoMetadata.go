package model

import "gorm.io/gorm"

type VideoMetaData struct {
	gorm.Model
	Title      string  `gorm:"unique; not null" json:"title" binding:"required"`
	Format     string  `json:"format" binding:"required"`
	UploadedBy int     `json:"uploadedBy" binding:"required"`
	Length     float64 `json:"length" binding:"required"`
	Size       float64 `json:"size" binding:"required"`
	Resolution string  `json:"resolution" binding:"required"`
	Topic      string  `json:"topic" binding:"required"`
}
