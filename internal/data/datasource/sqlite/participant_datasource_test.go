package sqlite

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupParticipantDB(t *testing.T) (*sql.DB, ParticipantDataSource) {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	createParticipantTable := `
	CREATE TABLE participant (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		server_id TEXT NOT NULL,
		discord_id TEXT NOT NULL,
		points INTEGER DEFAULT 0,
		botm_enabled BOOLEAN DEFAULT true,
		kots_enabled BOOLEAN DEFAULT true,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		UNIQUE(server_id, discord_id)
	);`

	createAccountTable := `
	CREATE TABLE account (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		participant_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (participant_id) REFERENCES participant(id),
		UNIQUE(participant_id, name)
	);`

	if _, err := db.Exec(createParticipantTable); err != nil {
		t.Fatalf("failed to create participant table: %v", err)
	}
	if _, err := db.Exec(createAccountTable); err != nil {
		t.Fatalf("failed to create account table: %v", err)
	}

	datasource := NewParticipantDataSource(db)
	return db, datasource
}

func TestAddAccount_NewParticipant(t *testing.T) {
	_, ds := setupParticipantDB(t)

	serverID := "server1"
	discordID := "discord1"
	accountName := "testaccount"

	err := ds.AddAccount(serverID, discordID, accountName)
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	var participantID int64
	var points int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT id, points FROM participant WHERE server_id = ? AND discord_id = ?",
		serverID, discordID).Scan(&participantID, &points)
	if err != nil {
		t.Fatalf("failed to query participant: %v", err)
	}

	if participantID == 0 {
		t.Error("participant ID should not be 0")
	}
	if points != 0 {
		t.Errorf("expected points to be 0, got %d", points)
	}

	var accountID int64
	err = ds.(*participantDS).db.QueryRow(
		"SELECT id FROM account WHERE participant_id = ? AND name = ?",
		participantID, accountName).Scan(&accountID)
	if err != nil {
		t.Fatalf("failed to query account: %v", err)
	}

	if accountID == 0 {
		t.Error("account ID should not be 0")
	}
}

func TestAddAccount_ExistingParticipant(t *testing.T) {
	_, ds := setupParticipantDB(t)

	serverID := "server1"
	discordID := "discord1"
	accountName1 := "testaccount1"
	accountName2 := "testaccount2"

	err := ds.AddAccount(serverID, discordID, accountName1)
	if err != nil {
		t.Fatalf("AddAccount failed for first account: %v", err)
	}

	err = ds.AddAccount(serverID, discordID, accountName2)
	if err != nil {
		t.Fatalf("AddAccount failed for second account: %v", err)
	}

	var participantID int64
	err = ds.(*participantDS).db.QueryRow(
		"SELECT id FROM participant WHERE server_id = ? AND discord_id = ?",
		serverID, discordID).Scan(&participantID)
	if err != nil {
		t.Fatalf("failed to query participant: %v", err)
	}

	var count int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT COUNT(*) FROM account WHERE participant_id = ?",
		participantID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count accounts: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 accounts, got %d", count)
	}
}

func TestAddAccount_DuplicateAccount(t *testing.T) {
	_, ds := setupParticipantDB(t)

	serverID := "server1"
	discordID := "discord1"
	accountName := "testaccount"

	err := ds.AddAccount(serverID, discordID, accountName)
	if err != nil {
		t.Fatalf("AddAccount failed for first attempt: %v", err)
	}

	err = ds.AddAccount(serverID, discordID, accountName)
	if err != nil {
		t.Fatalf("AddAccount failed for duplicate attempt: %v", err)
	}

	var participantID int64
	err = ds.(*participantDS).db.QueryRow(
		"SELECT id FROM participant WHERE server_id = ? AND discord_id = ?",
		serverID, discordID).Scan(&participantID)
	if err != nil {
		t.Fatalf("failed to query participant: %v", err)
	}

	var count int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT COUNT(*) FROM account WHERE participant_id = ? AND name = ?",
		participantID, accountName).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count accounts: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 account, got %d", count)
	}
}
