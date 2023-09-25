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
