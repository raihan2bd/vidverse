package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/raihan2bd/vidverse/config"
	"github.com/raihan2bd/vidverse/handlers"
	"github.com/raihan2bd/vidverse/handlers/websocket"
)

var app *config.Application

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app, err = config.LoadConfig()
	if err != nil {
		panic(err)
	}

	repo := handlers.NewAPP(app)
	handlers.NewHandler(repo)
	socketRepo := websocket.NewAPP(app)
	websocket.NewSocket(socketRepo)
	r := NewRouter()

	r.Run()
}
