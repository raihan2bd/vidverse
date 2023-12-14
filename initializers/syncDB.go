package initializers

import (
	"log"

	"errors"

	"github.com/raihan2bd/vidverse/models"
)

func SyncDatabase() error {
	err := DB.AutoMigrate(&models.User{}, &models.Channel{}, &models.Video{}, &models.Like{}, &models.Comment{}, &models.Subscription{}, &models.Notification{})

	if err != nil {
		log.Println(err)
		return errors.New("failed to sync database")
	}

	return nil

	// isSeeded := checkDatabaseSeed()
	// if !isSeeded {
	// 	seedUsers()
	// 	seedChannels()
	// 	seedVideos()
	// 	seedLikes()
	// 	seedComments()
	// }
}

func checkDatabaseSeed() bool {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	return count > 0
}

func seedUsers() {
	users := []models.User{
		{
			Name:     "Alice",
			Email:    "alice@example.com",
			Password: "Password@123",
		},
		{
			Name:     "Bob",
			Email:    "bob@example.com",
			Password: "Password@123",
		},
	}

	for _, user := range users {
		DB.Create(&user)
	}
}

func seedChannels() {
	channels := []models.Channel{
		{
			Title:       "Channel 1",
			Description: "Description for Channel 1",
			UserID:      1,
		},
		{
			Title:       "Channel 2",
			Description: "Description for Channel 2",
			UserID:      2,
		},
	}

	for _, channel := range channels {
		DB.Create(&channel)
	}
}

func seedVideos() {
	videos := []models.Video{
		{
			Title:       "Video 1",
			Description: "Description for Video 1. This is dummy video description.",
			SecureURL:   "https://res.cloudinary.com/dog87elav/video/upload/v1695878153/vidverse/uploads/videos/d8kukajtsusdsfihoeoq.mp4",
			ChannelID:   1,
		},
		{
			Title:       "Video 2",
			Description: "Description for Video 2. This is dummy video description.",
			SecureURL:   "https://res.cloudinary.com/dog87elav/video/upload/v1695878153/vidverse/uploads/videos/d8kukajtsusdsfihoeoq.mp4",
			ChannelID:   2,
		},
	}

	for _, video := range videos {
		DB.Create(&video)
	}
}

func seedLikes() {
	likes := []models.Like{
		{
			UserID:  1,
			VideoID: 1,
		},
		{
			UserID:  2,
			VideoID: 2,
		},
	}

	for _, like := range likes {
		DB.Create(&like)
	}
}

func seedComments() {
	comments := []models.Comment{
		{
			Text:    "Comment 1 on Video 1",
			UserID:  1,
			VideoID: 1,
		},
		{
			Text:    "Comment 2 on Video 2",
			UserID:  2,
			VideoID: 2,
		},
	}

	for _, comment := range comments {
		DB.Create(&comment)
	}
}
