package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/handlers"
	"github.com/raihan2bd/vidverse/handlers/websocket"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	r.Use(cors.New(config))
	r.Use(gin.Logger())

	// r.Use(handlers.Methods.IsLoggedIn)
	v1 := r.Group("/api/v1")
	v1.GET("/", handlers.Methods.GetStatus)

	v1.POST("/auth/login", handlers.Methods.LoginHandler)
	v1.POST("/auth/signup", handlers.Methods.SignupHandler)

	v1.GET("/videos", handlers.Methods.HandleGetAllVideos)
	v1.POST("/videos", IsAdmin, handlers.Methods.HandleCreateVideo)
	v1.GET("/get_videos/:channelID", handlers.Methods.HandleGetVideosByChannelID)
	v1.GET("/videos/:videoID", HasToken, handlers.Methods.HandleGetSingleVideo)
	v1.DELETE("/videos/:videoID", handlers.Methods.HandleDeleteVidoe)
	v1.GET("/related_videos/:channelID", handlers.Methods.HandleGetRelatedVideos)
	v1.GET("/file/video/:videoID", handlers.Methods.StreamVideoBuff)

	v1.GET("/subscribed_channels/:channelID", IsLoggedIn, handlers.Methods.HandleGetSubscribedChannels)

	v1.POST("/comments", IsLoggedIn, handlers.Methods.HandleCreateOrUpdateComment)
	v1.DELETE("/comments/:commentID", IsLoggedIn, handlers.Methods.HandleDeleteComment)
	v1.GET("/comments/:videoID", handlers.Methods.HandleGetComments)

	v1.GET("/likes/:videoID", IsLoggedIn, handlers.Methods.HandleVideoLike)

	v1.GET("/channels", handlers.Methods.HandleGetChannels)
	v1.POST("/channels", isAuthor, handlers.Methods.HandleCreateChannel)
	v1.GET("/channels/:channelID", handlers.Methods.HandleGetChannel)
	v1.DELETE("/channels/:channelID", isAuthor, handlers.Methods.HandleDeleteChannel)

	// websocket handler
	v1.GET("/ws", websocket.Methods.WSHandler)

	return r
}
