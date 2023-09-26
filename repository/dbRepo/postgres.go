package dbrepo

import (
	"errors"

	"github.com/raihan2bd/vidverse/models"
)

// Get all videos from the database
func (m *postgresDBRepo) GetAllVideos() ([]models.Video, error) {
	var videos []models.Video
	if err := m.DB.Find(&videos).Error; err != nil {
		return nil, errors.New("internal server error. Please try again")
	}

	return videos, nil
}

// Get single video by Id
func (m *postgresDBRepo) GetVidoeByID(id int) (*models.Video, error) {

	var video models.Video
	if id > 0 {
		m.DB.First(&video, "id = ?", id)

		if video.ID == 0 {
			return nil, errors.New("404 video not found")
		}
	}

	return &video, nil
}
