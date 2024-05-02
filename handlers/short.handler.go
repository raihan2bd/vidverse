package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	shortFile, fileInfo, err := c.Request.FormFile("short")
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "File is required."})
		return
	} else if fileInfo == nil {
		defer shortFile.Close()
		c.IndentedJSON(400, gin.H{"error": "short is required."})
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
	validator.IsVideoSize(fileInfo.Size, 100*1024*1024, "video")

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
}
