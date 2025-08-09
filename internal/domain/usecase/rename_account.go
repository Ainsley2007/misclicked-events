package usecase

import (
	"fmt"
	"misclicked-events/internal/data/repository"
	"misclicked-events/internal/utils"
)

type RenameAccountUseCase struct {
	participantRepo *repository.ParticipantRepository
	hiscoreRepo     *repository.HiscoreRepository
}

func NewRenameAccountUseCase(participantRepo *repository.ParticipantRepository, hiscoreRepo *repository.HiscoreRepository) *RenameAccountUseCase {
	return &RenameAccountUseCase{
		participantRepo: participantRepo,
		hiscoreRepo:     hiscoreRepo,
	}
}

func (uc *RenameAccountUseCase) Execute(serverID, discordID, oldUsername, newUsername string) error {
	utils.Debug("RenameAccountUseCase: Starting execution for old username %s, new username %s, participant %s, server %s", oldUsername, newUsername, discordID, serverID)

	utils.Debug("RenameAccountUseCase: Checking if new player %s exists in OSRS", newUsername)
	exists, err := uc.hiscoreRepo.CheckIfPlayerExists(newUsername)
	if err != nil {
		utils.Error("RenameAccountUseCase: Failed to check if new player %s exists: %v", newUsername, err)
		return fmt.Errorf("failed to check username")
	}
	if !exists {
		utils.Debug("RenameAccountUseCase: New player %s does not exist in OSRS", newUsername)
		return fmt.Errorf("player does not exist")
	}
	utils.Debug("RenameAccountUseCase: New player %s exists in OSRS", newUsername)

	utils.Debug("RenameAccountUseCase: Renaming account %s to %s in database", oldUsername, newUsername)
	err = uc.participantRepo.RenameAccount(serverID, discordID, oldUsername, newUsername)
	if err != nil {
		utils.Error("RenameAccountUseCase: Failed to rename account %s to %s for participant %s in server %s: %v", oldUsername, newUsername, discordID, serverID, err)
		return fmt.Errorf("failed to rename account")
	}
	utils.Info("RenameAccountUseCase: Successfully renamed account %s to %s for participant %s in server %s", oldUsername, newUsername, discordID, serverID)

	utils.Info("RenameAccountUseCase: Successfully completed execution for old username %s, new username %s, participant %s, server %s", oldUsername, newUsername, discordID, serverID)
	return nil
}
