package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/config"
)

type Repo struct {
	App *config.Application
}

var Methods *Repo

func NewAPP(a *config.Application) *Repo {
	return &Repo{
		App: a,
	}
}

func NewSocket(m *Repo) {
	Methods = m
}

// handle websocket request
func (m *Repo) WSHandler(c *gin.Context) {

}
