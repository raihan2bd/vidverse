package main

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	v1 := r.Group("/api/v1")
	v1.GET("/", GetStatus)

	return r
}
