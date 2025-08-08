package sqlite

import (
	"database/sql"
	"fmt"
	"misclicked-events/internal/utils"
)

type ConfigDataSource interface {
	GetConfig(serverID string) (*ConfigModel, error)
	UpsertConfig(cfg *ConfigModel) error
	UpdateRankingChannelID(serverID, channelID string) error
	UpdateHiscoreChannelID(serverID, channelID string) error
	UpdateCategoryChannelID(serverID, channelID string) error
	UpdateRankingMessageID(serverID, messageID string) error
	UpdateHiscoreMessageID(serverID, messageID string) error
}

func NewConfigDataSource(db *sql.DB) ConfigDataSource {
	return &configDS{db: db}
}

type configDS struct{ db *sql.DB }

func (ds *configDS) GetConfig(serverID string) (*ConfigModel, error) {
	query := `
		SELECT ranking_channel_id, hiscore_channel_id, category_channel_id,
		       ranking_message_id, hiscore_message_id
		FROM config WHERE server_id = ?`

	row := ds.db.QueryRow(query, serverID)
	cfg := &ConfigModel{ServerID: serverID}

	err := row.Scan(
		&cfg.RankingChannelID,
		&cfg.HiscoreChannelID,
		&cfg.CategoryChannelID,
		&cfg.RankingMessageID,
		&cfg.HiscoreMessageID,
	)

	if err == sql.ErrNoRows {
		utils.Debug("No config found for server %s, returning default config", serverID)
		return cfg, nil
	}

	if err != nil {
		utils.Error("Failed to scan config for server %s: %v", serverID, err)
		return nil, fmt.Errorf("failed to retrieve config for server %s: %w", serverID, err)
	}

	utils.Debug("Successfully retrieved config for server %s", serverID)
	return cfg, nil
}

func (ds *configDS) UpsertConfig(cfg *ConfigModel) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if cfg.ServerID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	sqlStmt := `
		INSERT INTO config(server_id, ranking_channel_id, hiscore_channel_id,
		                    category_channel_id, ranking_message_id, hiscore_message_id)
		VALUES(?, ?, ?, ?, ?, ?)
		ON CONFLICT(server_id) DO UPDATE SET
		  ranking_channel_id   = excluded.ranking_channel_id,
		  hiscore_channel_id   = excluded.hiscore_channel_id,
		  category_channel_id  = excluded.category_channel_id,
		  ranking_message_id   = excluded.ranking_message_id,
		  hiscore_message_id   = excluded.hiscore_message_id`

	_, err := ds.db.Exec(sqlStmt,
		cfg.ServerID,
		cfg.RankingChannelID,
		cfg.HiscoreChannelID,
		cfg.CategoryChannelID,
		cfg.RankingMessageID,
		cfg.HiscoreMessageID,
	)

	if err != nil {
		utils.Error("Failed to upsert config for server %s: %v", cfg.ServerID, err)
		return fmt.Errorf("failed to save config for server %s: %w", cfg.ServerID, err)
	}

	utils.Info("Successfully saved config for server %s", cfg.ServerID)
	return nil
}

func (ds *configDS) UpdateRankingChannelID(serverID, channelID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	_, err := ds.db.Exec(
		`UPDATE config SET ranking_channel_id = ? WHERE server_id = ?`,
		channelID, serverID,
	)

	if err != nil {
		utils.Error("Failed to update ranking channel ID for server %s: %v", serverID, err)
		return fmt.Errorf("failed to update ranking channel for server %s: %w", serverID, err)
	}

	utils.Debug("Updated ranking channel ID to %s for server %s", channelID, serverID)
	return nil
}

func (ds *configDS) UpdateHiscoreChannelID(serverID, channelID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	_, err := ds.db.Exec(
		`UPDATE config SET hiscore_channel_id = ? WHERE server_id = ?`,
		channelID, serverID,
	)

	if err != nil {
		utils.Error("Failed to update hiscore channel ID for server %s: %v", serverID, err)
		return fmt.Errorf("failed to update hiscore channel for server %s: %w", serverID, err)
	}

	utils.Debug("Updated hiscore channel ID to %s for server %s", channelID, serverID)
	return nil
}

func (ds *configDS) UpdateCategoryChannelID(serverID, channelID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	_, err := ds.db.Exec(
		`UPDATE config SET category_channel_id = ? WHERE server_id = ?`,
		channelID, serverID,
	)

	if err != nil {
		utils.Error("Failed to update category channel ID for server %s: %v", serverID, err)
		return fmt.Errorf("failed to update category channel for server %s: %w", serverID, err)
	}

	utils.Debug("Updated category channel ID to %s for server %s", channelID, serverID)
	return nil
}

func (ds *configDS) UpdateRankingMessageID(serverID, messageID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	_, err := ds.db.Exec(
		`UPDATE config SET ranking_message_id = ? WHERE server_id = ?`,
		messageID, serverID,
	)

	if err != nil {
		utils.Error("Failed to update ranking message ID for server %s: %v", serverID, err)
		return fmt.Errorf("failed to update ranking message for server %s: %w", serverID, err)
	}

	utils.Debug("Updated ranking message ID to %s for server %s", messageID, serverID)
	return nil
}

func (ds *configDS) UpdateHiscoreMessageID(serverID, messageID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	_, err := ds.db.Exec(
		`UPDATE config SET hiscore_message_id = ? WHERE server_id = ?`,
		messageID, serverID,
	)

	if err != nil {
		utils.Error("Failed to update hiscore message ID for server %s: %v", serverID, err)
		return fmt.Errorf("failed to update hiscore message for server %s: %w", serverID, err)
	}

	utils.Debug("Updated hiscore message ID to %s for server %s", messageID, serverID)
	return nil
}
