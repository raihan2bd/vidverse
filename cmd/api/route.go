package main

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	v1 := r.Group("/api/v1")
	v1.GET("/", GetStatus)
	v1.POST("/uploads/video", UploadVideo)
	v1.DELETE("/videos/:videoID", DeleteVidoe)
	v1.GET("/videos/:videoID", GetSingleVideo)
	v1.GET("/file/video/:videoID", StreamVideoBuff)

	return r
}
