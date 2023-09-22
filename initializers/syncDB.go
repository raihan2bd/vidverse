package initializers

import "github.com/raihan2bd/vidverse/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Video{})
}
