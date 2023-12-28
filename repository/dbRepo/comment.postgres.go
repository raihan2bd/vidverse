package dbrepo

import (
	"errors"

	"github.com/raihan2bd/vidverse/models"
)

// Get all comments
func (m *postgresDBRepo) GetCommentsByVideoID(id, page, limit int) ([]models.CommentDTO, int64, error) {
	var comments []models.CommentDTO
	var count int64
	offset := (page - 1) * limit

	err := m.DB.Table("comments").Select("comments.id, comments.text, comments.video_id, users.id as user_id, users.name as user_name, users.avatar as user_avatar, comments.created_at").
		Joins("left join users on users.id = comments.user_id").
		Where("comments.video_id = ?", id).
		Count(&count).
		Offset(offset).Limit(limit).
		Order("comments.created_at desc").
		Find(&comments).Error
	if err != nil {
		return nil, 0, errors.New("internal server error. Please try again")
	}

	return comments, count, nil
}

// Get comment by ID
func (m *postgresDBRepo) GetCommentByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := m.DB.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, errors.New("404 comment not found")
	}

	return &comment, nil
}

// Create new comment
func (m *postgresDBRepo) CreateComment(comment *models.Comment) (uint, error) {
	result := m.DB.Create(&comment)
	if result.Error != nil {
		return 0, errors.New("failed to create comment")
	}

	return comment.ID, nil
}

// update comment
func (m *postgresDBRepo) UpdateComment(comment *models.Comment) error {
	result := m.DB.Save(&comment)
	if result.Error != nil {
		return errors.New("failed to update comment")
	}

	return nil
}

// delete comment
func (m *postgresDBRepo) DeleteCommentByID(id uint) error {
	result := m.DB.Unscoped().Delete(&models.Comment{}, id)
	if result.Error != nil {
		return errors.New("something went wrong. failed to delete the comment")
	}

	return nil
}

func (m *postgresDBRepo) DeleteNotificationByCommentID(commentID uint) error {
	// delete all notifications where comment_id = commentID
	result := m.DB.Unscoped().Delete(&models.Notification{}, "comment_id = ?", commentID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *postgresDBRepo) DeleteNotificationByVideoID(videoID uint) error {
	// delete all notifications where video_id = videoID
	result := m.DB.Unscoped().Delete(&models.Notification{}, "video_id = ?", videoID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *postgresDBRepo) DeleteNotificationByChannelID(channelID uint) error {
	// delete all notifications where channel_id = channelID
	result := m.DB.Unscoped().Delete(&models.Notification{}, "channel_id = ?", channelID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// delete notification by likeID
func (m *postgresDBRepo) DeleteNotificationByLikeID(likeID uint) error {
	// delete all notifications where like_id = likeID
	result := m.DB.Unscoped().Delete(&models.Notification{}, "like_id = ?", likeID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// delete notification by senderID
func (m *postgresDBRepo) DeleteNotificationBySenderID(senderID uint) error {
	// delete all notifications where sender_id = senderID
	result := m.DB.Unscoped().Delete(&models.Notification{}, "sender_id = ?", senderID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// delete notification by receiverID
func (m *postgresDBRepo) DeleteNotificationByReceiverID(receiverID uint) error {
	// delete all notifications where receiver_id = receiverID
	result := m.DB.Unscoped().Delete(&models.Notification{}, "receiver_id = ?", receiverID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
