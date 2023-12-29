package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/initializers"
	"github.com/raihan2bd/vidverse/models"
	validator "github.com/raihan2bd/vidverse/validators"
)

// Get Channels
func (m *Repo) HandleGetChannels(c *gin.Context) {
	var channels []models.CustomChannel
	userID := 1
	channels, err := m.App.DBMethods.GetChannels(userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 Channel not found!",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"channels": channels,
	})
}

// get single channel with videos
func (m *Repo) HandleGetChannel(c *gin.Context) {
	chanID, err := strconv.Atoi(c.Params.ByName("channelID"))

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "404 Channel not found!",
		})
		return
	}

	var channel *models.CustomChannelDTO
	channel, err = m.App.DBMethods.GetChannelByID(chanID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 Channel not found!",
		})
		return
	}

	// Send the response
	c.IndentedJSON(http.StatusOK, gin.H{
		"channel": channel,
	})
}

// create new channel
func (m *Repo) HandleCreateChannel(c *gin.Context) {
	// get user id from context
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// convert user id to uint
	userID := uint(user_id.(float64))

	// fetch user from db
	user, err := m.App.DBMethods.GetUserByID(userID)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	fmt.Println(user.UserRole)

	// check user role
	if user.UserRole != "author" {
		if user.UserRole != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to create channel"})
			return
		}
	}

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

	channel := models.Channel{Title: title, Description: description, Logo: resp.SecureURL, UserID: user.ID}

	var id uint
	id, err = m.App.DBMethods.CreateChannel(&channel)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Channel created successfully", "channel_id": id})
}

// delete channel
func (m *Repo) HandleDeleteChannel(c *gin.Context) {
	fmt.Println("Delete channel handler")
	channelID, err := strconv.Atoi(c.Param("channelID"))

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid channel id"})
		return
	}

	// delete channel
	deleteError := m.App.DBMethods.DeleteChannelByID(channelID)

	if deleteError != nil {
		c.JSON(deleteError.Status, gin.H{"error": deleteError.Err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Channel deleted successfully"})
}
