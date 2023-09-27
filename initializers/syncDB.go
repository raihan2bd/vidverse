package initializers

import "github.com/raihan2bd/vidverse/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{}, &models.Channel{}, &models.Video{}, &models.Like{}, &models.Comment{})
}
