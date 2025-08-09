package usecase

import (
	"fmt"
	"misclicked-events/internal/data/repository"
	"misclicked-events/internal/utils"
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
	utils.Debug("AddAccountUseCase: Starting execution for account %s, participant %s, server %s", accountName, discordID, serverID)

	utils.Debug("AddAccountUseCase: Checking if player %s exists in OSRS", accountName)
	exists, err := uc.hiscoreRepo.CheckIfPlayerExists(accountName)
	if err != nil {
		utils.Error("AddAccountUseCase: Failed to check if player %s exists: %v", accountName, err)
		return fmt.Errorf("failed to check username")
	}
	if !exists {
		utils.Debug("AddAccountUseCase: Player %s does not exist in OSRS", accountName)
		return fmt.Errorf("player does not exist")
	}
	utils.Debug("AddAccountUseCase: Player %s exists in OSRS", accountName)

	utils.Debug("AddAccountUseCase: Adding account %s to database", accountName)
	err = uc.participantRepo.AddAccount(serverID, discordID, accountName)
	if err != nil {
		if strings.Contains(err.Error(), "already tracked") {
			utils.Debug("AddAccountUseCase: Account %s already tracked for participant %s", accountName, discordID)
			return err
		}
		utils.Error("AddAccountUseCase: Failed to add account %s for participant %s in server %s: %v", accountName, discordID, serverID, err)
		return fmt.Errorf("failed to add account")
	}
	utils.Info("AddAccountUseCase: Successfully added account %s for participant %s in server %s", accountName, discordID, serverID)

	utils.Debug("AddAccountUseCase: Checking for active BOTM competition in server %s", serverID)
	hasActiveBotm, err := uc.competitionRepo.HasRunningBotmCompetition(serverID)
	if err != nil {
		utils.Error("AddAccountUseCase: Failed to check for active BOTM competition in server %s: %v", serverID, err)
		return fmt.Errorf("failed to add account")
	}

	if hasActiveBotm {
		utils.Debug("AddAccountUseCase: Active BOTM competition found in server %s", serverID)
		botm, err := uc.competitionRepo.GetBotm(serverID)
		if err != nil {
			utils.Error("AddAccountUseCase: Failed to get active BOTM for server %s: %v", serverID, err)
			return fmt.Errorf("failed to add account")
		}
		utils.Debug("AddAccountUseCase: Retrieved BOTM %d with boss %s", botm.ID, botm.CurrentBoss)

		utils.Debug("AddAccountUseCase: Getting participant ID for %s in server %s", discordID, serverID)
		participantID, err := uc.participantRepo.GetParticipantID(serverID, discordID)
		if err != nil {
			utils.Error("AddAccountUseCase: Failed to get participant ID for %s in server %s: %v", discordID, serverID, err)
			return fmt.Errorf("failed to add account")
		}
		utils.Debug("AddAccountUseCase: Retrieved participant ID %d", participantID)

		utils.Debug("AddAccountUseCase: Getting current KC for %s on boss %s", accountName, botm.CurrentBoss)
		startingKC, err := uc.hiscoreRepo.GetBossKC(accountName, botm.CurrentBoss)
		if err != nil {
			utils.Error("AddAccountUseCase: Failed to get current KC for %s on boss %s: %v", accountName, botm.CurrentBoss, err)
			return fmt.Errorf("failed to add account")
		}
		utils.Debug("AddAccountUseCase: Retrieved starting KC %d for %s on boss %s", startingKC, accountName, botm.CurrentBoss)

		utils.Debug("AddAccountUseCase: Adding BOTM participation for participant %d in BOTM %d", participantID, botm.ID)
		err = uc.participantRepo.AddBotmParticipation(participantID, botm.ID, startingKC)
		if err != nil {
			utils.Error("AddAccountUseCase: Failed to add BOTM participation for participant %d in BOTM %d: %v", participantID, botm.ID, err)
			return fmt.Errorf("failed to add account")
		}
		utils.Info("AddAccountUseCase: Successfully added BOTM participation for participant %d in BOTM %d with starting KC %d", participantID, botm.ID, startingKC)
	} else {
		utils.Debug("AddAccountUseCase: No active BOTM competition found in server %s", serverID)
	}

	utils.Info("AddAccountUseCase: Successfully completed execution for account %s, participant %s, server %s", accountName, discordID, serverID)
	return nil
}
