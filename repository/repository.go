package repository

import "github.com/raihan2bd/vidverse/models"

type DatabaseRepo interface {
	CreateNewUser(user *models.User) (int, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)

	GetAllVideos(page, limit int, searchQuery string) ([]models.VideoDTO, int64, error)
	GetTotalVideosCount(searchQuery string) (int64, error)
	GetVideoByID(id int) (*models.Video, error)
	GetVideosByChannelID(id, page, limit int) ([]models.VideoDTO, int64, error)
	DeleteVideoByID(id int) error
	GetCommentsByVideoID(id, page, limit int) ([]models.CommentDTO, int64, error)
	GetChannels(id int) ([]models.CustomChannel, error)
	GetChannelByID(id int) (*models.CustomChannelDTO, error)
}
