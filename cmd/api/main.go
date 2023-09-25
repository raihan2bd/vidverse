package main

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/raihan2bd/vidverse/initializers"
	"github.com/raihan2bd/vidverse/repository"
	dbrepo "github.com/raihan2bd/vidverse/repository/dbRepo"
)

type application struct {
	DB  repository.DatabaseRepo
	CLD *cloudinary.Cloudinary
}

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.ConnectToCloudinary()
	initializers.SyncDatabase()
}

func main() {
	app := &application{}
	repo := dbrepo.NewPostgresRepo(initializers.DB, initializers.CLD)
	app.DB = repo
	app.CLD = initializers.CLD
	r := app.NewRouter()

	r.Run()
}
