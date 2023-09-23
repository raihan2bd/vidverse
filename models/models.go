package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string `json:"-"`
}

type Video struct {
	gorm.Model
	Title       string
	Description string
	PublicID    string
	SecureURL   string
}
