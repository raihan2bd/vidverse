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

func (m *postgresDBRepo) DeleteVideoFromCloudinary(publicID string) error {
	ctx := context.Background()
	result, err := m.CLD.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID, ResourceType: "video"})
	if err != nil {
		return errors.New("failed to delete image")
	}

	if result.Result != "ok" {
		return errors.New("failed to delete image")
	}

	return nil
}

// Delete video with its related data
func (m *postgresDBRepo) DeleteVideoWithRelatedData(videoID uint) error {
	// delete video from cloudinary
	// select publicID from videos where id = videoID
	var publicID string
	err := m.DB.Table("videos").Select("videos.public_id").Where("videos.id = ?", videoID).First(&publicID).Error
	if err != nil {
		return errors.New("failed to delete the video")
	}

	// delete video from db
	result := m.DB.Unscoped().Delete(&models.Video{}, videoID)
	if result.Error != nil {
		return errors.New("something went wrong. failed to delete the video")
	}

	// delete video from cloudinary
	_ = m.DeleteVideoFromCloudinary(publicID)

	// delete comments
	_ = m.DB.Unscoped().Where("video_id = ?", videoID).Delete(&models.Comment{}).Error

	// delete likes
	_ = m.DB.Unscoped().Where("video_id = ?", videoID).Delete(&models.Like{}).Error

	// delete notifications
	_ = m.DB.Unscoped().Where("video_id = ?", videoID).Delete(&models.Notification{}).Error

	return nil

}

func (m *postgresDBRepo) FindAllVideoIDByChannelID(id uint) ([]uint, error) {
	var videoIDs []uint
	err := m.DB.Table("videos").Select("videos.id").Where("channel_id = ?", id).Find(&videoIDs).Error
	if err != nil {
		return nil, errors.New("internal server error. Please try again")
	}

	return videoIDs, nil
}

func (m *postgresDBRepo) DeleteAllVideoIDByChannelID(id uint) error {
	var videoIDs []uint
	err := m.DB.Table("videos").Select("videos.id").Where("channel_id = ?", id).Find(&videoIDs).Error
	if err != nil {
		return errors.New("internal server error. Please try again")
	}

	for _, videoID := range videoIDs {
		go m.DeleteVideoWithRelatedData(videoID)
	}

	return nil
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

func (m *postgresDBRepo) ToggleSubscription(userID, channelID uint) (uint, error) {
	var subscription models.Subscription
	var subscribed uint = 0
	tsx := m.DB.Where("user_id = ? AND channel_id = ?", userID, channelID).First(&subscription)
	if tsx.Error != nil {
		subscription.UserID = userID
		subscription.ChannelID = channelID
		ts := m.DB.Create(&subscription)
		if ts.Error != nil {
			fmt.Println(ts.Error, "error")
			return 0, errors.New("failed to subscription the channel")
		}
		subscribed = 1
	} else {
		ts := m.DB.Unscoped().Delete(&subscription)
		if ts.Error != nil {
			return 0, errors.New("failed to  unsubscribe the channel")
		}
	}

	return subscribed, nil
}

func (m *postgresDBRepo) CreateContactUs(contactUs *models.ContactUs) error {
	result := m.DB.Create(&contactUs)
	if result.Error != nil {
		return errors.New("failed to create contact us")
	}

	return nil
}

func (m *postgresDBRepo) IsContactUsSubmitted(email string) bool {
	var contactUs models.ContactUs
	result := m.DB.Table("contact_us").Select("contact_us.id").Where("contact_us.email = ?", email).First(&contactUs)

	if result.Error != nil {
		return false
	}

	if contactUs.ID == 0 {
		return false
	}
	return true
}
