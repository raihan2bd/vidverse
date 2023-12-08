package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/config"
)

var Methods *Repo

type Repo struct {
	App *config.Application
}

func NewAPP(a *config.Application) *Repo {
	return &Repo{
		App: a,
	}
}

func NewHandler(m *Repo) {
	Methods = m
}

func (m *Repo) GetStatus(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"status":      "Available",
		"version":     "1.0.0",
		"environment": "Development",
	})
}
