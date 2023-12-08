package config

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/raihan2bd/vidverse/repository"
)

type Application struct {
	DB  repository.DatabaseRepo
	CLD *cloudinary.Cloudinary
}
