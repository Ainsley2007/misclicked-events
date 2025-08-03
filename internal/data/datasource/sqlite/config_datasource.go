package sqlite

import (
	"database/sql"
)

type ConfigDataSource interface {
	GetConfig(serverID string) (*Config, error)
	UpsertConfig(cfg *Config) error
	UpdateRankingChannelID(serverID, channelID string) error
	UpdateHiscoreChannelID(serverID, channelID string) error
	UpdateCategoryChannelID(serverID, channelID string) error
	UpdateRankingMessageID(serverID, messageID string) error
	UpdateHiscoreMessageID(serverID, messageID string) error
}

func NewConfigDataSource(db *sql.DB) ConfigDataSource {
	return &configDS{db}
}

type configDS struct{ db *sql.DB }

func (ds *configDS) GetConfig(serverID string) (*Config, error) {
	query := `
		SELECT ranking_channel_id, hiscore_channel_id, category_channel_id,
		       ranking_message_id, hiscore_message_id
		FROM config WHERE server_id = ?`
	row := ds.db.QueryRow(query, serverID)
	cfg := &Config{ServerID: serverID}
	err := row.Scan(
		&cfg.RankingChannelID,
		&cfg.HiscoreChannelID,
		&cfg.CategoryChannelID,
		&cfg.RankingMessageID,
		&cfg.HiscoreMessageID,
	)
	if err == sql.ErrNoRows {
		return cfg, nil
	}
	return cfg, err
}

func (ds *configDS) UpsertConfig(cfg *Config) error {
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
	return err
}

func (ds *configDS) UpdateRankingChannelID(serverID, channelID string) error {
	_, err := ds.db.Exec(
		`UPDATE config SET ranking_channel_id = ? WHERE server_id = ?`,
		channelID, serverID,
	)
	return err
}

func (ds *configDS) UpdateHiscoreChannelID(serverID, channelID string) error {
	_, err := ds.db.Exec(
		`UPDATE config SET hiscore_channel_id = ? WHERE server_id = ?`,
		channelID, serverID,
	)
	return err
}

func (ds *configDS) UpdateCategoryChannelID(serverID, channelID string) error {
	_, err := ds.db.Exec(
		`UPDATE config SET category_channel_id = ? WHERE server_id = ?`,
		channelID, serverID,
	)
	return err
}

func (ds *configDS) UpdateRankingMessageID(serverID, messageID string) error {
	_, err := ds.db.Exec(
		`UPDATE config SET ranking_message_id = ? WHERE server_id = ?`,
		messageID, serverID,
	)
	return err
}

func (ds *configDS) UpdateHiscoreMessageID(serverID, messageID string) error {
	_, err := ds.db.Exec(
		`UPDATE config SET hiscore_message_id = ? WHERE server_id = ?`,
		messageID, serverID,
	)
	return err
}
