package videoMetadata

import (
	// "vod/model"
	"fmt"
	"vod/model"

	"gorm.io/gorm"
)

type Repository interface {
	// Create(user model.User) (model.User, error)

	Find(videoMDToPlay model.VideoMetaData) (model.VideoMetaData, error)
}

type repository struct {
	client *gorm.DB
}

func NewRepository(client *gorm.DB) Repository {
	return &repository{client: client}
}

func (r *repository) Find(videoMDToPlay model.VideoMetaData) (model.VideoMetaData, error) {
	fmt.Println("Obj: ", videoMDToPlay)
	var videoMD model.VideoMetaData

	res := r.client.Where(
		"video_meta_data.title=? and video_meta_data.resolution=?",
		videoMDToPlay.Title, videoMDToPlay.Resolution,
	).Find(&videoMD)

	fmt.Println("VideoMD: ", videoMD)

	if res.Error != nil {
		return model.VideoMetaData{}, res.Error
	}

	return videoMD, nil
}

// 	encryptPassword, err := generateFromPassword(user.Password)
// 	if err != nil {
// 		return model.User{}, err
// 	}
// 	user.Password = encryptPassword

// 	role, err := setRole(user.Role)
// 	if err != nil {
// 		return model.User{}, err
// 	}
// 	user.Role = role

// 	result := r.client.Create(&user)

// 	if result.Error != nil {
// 		return model.User{}, result.Error
// 	}

// 	return user, nil
// }
