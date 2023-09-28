package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *application) NewRouter() *gin.Engine {
	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.Use(gin.Logger())

	v1 := r.Group("/api/v1")
	v1.GET("/", GetStatus)
	v1.GET("/videos", app.HandleGetAllVideos)
	v1.POST("/videos", UploadVideo)
	v1.GET("/videos/:videoID", app.HandleGetSingleVideo)
	v1.DELETE("/videos/:videoID", app.HandleDeleteVidoe)
	v1.GET("/file/video/:videoID", StreamVideoBuff)

	return r
}
