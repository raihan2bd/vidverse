package dbrepo

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/raihan2bd/vidverse/repository"
	"gorm.io/gorm"
)

type postgresDBRepo struct {
	DB  *gorm.DB
	CLD *cloudinary.Cloudinary
}

func NewPostgresRepo(db *gorm.DB, cld *cloudinary.Cloudinary) repository.DatabaseRepo {
	return &postgresDBRepo{
		DB:  db,
		CLD: cld,
	}
}
