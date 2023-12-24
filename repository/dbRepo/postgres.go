package dbrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/raihan2bd/vidverse/models"
	"gorm.io/gorm"
)

// Get user by username
func (m *postgresDBRepo) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	var ErrUserNotFound = errors.New("user not found")

	result := m.DB.First(&user, "user_name = ?", username)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	if result.Error != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

// Get user by email
func (m *postgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	m.DB.First(&user, "email = ?", email)
	if user.ID > 0 {
		return &user, nil
	}
	return nil, errors.New("404 user not found")
}

// Get user by ID
func (m *postgresDBRepo) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := m.DB.First(&user, id)
	if result.Error != nil {
		return nil, errors.New("404 user not found")
	}

	user.Password = ""

	return &user, nil
}

// create new user
func (m *postgresDBRepo) CreateNewUser(user *models.User) (int, error) {
	result := m.DB.Create(&user)
	if result.Error != nil {
		return 0, errors.New("failed to create the user. please try again later")
	}

	return int(user.ID), nil
}

// Get all videos from the database
func (m *postgresDBRepo) GetAllVideos(page, limit int, searchQuery string) ([]models.VideoDTO, int64, error) {
	var videos []models.VideoDTO
	var count int64

	offset := (page - 1) * limit

	err := m.DB.Table("videos").Select("videos.id, videos.title, videos.thumb, videos.views, channels.id as channel_id, channels.title as channel_title, channels.logo as channel_logo").
		Joins("left join channels on channels.id = videos.channel_id").
		Where("videos.title ILIKE ? OR videos.description ILIKE ? OR channels.title ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%").
		Count(&count).
		Offset(offset).Limit(limit).
		Order("videos.created_at asc").
		Find(&videos).Error
	if err != nil {
		return nil, 0, errors.New("internal server error. Please try again")
	}

	return videos, count, nil
}

// Get total videos count
func (m *postgresDBRepo) GetTotalVideosCount(searchQuery string) (int64, error) {
	var count int64
	// return only videos count from the database with search query (videos title, description, channel title)
	err := m.DB.Table("videos").Select("videos.id").
		Joins("left join channels on channels.id = videos.channel_id").
		Where("videos.title ILIKE ? OR videos.description ILIKE ? OR channels.title ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%").
		Count(&count).Error

	if err != nil {
		return 0, errors.New("internal server error. Please try again")
	}

	return count, nil
}

// Get single video by Id
func (m *postgresDBRepo) GetVideoByID(id int) (*models.Video, error) {

	var video models.Video
	err := m.DB.Preload("Likes").Preload("Comments").Preload("Channel").First(&video, "id = ?", id).Error
	if err != nil {
		return nil, errors.New("404 video not found")
	}
	// update view count
	m.DB.Model(&video).Update("views", video.Views+1)
	video.Views = video.Views + 1

	video.Channel.Subscriptions = 0

	var count int64 = 0
	tx := m.DB.Table("subscriptions").Select("subscriptions.id").
		Joins("left join channels on channels.id = subscriptions.channel_id").
		Where("subscriptions.channel_id = ?", video.ChannelID).
		Count(&count)

	if tx.Error != nil {
		return nil, errors.New("internal server error. Please try again")
	}

	video.Channel.Subscriptions = count
	return &video, nil
}

// Delete video by ID
func (m *postgresDBRepo) DeleteVideoByID(id int) error {
	result := m.DB.Unscoped().Delete(&models.Video{}, id)
	if result.Error != nil {
		return errors.New("something went wrong. failed to delete the video")
	}

	return nil
}

// Get videos by channelID including pagination
func (m *postgresDBRepo) GetVideosByChannelID(id, page, limit int) ([]models.VideoDTO, int64, error) {
	var videos []models.VideoDTO
	var count int64
	offset := (page - 1) * limit

	err := m.DB.Table("videos").Select("videos.id, videos.title, videos.thumb, videos.views, channels.id as channel_id, channels.title as channel_title, channels.logo as channel_logo").
		Joins("left join channels on channels.id = videos.channel_id").
		Where("videos.channel_id = ?", id).
		Count(&count).
		Offset(offset).Limit(limit).
		Order("videos.created_at asc").
		Find(&videos).Error
	if err != nil {
		return nil, 0, errors.New("internal server error. Please try again")
	}

	return videos, count, nil
}

// Get all comments
func (m *postgresDBRepo) GetCommentsByVideoID(id, page, limit int) ([]models.CommentDTO, int64, error) {
	var comments []models.CommentDTO
	var count int64
	offset := (page - 1) * limit

	err := m.DB.Table("comments").Select("comments.id, comments.text, comments.video_id, users.id as user_id, users.name as user_name, users.avatar as user_avatar").
		Joins("left join users on users.id = comments.user_id").
		Where("comments.video_id = ?", id).
		Count(&count).
		Offset(offset).Limit(limit).
		Order("comments.created_at desc").
		Find(&comments).Error
	if err != nil {
		return nil, 0, errors.New("internal server error. Please try again")
	}

	return comments, count, nil
}

// Get comment by ID
func (m *postgresDBRepo) GetCommentByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := m.DB.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, errors.New("404 comment not found")
	}

	return &comment, nil
}

// Create new comment
func (m *postgresDBRepo) CreateComment(comment *models.Comment) (uint, error) {
	result := m.DB.Create(&comment)
	if result.Error != nil {
		return 0, errors.New("failed to create comment")
	}

	return comment.ID, nil
}

// update comment
func (m *postgresDBRepo) UpdateComment(comment *models.Comment) error {
	result := m.DB.Save(&comment)
	if result.Error != nil {
		return errors.New("failed to update comment")
	}

	return nil
}

// delete comment
func (m *postgresDBRepo) DeleteCommentByID(id uint) error {
	result := m.DB.Unscoped().Delete(&models.Comment{}, id)
	if result.Error != nil {
		return errors.New("something went wrong. failed to delete the comment")
	}

	return nil
}

// Get All the channels
func (m *postgresDBRepo) GetChannels(userID int) ([]models.CustomChannel, error) {
	var channels []models.CustomChannel
	err := m.DB.Table("channels").Select("channels.id, channels.title, channels.logo").Where("channels.user_id = ?", userID).
		Find(&channels).Error
	if err != nil {
		return nil, errors.New("internal server error. Please try again")
	}

	return channels, nil
}

// get channel details
func (m *postgresDBRepo) GetChannelByID(id int) (*models.CustomChannelDTO, error) {
	var channel models.CustomChannelDTO

	err := m.DB.Table("channels").Select("channels.id, channels.title, channels.logo, channels.description, count(videos.id) as total_videos, count(subscriptions.id) as total_subscribers").
		Joins("left join videos on videos.channel_id = channels.id").
		Joins("left join subscriptions on subscriptions.channel_id = channels.id").
		Where("channels.id = ?", id).
		Group("channels.id").
		Find(&channel).Error

	if err != nil {
		return nil, errors.New("internal server error. Please try again")
	}

	fmt.Println("I'm working")

	return &channel, nil
}

// delete channel By Id
func (m *postgresDBRepo) DeleteChannelByID(id int) *models.CustomError {
	// get channel by id
	var channel models.Channel
	result := m.DB.First(&channel, id)

	if result.Error != nil {
		return &models.CustomError{Status: 404, Err: errors.New("the channel you want to delete is not found")}
	}

	// delete channel with transaction
	tx := m.DB.Begin()
	tx.Delete(&channel)

	if tx.Error != nil {
		tx.Rollback()
		return &models.CustomError{Status: 500, Err: errors.New("failed to delete the channel")}
	}

	// delete the logoImage
	_, err := m.CLD.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: channel.LogoPublicID, ResourceType: "image"})

	if err != nil {
		tx.Rollback()
		return &models.CustomError{Status: 500, Err: errors.New("failed to delete the channel")}
	}

	tx.Commit()

	return nil
}

// Get like by videoID and userID
func (m *postgresDBRepo) GetLikeByVideoIDAndUserID(videoID, userID uint) (*models.Like, error) {
	var like models.Like
	err := m.DB.Where("video_id = ? AND user_id = ?", videoID, userID).First(&like).Error
	if err != nil {
		return nil, errors.New("404 like not found")
	}

	if like.ID == 0 {
		return nil, errors.New("404 like not found")
	}

	return &like, nil
}

// create like
func (m *postgresDBRepo) CreateLike(like *models.Like) (uint, error) {
	result := m.DB.Create(&like)
	if result.Error != nil {
		return 0, errors.New("failed to create like")
	}

	return like.ID, nil
}

// Delete like by ID
func (m *postgresDBRepo) DeleteLikeByID(id uint) error {
	result := m.DB.Unscoped().Delete(&models.Like{}, id)
	if result.Error != nil {
		return errors.New("something went wrong. failed to delete the like")
	}

	return nil
}

// Check user subscription status for a channel
func (m *postgresDBRepo) IsSubscribed(userID, channelID uint) bool {
	var subscription models.Subscription
	err := m.DB.Where("user_id = ? AND channel_id = ?", userID, channelID).First(&subscription).Error
	if err != nil {
		return false
	}

	if subscription.ID == 0 {
		return false
	}

	return true
}

func (m *postgresDBRepo) ToggleSubscription(userID, channelID uint) error {
	var subscription models.Subscription
	tsx := m.DB.Where("user_id = ? AND channel_id = ?", userID, channelID).First(&subscription)
	if tsx.Error != nil {
		subscription.UserID = userID
		subscription.ChannelID = channelID
		ts := m.DB.Create(&subscription)
		if ts.Error != nil {
			fmt.Println(ts.Error, "error")
			return errors.New("failed to subscription the channel")
		}
	} else {
		ts := m.DB.Unscoped().Delete(&subscription)
		if ts.Error != nil {
			return errors.New("failed to  unsubscribe the channel")
		}
	}

	return nil
}
