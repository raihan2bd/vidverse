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

	channel, err := m.App.DBMethods.GetChannelByID(channelID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	id, err := m.App.DBMethods.ToggleSubscription(userIDUint, uint(channelID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"channels": channel})

	if id == 0 {
		return
	}

	if user.ID == channel.UserID {
		return
	}

	// create notification
	notification := &models.Notification{
		ReceiverID: channel.UserID,
		IsRead:     false,
		SenderID:   userIDUint,
		SenderName: user.Name,
		ChannelID:  channel.ID,
		Type:       "subscribe",
	}

	nID, err := m.App.DBMethods.CreateNotification(notification)
	if err != nil {
		return
	}

	notification.SenderAvatar = user.Avatar
	notification.Thumb = channel.Logo
	notification.ID = nID

	// send notification to the user
	m.App.NotificationChan <- &config.NotificationEvent{
		BroadcasterID: channel.UserID,
		Action:        "a_new_notification",
		Data:          notification,
	}

}

// HandleUpdateNotification update notification by notification ID
func (m *Repo) HandleUpdateNotification(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	notificationInt, err := strconv.Atoi(c.Param("notificationID"))
	if err != nil {
		c.JSON(404, gin.H{"error": "the notification you are trying to update is invalid"})
		return
	}

	notificationID := uint(notificationInt)

	userIDUint := uint(userID.(float64))

	// get notification by id
	notification, err := m.App.DBMethods.GetNotificationByID(notificationID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	if notification.IsRead {
		c.JSON(200, gin.H{"message": "notification is already read"})
		return
	}

	if notification.ReceiverID != userIDUint {
		c.JSON(403, gin.H{"error": "you are not authorized to update this notification"})
		return
	}

	// update notification
	err = m.App.DBMethods.UpdateNotificationByID(notification.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{"notification": notification})

	// send websocket signal to the user
	m.App.NotificationChan <- &config.NotificationEvent{
		BroadcasterID: userIDUint,
		Action:        "a_notification_is_read",
		Data:          notificationID,
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

// HandleGetNotificationsByUserID get all notifications by user ID
func (m *Repo) HandleGetNotifications(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var (
		err         error
		page, limit int
	)

	page, err = strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid page number"})
		return
	}

	limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid limit number"})
		return
	}

	userIDUint := uint(userID.(float64))

	var (
		notifications []models.Notification
		total         int64
	)

	notifications, total, err = m.App.DBMethods.GetNotificationsByUserID(userIDUint, page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	var has_next_page bool
	if total > int64(page*limit) {
		has_next_page = true
	}

	c.JSON(200, gin.H{"notifications": notifications, "total": total, "has_next_page": has_next_page, "page": page})
}

// HandleGetLikedVideos get all liked videos by user ID
func (m *Repo) HandleGetLikedVideos(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var (
		err         error
		page, limit int
	)

	page, err = strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid page number"})
		return
	}

	limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid limit number"})
		return
	}

	userIDUint := uint(userID.(float64))

	var (
		videos []models.VideoDTO
		total  int64
	)

	videos, total, err = m.App.DBMethods.GetLikedVideos(userIDUint, page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	var has_next_page bool
	if total > int64(page*limit) {
		has_next_page = true
	}

	c.JSON(200, gin.H{"videos": videos, "total": total, "has_next_page": has_next_page, "page": page})
}
