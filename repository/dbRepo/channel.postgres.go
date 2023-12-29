package dbrepo

import (
	"errors"

	"github.com/raihan2bd/vidverse/models"
)

// Create new channel
func (m *postgresDBRepo) CreateChannel(channel *models.Channel) (uint, error) {
	result := m.DB.Create(&channel)
	if result.Error != nil {
		return 0, errors.New("failed to create the channel. please try again later")
	}
	return channel.ID, nil
}

// Update channel
func (m *postgresDBRepo) UpdateChannel(channel *models.Channel) error {
	result := m.DB.Save(&channel)
	if result.Error != nil {
		return errors.New("failed to update the channel. please try again later")
	}
	return nil
}
