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
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	// UserName string `gorm:"type:varchar(100);unique;not null" json:"username"`
	Email    string `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Avatar   string `gorm:"type:varchar(255);not null;default:'https://upload.wikimedia.org/wikipedia/commons/5/59/User-avatar.svg'" json:"avatar"`
	IsActive bool   `gorm:"type:boolean;not null;default:false" json:"is_active"`
	UserRole string `gorm:"type:varchar(150);not null;default:'user'" json:"user_role"`
}

type UserPayload struct {
	CustomModel
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	// UserName string `gorm:"type:varchar(100);unique;not null" json:"username"`
	Email    string `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"password"`
	Avatar   string `gorm:"type:varchar(255);not null;default:'https://upload.wikimedia.org/wikipedia/commons/5/59/User-avatar.svg'" json:"avatar"`
	IsActive bool   `gorm:"type:boolean;not null;default:false" json:"is_active"`
	UserRole string `gorm:"type:varchar(150);not null;default:'user'" json:"user_role"`
}

type Channel struct {
	CustomModel
	Title        string  `gorm:"type:varchar(100)" json:"title"`
	Description  string  `gorm:"type:text;size:500" json:"description"`
	Logo         string  `gorm:"type:varchar(255);" json:"logo"`
	UserID       uint    `json:"user_id"`
	LogoPublicID string  `gorm:"type:varchar(255);not null" json:"-"`
	User         User    `gorm:"foreignKey:UserID" json:"user"`
	Videos       []Video `gorm:"foreignKey:ChannelID" json:"videos"`
}

type Subscription struct {
	CustomModel
	UserID    uint `gorm:"foreignKey:UserID" json:"user_id"`
	ChannelID uint `gorm:"foreignKey:ChannelID" json:"channel_id"`
}

type ChannelPayload struct {
	ID          uint   `json:"id"`
	Title       string `gorm:"type:varchar(100)" json:"title,omitempty"`
	Description string `gorm:"type:text;size:500" json:"description,omitempty"`
	Logo        string `gorm:"type:varchar(255);" json:"logo,omitempty"`
	UserID      uint   `json:"user_id,omitempty"`
}

type Video struct {
	CustomModel
	Title       string    `gorm:"type:varchar(255);not null" json:"title,omitempty" binding:"required,min=2,max=255"`
	Description string    `gorm:"type:text;size:500;not null" json:"description,omitempty" binding:"required,min=2,max=500"`
	PublicID    string    `gorm:"type:varchar(255);not null" json:"-"`
	SecureURL   string    `gorm:"type:varchar(255);not null" json:"secure_url,omitempty"`
	ChannelID   uint      `json:"channel_id,omitempty"`
	Channel     Channel   `gorm:"foreignKey:ChannelID" json:"channel,omitempty"`
	Thumb       string    `gorm:"type:varchar(255)" json:"tumb,omitempty"`
	Likes       []Like    `json:"likes,omitempty"`
	Comments    []Comment `json:"comments,omitempty"`
	Views       int64     `gorm:"type:bigint;not null;default:0" json:"views,omitempty"`
}

type VideoDTO struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	Thumb        string `json:"thumb"`
	Views        int64  `json:"views"`
	ChannelID    uint   `json:"channel_id"`
	ChannelTitle string `json:"channel_title"`
	ChannelLogo  string `json:"channel_logo"`
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

type CommentDTO struct {
	ID         uint   `json:"id"`
	Text       string `json:"text"`
	VideoID    uint   `json:"video_id"`
	UserID     uint   `json:"user_id"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar"`
}

type CustomChannel struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Logo  string `json:"logo"`
}

type CustomChannelDTO struct {
	ID              uint   `json:"id,omitempty"`
	Title           string `json:"title,omitempty"`
	Logo            string `json:"logo,omitempty"`
	Description     string `json:"description,omitempty"`
	TotalSubscriber int64  `json:"total_subscriber,omitempty"`
	TotalVideo      int64  `json:"total_video,omitempty"`
}

type Notification struct {
	CustomModel
	IsRead     bool   `gorm:"type:boolean;not null;default:false" json:"is_read"`
	ReceiverID uint   `gorm:"foreignKey:UserID; not null;" json:"receiver_id"`
	SenderID   uint   `gorm:"foreignKey:UserID; not null;" json:"sender_id"`
	SenderName string `gorm:"type:varchar(100);not null;" json:"sender_name"`
	VideoID    uint   `gorm:"foreignKey:VideoID" json:"video_id,omitempty"`
	ChannelID  uint   `gorm:"foreignKey:ChannelID" json:"channel_id,omitempty"`
	CommentID  uint   `gorm:"foreignKey:CommentID" json:"comment_id,omitempty"`
	LikeID     uint   `gorm:"foreignKey:LikeID" json:"like_id,omitempty"`
	Type       string `gorm:"type:varchar(100);not null;" json:"type"`
}
