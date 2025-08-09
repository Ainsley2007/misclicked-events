package repository

import (
	"fmt"
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/data/mappers"
	"misclicked-events/internal/utils"
	"strings"
)

type ParticipantRepository struct {
	ds     sqlite.ParticipantDataSource
	mapper *mappers.AccountMapper
}

func NewParticipantRepository(ds sqlite.ParticipantDataSource) *ParticipantRepository {
	return &ParticipantRepository{
		ds:     ds,
		mapper: mappers.NewAccountMapper(),
	}
}

func (r *ParticipantRepository) AddAccount(serverID, discordID, accountName string) error {
	if serverID == "" {
		utils.Error("AddAccount called with empty server ID")
		return fmt.Errorf("server ID cannot be empty")
	}

	if discordID == "" {
		utils.Error("AddAccount called with empty discord ID")
		return fmt.Errorf("discord ID cannot be empty")
	}

	if accountName == "" {
		utils.Error("AddAccount called with empty account name")
		return fmt.Errorf("account name cannot be empty")
	}

	utils.Debug("Adding account %s for participant %s in server %s", accountName, discordID, serverID)
	err := r.ds.AddAccount(serverID, discordID, accountName)
	if err != nil {
		if strings.Contains(err.Error(), "already tracked") {
			utils.Debug("Account %s already tracked for participant %s in server %s", accountName, discordID, serverID)
			return err
		}
		utils.Error("Failed to add account %s for participant %s in server %s: %v", accountName, discordID, serverID, err)
		return fmt.Errorf("failed to add account")
	}

	utils.Info("Successfully added account %s for participant %s in server %s", accountName, discordID, serverID)
	return nil
}

func (r *ParticipantRepository) RemoveAccount(serverID, discordID, accountName string) error {
	if serverID == "" {
		utils.Error("RemoveAccount called with empty server ID")
		return fmt.Errorf("server ID cannot be empty")
	}

	if discordID == "" {
		utils.Error("RemoveAccount called with empty discord ID")
		return fmt.Errorf("discord ID cannot be empty")
	}

	if accountName == "" {
		utils.Error("RemoveAccount called with empty account name")
		return fmt.Errorf("account name cannot be empty")
	}

	utils.Debug("Removing account %s for participant %s in server %s", accountName, discordID, serverID)
	err := r.ds.RemoveAccount(serverID, discordID, accountName)
	if err != nil {
		utils.Error("Failed to remove account %s for participant %s in server %s: %v", accountName, discordID, serverID, err)
		return fmt.Errorf("failed to remove account")
	}

	utils.Info("Successfully removed account %s for participant %s in server %s", accountName, discordID, serverID)
	return nil
}

func (r *ParticipantRepository) RenameAccount(serverID, discordID, oldUsername, newUsername string) error {
	if serverID == "" {
		utils.Error("RenameAccount called with empty server ID")
		return fmt.Errorf("server ID cannot be empty")
	}

	if discordID == "" {
		utils.Error("RenameAccount called with empty discord ID")
		return fmt.Errorf("discord ID cannot be empty")
	}

	if oldUsername == "" {
		utils.Error("RenameAccount called with empty old username")
		return fmt.Errorf("old username cannot be empty")
	}

	if newUsername == "" {
		utils.Error("RenameAccount called with empty new username")
		return fmt.Errorf("new username cannot be empty")
	}

	if oldUsername == newUsername {
		utils.Debug("RenameAccount called with same old and new username: %s", oldUsername)
		return fmt.Errorf("old and new usernames cannot be the same")
	}

	utils.Debug("Renaming account %s to %s for participant %s in server %s", oldUsername, newUsername, discordID, serverID)
	err := r.ds.RenameAccount(serverID, discordID, oldUsername, newUsername)
	if err != nil {
		utils.Error("Failed to rename account %s to %s for participant %s in server %s: %v", oldUsername, newUsername, discordID, serverID, err)
		return fmt.Errorf("failed to rename account")
	}

	utils.Info("Successfully renamed account %s to %s for participant %s in server %s", oldUsername, newUsername, discordID, serverID)
	return nil
}

func (r *ParticipantRepository) GetParticipantID(serverID, discordID string) (int64, error) {
	if serverID == "" {
		utils.Error("GetParticipantID called with empty server ID")
		return 0, fmt.Errorf("server ID cannot be empty")
	}

	if discordID == "" {
		utils.Error("GetParticipantID called with empty discord ID")
		return 0, fmt.Errorf("discord ID cannot be empty")
	}

	utils.Debug("Getting participant ID for discord ID %s in server %s", discordID, serverID)
	participantID, err := r.ds.GetParticipantID(serverID, discordID)
	if err != nil {
		utils.Error("Failed to get participant ID for server %s and discord ID %s: %v", serverID, discordID, err)
		return 0, fmt.Errorf("failed to get participant")
	}

	utils.Debug("Retrieved participant ID %d for discord ID %s in server %s", participantID, discordID, serverID)
	return participantID, nil
}

func (r *ParticipantRepository) GetTrackedAccounts(serverID, discordID string) ([]string, error) {
	if serverID == "" {
		utils.Error("GetTrackedAccounts called with empty server ID")
		return nil, fmt.Errorf("server ID cannot be empty")
	}

	if discordID == "" {
		utils.Error("GetTrackedAccounts called with empty discord ID")
		return nil, fmt.Errorf("discord ID cannot be empty")
	}

	utils.Debug("Getting tracked accounts for discord ID %s in server %s", discordID, serverID)
	accounts, err := r.ds.GetTrackedAccounts(serverID, discordID)
	if err != nil {
		utils.Error("Failed to get tracked accounts for server %s and discord ID %s: %v", serverID, discordID, err)
		return nil, fmt.Errorf("failed to get accounts")
	}

	utils.Debug("Retrieved %d tracked accounts for discord ID %s in server %s: %v", len(accounts), discordID, serverID, accounts)
	return accounts, nil
}

func (r *ParticipantRepository) AddBotmParticipation(participantID, botmID int64, startingKC int) error {
	if participantID <= 0 {
		utils.Error("AddBotmParticipation called with invalid participant ID: %d", participantID)
		return fmt.Errorf("participant ID must be positive")
	}

	if botmID <= 0 {
		utils.Error("AddBotmParticipation called with invalid BOTM ID: %d", botmID)
		return fmt.Errorf("BOTM ID must be positive")
	}

	utils.Debug("Adding BOTM participation for participant %d in BOTM %d with starting KC %d", participantID, botmID, startingKC)
	err := r.ds.AddBotmParticipation(participantID, botmID, startingKC)
	if err != nil {
		utils.Error("Failed to add BOTM participation for participant %d in BOTM %d: %v", participantID, botmID, err)
		return fmt.Errorf("failed to add participation")
	}

	utils.Info("Successfully added BOTM participation for participant %d in BOTM %d with starting KC %d", participantID, botmID, startingKC)
	return nil
}
