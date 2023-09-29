package videoEncryption

import (
	// "strconv"
	"vod/model"

	"gorm.io/gorm"
)

type Repository interface {
	Find(videoID int) (model.VideoMetaData, error)
	Update(encodeToEncrypt map[model.VideoLinks]model.VideoLinks) (bool, error)
}

type repository struct {
	client *gorm.DB
}

func NewRepository(client *gorm.DB) Repository {
	return &repository{client: client}
}

func (r *repository) Update(
	encodeToEncrypt map[model.VideoLinks]model.VideoLinks,
) (bool, error,
) {
	
	for key, val := range encodeToEncrypt {
		// fmt.Printf("key: %v, \nval: %v\n", key.ToString(), val.ToString())

		res := r.client.Model(
			&model.VideoLinks{},
		).Where(
			"video_meta_data_id=? and link=?",
			key.VideoMetaDataID,
			key.Link,
		).Update(
			"link",
			val.Link,
		)

		if res.Error != nil {
			return false, res.Error
		}
	}
	return true, nil
}

func (r *repository) Find(videoID int) (model.VideoMetaData, error) {
	var videoMD model.VideoMetaData

	res := r.client.Where(
		"id=?", videoID,
	).Find(&videoMD)

	if res.Error != nil {
		return model.VideoMetaData{}, res.Error
	}

	return videoMD, nil
}
