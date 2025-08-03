package repository

import (
	"misclicked-events/internal/data/datasource/sqlite"
)

type ConfigRepository struct {
	ds sqlite.ConfigDataSource
}

func NewConfigRepository(ds sqlite.ConfigDataSource) *ConfigRepository {
	return &ConfigRepository{ds}
}

func (r *ConfigRepository) FetchConfig(serverID string) (*sqlite.Config, error) {
	cfg, err := r.ds.GetConfig(serverID)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (r *ConfigRepository) SaveConfig(cfg *sqlite.Config) error {
	return r.ds.UpsertConfig(cfg)
}

func (r *ConfigRepository) EditRankingChannelID(serverID, channelID string) error {
	return r.ds.UpdateRankingChannelID(serverID, channelID)
}

func (r *ConfigRepository) EditHiscoreChannelID(serverID, channelID string) error {
	return r.ds.UpdateHiscoreChannelID(serverID, channelID)
}

func (r *ConfigRepository) EditCategoryChannelID(serverID, channelID string) error {
	return r.ds.UpdateCategoryChannelID(serverID, channelID)
}

func (r *ConfigRepository) EditRankingMessageID(serverID, messageID string) error {
	return r.ds.UpdateRankingMessageID(serverID, messageID)
}

func (r *ConfigRepository) EditHiscoreMessageID(serverID, messageID string) error {
	return r.ds.UpdateHiscoreMessageID(serverID, messageID)
}
