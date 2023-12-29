package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/helpers"
)

// pass the user_id if logged in
func HasToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.Next()
		return
	}

	claims, err := helpers.DecodeToken(token)

	if err != nil {
		fmt.Println(err)
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

func IsLoggedIn(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.Next()
		return
	}

	claims, err := helpers.DecodeToken(token)

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

func IsAdmin(c *gin.Context) {
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

func isAuthor(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.Next()
		return
	}
	claims, err := helpers.DecodeToken(token)

	if err != nil {
		c.IndentedJSON(401, gin.H{
			"error": "unauthorized",
		})
		return
	}

	if !helpers.ValidateToken(claims) {
		if err != nil {
			c.IndentedJSON(401, gin.H{
				"error": "unauthorized",
			})
		}
		return
	}

	if claims["user_role"] != "author" {
		if claims["user_role"] != "admin" {
			c.IndentedJSON(403, gin.H{
				"error": "forbidden",
			})
			return
		}
	}

	c.Set("user_id", claims["sub"])
	c.Next()
}
