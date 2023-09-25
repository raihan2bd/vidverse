package repository

import "github.com/raihan2bd/vidverse/models"

type DatabaseRepo interface {
	GetAllVideos() ([]models.Video, error)
}
