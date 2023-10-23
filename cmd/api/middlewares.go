package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (app *application) IsLoggedIn(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Logged-In user", tokenString)
	c.Next()
}
