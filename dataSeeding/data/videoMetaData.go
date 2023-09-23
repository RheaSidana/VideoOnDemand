package data

import (
	"strconv"
	"vod/initializer"
	"vod/model"
	videoEncoding "vod/modules/video/encoding"
)

func videoMDData(users []model.User) []model.VideoMetaData {
	genres := []string{
		"Science", "Math", "Hindi",
		"English", "Social Science",
	}
	videoMDToAdd := model.VideoMetaData{
		Title:      "Title",
		UploadedBy: int(users[0].ID),
		Length:     15,
		Size:       950,
		Resolution: videoEncoding.Resolution1920x1080(),
	}

	var videoMDList []model.VideoMetaData
	for i := 1; i <= 5; i++ {
		videoMD := videoMDToAdd
		videoMD.Title = videoMD.Title + " " + strconv.Itoa(i)
		videoMD.Topic = genres[i-1]
		if i >= 4 {
			videoMD.UploadedBy = int(users[1].ID)
		}
		videoMDList = append(videoMDList, videoMD)
	}

	return videoMDList
}

func AddVideoMDToDB(users []model.User) ([]model.VideoMetaData) {
	var videoMDList []model.VideoMetaData
	for _, videoMD := range videoMDData(users) {
		if initializer.Db.Where(
			"title=? and resolution=?",
			videoMD.Title,
			videoMD.Resolution).Find(&videoMD).RowsAffected == 1 {
			continue
		}
		initializer.Db.Create(&videoMD)

		videoMDList = append(videoMDList, videoMD)
	}

	return videoMDList
}
