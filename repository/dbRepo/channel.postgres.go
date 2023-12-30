package dbrepo

import (
	"errors"

	"github.com/raihan2bd/vidverse/models"
)

// get channel details
func (m *postgresDBRepo) GetChannelByID(id int) (*models.CustomChannelDTO, error) {
	var channel models.CustomChannelDTO

	err := m.DB.Table("channels").Select("channels.id, channels.title, channels.logo, channels.description, count(videos.id) as total_videos, count(subscriptions.id) as total_subscribers, channels.user_id, channels.logo_public_id, channels.cover, channels.cover_public_id").
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
