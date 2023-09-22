package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"status":      "Available",
		"version":     "1.0.0",
		"environment": "Development",
	})
}
