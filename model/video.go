package model

import "gorm.io/gorm"

type VideoMetaData struct {
	gorm.Model
	Title      string  `gorm:"primaryKey" json:"title" binding:"required"`
	UploadedBy int     `json:"uploadedBy" binding:"required"`
	Length     float64 `json:"length" binding:"required"`
	Size       float64 `json:"size" binding:"required"`
	Resolution string  `gorm:"primaryKey" json:"resolution" binding:"required"`
	Topic      string  `json:"topic" binding:"required"`
}
