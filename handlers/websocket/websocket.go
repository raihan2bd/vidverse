package websocket

import "github.com/raihan2bd/vidverse/config"

var app *config.Application

func NewAPP(a *config.Application) {
	app = a
}
