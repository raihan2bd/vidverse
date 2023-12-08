package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/config"
)

var app *config.Application

func NewAPP(a *config.Application) {
	app = a
}

func GetStatus(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"status":      "Available",
		"version":     "1.0.0",
		"environment": "Development",
	})
}
