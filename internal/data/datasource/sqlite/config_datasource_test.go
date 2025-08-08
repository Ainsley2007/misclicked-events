package sqlite

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, ConfigDataSource) {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	createStmt := `
CREATE TABLE config (
    server_id TEXT PRIMARY KEY,
    ranking_channel_id TEXT,
    hiscore_channel_id TEXT,
    category_channel_id TEXT,
    ranking_message_id TEXT,
    hiscore_message_id TEXT
);`
	if _, err := db.Exec(createStmt); err != nil {
		t.Fatalf("failed to create config table: %v", err)
	}
	datasource := NewConfigDataSource(db)
	return db, datasource
}

func TestGetConfig_Default(t *testing.T) {
	_, ds := setupTestDB(t)
	serverID := "server1"
	cfg, err := ds.GetConfig(serverID)
	if err != nil {
		t.Fatalf("GetConfig returned error: %v", err)
	}
	expected := &ConfigModel{ServerID: serverID}
	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("expected %+v, got %+v", expected, cfg)
	}
}

func TestUpsertAndGetConfig(t *testing.T) {
	_, ds := setupTestDB(t)
	serverID := "server1"
	input := &ConfigModel{
		ServerID:          serverID,
		RankingChannelID:  "rank1",
		HiscoreChannelID:  "hiscore1",
		CategoryChannelID: "cat1",
		RankingMessageID:  "msg1",
		HiscoreMessageID:  "msg2",
	}
	if err := ds.UpsertConfig(input); err != nil {
		t.Fatalf("UpsertConfig error: %v", err)
	}
	cfg, err := ds.GetConfig(serverID)
	if err != nil {
		t.Fatalf("GetConfig returned error after upsert: %v", err)
	}
	if !reflect.DeepEqual(cfg, input) {
		t.Errorf("expected %+v, got %+v", input, cfg)
	}
}

func TestUpdateFields(t *testing.T) {
	_, ds := setupTestDB(t)
	serverID := "server1"
	_ = ds.UpsertConfig(&ConfigModel{
		ServerID:          serverID,
		RankingChannelID:  "rank1",
		HiscoreChannelID:  "hiscore1",
		CategoryChannelID: "cat1",
		RankingMessageID:  "msg1",
		HiscoreMessageID:  "msg2",
	})
	// Update Ranking Channel ID
	if err := ds.UpdateRankingChannelID(serverID, "newRank"); err != nil {
		t.Fatalf("UpdateRankingChannelID error: %v", err)
	}
	cfg, _ := ds.GetConfig(serverID)
	if cfg.RankingChannelID != "newRank" {
		t.Errorf("expected ranking_channel_id 'newRank', got '%s'", cfg.RankingChannelID)
	}
	// Update Hiscore Channel ID
	if err := ds.UpdateHiscoreChannelID(serverID, "newHiscore"); err != nil {
		t.Fatalf("UpdateHiscoreChannelID error: %v", err)
	}
	cfg, _ = ds.GetConfig(serverID)
	if cfg.HiscoreChannelID != "newHiscore" {
		t.Errorf("expected hiscore_channel_id 'newHiscore', got '%s'", cfg.HiscoreChannelID)
	}
	// Update Category Channel ID
	if err := ds.UpdateCategoryChannelID(serverID, "newCat"); err != nil {
		t.Fatalf("UpdateCategoryChannelID error: %v", err)
	}
	cfg, _ = ds.GetConfig(serverID)
	if cfg.CategoryChannelID != "newCat" {
		t.Errorf("expected category_channel_id 'newCat', got '%s'", cfg.CategoryChannelID)
	}
	// Update Ranking Message ID
	if err := ds.UpdateRankingMessageID(serverID, "newMsg"); err != nil {
		t.Fatalf("UpdateRankingMessageID error: %v", err)
	}
	cfg, _ = ds.GetConfig(serverID)
	if cfg.RankingMessageID != "newMsg" {
		t.Errorf("expected ranking_message_id 'newMsg', got '%s'", cfg.RankingMessageID)
	}
	// Update Hiscore Message ID
	if err := ds.UpdateHiscoreMessageID(serverID, "newMsg2"); err != nil {
		t.Fatalf("UpdateHiscoreMessageID error: %v", err)
	}
	cfg, _ = ds.GetConfig(serverID)
	if cfg.HiscoreMessageID != "newMsg2" {
		t.Errorf("expected hiscore_message_id 'newMsg2', got '%s'", cfg.HiscoreMessageID)
	}
}
