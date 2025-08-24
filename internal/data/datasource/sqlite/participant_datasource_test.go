package sqlite

import (
	"database/sql"
	"sort"
	"strings"
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
		discord_id TEXT PRIMARY KEY,
		server_id TEXT NOT NULL,
		botm_points INTEGER DEFAULT 0,
		kots_points INTEGER DEFAULT 0
	);`

	createAccountTable := `
	CREATE TABLE account (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		participant_id TEXT NOT NULL,
		username TEXT NOT NULL,
		failed_fetch_count INTEGER DEFAULT 0,
		FOREIGN KEY (participant_id) REFERENCES participant(discord_id),
		UNIQUE(participant_id, username)
	);`

	createBotmParticipationTable := `
	CREATE TABLE botm_participation (
		account_id INTEGER NOT NULL,
		botm_id INTEGER NOT NULL,
		start_amount INTEGER NOT NULL,
		current_amount INTEGER NOT NULL,
		PRIMARY KEY (account_id, botm_id)
	);`

	createKotsParticipationTable := `
	CREATE TABLE kots_participation (
		account_id INTEGER NOT NULL,
		kots_id INTEGER NOT NULL,
		start_amount INTEGER NOT NULL,
		current_amount INTEGER NOT NULL,
		PRIMARY KEY (account_id, kots_id)
	);`

	if _, err := db.Exec(createParticipantTable); err != nil {
		t.Fatalf("failed to create participant table: %v", err)
	}
	if _, err := db.Exec(createAccountTable); err != nil {
		t.Fatalf("failed to create account table: %v", err)
	}
	if _, err := db.Exec(createBotmParticipationTable); err != nil {
		t.Fatalf("failed to create botm_participation table: %v", err)
	}
	if _, err := db.Exec(createKotsParticipationTable); err != nil {
		t.Fatalf("failed to create kots_participation table: %v", err)
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

	var botmPoints, kotsPoints int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT botm_points, kots_points FROM participant WHERE server_id = ? AND discord_id = ?",
		serverID, discordID).Scan(&botmPoints, &kotsPoints)
	if err != nil {
		t.Fatalf("failed to query participant: %v", err)
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
		discordID, accountName).Scan(&accountID)
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

	var count int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT COUNT(*) FROM account WHERE participant_id = ?",
		discordID).Scan(&count)
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

	var count int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT COUNT(*) FROM account WHERE participant_id = ? AND username = ?",
		discordID, accountName).Scan(&count)
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

	var count int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT COUNT(*) FROM account WHERE participant_id = ? AND username = ?",
		discordID, accountName).Scan(&count)
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

	var count int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT COUNT(*) FROM account WHERE participant_id = ? AND username = ?",
		discordID, accountName1).Scan(&count)
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

	var count int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT COUNT(*) FROM account WHERE participant_id = ?",
		discordID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count accounts: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 account, got %d", count)
	}

	var remainingAccount string
	err = ds.(*participantDS).db.QueryRow(
		"SELECT username FROM account WHERE participant_id = ?",
		discordID).Scan(&remainingAccount)
	if err != nil {
		t.Fatalf("failed to query remaining account: %v", err)
	}

	if remainingAccount != accountName2 {
		t.Errorf("expected remaining account to be %s, got %s", accountName2, remainingAccount)
	}
}

func TestCreateBotmParticipation_Success(t *testing.T) {
	_, ds := setupParticipantDB(t)

	// First create a participant and account
	serverID := "server1"
	discordID := "discord1"
	accountName := "testaccount"

	err := ds.AddAccount(serverID, discordID, accountName)
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	// Get the account ID
	var accountID int64
	err = ds.(*participantDS).db.QueryRow(
		"SELECT id FROM account WHERE participant_id = ? AND username = ?",
		discordID, accountName).Scan(&accountID)
	if err != nil {
		t.Fatalf("failed to query account: %v", err)
	}

	botmID := int64(1)
	startingKC := 100

	err = ds.CreateBotmParticipation(accountID, botmID, startingKC)
	if err != nil {
		t.Fatalf("CreateBotmParticipation failed: %v", err)
	}

	// Verify the participation was created
	var startAmount, currentAmount int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT start_amount, current_amount FROM botm_participation WHERE account_id = ? AND botm_id = ?",
		accountID, botmID).Scan(&startAmount, &currentAmount)
	if err != nil {
		t.Fatalf("failed to query botm_participation: %v", err)
	}

	if startAmount != startingKC {
		t.Errorf("expected start_amount to be %d, got %d", startingKC, startAmount)
	}
	if currentAmount != startingKC {
		t.Errorf("expected current_amount to be %d, got %d", startingKC, currentAmount)
	}
}

func TestGetBotmParticipation_Success(t *testing.T) {
	_, ds := setupParticipantDB(t)

	// First create a participant and account
	serverID := "server1"
	discordID := "discord1"
	accountName := "testaccount"

	err := ds.AddAccount(serverID, discordID, accountName)
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	// Get the account ID
	var accountID int64
	err = ds.(*participantDS).db.QueryRow(
		"SELECT id FROM account WHERE participant_id = ? AND username = ?",
		discordID, accountName).Scan(&accountID)
	if err != nil {
		t.Fatalf("failed to query account: %v", err)
	}

	botmID := int64(1)
	startingKC := 100

	// Create participation
	err = ds.CreateBotmParticipation(accountID, botmID, startingKC)
	if err != nil {
		t.Fatalf("CreateBotmParticipation failed: %v", err)
	}

	// Get participation
	participation, err := ds.GetBotmParticipation(accountID, botmID)
	if err != nil {
		t.Fatalf("GetBotmParticipation failed: %v", err)
	}

	if participation == nil {
		t.Fatal("expected participation to not be nil")
	}

	if participation.AccountID != accountID {
		t.Errorf("expected AccountID to be %d, got %d", accountID, participation.AccountID)
	}
	if participation.BotmID != botmID {
		t.Errorf("expected BotmID to be %d, got %d", botmID, participation.BotmID)
	}
	if participation.StartAmount != startingKC {
		t.Errorf("expected StartAmount to be %d, got %d", startingKC, participation.StartAmount)
	}
	if participation.CurrentAmount != startingKC {
		t.Errorf("expected CurrentAmount to be %d, got %d", startingKC, participation.CurrentAmount)
	}
}

func TestGetBotmParticipation_NotFound(t *testing.T) {
	_, ds := setupParticipantDB(t)

	accountID := int64(999)
	botmID := int64(1)

	participation, err := ds.GetBotmParticipation(accountID, botmID)
	if err != nil {
		t.Fatalf("GetBotmParticipation failed: %v", err)
	}

	if participation != nil {
		t.Error("expected participation to be nil when not found")
	}
}

func TestUpdateBotmParticipation_Success(t *testing.T) {
	_, ds := setupParticipantDB(t)

	// First create a participant and account
	serverID := "server1"
	discordID := "discord1"
	accountName := "testaccount"

	err := ds.AddAccount(serverID, discordID, accountName)
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	// Get the account ID
	var accountID int64
	err = ds.(*participantDS).db.QueryRow(
		"SELECT id FROM account WHERE participant_id = ? AND username = ?",
		discordID, accountName).Scan(&accountID)
	if err != nil {
		t.Fatalf("failed to query account: %v", err)
	}

	botmID := int64(1)
	startingKC := 100

	// Create participation
	err = ds.CreateBotmParticipation(accountID, botmID, startingKC)
	if err != nil {
		t.Fatalf("CreateBotmParticipation failed: %v", err)
	}

	// Update participation
	newStartAmount := 150
	newCurrentAmount := 200

	err = ds.UpdateBotmParticipation(accountID, botmID, newStartAmount, newCurrentAmount)
	if err != nil {
		t.Fatalf("UpdateBotmParticipation failed: %v", err)
	}

	// Verify the participation was updated
	var startAmount, currentAmount int
	err = ds.(*participantDS).db.QueryRow(
		"SELECT start_amount, current_amount FROM botm_participation WHERE account_id = ? AND botm_id = ?",
		accountID, botmID).Scan(&startAmount, &currentAmount)
	if err != nil {
		t.Fatalf("failed to query botm_participation: %v", err)
	}

	if startAmount != newStartAmount {
		t.Errorf("expected start_amount to be %d, got %d", newStartAmount, startAmount)
	}
	if currentAmount != newCurrentAmount {
		t.Errorf("expected current_amount to be %d, got %d", newCurrentAmount, currentAmount)
	}
}

func TestGetBotmParticipationByParticipant_Success(t *testing.T) {
	_, ds := setupParticipantDB(t)

	// First create a participant and multiple accounts
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

	// Get the account IDs
	var accountID1, accountID2 int64
	err = ds.(*participantDS).db.QueryRow(
		"SELECT id FROM account WHERE participant_id = ? AND username = ?",
		discordID, accountName1).Scan(&accountID1)
	if err != nil {
		t.Fatalf("failed to query account1: %v", err)
	}

	err = ds.(*participantDS).db.QueryRow(
		"SELECT id FROM account WHERE participant_id = ? AND username = ?", discordID, accountName2).Scan(&accountID2)
	if err != nil {
		t.Fatalf("failed to query account2: %v", err)
	}

	botmID := int64(1)

	// Create participations for both accounts
	err = ds.CreateBotmParticipation(accountID1, botmID, 100)
	if err != nil {
		t.Fatalf("CreateBotmParticipation for account1 failed: %v", err)
	}

	err = ds.CreateBotmParticipation(accountID2, botmID, 200)
	if err != nil {
		t.Fatalf("CreateBotmParticipation for account2 failed: %v", err)
	}

	// Get participations by participant
	participations, err := ds.GetBotmParticipationByParticipant(discordID, botmID)
	if err != nil {
		t.Fatalf("GetBotmParticipationByParticipant failed: %v", err)
	}

	if len(participations) != 2 {
		t.Errorf("expected 2 participations, got %d", len(participations))
	}

	// Verify both participations exist
	foundAccount1 := false
	foundAccount2 := false
	for _, p := range participations {
		if p.AccountID == accountID1 {
			foundAccount1 = true
			if p.StartAmount != 100 || p.CurrentAmount != 100 {
				t.Errorf("expected account1 to have start=100, current=100, got start=%d, current=%d", p.StartAmount, p.CurrentAmount)
			}
		}
		if p.AccountID == accountID2 {
			foundAccount2 = true
			if p.StartAmount != 200 || p.CurrentAmount != 200 {
				t.Errorf("expected account2 to have start=200, current=200, got start=%d, current=%d", p.StartAmount, p.CurrentAmount)
			}
		}
	}

	if !foundAccount1 {
		t.Error("expected to find participation for account1")
	}
	if !foundAccount2 {
		t.Error("expected to find participation for account2")
	}
}

func TestGetAccountID_Success(t *testing.T) {
	_, ds := setupParticipantDB(t)

	// First create a participant and account
	serverID := "server1"
	discordID := "discord1"
	accountName := "testaccount"

	err := ds.AddAccount(serverID, discordID, accountName)
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	// Get the account ID
	accountID, err := ds.GetAccountID(discordID, accountName)
	if err != nil {
		t.Fatalf("GetAccountID failed: %v", err)
	}

	if accountID == 0 {
		t.Error("expected account ID to not be 0")
	}

	// Verify the account ID matches what we expect
	var expectedAccountID int64
	err = ds.(*participantDS).db.QueryRow(
		"SELECT id FROM account WHERE participant_id = ? AND username = ?",
		discordID, accountName).Scan(&expectedAccountID)
	if err != nil {
		t.Fatalf("failed to query account: %v", err)
	}

	if accountID != expectedAccountID {
		t.Errorf("expected account ID to be %d, got %d", expectedAccountID, accountID)
	}
}

func TestGetAccountID_NotFound(t *testing.T) {
	_, ds := setupParticipantDB(t)

	discordID := "discord1"
	accountName := "nonexistent"

	_, err := ds.GetAccountID(discordID, accountName)
	if err == nil {
		t.Error("expected GetAccountID to fail when account not found")
	}

	if !strings.Contains(err.Error(), "account not found") {
		t.Errorf("expected error to contain 'account not found', got: %v", err)
	}
}

func TestGetTrackedAccountsWithIDs_Success(t *testing.T) {
	_, ds := setupParticipantDB(t)

	// Add a participant with multiple accounts
	err := ds.AddAccount("server1", "discord1", "account1")
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	err = ds.AddAccount("server1", "discord1", "account2")
	if err != nil {
		t.Fatalf("AddAccount failed: %v", err)
	}

	// Get accounts with IDs
	accounts, err := ds.GetTrackedAccountsWithIDs("server1", "discord1")
	if err != nil {
		t.Fatalf("GetTrackedAccountsWithIDs failed: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts))
	}

	// Verify account details
	if accounts[0].ParticipantID != "discord1" {
		t.Errorf("expected first account ParticipantID to be 'discord1', got %s", accounts[0].ParticipantID)
	}
	if accounts[1].ParticipantID != "discord1" {
		t.Errorf("expected second account ParticipantID to be 'discord1', got %s", accounts[1].ParticipantID)
	}
	if accounts[0].FailedFetchCount != 0 {
		t.Errorf("expected first account FailedFetchCount to be 0, got %d", accounts[0].FailedFetchCount)
	}
	if accounts[1].FailedFetchCount != 0 {
		t.Errorf("expected second account FailedFetchCount to be 0, got %d", accounts[1].FailedFetchCount)
	}

	// Verify usernames (order should be alphabetical due to ORDER BY LOWER(username))
	usernames := []string{accounts[0].Username, accounts[1].Username}
	sort.Strings(usernames)
	if usernames[0] != "account1" {
		t.Errorf("expected first username to be 'account1', got %s", usernames[0])
	}
	if usernames[1] != "account2" {
		t.Errorf("expected second username to be 'account2', got %s", usernames[1])
	}

	// Verify IDs are unique and positive
	if accounts[0].ID <= 0 {
		t.Errorf("expected first account ID to be positive, got %d", accounts[0].ID)
	}
	if accounts[1].ID <= 0 {
		t.Errorf("expected second account ID to be positive, got %d", accounts[1].ID)
	}
	if accounts[0].ID == accounts[1].ID {
		t.Errorf("expected account IDs to be different, both got %d", accounts[0].ID)
	}
}

func TestGetTrackedAccountsWithIDs_NoAccounts(t *testing.T) {
	_, ds := setupParticipantDB(t)

	// Try to get accounts for a participant that doesn't exist
	accounts, err := ds.GetTrackedAccountsWithIDs("server1", "nonexistent")
	if err != nil {
		t.Fatalf("GetTrackedAccountsWithIDs failed: %v", err)
	}
	if len(accounts) != 0 {
		t.Fatalf("expected 0 accounts, got %d", len(accounts))
	}
}

func TestGetTrackedAccountsWithIDs_EmptyServerID(t *testing.T) {
	_, ds := setupParticipantDB(t)

	// The datasource layer doesn't validate empty parameters, it just executes the SQL
	// Empty server ID will result in no results, not an error
	accounts, err := ds.GetTrackedAccountsWithIDs("", "discord1")
	if err != nil {
		t.Fatalf("GetTrackedAccountsWithIDs failed: %v", err)
	}
	if len(accounts) != 0 {
		t.Fatalf("expected 0 accounts for empty server ID, got %d", len(accounts))
	}
}

func TestGetTrackedAccountsWithIDs_EmptyDiscordID(t *testing.T) {
	_, ds := setupParticipantDB(t)

	// The datasource layer doesn't validate empty parameters, it just executes the SQL
	// Empty discord ID will result in no results, not an error
	accounts, err := ds.GetTrackedAccountsWithIDs("server1", "")
	if err != nil {
		t.Fatalf("GetTrackedAccountsWithIDs failed: %v", err)
	}
	if len(accounts) != 0 {
		t.Fatalf("expected 0 accounts for empty discord ID, got %d", len(accounts))
	}
}
