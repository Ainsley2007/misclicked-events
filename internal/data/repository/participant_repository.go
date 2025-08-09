package repository

import (
	"fmt"
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/data/mappers"
	"misclicked-events/internal/utils"
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
		return fmt.Errorf("server ID cannot be empty")
	}

	if discordID == "" {
		return fmt.Errorf("discord ID cannot be empty")
	}

	if accountName == "" {
		return fmt.Errorf("account name cannot be empty")
	}

	utils.Debug("Adding account %s for participant %s in server %s", accountName, discordID, serverID)
	err := r.ds.AddAccount(serverID, discordID, accountName)
	if err != nil {
		utils.Error("Failed to add account %s for participant %s in server %s: %v", accountName, discordID, serverID, err)
		return fmt.Errorf("failed to add account %s for participant %s: %w", accountName, discordID, err)
	}

	utils.Info("Successfully added account %s for participant %s in server %s", accountName, discordID, serverID)
	return nil
}

func (r *ParticipantRepository) RemoveAccount(serverID, discordID, accountName string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}

	if discordID == "" {
		return fmt.Errorf("discord ID cannot be empty")
	}

	if accountName == "" {
		return fmt.Errorf("account name cannot be empty")
	}

	utils.Debug("Removing account %s for participant %s in server %s", accountName, discordID, serverID)
	err := r.ds.RemoveAccount(serverID, discordID, accountName)
	if err != nil {
		utils.Error("Failed to remove account %s for participant %s in server %s: %v", accountName, discordID, serverID, err)
		return fmt.Errorf("failed to remove account %s for participant %s: %w", accountName, discordID, err)
	}

	utils.Info("Successfully removed account %s for participant %s in server %s", accountName, discordID, serverID)
	return nil
}
