package repository

import "github.com/raihan2bd/vidverse/models"

type DatabaseRepo interface {
	CreateNewUser(user *models.User) (int, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)

	GetAllVideos(page, limit int, searchQuery string) ([]models.VideoDTO, int64, error)
	GetTotalVideosCount(searchQuery string) (int64, error)
	GetVideoByID(id int) (*models.Video, error)
	DeleteVideoByID(id int) error
}
