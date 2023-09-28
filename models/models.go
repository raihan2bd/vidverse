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
	Name     string `gorm:"type:varchar(100);not null" json:"name" binding:"required,min=3,max=100"`
	Email    string `gorm:"type:varchar(255);unique;not null" json:"email" binding:"required,email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
}

type Channel struct {
	CustomModel
	Title       string  `gorm:"type:varchar(100)" json:"title" binding:"required,min=5,max=100"`
	Description string  `gorm:"type:text;size:500" json:"description" binding:"required,min=25,max=100"`
	UserID      uint    `json:"user_id"`
	User        User    `gorm:"foreignKey:UserID" json:"user"`
	Videos      []Video `gorm:"foreignKey:ChannelID" json:"videos"`
}

type Video struct {
	CustomModel
	Title       string    `gorm:"type:varchar(255);not null" json:"title" binding:"required,min=2,max=255"`
	Description string    `gorm:"type:text;size:500;not null" json:"description" binding:"required,min=2,max=500"`
	PublicID    string    `gorm:"type:varchar(255);not null" json:"-"`
	SecureURL   string    `gorm:"type:varchar(255);not null" json:"secure_url"`
	ChannelID   uint      `json:"channel_id"`
	Channel     Channel   `gorm:"foreignKey:ChannelID" json:"channel"`
	Thumb       string    `gorm:"type:varchar(255)" json:"tumb"`
	Likes       []Like    `json:"likes"`
	Comments    []Comment `json:"comments"`
}

type Like struct {
	CustomModel
	UserID  uint  `json:"user_id"`
	VideoID uint  `json:"video_id"`
	Video   Video `gorm:"foreignKey:VideoID"`
	User    User  `gorm:"foreignKey:UserID"`
}

type Comment struct {
	CustomModel
	Text    string `gorm:"type:text;size:500" json:"text"`
	UserID  uint   `json:"user_id"`
	VideoID uint   `json:"video_id"`
	Video   Video  `gorm:"foreignKey:VideoID"`
	User    User   `gorm:"foreignKey:UserID"`
}
