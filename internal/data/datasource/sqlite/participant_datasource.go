package sqlite

import (
	"database/sql"
	"fmt"
	"misclicked-events/internal/utils"
	"strings"
)

type ParticipantDataSource interface {
	AddAccount(serverID, discordID, accountName string) error
	RemoveAccount(serverID, discordID, accountName string) error
	RenameAccount(serverID, discordID, oldUsername, newUsername string) error
	GetParticipantID(serverID, discordID string) (int64, error)
	GetTrackedAccounts(serverID, discordID string) ([]string, error)
	AddBotmParticipation(participantID, botmID int64, startingKC int) error
}

func NewParticipantDataSource(db *sql.DB) ParticipantDataSource {
	return &participantDS{db: db}
}

type participantDS struct{ db *sql.DB }

func (ds *participantDS) AddAccount(serverID, discordID, accountName string) error {
	utils.Debug("Adding account %s for participant %s in server %s", accountName, discordID, serverID)

	result, err := ds.db.Exec(`
		INSERT OR IGNORE INTO participant (server_id, discord_id, botm_points, kots_points)
		VALUES (?, ?, 0, 0)`,
		serverID, discordID)
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
		SELECT id, ?, 0 FROM participant WHERE server_id = ? AND discord_id = ?`,
		accountName, serverID, discordID)
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
		WHERE account.participant_id = (
			SELECT id FROM participant 
			WHERE server_id = ? AND discord_id = ?
		) AND LOWER(account.username) = LOWER(?)`, serverID, discordID, accountName)
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
	utils.Debug("Renaming account %s to %s for participant %s in server %s", oldUsername, newUsername, discordID, serverID)

	result, err := ds.db.Exec(`
		UPDATE account 
		SET username = ?
		WHERE account.participant_id = (
			SELECT id FROM participant 
			WHERE server_id = ? AND discord_id = ?
		) AND LOWER(account.username) = LOWER(?)`, newUsername, serverID, discordID, oldUsername)
	if err != nil {
		utils.Error("Failed to rename account %s to %s for participant %s in server %s: %v", oldUsername, newUsername, discordID, serverID, err)
		return fmt.Errorf("failed to rename account")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.Error("Failed to get rows affected for account rename %s to %s for participant %s in server %s: %v", oldUsername, newUsername, discordID, serverID, err)
		return fmt.Errorf("failed to rename account")
	}

	if rowsAffected == 0 {
		utils.Error("Account %s not found for participant %s in server %s", oldUsername, discordID, serverID)
		return fmt.Errorf("account not found")
	}

	utils.Info("Successfully renamed account %s to %s for participant %s in server %s", oldUsername, newUsername, discordID, serverID)
	return nil
}

func (ds *participantDS) GetParticipantID(serverID, discordID string) (int64, error) {
	utils.Debug("Getting participant ID for discord ID %s in server %s", discordID, serverID)

	var participantID int64
	err := ds.db.QueryRow(`
		SELECT id FROM participant WHERE server_id = ? AND discord_id = ?`,
		serverID, discordID).Scan(&participantID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error("Participant not found for discord ID %s in server %s", discordID, serverID)
			return 0, fmt.Errorf("participant not found")
		}
		utils.Error("Failed to get participant ID for discord ID %s in server %s: %v", discordID, serverID, err)
		return 0, fmt.Errorf("failed to get participant")
	}

	utils.Debug("Found participant ID %d for discord ID %s in server %s", participantID, discordID, serverID)
	return participantID, nil
}

func (ds *participantDS) GetTrackedAccounts(serverID, discordID string) ([]string, error) {
	utils.Debug("Getting tracked accounts for participant %s in server %s", discordID, serverID)

	rows, err := ds.db.Query(`
		SELECT username FROM account 
		WHERE participant_id = (
			SELECT id FROM participant 
			WHERE server_id = ? AND discord_id = ?
		) ORDER BY username`, serverID, discordID)
	if err != nil {
		utils.Error("Failed to get tracked accounts for participant %s in server %s: %v", discordID, serverID, err)
		return nil, fmt.Errorf("failed to get accounts")
	}
	defer rows.Close()

	var accounts []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			utils.Error("Failed to scan username for participant %s in server %s: %v", discordID, serverID, err)
			return nil, fmt.Errorf("failed to get accounts")
		}
		accounts = append(accounts, username)
	}

	if err := rows.Err(); err != nil {
		utils.Error("Error iterating over tracked accounts for participant %s in server %s: %v", discordID, serverID, err)
		return nil, fmt.Errorf("failed to get accounts")
	}

	utils.Debug("Found %d tracked accounts for participant %s in server %s: %v", len(accounts), discordID, serverID, accounts)
	return accounts, nil
}

func (ds *participantDS) AddBotmParticipation(participantID, botmID int64, startingKC int) error {
	utils.Debug("Adding BOTM participation for participant %d in BOTM %d with starting KC %d", participantID, botmID, startingKC)

	result, err := ds.db.Exec(`
		INSERT OR IGNORE INTO botm_participation (participant_id, botm_id, start_amount, current_amount)
		VALUES (?, ?, ?, ?)`,
		participantID, botmID, startingKC, startingKC)
	if err != nil {
		utils.Error("Failed to add BOTM participation for participant %d in BOTM %d: %v", participantID, botmID, err)
		return fmt.Errorf("failed to add participation")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.Error("Failed to get rows affected for BOTM participation creation for participant %d in BOTM %d: %v", participantID, botmID, err)
		return fmt.Errorf("failed to add participation")
	}

	if rowsAffected > 0 {
		utils.Info("Successfully added BOTM participation for participant %d in BOTM %d with starting KC %d", participantID, botmID, startingKC)
	} else {
		utils.Debug("BOTM participation already exists for participant %d in BOTM %d", participantID, botmID)
	}

	return nil
}
