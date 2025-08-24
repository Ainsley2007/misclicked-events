package usecase

import (
	"fmt"
	"misclicked-events/internal/data/repository"
	"strings"
)

type AddAccountUseCase struct {
	participantRepo *repository.ParticipantRepository
	hiscoreRepo     *repository.HiscoreRepository
	competitionRepo *repository.CompetitionRepository
}

func NewAddAccountUseCase(participantRepo *repository.ParticipantRepository, hiscoreRepo *repository.HiscoreRepository, competitionRepo *repository.CompetitionRepository) *AddAccountUseCase {
	return &AddAccountUseCase{
		participantRepo: participantRepo,
		hiscoreRepo:     hiscoreRepo,
		competitionRepo: competitionRepo,
	}
}

func (uc *AddAccountUseCase) Execute(serverID, discordID, accountName string) error {
	exists, err := uc.hiscoreRepo.CheckIfPlayerExists(accountName)
	if err != nil {
		return fmt.Errorf("failed to check username")
	}
	if !exists {
		return fmt.Errorf("player does not exist")
	}

	err = uc.participantRepo.AddAccount(serverID, discordID, accountName)
	if err != nil {
		if strings.Contains(err.Error(), "already tracked") {
			return err
		}
		return fmt.Errorf("failed to add account")
	}

	hasActiveBotm, err := uc.competitionRepo.HasRunningBotmCompetition(serverID)
	if err != nil {
		return fmt.Errorf("failed to add account")
	}

	if hasActiveBotm {
		botm, err := uc.competitionRepo.GetBotm(serverID)
		if err != nil {
			return fmt.Errorf("failed to add account")
		}

		startingKC, err := uc.hiscoreRepo.GetBossKCCombined(accountName, botm.Activity.HiscoreNames)
		if err != nil {
			return fmt.Errorf("failed to add account")
		}

		// Create a map with the single account's KC
		accountStartingKC := map[string]int{accountName: startingKC}
		err = uc.participantRepo.AddBotmParticipation(discordID, botm.ID, accountStartingKC)
		if err != nil {
			return fmt.Errorf("failed to add account")
		}
	}

	return nil
}
