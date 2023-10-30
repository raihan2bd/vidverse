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
func (m *postgresDBRepo) GetAllVideos() ([]models.Video, error) {
	var videos []models.Video
	if err := m.DB.Find(&videos).Error; err != nil {
		return nil, errors.New("internal server error. Please try again")
	}

	return videos, nil
}

// Get single video by Id
func (m *postgresDBRepo) GetVidoeByID(id int) (*models.Video, error) {

	var video models.Video
	if id > 0 {
		m.DB.First(&video, "id = ?", id)

		if video.ID == 0 {
			return nil, errors.New("404 video not found")
		}
	}

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
