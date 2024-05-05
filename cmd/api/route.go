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
	// config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowOrigins = []string{"*"}

	config.AllowCredentials = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	r.Use(cors.New(config))
	r.Use(gin.Logger())

	// r.Use(handlers.Methods.IsLoggedIn)
	v1 := r.Group("/api/v1")
	v1.GET("/", handlers.Methods.GetStatus)

	v1.POST("/auth/login", handlers.Methods.LoginHandler)
	v1.POST("/auth/signup", handlers.Methods.SignupHandler)
	v1.POST("/auth/request-password-reset", handlers.Methods.RequestPasswordReset)
	v1.POST("/auth/verify-password-reset", handlers.Methods.VerifyPasswordReset)
	v1.POST("/auth/reset-password", handlers.Methods.ResetPassword)
	// v1.POST("/auth/request_forgot_password", handlers.Methods.RequestForgotPassword)

	// social auth
	v1.POST("/auth/social_login", handlers.Methods.SocialLogin)

	v1.GET("/videos", handlers.Methods.HandleGetAllVideos)
	v1.POST("/videos", isAuthor, handlers.Methods.HandleCreateVideo)
	v1.POST("/videos/:videoID", isAuthor, handlers.Methods.HandleUpdateVideo)
	v1.GET("/get_videos/:channelID", handlers.Methods.HandleGetVideosByChannelID)
	v1.GET("/videos/:videoID", HasToken, handlers.Methods.HandleGetSingleVideo)
	v1.DELETE("/videos/:videoID", isAuthor, handlers.Methods.HandleDeleteVideo)
	v1.GET("/related_videos/:channelID", handlers.Methods.HandleGetRelatedVideos)
	v1.GET("/file/video/:videoID", handlers.Methods.StreamVideoBuff)
	// insert watch later video
	v1.POST("/watch_later", IsLoggedIn, handlers.Methods.HandleWatchLater)
	v1.GET("/watch_later", IsLoggedIn, handlers.Methods.HandleGetWatchLater)

	// routes for shots
	v1.POST("/shots", isAuthor, handlers.Methods.HandleCreateShot)
	// v1.GET("/shots", handlers.Methods.HandleGetAllShots)
	// v1.GET("/shots/:shortID", handlers.Methods.HandleGetSingleShot)
	// v1.PATCH("/shots/:shotID", isAuthor, handlers.Methods.HandleUpdateShot)
	// v1.DELETE("/shots/:shotID", isAuthor, handlers.Methods.HandleDeleteShot)

	v1.GET("/subscribed_channels/:channelID", IsLoggedIn, handlers.Methods.HandleGetSubscribedChannels)
	v1.GET("/notifications", IsLoggedIn, handlers.Methods.HandleGetNotifications)
	v1.PATCH("/notifications/:notificationID", IsLoggedIn, handlers.Methods.HandleUpdateNotification)

	v1.POST("/comments", IsLoggedIn, handlers.Methods.HandleCreateOrUpdateComment)
	v1.DELETE("/comments/:commentID", IsLoggedIn, handlers.Methods.HandleDeleteComment)
	v1.GET("/comments/:videoID", handlers.Methods.HandleGetComments)

	v1.GET("/likes/:videoID", IsLoggedIn, handlers.Methods.HandleVideoLike)
	v1.GET("/liked_videos", IsLoggedIn, handlers.Methods.HandleGetLikedVideos)

	v1.GET("/channels", isAuthor, handlers.Methods.HandleGetChannels)
	v1.GET("/channels_by_user_with_details", isAuthor, handlers.Methods.HandleGetChannelsWithDetailsByUserID)
	v1.POST("/channels", isAuthor, handlers.Methods.HandleCreateChannel)
	v1.PATCH("/channels/:channelID", isAuthor, handlers.Methods.HandleEditChannel)
	v1.GET("/channels/:channelID", handlers.Methods.HandleGetChannel)
	v1.DELETE("/channels/:channelID", isAuthor, handlers.Methods.HandleDeleteChannel)
	v1.GET("/get_channel_videos/:channelID", handlers.Methods.HandleGetChannelsVideos)
	v1.GET("/get_channel_with_details/:channelID", HasToken, handlers.Methods.HandleGetChannelWithDetails)

	v1.POST("/contact_us", HasToken, handlers.Methods.HandleContactUs)

	// websocket handler
	v1.GET("/ws", websocket.Methods.WSHandler)

	return r
}
