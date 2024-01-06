package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/config"
	"github.com/raihan2bd/vidverse/models"
	validator "github.com/raihan2bd/vidverse/validators"
)

func (m *Repo) HandleGetSubscribedChannels(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userIDUint := uint(userID.(float64))
	user, err := m.App.DBMethods.GetUserByID(userIDUint)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	channelID, err := strconv.Atoi(c.Param("channelID"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid channel id"})
		return
	}

	channels, err := m.App.DBMethods.GetChannelByID(channelID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	id, err := m.App.DBMethods.ToggleSubscription(userIDUint, uint(channelID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"channels": channels})

	if id == 0 {
		return
	}

	if user.ID == channels.UserID {
		return
	}

	// create notification
	notification := &models.Notification{
		ReceiverID: channels.UserID,
		IsRead:     false,
		SenderID:   userIDUint,
		SenderName: user.Name,
		ChannelID:  channels.ID,
		Type:       "subscribe",
	}

	err = m.App.DBMethods.CreateNotification(notification)
	if err != nil {
		return
	}

	// send notification to the user
	m.App.NotificationChan <- &config.NotificationEvent{
		BroadcasterID: channels.UserID,
		Action:        "a_new_notification",
		Data:          notification,
	}

}

func (m *Repo) HandleContactUs(c *gin.Context) {
	userID, ok := c.Get("user_id")

	// request data from client
	var contactUs models.ContactUs
	err := c.ShouldBindJSON(&contactUs)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	v := validator.New()
	v.IsLength(contactUs.Name, "name", 3, 100)
	v.IsEmail(contactUs.Email, "email", "invalid email")
	v.IsLength(contactUs.Message, "message", 10, 500)

	if !v.Valid() {
		c.JSON(400, gin.H{"error": v.GetErrMsg()})
		return
	}

	// get user from db
	if ok && userID != nil {
		var user *models.User
		user, err = m.App.DBMethods.GetUserByID(uint(userID.(float64)))
		if err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}

		if contactUs.IsForAuthor {
			if user.UserRole == "author" || user.UserRole == "admin" {
				c.JSON(400, gin.H{"error": "you are an author"})
				return
			}
		}

		contactUs.UserID = user.ID
	} else {
		isSubmitted := m.App.DBMethods.IsContactUsSubmitted(contactUs.Email)
		if isSubmitted {
			c.JSON(401, gin.H{"error": "Your message is already pending. Please sign up to submit more messages"})
			return
		}
	}

	// create contact us
	err = m.App.DBMethods.CreateContactUs(&contactUs)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

}
