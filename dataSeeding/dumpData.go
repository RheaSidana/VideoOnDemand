package main

import (
	"vod/dataSeeding/data"
	"vod/initializer"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectToDB()
}

func main() {
	usersData := data.AddAdminUsersToDB()
	videoMDs :=  data.AddVideoMDToDB(usersData)

	data.SetInRedis(usersData, videoMDs)
}