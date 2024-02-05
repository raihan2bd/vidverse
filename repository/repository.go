package repository

import "github.com/raihan2bd/vidverse/models"

type DatabaseRepo interface {
	CreateNewUser(user *models.User) (int, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	AddForgotPasswordToken(token *models.Token) error
	UpdateUserPassword(user *models.User) error

	GetAllVideos(page, limit int, searchQuery string) ([]models.VideoDTO, int64, error)
	GetTotalVideosCount(searchQuery string) (int64, error)
	GetVideoByID(id int) (*models.Video, error)
	GetVideosByChannelID(id, page, limit int) ([]models.VideoDTO, int64, error)
	DeleteVideoModel(*models.Video) error
	FindAllVideoIDByChannelID(id uint) ([]uint, error)
	DeleteAllVideoIDByChannelID(id uint) error
	DeleteVideoFromCloudinary(publicID string) error
	CreateVideo(video *models.Video) (uint, error)
	UpdateVideo(video *models.Video) error

	GetCommentsByVideoID(id, page, limit int) ([]models.CommentDTO, int64, error)
	GetCommentByID(id uint) (*models.Comment, error)
	CreateComment(comment *models.Comment) (uint, error)
	UpdateComment(comment *models.Comment) error
	DeleteCommentByID(id uint) error
	DeleteNotificationByCommentID(commentID uint) error

	CreateChannel(channel *models.Channel) (uint, error)
	GetChannels(id int) ([]models.CustomChannel, error)
	GetChannelByID(id int) (*models.CustomChannelDTO, error)
	UpdateChannel(channel *models.CustomChannelDTO) error
	DeleteChannelByID(id int) *models.CustomError
	GetChannelsWithDetailsByUserID(userID uint) ([]models.CustomChannelDTO, error)
	GetVideosByChannelIDWithPagination(channelID uint, page, limit int) ([]models.VideoDTO, int64, error)
	GetChannelWithDetails(channelID uint, userID uint) (*models.CustomChannelDTO, error)

	GetLikeByVideoIDAndUserID(videoID, userID uint) (*models.Like, error)
	CreateLike(like *models.Like) (uint, error)
	DeleteLikeByID(id uint) error
	GetLikedVideos(userIDUint uint, page, limit int) ([]models.VideoDTO, int64, error)

	IsSubscribed(userID, channelID uint) bool
	ToggleSubscription(userID, channelID uint) (uint, error)

	GetNotificationsByUserID(userID uint, page, limit int) ([]models.Notification, int64, error)
	GetUnreadNotificationsByUserID(userID uint) ([]models.Notification, error)
	GetNotificationByID(id uint) (*models.Notification, error)
	CreateNotification(notification *models.Notification) (uint, error)
	DeleteNotificationsByChannelID(id uint) error
	GetUnreadNotificationsCountByUserID(userID uint) (int64, error)
	UpdateNotificationByID(id uint) error

	CreateContactUs(contactUs *models.ContactUs) error
	IsContactUsSubmitted(email string) bool
}
