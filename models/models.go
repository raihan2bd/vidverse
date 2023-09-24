package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomModel struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type User struct {
	CustomModel
	Name     string
	Email    string
	Password string `json:"-"`
}

type Video struct {
	CustomModel
	Title       string `json:"title"`
	Description string `json:"description"`
	PublicID    string `json:"-"`
	SecureURL   string `json:"secure_url"`
}
