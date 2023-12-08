package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/handlers"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true

	r.Use(cors.New(config))
	r.Use(gin.Logger())

	// r.Use(handlers.IsLoggedIn)
	v1 := r.Group("/api/v1")
	v1.GET("/", handlers.GetStatus)
	v1.POST("/auth/login", handlers.LoginHandler)
	v1.POST("/auth/signup", handlers.SignupHandler)
	v1.GET("/videos", handlers.HandleGetAllVideos)
	// v1.POST("/videos", handlers.IsAdmin, handlers.HandleCreateVideo)
	v1.POST("/videos", handlers.HandleCreateVideo)
	v1.GET("/get_videos/:channelID", handlers.HandleGetVideosByChannelID)
	v1.GET("/videos/:videoID", handlers.HandleGetSingleVideo)
	v1.DELETE("/videos/:videoID", handlers.HandleDeleteVidoe)
	v1.GET("/related_videos/:channelID", handlers.HandleGetRelatedVideos)
	v1.GET("/comments/:videoID", handlers.HandleGetComments)
	v1.GET("/file/video/:videoID", handlers.StreamVideoBuff)

	v1.GET("/channels", handlers.HandleGetChannels)
	v1.POST("/channels", handlers.HandleCreateChannel)
	v1.GET("/channels/:channelID", handlers.HandleGetChannel)
	v1.DELETE("/channels/:channelID", handlers.HandleDeleteChannel)

	return r
}
