package main

import (
	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/helpers"
)

// pass the user_id if logged in
func (app *application) HasToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.Next()
		return
	}

	claims, err := helpers.DecodeToken(token)

	if err != nil {
		c.Next()
		return
	}

	if !helpers.ValidateToken(claims) {
		c.Next()
		return
	}

	c.Set("user_id", claims["sub"])
	c.Next()
}

func (app *application) IsLoggedIn(c *gin.Context) {
	claims, err := helpers.DecodeToken(c.Request.Header.Get("Authorization"))

	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	if !helpers.ValidateToken(claims) {
		c.AbortWithStatus(401)
		return
	}

	c.Set("user_id", claims["sub"])
	c.Next()
}

func (app *application) IsAdmin(c *gin.Context) {
	claims, err := helpers.DecodeToken(c.Request.Header.Get("Authorization"))

	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	if !helpers.ValidateToken(claims) {
		c.AbortWithStatus(401)
		return
	}

	if claims["user_role"] != "admin" {
		c.AbortWithStatus(403)
		return
	}

	c.Set("user_id", claims["sub"])
	c.Next()
}
