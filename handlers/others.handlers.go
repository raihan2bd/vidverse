package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (m *Repo) HandleGetSubscribedChannels(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userIDUint := uint(userID.(float64))
	channelID, err := strconv.Atoi(c.Param("channelID"))

	channels, err := m.App.DBMethods.GetChannelByID(channelID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	err = m.App.DBMethods.ToggleSubscription(userIDUint, uint(channelID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"channels": channels})
}
