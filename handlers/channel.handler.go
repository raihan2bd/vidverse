package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/helpers"
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
	if user.UserRole != "admin" {
		if user.UserRole != "author" {
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

	ctx := context.Background()

	// upload logo to cloudinary
	var (
		secureURL string
		publicID  string
	)

	secureURL, publicID, err = helpers.UploadImageToCloudinary(ctx, m.App.CLD, channelLogo)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload logo"})
		return
	}

	channel := models.Channel{Title: title, Description: description, Logo: secureURL, UserID: user.ID, LogoPublicID: publicID}

	var id uint
	id, err = m.App.DBMethods.CreateChannel(&channel)
	if err != nil {
		// delete logo from cloudinary if exists
		ctx := context.Background()
		_ = helpers.DeleteImageFromCloudinary(ctx, m.App.CLD, publicID)

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Channel created successfully", "channel_id": id})
}

// edit channel
func (m *Repo) HandleEditChannel(c *gin.Context) {
	channelID, err := strconv.Atoi(c.Param("channelID"))
	if err != nil {
		c.JSON(400, gin.H{"error": "The channel you are trying to edit does not exist."})
		return
	}

	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// channel data from form
	title := c.PostForm("title")
	description := c.PostForm("description")
	logoURL := c.PostForm("logo_url")

	// channel logo from form
	channelLogo, logoHeader, err := c.Request.FormFile("logo")
	if err != nil {
		fmt.Println(err)
		if logoURL == "" {
			c.JSON(400, gin.H{"error": "Logo is required."})
			return
		}
	} else {
		defer channelLogo.Close()
	}

	// validate inputs
	validator := validator.New()
	validator.IsLength(title, "title", 5, 255)
	validator.IsLength(description, "description", 25, 500)

	// validate logo if exists
	if logoHeader != nil {
		validator.IsImage(logoHeader.Header.Get("Content-Type"), "logo")
		validator.IsImageSize(logoHeader.Size, 5*1024*1024, "logo")
	}

	if !validator.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validator.GetErrMsg(),
		})
		return
	}

	// convert user id to uint fetch user from db
	userID := uint(user_id.(float64))
	var user *models.User
	user, err = m.App.DBMethods.GetUserByID(userID)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// get channel by id
	channel, err := m.App.DBMethods.GetChannelByID(channelID)
	if err != nil {
		c.JSON(500, gin.H{"error": "The channel you are trying to edit does not exist."})
		return
	}

	// check user role
	if user.UserRole != "admin" {
		if user.ID != channel.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to edit channel"})
			return
		}
	}

	// check if data is the same of not
	if title == channel.Title && description == channel.Description && logoURL == channel.Logo {
		if logoHeader == nil {
			c.JSON(200, gin.H{"message": "Your channel is already up to date!"})
			return
		}
	}

	// upload image to cloudinary
	var (
		secureURL string
		publicID  string
	)
	oldPublicID := channel.LogoPublicID

	if logoHeader != nil && channelLogo != nil {
		ctx := context.Background()
		uploadPath := "vidverse/uploads/channel_logos"
		secureURL, publicID, err = helpers.UploadImageToCloudinary(ctx, m.App.CLD, channelLogo, uploadPath)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to upload logo"})
			return
		}
	}

	if publicID != "" && secureURL != "" {
		channel.LogoPublicID = publicID
		channel.Logo = secureURL
	}

	if title != "" && title != channel.Title {
		channel.Title = title
	}

	if description != "" && description != channel.Description {
		channel.Description = description
	}

	err = m.App.DBMethods.UpdateChannel(channel)
	if err != nil {
		// delete logo from cloudinary if exists
		if secureURL != "" {
			ctx := context.Background()
			_ = helpers.DeleteImageFromCloudinary(ctx, m.App.CLD, publicID)
		}

		c.JSON(500, gin.H{"error": "Failed to update channel"})
		return
	}

	c.JSON(201, gin.H{"message": "Channel updated successfully!"})

	if oldPublicID != channel.LogoPublicID {
		ctx := context.Background()
		err = helpers.DeleteImageFromCloudinary(ctx, m.App.CLD, oldPublicID)
		if err != nil {
			return
		}
	}

}

// delete channel
func (m *Repo) HandleDeleteChannel(c *gin.Context) {
	channelID, err := strconv.Atoi(c.Param("channelID"))

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid channel id"})
		return
	}

	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// convert user id to uint
	userID := uint(user_id.(float64))
	var user *models.User
	user, err = m.App.DBMethods.GetUserByID(userID)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// get channel by id
	channel, err := m.App.DBMethods.GetChannelByID(channelID)
	if err != nil {
		c.JSON(500, gin.H{"error": "The channel you are trying to delete does not exist."})
		return
	}

	// check user role
	if user.UserRole != "admin" {
		if user.ID != channel.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to delete channel"})
			return
		}
	}

	// delete channel
	deleteError := m.App.DBMethods.DeleteChannelByID(channelID)

	if deleteError != nil {
		c.JSON(deleteError.Status, gin.H{"error": deleteError.Err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Channel deleted successfully"})
}
