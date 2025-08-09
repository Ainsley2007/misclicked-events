package usecase

import (
	"fmt"
	"misclicked-events/internal/data/repository"
)

type AddAccountUseCase struct {
	participantRepo *repository.ParticipantRepository
	hiscoreRepo     *repository.HiscoreRepository
}

func NewAddAccountUseCase(participantRepo *repository.ParticipantRepository, hiscoreRepo *repository.HiscoreRepository) *AddAccountUseCase {
	return &AddAccountUseCase{
		participantRepo: participantRepo,
		hiscoreRepo:     hiscoreRepo,
	}
}

func (uc *AddAccountUseCase) Execute(serverID, discordID, accountName string) error {
	exists, err := uc.hiscoreRepo.CheckIfPlayerExists(accountName)
	if err != nil {
		return fmt.Errorf("failed to check if player exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("player '%s' does not exist in OSRS", accountName)
	}

	return uc.participantRepo.AddAccount(serverID, discordID, accountName)
}
