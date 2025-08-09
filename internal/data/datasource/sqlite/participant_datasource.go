package sqlite

import (
	"database/sql"
	"fmt"
	"misclicked-events/internal/utils"
)

type ParticipantDataSource interface {
	AddAccount(serverID, discordID, accountName string) error
	RemoveAccount(serverID, discordID, accountName string) error
}

func NewParticipantDataSource(db *sql.DB) ParticipantDataSource {
	return &participantDS{db: db}
}

type participantDS struct{ db *sql.DB }

func (ds *participantDS) AddAccount(serverID, discordID, accountName string) error {
	result, err := ds.db.Exec(`
		INSERT OR IGNORE INTO participant (server_id, discord_id, botm_points, kots_points)
		VALUES (?, ?, 0, 0)`,
		serverID, discordID)
	if err != nil {
		utils.Error("Failed to create participant for discord ID %s: %v", discordID, err)
		return fmt.Errorf("failed to create participant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.Error("Failed to get rows affected for participant creation: %v", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected > 0 {
		utils.Debug("Created new participant for discord ID %s", discordID)
	}

	result, err = ds.db.Exec(`
		INSERT OR IGNORE INTO account (participant_id, username, failed_fetch_count)
		SELECT id, ?, 0 FROM participant WHERE server_id = ? AND discord_id = ?`,
		accountName, serverID, discordID)
	if err != nil {
		utils.Error("Failed to add account %s for participant %s: %v", accountName, discordID, err)
		return fmt.Errorf("failed to add account: %w", err)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		utils.Error("Failed to get rows affected for account creation: %v", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected > 0 {
		utils.Info("Successfully added account %s for participant %s", accountName, discordID)
	} else {
		utils.Debug("Account %s already exists for participant %s", accountName, discordID)
	}

	return nil
}

func (ds *participantDS) RemoveAccount(serverID, discordID, accountName string) error {
	result, err := ds.db.Exec(`
		DELETE FROM account 
		WHERE account.participant_id = (
			SELECT id FROM participant 
			WHERE server_id = ? AND discord_id = ?
		) AND account.username = ?`, serverID, discordID, accountName)
	if err != nil {
		utils.Error("Failed to remove account %s for participant %s: %v", accountName, discordID, err)
		return fmt.Errorf("failed to remove account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.Error("Failed to get rows affected for account removal %s: %v", accountName, err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		utils.Error("Account %s not found for participant %s in server %s", accountName, discordID, serverID)
		return fmt.Errorf("account %s not found for participant %s", accountName, discordID)
	}

	utils.Info("Successfully removed account %s for participant %s", accountName, discordID)
	return nil
}
