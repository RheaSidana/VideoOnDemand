package main

import (
	"vod/initializer"
	"vod/model"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectToDB()
}

func main() {
	initializer.Db.AutoMigrate(&model.User{})
	initializer.Db.AutoMigrate(&model.VideoMetaData{})
}
