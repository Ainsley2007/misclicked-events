package sqlite

import (
	"database/sql"
	"fmt"
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
	GetAllParticipantsWithAccounts(serverID string) ([]ParticipantWithAccounts, error)
	AddBotmParticipation(participantID string, botmID int64, startingKC int) error
	GetBotmParticipation(participantID string, botmID int64) (*BotmParticipationModel, error)
	CreateBotmParticipation(participantID string, botmID int64, startingKC int) error
	UpdateBotmParticipation(participantID string, botmID int64, startAmount, currentAmount int) error
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

func (ds *participantDS) AddBotmParticipation(participantID string, botmID int64, startingKC int) error {
	// This method is now deprecated - the repository handles the logic
	// Keeping for backward compatibility but it just delegates to CreateBotmParticipation
	return ds.CreateBotmParticipation(participantID, botmID, startingKC)
}

func (ds *participantDS) GetBotmParticipation(participantID string, botmID int64) (*BotmParticipationModel, error) {
	var participation BotmParticipationModel
	err := ds.db.QueryRow(`
		SELECT participant_id, botm_id, start_amount, current_amount
		FROM botm_participation 
		WHERE participant_id = ? AND botm_id = ?`, participantID, botmID).Scan(
		&participation.ParticipantID, &participation.BotmID,
		&participation.StartAmount, &participation.CurrentAmount)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		utils.Error("Failed to get BOTM participation for participant %s in competition %d: %v", participantID, botmID, err)
		return nil, fmt.Errorf("failed to get participation")
	}

	return &participation, nil
}

func (ds *participantDS) CreateBotmParticipation(participantID string, botmID int64, startingKC int) error {
	_, err := ds.db.Exec(`
		INSERT INTO botm_participation (participant_id, botm_id, start_amount, current_amount)
		VALUES (?, ?, ?, ?)`,
		participantID, botmID, startingKC, startingKC)
	if err != nil {
		utils.Error("Failed to create BOTM participation for participant %s in competition %d: %v", participantID, botmID, err)
		return fmt.Errorf("failed to create participation")
	}
	return nil
}

func (ds *participantDS) UpdateBotmParticipation(participantID string, botmID int64, startAmount, currentAmount int) error {
	_, err := ds.db.Exec(`
		UPDATE botm_participation 
		SET start_amount = ?, current_amount = ?
		WHERE participant_id = ? AND botm_id = ?`,
		startAmount, currentAmount, participantID, botmID)
	if err != nil {
		utils.Error("Failed to update BOTM participation for participant %s in competition %d: %v", participantID, botmID, err)
		return fmt.Errorf("failed to update participation")
	}
	return nil
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
