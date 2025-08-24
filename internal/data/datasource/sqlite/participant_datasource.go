package sqlite

import (
	"database/sql"
	"fmt"
	"misclicked-events/internal/domain"
	"misclicked-events/internal/utils"
	"strings"
)

type ParticipantWithAccounts struct {
	DiscordID string
	Accounts  []string
}

type ParticipantDataSource interface {
	AddAccount(serverID, discordID, accountName string) error
	RemoveAccount(serverID, discordID, accountName string) error
	RenameAccount(serverID, discordID, oldUsername, newUsername string) error
	GetTrackedAccounts(serverID, discordID string) ([]string, error)
	GetTrackedAccountsWithIDs(serverID, discordID string) ([]AccountModel, error)
	GetAllParticipantsWithAccounts(serverID string) ([]ParticipantWithAccounts, error)
	GetParticipantWithAccountKC(participantID string, botmID int64) (*domain.ParticipantWithAccountKC, error)
	GetAccountID(participantID, username string) (int64, error) // Helper method for internal use
	GetBotmParticipation(accountID int64, botmID int64) (*BotmParticipationModel, error)
	CreateBotmParticipation(accountID int64, botmID int64, startingKC int) error
	UpdateBotmParticipation(accountID int64, botmID int64, startAmount, currentAmount int) error
	GetBotmParticipationByParticipant(participantID string, botmID int64) ([]*BotmParticipationModel, error)
}

func NewParticipantDataSource(db *sql.DB) ParticipantDataSource {
	return &participantDS{db: db}
}

type participantDS struct{ db *sql.DB }

func (ds *participantDS) AddAccount(serverID, discordID, accountName string) error {
	utils.Debug("Adding account %s for participant %s in server %s", accountName, discordID, serverID)

	result, err := ds.db.Exec(`
		INSERT OR IGNORE INTO participant (discord_id, server_id, botm_points, kots_points)
		VALUES (?, ?, 0, 0)`,
		discordID, serverID)
	if err != nil {
		utils.Error("Failed to create participant for discord ID %s in server %s: %v", discordID, serverID, err)
		return fmt.Errorf("failed to create participant")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.Error("Failed to get rows affected for participant creation for discord ID %s in server %s: %v", discordID, serverID, err)
		return fmt.Errorf("failed to create participant")
	}

	if rowsAffected > 0 {
		utils.Info("Created new participant for discord ID %s in server %s", discordID, serverID)
	} else {
		utils.Debug("Participant already exists for discord ID %s in server %s", discordID, serverID)
	}

	result, err = ds.db.Exec(`
		INSERT INTO account (participant_id, username, failed_fetch_count)
		VALUES (?, ?, 0)`,
		discordID, accountName)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.Debug("Account %s already exists for participant %s in server %s (case-insensitive)", accountName, discordID, serverID)
			return fmt.Errorf("account already tracked")
		}
		utils.Error("Failed to add account %s for participant %s in server %s: %v", accountName, discordID, serverID, err)
		return fmt.Errorf("failed to add account")
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		utils.Error("Failed to get rows affected for account creation for %s in server %s: %v", accountName, serverID, err)
		return fmt.Errorf("failed to add account")
	}

	if rowsAffected > 0 {
		utils.Info("Successfully added account %s for participant %s in server %s", accountName, discordID, serverID)
	} else {
		utils.Error("Failed to add account %s for participant %s in server %s: no rows affected", accountName, discordID, serverID)
		return fmt.Errorf("failed to add account")
	}

	return nil
}

func (ds *participantDS) RemoveAccount(serverID, discordID, accountName string) error {
	utils.Debug("Removing account %s for participant %s in server %s", accountName, discordID, serverID)

	result, err := ds.db.Exec(`
		DELETE FROM account 
		WHERE participant_id = ? AND LOWER(username) = LOWER(?)`, discordID, accountName)
	if err != nil {
		utils.Error("Failed to remove account %s for participant %s in server %s: %v", accountName, discordID, serverID, err)
		return fmt.Errorf("failed to remove account")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.Error("Failed to get rows affected for account removal %s for participant %s in server %s: %v", accountName, discordID, serverID, err)
		return fmt.Errorf("failed to remove account")
	}

	if rowsAffected == 0 {
		utils.Error("Account %s not found for participant %s in server %s", accountName, discordID, serverID)
		return fmt.Errorf("account not found")
	}

	utils.Info("Successfully removed account %s for participant %s in server %s", accountName, discordID, serverID)
	return nil
}

func (ds *participantDS) RenameAccount(serverID, discordID, oldUsername, newUsername string) error {
	utils.Debug("Renaming account from %s to %s for participant %s in server %s", oldUsername, newUsername, discordID, serverID)

	result, err := ds.db.Exec(`
		UPDATE account 
		SET username = ? 
		WHERE participant_id = ? AND LOWER(username) = LOWER(?)`, newUsername, discordID, oldUsername)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.Debug("Account %s already exists for participant %s in server %s (case-insensitive)", newUsername, discordID, serverID)
			return fmt.Errorf("account already tracked")
		}
		utils.Error("Failed to rename account from %s to %s for participant %s in server %s: %v", oldUsername, newUsername, discordID, serverID, err)
		return fmt.Errorf("failed to rename account")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.Error("Failed to get rows affected for account rename from %s to %s for participant %s in server %s: %v", oldUsername, newUsername, discordID, serverID, err)
		return fmt.Errorf("failed to rename account")
	}

	if rowsAffected == 0 {
		utils.Error("Account %s not found for participant %s in server %s", oldUsername, discordID, serverID)
		return fmt.Errorf("account not found")
	}

	utils.Info("Successfully renamed account from %s to %s for participant %s in server %s", oldUsername, newUsername, discordID, serverID)
	return nil
}

func (ds *participantDS) GetTrackedAccounts(serverID, discordID string) ([]string, error) {
	utils.Debug("Getting tracked accounts for participant %s in server %s", discordID, serverID)

	rows, err := ds.db.Query(`
		SELECT username 
		FROM account 
		WHERE participant_id = ? 
		ORDER BY LOWER(username)`, discordID)
	if err != nil {
		utils.Error("Failed to get tracked accounts for participant %s in server %s: %v", discordID, serverID, err)
		return nil, fmt.Errorf("failed to get tracked accounts")
	}
	defer rows.Close()

	var accounts []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			utils.Error("Failed to scan username for participant %s in server %s: %v", discordID, serverID, err)
			return nil, fmt.Errorf("failed to get tracked accounts")
		}
		accounts = append(accounts, username)
	}

	if err := rows.Err(); err != nil {
		utils.Error("Error iterating over tracked accounts for participant %s in server %s: %v", discordID, serverID, err)
		return nil, fmt.Errorf("failed to get tracked accounts")
	}

	utils.Debug("Found %d tracked accounts for participant %s in server %s", len(accounts), discordID, serverID)
	return accounts, nil
}

func (ds *participantDS) GetTrackedAccountsWithIDs(serverID, discordID string) ([]AccountModel, error) {
	utils.Debug("Getting tracked accounts with IDs for participant %s in server %s", discordID, serverID)

	rows, err := ds.db.Query(`
		SELECT id, username, failed_fetch_count 
		FROM account 
		WHERE participant_id = ? 
		ORDER BY LOWER(username)`, discordID)
	if err != nil {
		utils.Error("Failed to get tracked accounts with IDs for participant %s in server %s: %v", discordID, serverID, err)
		return nil, fmt.Errorf("failed to get tracked accounts with IDs")
	}
	defer rows.Close()

	var accounts []AccountModel
	for rows.Next() {
		var account AccountModel
		if err := rows.Scan(&account.ID, &account.Username, &account.FailedFetchCount); err != nil {
			utils.Error("Failed to scan account for participant %s in server %s: %v", discordID, serverID, err)
			continue
		}
		account.ParticipantID = discordID
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		utils.Error("Error iterating over tracked accounts with IDs for participant %s in server %s: %v", discordID, serverID, err)
		return nil, fmt.Errorf("failed to get tracked accounts with IDs")
	}

	utils.Debug("Found %d tracked accounts with IDs for participant %s in server %s", len(accounts), discordID, serverID)
	return accounts, nil
}

func (ds *participantDS) GetAccountID(participantID, username string) (int64, error) {
	utils.Debug("Getting account ID for participant %s with username %s", participantID, username)

	var accountID int64
	err := ds.db.QueryRow(`
		SELECT id 
		FROM account 
		WHERE participant_id = ? AND LOWER(username) = LOWER(?)`, participantID, username).Scan(&accountID)

	if err == sql.ErrNoRows {
		utils.Error("Account %s not found for participant %s", username, participantID)
		return 0, fmt.Errorf("account not found")
	} else if err != nil {
		utils.Error("Failed to get account ID for participant %s with username %s: %v", participantID, username, err)
		return 0, fmt.Errorf("failed to get account ID")
	}

	utils.Debug("Found account ID %d for participant %s with username %s", accountID, participantID, username)
	return accountID, nil
}

func (ds *participantDS) GetBotmParticipation(accountID int64, botmID int64) (*BotmParticipationModel, error) {
	var participation BotmParticipationModel
	err := ds.db.QueryRow(`
		SELECT account_id, botm_id, start_amount, current_amount
		FROM botm_participation 
		WHERE account_id = ? AND botm_id = ?`, accountID, botmID).Scan(
		&participation.AccountID, &participation.BotmID,
		&participation.StartAmount, &participation.CurrentAmount)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		utils.Error("Failed to get BOTM participation for account %d in competition %d: %v", accountID, botmID, err)
		return nil, fmt.Errorf("failed to get participation")
	}

	return &participation, nil
}

func (ds *participantDS) CreateBotmParticipation(accountID int64, botmID int64, startingKC int) error {
	_, err := ds.db.Exec(`
		INSERT INTO botm_participation (account_id, botm_id, start_amount, current_amount)
		VALUES (?, ?, ?, ?)`,
		accountID, botmID, startingKC, startingKC)
	if err != nil {
		utils.Error("Failed to create BOTM participation for account %d in competition %d: %v", accountID, botmID, err)
		return fmt.Errorf("failed to create participation")
	}
	return nil
}

func (ds *participantDS) UpdateBotmParticipation(accountID int64, botmID int64, startAmount, currentAmount int) error {
	_, err := ds.db.Exec(`
		UPDATE botm_participation 
		SET start_amount = ?, current_amount = ?
		WHERE account_id = ? AND botm_id = ?`,
		startAmount, currentAmount, accountID, botmID)
	if err != nil {
		utils.Error("Failed to update BOTM participation for account %d in competition %d: %v", accountID, botmID, err)
		return fmt.Errorf("failed to update participation")
	}
	return nil
}

func (ds *participantDS) GetParticipantWithAccountKC(participantID string, botmID int64) (*domain.ParticipantWithAccountKC, error) {
	utils.Debug("Getting participant %s with account KC for BOTM %d", participantID, botmID)

	// Single query to get all accounts with their KC for this participant
	rows, err := ds.db.Query(`
		SELECT a.id, a.username, 
		       COALESCE(bp.current_amount - bp.start_amount, 0) as kc_gained
		FROM account a
		LEFT JOIN botm_participation bp ON a.id = bp.account_id AND bp.botm_id = ?
		WHERE a.participant_id = ?
		ORDER BY LOWER(a.username)`, botmID, participantID)
	if err != nil {
		utils.Error("Failed to get participant with account KC for participant %s in BOTM %d: %v", participantID, botmID, err)
		return nil, fmt.Errorf("failed to get participant with account KC")
	}
	defer rows.Close()

	var accounts []domain.AccountWithKC
	for rows.Next() {
		var account domain.AccountWithKC
		if err := rows.Scan(&account.ID, &account.Username, &account.KCGained); err != nil {
			utils.Error("Failed to scan account KC data for participant %s in BOTM %d: %v", participantID, botmID, err)
			return nil, fmt.Errorf("failed to get participant with account KC")
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		utils.Error("Error iterating over account KC data for participant %s in BOTM %d: %v", participantID, botmID, err)
		return nil, fmt.Errorf("failed to get participant with account KC")
	}

	participant := &domain.ParticipantWithAccountKC{
		DiscordID: participantID,
		Accounts:  accounts,
	}

	utils.Debug("Found %d accounts with KC for participant %s in BOTM %d", len(accounts), participantID, botmID)
	return participant, nil
}

func (ds *participantDS) GetAllParticipantsWithAccounts(serverID string) ([]ParticipantWithAccounts, error) {
	utils.Debug("Getting all participants with accounts for server %s", serverID)

	rows, err := ds.db.Query(`
		SELECT p.discord_id, a.username
		FROM participant p
		JOIN account a ON p.discord_id = a.participant_id
		WHERE p.server_id = ?
		ORDER BY p.discord_id, LOWER(a.username)`, serverID)
	if err != nil {
		utils.Error("Failed to get participants with accounts for server %s: %v", serverID, err)
		return nil, fmt.Errorf("failed to get participants with accounts")
	}
	defer rows.Close()

	participantsMap := make(map[string][]string)
	for rows.Next() {
		var discordID, username string
		if err := rows.Scan(&discordID, &username); err != nil {
			utils.Error("Failed to scan participant data for server %s: %v", serverID, err)
			return nil, fmt.Errorf("failed to get participants with accounts")
		}
		participantsMap[discordID] = append(participantsMap[discordID], username)
	}

	if err := rows.Err(); err != nil {
		utils.Error("Error iterating over participants for server %s: %v", serverID, err)
		return nil, fmt.Errorf("failed to get participants with accounts")
	}

	var participants []ParticipantWithAccounts
	for discordID, accounts := range participantsMap {
		participants = append(participants, ParticipantWithAccounts{
			DiscordID: discordID,
			Accounts:  accounts,
		})
	}

	utils.Debug("Found %d participants with accounts for server %s", len(participants), serverID)
	return participants, nil
}

func (ds *participantDS) GetBotmParticipationByParticipant(participantID string, botmID int64) ([]*BotmParticipationModel, error) {
	utils.Debug("Getting BOTM participation for participant %s in competition %d", participantID, botmID)

	rows, err := ds.db.Query(`
		SELECT bp.account_id, bp.botm_id, bp.start_amount, bp.current_amount
		FROM botm_participation bp
		JOIN account a ON bp.account_id = a.id
		WHERE a.participant_id = ? AND bp.botm_id = ?`, participantID, botmID)
	if err != nil {
		utils.Error("Failed to get BOTM participation for participant %s in competition %d: %v", participantID, botmID, err)
		return nil, fmt.Errorf("failed to get participation")
	}
	defer rows.Close()

	var participations []*BotmParticipationModel
	for rows.Next() {
		var participation BotmParticipationModel
		if err := rows.Scan(&participation.AccountID, &participation.BotmID,
			&participation.StartAmount, &participation.CurrentAmount); err != nil {
			utils.Error("Failed to scan BOTM participation for participant %s in competition %d: %v", participantID, botmID, err)
			return nil, fmt.Errorf("failed to get participation")
		}
		participations = append(participations, &participation)
	}

	if err := rows.Err(); err != nil {
		utils.Error("Error iterating over BOTM participation for participant %s in competition %d: %v", participantID, botmID, err)
		return nil, fmt.Errorf("failed to get participation")
	}

	utils.Debug("Found %d BOTM participations for participant %s in competition %d", len(participations), participantID, botmID)
	return participations, nil
}
