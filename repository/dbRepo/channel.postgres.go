package dbrepo

import (
	"errors"
	"log"

	"github.com/raihan2bd/vidverse/models"
	"gorm.io/gorm/clause"
)

// get channel details
func (m *postgresDBRepo) GetChannelByID(id int) (*models.CustomChannelDTO, error) {
	var channel models.CustomChannelDTO

	err := m.DB.Table("channels").Select("channels.id, channels.title, channels.logo, channels.description, count(DISTINCT videos.id) as total_videos, count(DISTINCT subscriptions.id) as total_subscribers, channels.user_id, channels.logo_public_id, channels.cover, channels.cover_public_id").
		Joins("left join videos on videos.channel_id = channels.id").
		Joins("left join subscriptions on subscriptions.channel_id = channels.id").
		Where("channels.id = ?", id).
		Group("channels.id").
		Find(&channel).Error

	if err != nil {
		return nil, errors.New("internal server error. Please try again")
	}
	return &channel, nil
}

// Create new channel
func (m *postgresDBRepo) CreateChannel(channel *models.Channel) (uint, error) {
	result := m.DB.Create(&channel)
	if result.Error != nil {
		return 0, errors.New("failed to create the channel. please try again later")
	}
	return channel.ID, nil
}

// Get All the channels
func (m *postgresDBRepo) GetChannels(userID int) ([]models.CustomChannel, error) {
	var channels []models.CustomChannel
	err := m.DB.Table("channels").Select("channels.id, channels.title, channels.logo").Where("channels.user_id = ?", userID).
		Find(&channels).Error
	if err != nil {
		return nil, errors.New("internal server error. Please try again")
	}

	return channels, nil
}

// Get channel by user id with details
func (m *postgresDBRepo) GetChannelsWithDetailsByUserID(userID uint) ([]models.CustomChannelDTO, error) {
	var channels []models.CustomChannelDTO

	err := m.DB.Table("channels").Select("channels.id, channels.title, channels.logo, channels.description, channels.cover, count(DISTINCT videos.id) as total_video, count(DISTINCT subscriptions.id) as total_subscriber, channels.user_id").
		Joins("left join videos on videos.channel_id = channels.id").
		Joins("left join subscriptions on subscriptions.channel_id = channels.id").
		Where("channels.user_id = ?", userID).
		Group("channels.id").
		Order("channels.created_at asc").
		Find(&channels).Error

	if err != nil {
		return nil, errors.New("internal server error. Please try again")
	}

	return channels, nil
}

// delete channel By Id
func (m *postgresDBRepo) DeleteChannelByID(id int) *models.CustomError {
	// get channel by id
	var channel models.Channel
	result := m.DB.Table("channels").Where("id = ?", id).First(&channel)

	if result.Error != nil {
		return &models.CustomError{Status: 404, Err: errors.New("the channel you want to delete is not found")}
	}

	var videos []models.Video
	result = m.DB.Table("videos").Select("id, public_id, thumb_public_id").Where("channel_id = ?", id).Find(&videos)

	if result.Error != nil {
		videos = nil
		log.Println(result.Error)
	}

	tsx := m.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tsx.Rollback()
		}
	}()

	// delete all videos related to this channel
	for _, video := range videos {
		err := m.DeleteVideoModel(&video)
		if err != nil {
			log.Println(err)
		}
	}

	// delete the channel
	channelPublicID := channel.LogoPublicID
	channelCoverPublicID := channel.CoverPublicID

	err := m.DB.Select(clause.Associations).Unscoped().Delete(&channel).Error
	if err != nil {
		tsx.Rollback()
		return &models.CustomError{Status: 500, Err: errors.New("failed to delete the channel")}
	}

	err = tsx.Commit().Error
	if err != nil {
		return &models.CustomError{Status: 500, Err: errors.New("failed to delete the channel")}
	}

	go func() {
		// delete the logo from cloudinary
		if channelPublicID != "" {
			err = m.DeleteImageFromCloudinary(channelPublicID)
			if err != nil {
				log.Println(err)
			}
		}

		// delete cover from cloudinary
		if channel.CoverPublicID != "" {
			err = m.DeleteImageFromCloudinary(channelCoverPublicID)
			if err != nil {
				log.Println(err)
			}
		}
		// delete all notification related to this channel
		err = m.DeleteNotificationsByChannelID(uint(id))
		if err != nil {
			log.Println(err)
		}
		// delete all subscriptions related to this channel
		err = m.DeleteAllSubscriptionByChannelID(uint(id))
		if err != nil {
			log.Println(err)
		}
	}()

	return nil

}

// Update channel
func (m *postgresDBRepo) UpdateChannel(channel *models.CustomChannelDTO) error {

	result := m.DB.Table("channels").Where("id = ?", channel.ID).Updates(map[string]interface{}{
		"title":          channel.Title,
		"description":    channel.Description,
		"logo_public_id": channel.LogoPublicID,
		"logo":           channel.Logo,
	})

	if result.Error != nil {
		return errors.New("failed to update the channel")
	}
	return nil
}

// Get videos by channel id with pagination
func (m *postgresDBRepo) GetVideosByChannelIDWithPagination(channelID uint, page, limit int) ([]models.VideoDTO, int64, error) {
	var videos []models.VideoDTO
	var count int64

	err := m.DB.Table("videos").Select("videos.id, videos.title, videos.thumb, videos.views, videos.channel_id, channels.title as channel_title, channels.logo as channel_logo").
		Joins("left join channels on channels.id = videos.channel_id").
		Where("videos.channel_id = ?", channelID).
		Order("videos.created_at desc").
		Count(&count).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&videos).Error

	if err != nil {
		return nil, 0, errors.New("internal server error. Please try again")
	}

	return videos, count, nil
}

func (m *postgresDBRepo) GetChannelWithDetails(channelID, userID uint) (*models.CustomChannelDTO, error) {
	var channel models.CustomChannelDTO

	err := m.DB.Table("channels").Select("channels.id, channels.title, channels.logo, channels.description, channels.cover, count(DISTINCT videos.id) as total_video, count(DISTINCT subscriptions.id) as total_subscriber, channels.user_id").
		Joins("left join videos on videos.channel_id = channels.id").
		Joins("left join subscriptions on subscriptions.channel_id = channels.id").
		Where("channels.id = ?", channelID).
		Group("channels.id").
		First(&channel).Error

	if err != nil {
		return nil, errors.New("internal server error. Please try again")
	}

	// check if the user is subscribed to this channel
	channel.IsSubscribed = m.IsSubscribed(userID, channelID)

	return &channel, nil
}
