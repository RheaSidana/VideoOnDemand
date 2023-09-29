package videoEncoding

import (
	"strconv"
	"vod/model"

	"gorm.io/gorm"
)

type Repository interface {
	Create(videoLinks []string, videoId string) ([]model.VideoLinks, model.VideoMetaData, error)
}

type repository struct {
	client *gorm.DB
}

func NewRepository(client *gorm.DB) Repository {
	return &repository{client: client}
}

func (r *repository) Create(videoLinks []string, videoId string) ([]model.VideoLinks, model.VideoMetaData, error) {
	//convert videoId
	vidId, err := strconv.Atoi(videoId)
	if err != nil {
		return []model.VideoLinks{}, model.VideoMetaData{}, err
	}

	//valid videoMD : find
	videoMD, err := findVideoMetadata(*r, uint(vidId))
	if err != nil {
		return []model.VideoLinks{}, model.VideoMetaData{}, err
	}

	//create videolinks
	createdVideoLinks, err := createVideoLinks(
		*r, uint(vidId), videoLinks)
	if err != nil {
		return []model.VideoLinks{}, model.VideoMetaData{}, err
	}

	return createdVideoLinks, videoMD, nil
}

func createVideoLinks(repo repository, videoId uint, videoLinks []string) ([]model.VideoLinks, error) {
	var createdVideoLinks []model.VideoLinks
	for _, link := range videoLinks {
		videoLink := model.VideoLinks{
			VideoMetaDataID: videoId,
			Link:            link,
		}

		res := repo.client.Create(&videoLink)

		if res.Error != nil {
			return []model.VideoLinks{}, res.Error
		}

		createdVideoLinks = append(createdVideoLinks, videoLink)
	}

	return createdVideoLinks, nil
}

func findVideoMetadata(repo repository, videoId uint) (model.VideoMetaData, error) {
	var videoMD model.VideoMetaData

	res := repo.client.Where(
		"id=?", videoId,
	).Find(&videoMD)

	if res.Error != nil {
		return model.VideoMetaData{}, res.Error
	}

	return videoMD, nil
}
