package sqlite

import (
	"database/sql"
	"fmt"
	"misclicked-events/internal/utils"
	"time"
)

type ParticipantDataSource interface {
	AddAccount(serverID, discordID, accountName string) error
}

func NewParticipantDataSource(db *sql.DB) ParticipantDataSource {
	return &participantDS{db: db}
}

type participantDS struct{ db *sql.DB }

func (ds *participantDS) AddAccount(serverID, discordID, accountName string) error {
	tx, err := ds.db.Begin()
	if err != nil {
		utils.Error("Failed to begin transaction for adding account %s to participant %s: %v", accountName, discordID, err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var participantID int64
	query := `SELECT id FROM participant WHERE server_id = ? AND discord_id = ?`
	err = tx.QueryRow(query, serverID, discordID).Scan(&participantID)
	if err == sql.ErrNoRows {
		result, err := tx.Exec(`
			INSERT INTO participant (server_id, discord_id, points, botm_enabled, kots_enabled, created_at, updated_at)
			VALUES (?, ?, 0, true, true, ?, ?)`,
			serverID, discordID, time.Now(), time.Now())
		if err != nil {
			utils.Error("Failed to create participant for discord ID %s: %v", discordID, err)
			return fmt.Errorf("failed to create participant: %w", err)
		}

		participantID, err = result.LastInsertId()
		if err != nil {
			utils.Error("Failed to get participant ID for discord ID %s: %v", discordID, err)
			return fmt.Errorf("failed to get participant ID: %w", err)
		}

		utils.Debug("Created new participant with ID %d for discord ID %s", participantID, discordID)
	} else if err != nil {
		utils.Error("Failed to check if participant exists for discord ID %s: %v", discordID, err)
		return fmt.Errorf("failed to check if participant exists: %w", err)
	}

	var existingAccountID int64
	accountQuery := `SELECT id FROM account WHERE participant_id = ? AND name = ?`
	err = tx.QueryRow(accountQuery, participantID, accountName).Scan(&existingAccountID)
	if err == nil {
		utils.Debug("Account %s already exists for participant %s", accountName, discordID)
		return tx.Commit()
	} else if err != sql.ErrNoRows {
		utils.Error("Failed to check if account %s exists for participant %s: %v", accountName, discordID, err)
		return fmt.Errorf("failed to check if account exists: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO account (participant_id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)`,
		participantID, accountName, time.Now(), time.Now())
	if err != nil {
		utils.Error("Failed to add account %s for participant %s: %v", accountName, discordID, err)
		return fmt.Errorf("failed to add account: %w", err)
	}

	utils.Info("Successfully added account %s for participant %s", accountName, discordID)
	return tx.Commit()
}
