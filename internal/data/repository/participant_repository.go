package repository

import (
	"fmt"
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/data/mappers"
	"misclicked-events/internal/domain"
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

	err := r.ds.AddAccount(serverID, discordID, accountName)
	if err != nil {
		if strings.Contains(err.Error(), "already tracked") {
			return err
		}
		utils.Error("Failed to add account %s for participant %s in server %s: %v", accountName, discordID, serverID, err)
		return fmt.Errorf("failed to add account")
	}

	utils.Debug("Added account %s for participant %s in server %s", accountName, discordID, serverID)
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

	err := r.ds.RemoveAccount(serverID, discordID, accountName)
	if err != nil {
		utils.Error("Failed to remove account %s for participant %s in server %s: %v", accountName, discordID, serverID, err)
		return fmt.Errorf("failed to remove account")
	}

	utils.Debug("Removed account %s for participant %s in server %s", accountName, discordID, serverID)
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
		return fmt.Errorf("old and new usernames cannot be the same")
	}

	err := r.ds.RenameAccount(serverID, discordID, oldUsername, newUsername)
	if err != nil {
		utils.Error("Failed to rename account %s to %s for participant %s in server %s: %v", oldUsername, newUsername, discordID, serverID, err)
		return fmt.Errorf("failed to rename account")
	}

	utils.Debug("Renamed account %s to %s for participant %s in server %s", oldUsername, newUsername, discordID, serverID)
	return nil
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

	accounts, err := r.ds.GetTrackedAccounts(serverID, discordID)
	if err != nil {
		utils.Error("Failed to get tracked accounts for server %s and discord ID %s: %v", serverID, discordID, err)
		return nil, fmt.Errorf("failed to get accounts")
	}

	return accounts, nil
}

func (r *ParticipantRepository) AddBotmParticipation(participantID string, botmID int64, accountStartingKC map[string]int) error {
	if participantID == "" {
		utils.Error("AddBotmParticipation called with empty participant ID")
		return fmt.Errorf("participant ID cannot be empty")
	}

	if botmID <= 0 {
		utils.Error("AddBotmParticipation called with invalid BOTM ID: %d", botmID)
		return fmt.Errorf("BOTM ID must be positive")
	}

	// Get all accounts with IDs for this participant
	accounts, err := r.ds.GetTrackedAccountsWithIDs("", participantID) // serverID not needed for this query
	if err != nil {
		utils.Error("Failed to get tracked accounts with IDs for participant %s: %v", participantID, err)
		return fmt.Errorf("failed to get tracked accounts with IDs")
	}

	if len(accounts) == 0 {
		utils.Error("Participant %s has no tracked accounts", participantID)
		return fmt.Errorf("participant has no tracked accounts")
	}

	// For each account, create participation
	for _, account := range accounts {
		accountID := account.ID
		accountName := account.Username

		existingParticipation, err := r.ds.GetBotmParticipation(accountID, botmID)
		if err != nil {
			utils.Error("Failed to check existing BOTM participation for account %s in BOTM %d: %v", accountName, botmID, err)
			continue
		}

		if existingParticipation == nil {
			accountStartKC, exists := accountStartingKC[accountName]
			if !exists {
				utils.Error("No starting KC provided for account %s", accountName)
				continue
			}

			err = r.ds.CreateBotmParticipation(accountID, botmID, accountStartKC)
			if err != nil {
				utils.Error("Failed to create BOTM participation for account %s in BOTM %d: %v", accountName, botmID, err)
				continue
			}
		}
	}

	return nil
}

func (r *ParticipantRepository) GetAllParticipantsWithAccounts(serverID string) ([]sqlite.ParticipantWithAccounts, error) {
	if serverID == "" {
		utils.Error("GetAllParticipantsWithAccounts called with empty server ID")
		return nil, fmt.Errorf("server ID cannot be empty")
	}

	participants, err := r.ds.GetAllParticipantsWithAccounts(serverID)
	if err != nil {
		utils.Error("Failed to get participants with accounts for server %s: %v", serverID, err)
		return nil, fmt.Errorf("failed to get participants with accounts")
	}

	return participants, nil
}

func (r *ParticipantRepository) GetParticipantWithAccountKC(participantID string, botmID int64) (*domain.ParticipantWithAccountKC, error) {
	if participantID == "" {
		utils.Error("GetParticipantWithAccountKC called with empty participant ID")
		return nil, fmt.Errorf("participant ID cannot be empty")
	}

	if botmID <= 0 {
		utils.Error("GetParticipantWithAccountKC called with invalid BOTM ID: %d", botmID)
		return nil, fmt.Errorf("BOTM ID must be positive")
	}

	participant, err := r.ds.GetParticipantWithAccountKC(participantID, botmID)
	if err != nil {
		utils.Error("Failed to get participant with account KC for participant %s in BOTM %d: %v", participantID, botmID, err)
		return nil, fmt.Errorf("failed to get participant with account KC")
	}

	return participant, nil
}
