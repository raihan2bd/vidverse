package main

import (
	"github.com/raihan2bd/vidverse/config"
	"github.com/raihan2bd/vidverse/handlers"
	"github.com/raihan2bd/vidverse/initializers"
	dbrepo "github.com/raihan2bd/vidverse/repository/dbRepo"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.ConnectToCloudinary()
	initializers.SyncDatabase()
}

func main() {
	app := config.Application{}
	repo := dbrepo.NewPostgresRepo(initializers.DB, initializers.CLD)
	app.DB = repo
	app.CLD = initializers.CLD
	handlers.NewAPP(&app)
	r := NewRouter()

	r.Run()
}
