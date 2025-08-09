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
		botm_points INTEGER DEFAULT 0,
		kots_points INTEGER DEFAULT 0,
		UNIQUE(server_id, discord_id)
	);`

	createAccountTable := `
	CREATE TABLE account (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		participant_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		failed_fetch_count INTEGER DEFAULT 0,
		FOREIGN KEY (participant_id) REFERENCES participant(id),
		UNIQUE(participant_id, username)
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
	var botmPoints, kotsPoints int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT id, botm_points, kots_points FROM participant WHERE server_id = ? AND discord_id = ?",
		serverID, discordID).Scan(&participantID, &botmPoints, &kotsPoints)
	if err != nil {
		t.Fatalf("failed to query participant: %v", err)
	}

	if participantID == 0 {
		t.Error("participant ID should not be 0")
	}
	if botmPoints != 0 {
		t.Errorf("expected botm_points to be 0, got %d", botmPoints)
	}
	if kotsPoints != 0 {
		t.Errorf("expected kots_points to be 0, got %d", kotsPoints)
	}

	var accountID int64
	err = ds.(*participantDS).db.QueryRow(
		"SELECT id FROM account WHERE participant_id = ? AND username = ?",
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
		t.Fatalf("AddAccount failed: %v", err)
	}

	err = ds.AddAccount(serverID, discordID, accountName2)
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
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
		t.Fatalf("AddAccount failed: %v", err)
	}

	err = ds.AddAccount(serverID, discordID, accountName)
	if err == nil {
		t.Error("AddAccount should fail when account already tracked")
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
		"SELECT COUNT(*) FROM account WHERE participant_id = ? AND username = ?",
		participantID, accountName).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count accounts: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 account, got %d", count)
	}
}

func TestRemoveAccount_Success(t *testing.T) {
	_, ds := setupParticipantDB(t)

	serverID := "server1"
	discordID := "discord1"
	accountName := "testaccount"

	err := ds.AddAccount(serverID, discordID, accountName)
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	err = ds.RemoveAccount(serverID, discordID, accountName)
	if err != nil {
		t.Fatalf("RemoveAccount failed: %v", err)
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
		"SELECT COUNT(*) FROM account WHERE participant_id = ? AND username = ?",
		participantID, accountName).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count accounts: %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 accounts, got %d", count)
	}
}

func TestRemoveAccount_ParticipantNotFound(t *testing.T) {
	_, ds := setupParticipantDB(t)

	serverID := "server1"
	discordID := "discord1"
	accountName := "testaccount"

	err := ds.RemoveAccount(serverID, discordID, accountName)
	if err == nil {
		t.Error("RemoveAccount should fail when participant not found")
	}
}

func TestRemoveAccount_AccountNotFound(t *testing.T) {
	_, ds := setupParticipantDB(t)

	serverID := "server1"
	discordID := "discord1"
	accountName1 := "testaccount1"
	accountName2 := "testaccount2"

	err := ds.AddAccount(serverID, discordID, accountName1)
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	err = ds.RemoveAccount(serverID, discordID, accountName2)
	if err == nil {
		t.Error("RemoveAccount should fail when account not found")
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
		"SELECT COUNT(*) FROM account WHERE participant_id = ? AND username = ?",
		participantID, accountName1).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count accounts: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 account, got %d", count)
	}
}

func TestRemoveAccount_MultipleAccounts(t *testing.T) {
	_, ds := setupParticipantDB(t)

	serverID := "server1"
	discordID := "discord1"
	accountName1 := "testaccount1"
	accountName2 := "testaccount2"

	err := ds.AddAccount(serverID, discordID, accountName1)
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	err = ds.AddAccount(serverID, discordID, accountName2)
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	err = ds.RemoveAccount(serverID, discordID, accountName1)
	if err != nil {
		t.Fatalf("RemoveAccount failed: %v", err)
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

	if count != 1 {
		t.Errorf("expected 1 account, got %d", count)
	}

	var remainingAccount string
	err = ds.(*participantDS).db.QueryRow(
		"SELECT username FROM account WHERE participant_id = ?",
		participantID).Scan(&remainingAccount)
	if err != nil {
		t.Fatalf("failed to query remaining account: %v", err)
	}

	if remainingAccount != accountName2 {
		t.Errorf("expected remaining account to be %s, got %s", accountName2, remainingAccount)
	}
}
