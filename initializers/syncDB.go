package initializers

import (
	"log"
	"os"

	"errors"

	"github.com/raihan2bd/vidverse/models"
)

func SyncDatabase() error {
	err := DB.AutoMigrate(&models.User{}, &models.Channel{}, &models.Video{}, &models.Like{}, &models.Comment{}, &models.Subscription{}, &models.Notification{}, &models.ContactUs{}, &models.Token{})

	if err != nil {
		log.Println(err)
		return errors.New("failed to sync database")
	}

	env := os.Getenv("ENVIRONMENT")

	if env == "development" {
		if !checkUserSeed() {
			seedUsers()
		}
		if !checkChannelSeed() {
			seedChannels()
		}
		if !checkVideoSeed() {
			seedVideos()
		}
		if !checkLikeSeed() {
			seedLikes()
		}
		if !checkCommentSeed() {
			seedComments()
		}
	}

	return nil

}

func checkUserSeed() bool {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	return count > 0
}

func checkChannelSeed() bool {
	var count int64
	DB.Model(&models.Channel{}).Count(&count)
	return count > 0
}

func checkVideoSeed() bool {
	var count int64
	DB.Model(&models.Video{}).Count(&count)
	return count > 0
}

func checkLikeSeed() bool {
	var count int64
	DB.Model(&models.Like{}).Count(&count)
	return count > 0
}

func checkCommentSeed() bool {
	var count int64
	DB.Model(&models.Comment{}).Count(&count)
	return count > 0
}

func seedUsers() {
	users := []models.User{
		{
			Name:     "Admin",
			Email:    "admin@test.com",
			Password: "$2a$12$tNErjg8dC6nRDPE9jU5Vj.nupSFbl0l6Hc4rCkQNVcUoKapiSkug2", // Admin@123
			UserRole: "admin",
		},
		{
			Name:     "Author",
			Email:    "author@test.com",
			Password: "$2a$12$OFZmsYtt7chRQ8zl8Swt/OHiWyAiFT.yREQGUSKBMMFnjSh2g6quW", // Pass@123
			UserRole: "author",
		},
		{
			Name:     "User",
			Email:    "user@test.com",
			Password: "$2a$12$OFZmsYtt7chRQ8zl8Swt/OHiWyAiFT.yREQGUSKBMMFnjSh2g6quW", // Pass@123
			UserRole: "user",
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
			SecureURL:   "https://res.cloudinary.com/dog87elav/video/upload/v1695878153/vidverse/uploads/seed-videos/robomydj09pndkde2iuk.mp4",
			ChannelID:   1,
			Thumb:       "https://res.cloudinary.com/dog87elav/video/upload/vidverse/uploads/seed-videos/robomydj09pndkde2iuk.jpeg",
			PublicID:    "vidverse/uploads/videos/robomydj09pndkde2iuk",
		},
		{
			Title:       "Video 2",
			Description: "Description for Video 2. This is dummy video description.",
			SecureURL:   "https://res.cloudinary.com/dog87elav/video/upload/v1695878153/vidverse/uploads/seed-videos/cx0g25tcjzqursnxyz2k.mp4",
			ChannelID:   1,
			Thumb:       "https://res.cloudinary.com/dog87elav/video/upload/vidverse/uploads/seed-videos/cx0g25tcjzqursnxyz2k.jpeg",
			PublicID:    "vidverse/uploads/videos/cx0g25tcjzqursnxyz2k",
		},
		{
			Title:       "Video 3",
			Description: "Description for Video 3. This is dummy video description.",
			SecureURL:   "https://res.cloudinary.com/dog87elav/video/upload/v1695878153/vidverse/uploads/seed-videos/robomydj09pndkde2iuk.mp4",
			ChannelID:   1,
			Thumb:       "https://res.cloudinary.com/dog87elav/video/upload/vidverse/uploads/seed-videos/robomydj09pndkde2iuk.jpeg",
			PublicID:    "vidverse/uploads/videos/robomydj09pndkde2iuk",
		},
		{
			Title:       "Video 4",
			Description: "Description for Video 4. This is dummy video description.",
			SecureURL:   "https://res.cloudinary.com/dog87elav/video/upload/v1695878153/vidverse/uploads/seed-videos/jvvcnicn3u0mvmhn64gu_ldllzy.mp4",
			ChannelID:   2,
			Thumb:       "https://res.cloudinary.com/dog87elav/video/upload/vidverse/uploads/seed-videos/jvvcnicn3u0mvmhn64gu_ldllzy.jpeg",
			PublicID:    "vidverse/uploads/videos/jvvcnicn3u0mvmhn64gu_ldllzy",
		},
		{
			Title:       "Video 5",
			Description: "Description for Video 5. This is dummy video description.",
			SecureURL:   "https://res.cloudinary.com/dog87elav/video/upload/v1695878153/vidverse/uploads/seed-videos/s1u8h4hqwkkxmoovi1s0.mp4",
			ChannelID:   2,
			Thumb:       "https://res.cloudinary.com/dog87elav/video/upload/vidverse/uploads/seed-videos/s1u8h4hqwkkxmoovi1s0.jpeg",
			PublicID:    "vidverse/uploads/videos/s1u8h4hqwkkxmoovi1s0",
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
