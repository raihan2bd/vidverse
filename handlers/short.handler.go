package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/helpers"
	"github.com/raihan2bd/vidverse/models"
	validator "github.com/raihan2bd/vidverse/validators"
)

func (m *Repo) HandleCreateShot(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Access denied! Please login first",
		})
		return
	}

	userID := uint(user_id.(float64))
	var user *models.User
	user, err := m.App.DBMethods.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Access denied! Please login first",
		})
		return
	}

	// check the user role
	if user.UserRole != "admin" {
		if user.UserRole != "author" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied! You are not allowed to upload video",
			})
			return
		}
	}

	videoFile, fileInfo, err := c.Request.FormFile("video")
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "File is required."})
		return
	} else if fileInfo == nil {
		defer videoFile.Close()
		c.IndentedJSON(400, gin.H{"error": "video is required."})
		return
	}

	var thumbSecureURL, thumbPublicID string
	thumbFile, thumbFileInfo, err := c.Request.FormFile("thumb")
	if err != nil {
	} else if thumbFileInfo != nil {
		defer thumbFile.Close()
	}

	validator := validator.New()

	// validate video
	validator.IsVideo(fileInfo.Header.Get("Content-Type"), "video")
	validator.IsVideoSize(fileInfo.Size, 50*1024*1024, "video")

	// Todo: add image upload system later

	title := c.PostForm("title")
	description := c.PostForm("description")
	channel_id := c.PostForm("channel_id")

	channelID, err := strconv.Atoi(channel_id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid channel id",
		})
		return
	}

	validator.Required(title, "title", "title is required.")
	validator.IsLength(title, "title", 5, 255)
	validator.Required(description, "description", "description is required")
	validator.IsLength(description, "description", 25, 500)
	validator.Required(channel_id, "channel_id", "channel_id is required")

	if thumbFileInfo != nil && thumbFile != nil {
		validator.IsImage(thumbFileInfo.Header.Get("Content-Type"), "thumb")
		validator.IsImageSize(thumbFileInfo.Size, 5*1024*1024, "thumb")
	}

	if !validator.Valid() {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": validator.GetErrMsg(),
		})
		return
	}

	// check the channel is available or not
	var channel *models.CustomChannelDTO
	channel, err = m.App.DBMethods.GetChannelByID(channelID)
	if err != nil || channel.ID == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "The channel you want to upload video is not found!",
		})
		return
	}

	// check if the channel user is the same or not
	if channel.UserID != userID {
		if user.UserRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied! You are not allowed to upload video to this channel",
			})
			return
		}
	}

	// upload video to cloudinary
	ctx := context.Background()
	var secureURL, videoPublicID string
	secureURL, videoPublicID, err = helpers.UploadShotToCloudinary(ctx, m.App.CLD, videoFile)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload the video",
		})
		return
	}

	if thumbFileInfo != nil && thumbFile != nil {
		// upload thumb to cloudinary
		thumbSecureURL, thumbPublicID, err = helpers.UploadImageToCloudinary(ctx, m.App.CLD, thumbFile, "vidverse/uploads/thumbs")
		if err != nil {
			thumbSecureURL = m.generateThumbURL(videoPublicID)
		}
	} else {
		// generate thumb url
		thumbSecureURL = m.generateThumbURL(videoPublicID)
	}

	video := models.Shot{Title: title, Description: description, PublicID: videoPublicID, SecureURL: secureURL, ChannelID: channel.ID, Thumb: thumbSecureURL, ThumbPublicID: thumbPublicID}

	videoID, err := m.App.DBMethods.CreateShot(&video)
	if err != nil {
		// delete thumbnail from cloudinary
		_ = helpers.DeleteImageFromCloudinary(ctx, m.App.CLD, thumbPublicID)
		// delete video from cloudinary
		_ = helpers.DeleteVideoFromCloudinary(ctx, m.App.CLD, videoPublicID)

		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create the video",
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"message":  "Successfully created the video",
		"video_id": videoID,
	})
}
