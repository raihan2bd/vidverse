package repository

import "github.com/raihan2bd/vidverse/models"

type DatabaseRepo interface {
	CreateNewUser(*models.User) (int, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)

	GetAllVideos() ([]models.VideoDTO, error)
	GetVidoeByID(id int) (*models.Video, error)
	DeleteVideoByID(id int) error
}
