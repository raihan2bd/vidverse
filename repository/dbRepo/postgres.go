package dbrepo

import (
	"errors"

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
