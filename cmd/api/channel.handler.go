package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/initializers"
	"github.com/raihan2bd/vidverse/models"
	validator "github.com/raihan2bd/vidverse/validators"
)

// create new channel
func (app *application) HandleCreateChannel(c *gin.Context) {
	fmt.Println("Create channel handler")
	// get channel logo from form
	channelLogo, logoHeader, err := c.Request.FormFile("logo")
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(400, gin.H{"error": "Logo is required."})
		return
	}

	defer channelLogo.Close()

	// validate logo
	validator := validator.New()
	validator.IsImage(logoHeader.Header.Get("Content-Type"), "logo")
	validator.IsImageSize(logoHeader.Size, 5*1024*1024, "logo")

	// get channel title from form
	title := c.PostForm("title")
	validator.Required(title, "title", "title is required.")
	validator.IsLength(title, "title", 5, 255)

	// get channel description from form
	description := c.PostForm("description")
	validator.Required(description, "description", "description is required")
	validator.IsLength(description, "description", 25, 500)

	if !validator.Valid() {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": validator.GetErrMsg(),
		})
		return
	}

	// Define the upload directory
	folder := "vidverse/uploads/images"

	ctx := context.Background()

	// Upload the file to Cloudinary with specified folder and transformations
	resp, err := initializers.CLD.Upload.Upload(ctx, channelLogo, uploader.UploadParams{
		Folder: folder,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	channel := models.Channel{Title: title, Description: description, Logo: resp.SecureURL, UserID: 1}

	result := initializers.DB.Create(&channel)

	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create the channel",
		})

		return
	}

	c.JSON(200, gin.H{"message": "Channel created successfully", "channel_id": channel.ID})

}
