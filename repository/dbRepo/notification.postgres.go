package dbrepo

import (
	"errors"
	"fmt"
	"log"

	"github.com/raihan2bd/vidverse/models"
)

// create notification
func (m *postgresDBRepo) CreateNotification(notification *models.Notification) error {
	var n models.Notification
	err := m.DB.Where("sender_id = ? AND receiver_id = ? AND channel_id = ? AND video_id = ? AND comment_id = ? AND type = ?", notification.SenderID, notification.ReceiverID, notification.ChannelID, notification.VideoID, notification.CommentID, notification.Type).First(&n).Error
	if err != nil {
		// do nothing
		log.Println(err)
	}

	if n.ID == 0 {
		result := m.DB.Create(&notification)
		if result.Error != nil {
			return errors.New("failed to create notification")
		}
	}

	return nil
}

// get all notifications by user ID
func (m *postgresDBRepo) GetNotificationsByUserID(userID uint, page, limit int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64
	offset := (page - 1) * limit

	err := m.DB.Table("notifications").Select("notifications.id, notifications.is_read, notifications.receiver_id, notifications.sender_id, notifications.sender_name, notifications.sender_avatar, notifications.thumb, notifications.video_id, notifications.channel_id, notifications.comment_id, notifications.like_id, notifications.type, notifications.created_at, users.avatar as sender_avatar, videos.thumb as thumb, channels.logo as thumb").
		Joins("left join users on users.id = notifications.sender_id").
		Joins("left join videos on videos.id = notifications.video_id").
		Joins("left join channels on channels.id = notifications.channel_id").
		Where("notifications.receiver_id = ?", userID).
		Order("notifications.is_read asc, notifications.created_at desc").
		Count(&total).
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error

	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("internal server error. Please try again")
	}

	return notifications, total, nil
}

// Get all unread notifications by user ID
func (m *postgresDBRepo) GetUnreadNotificationsByUserID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := m.DB.Where("receiver_id = ? AND is_read = ?", userID, false).Order("created_at desc").Find(&notifications).Error
	if err != nil {
		return nil, errors.New("internal server error. Please try again")
	}

	return notifications, nil
}

// Get all unread notifications count by user ID
func (m *postgresDBRepo) GetUnreadNotificationsCountByUserID(userID uint) (int64, error) {
	var count int64
	err := m.DB.Model(&models.Notification{}).Where("receiver_id = ? AND is_read = ?", userID, false).Count(&count).Error
	if err != nil {
		return 0, errors.New("internal server error. Please try again")
	}

	return count, nil
}

// Get notification by ID
func (m *postgresDBRepo) GetNotificationByID(id uint) (*models.Notification, error) {
	var notification models.Notification
	err := m.DB.First(&notification, id).Error
	if err != nil {
		return nil, errors.New("404 notification not found")
	}

	return &notification, nil
}

// Delete notification by ID
func (m *postgresDBRepo) DeleteNotificationByID(id uint) error {
	result := m.DB.Unscoped().Delete(&models.Notification{}, id)
	if result.Error != nil {
		return errors.New("something went wrong. failed to delete the notification")
	}

	return nil
}

// Delete notification by ID
func (m *postgresDBRepo) DeleteNotificationsByChannelID(id uint) error {
	result := m.DB.Unscoped().Where("channel_id = ?", id).Delete(&models.Notification{})
	if result.Error != nil {
		return errors.New("something went wrong. failed to delete the notifications")
	}

	return nil
}

// delete all notifications by user ID and is_read = true and created_at < 30 days
func (m *postgresDBRepo) DeleteAllNotificationsByUserID(userID int) error {
	result := m.DB.Unscoped().Where("user_id = ? AND is_read = ? AND created_at < now() - interval '30 days'", userID, true).Delete(&models.Notification{})
	if result.Error != nil {
		return errors.New("something went wrong. failed to delete the notifications")
	}

	return nil
}
