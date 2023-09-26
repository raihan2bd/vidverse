package repository

import "github.com/raihan2bd/vidverse/models"

type DatabaseRepo interface {
	GetAllVideos() ([]models.Video, error)
	GetVidoeByID(id int) (*models.Video, error)
	DeleteVideoByID(id int) error
}
