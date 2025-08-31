package repository

import (
	"fmt"
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/data/mappers"
	"misclicked-events/internal/domain"
	"misclicked-events/internal/utils"
)

type ConfigRepository struct {
	ds     sqlite.ConfigDataSource
	mapper *mappers.ConfigMapper
}

func NewConfigRepository(ds sqlite.ConfigDataSource) *ConfigRepository {
	return &ConfigRepository{
		ds:     ds,
		mapper: mappers.NewConfigMapper(),
	}
}

func (r *ConfigRepository) FetchConfig(serverID string) (*domain.Config, error) {
	if serverID == "" {
		return nil, fmt.Errorf("server ID cannot be empty")
	}

	cfgModel, err := r.ds.GetConfig(serverID)
	if err != nil {
		utils.Error("Failed to fetch config for server %s: %v", serverID, err)
		return nil, fmt.Errorf("failed to fetch config for server %s: %w", serverID, err)
	}

	cfg := r.mapper.ToDomain(cfgModel, serverID)
	return cfg, nil
}

func (r *ConfigRepository) SaveConfig(cfg *domain.Config, serverID string) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	cfgModel := r.mapper.ToModel(cfg, serverID)
	err := r.ds.UpsertConfig(cfgModel)
	if err != nil {
		utils.Error("Failed to save config for server %s: %v", serverID, err)
		return fmt.Errorf("failed to save config for server %s: %w", serverID, err)
	}

	utils.Info("Saved config for server %s", serverID)
	return nil
}

func (r *ConfigRepository) EditRankingChannelID(serverID, channelID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	if channelID == "" {
		return fmt.Errorf("channel ID cannot be empty")
	}

	err := r.ds.UpdateRankingChannelID(serverID, channelID)
	if err != nil {
		utils.Error("Failed to update ranking channel ID for server %s: %v", serverID, err)
		return fmt.Errorf("failed to update ranking channel for server %s: %w", serverID, err)
	}

	utils.Info("Updated ranking channel ID for server %s", serverID)
	return nil
}

func (r *ConfigRepository) EditHiscoreChannelID(serverID, channelID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	if channelID == "" {
		return fmt.Errorf("channel ID cannot be empty")
	}

	err := r.ds.UpdateHiscoreChannelID(serverID, channelID)
	if err != nil {
		utils.Error("Failed to update hiscore channel ID for server %s: %v", serverID, err)
		return fmt.Errorf("failed to update hiscore channel for server %s: %w", serverID, err)
	}

	utils.Info("Updated hiscore channel ID for server %s", serverID)
	return nil
}

func (r *ConfigRepository) EditCategoryChannelID(serverID, channelID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	// Category channel ID can be empty (optional field)
	err := r.ds.UpdateCategoryChannelID(serverID, channelID)
	if err != nil {
		utils.Error("Failed to update category channel ID for server %s: %v", serverID, err)
		return fmt.Errorf("failed to update category channel for server %s: %w", serverID, err)
	}

	utils.Info("Updated category channel ID for server %s", serverID)
	return nil
}

func (r *ConfigRepository) EditRankingMessageID(serverID, messageID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	if messageID == "" {
		return fmt.Errorf("message ID cannot be empty")
	}

	err := r.ds.UpdateRankingMessageID(serverID, messageID)
	if err != nil {
		utils.Error("Failed to update ranking message ID for server %s: %v", serverID, err)
		return fmt.Errorf("failed to update ranking message for server %s: %w", serverID, err)
	}

	return nil
}

func (r *ConfigRepository) EditHiscoreMessageID(serverID, messageID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	if messageID == "" {
		return fmt.Errorf("message ID cannot be empty")
	}

	err := r.ds.UpdateHiscoreMessageID(serverID, messageID)
	if err != nil {
		utils.Error("Failed to update hiscore message ID for server %s: %v", serverID, err)
		return fmt.Errorf("failed to update hiscore message for server %s: %w", serverID, err)
	}

	return nil
}
