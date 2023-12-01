package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *application) NewRouter() *gin.Engine {
	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true

	r.Use(cors.New(config))
	r.Use(gin.Logger())

	// r.Use(app.IsLoggedIn)
	v1 := r.Group("/api/v1")
	v1.GET("/", GetStatus)
	v1.POST("/auth/login", app.LoginHandler)
	v1.POST("/auth/signup", app.SignupHandler)
	v1.GET("/videos", app.HandleGetAllVideos)
	v1.POST("/videos", app.IsAdmin, UploadVideo)
	v1.GET("/videos/:videoID", app.HandleGetSingleVideo)
	v1.DELETE("/videos/:videoID", app.HandleDeleteVidoe)
	v1.GET("/related_videos/:channelID", app.HandleGetRelatedVideos)
	v1.GET("/comments/:videoID", app.HandleGetComments)
	v1.GET("/file/video/:videoID", StreamVideoBuff)
	v1.GET("/channels", app.HandleGetChannels)
	v1.GET("/channels/:channelID", app.HandleGetChannel)

	return r
}
