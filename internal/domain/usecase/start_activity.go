package usecase

import (
	"fmt"
	"misclicked-events/internal/domain"
)

type StartActivityUseCase struct {
	competitionRepo CompetitionRepository
	participantRepo ParticipantRepository
	activityRepo    ActivityRepository
	hiscoreRepo     HiscoreRepository
}

func NewStartActivityUseCase(competitionRepo CompetitionRepository, participantRepo ParticipantRepository, activityRepo ActivityRepository, hiscoreRepo HiscoreRepository) *StartActivityUseCase {
	return &StartActivityUseCase{
		competitionRepo: competitionRepo,
		participantRepo: participantRepo,
		activityRepo:    activityRepo,
		hiscoreRepo:     hiscoreRepo,
	}
}

func (uc *StartActivityUseCase) Execute(serverID, activityName, password string) (*domain.BotmWithActivity, error) {
	currentBotm, err := uc.competitionRepo.GetBotm(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to check current activity: %v", err)
	}
	if currentBotm != nil {
		return currentBotm, fmt.Errorf("an activity has already been selected: \"%s\", you need to end this activity before starting a new one", currentBotm.Activity.Name)
	}

	selectedActivity, err := uc.activityRepo.GetActivityByName(activityName)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity: %v", err)
	}

	if selectedActivity == nil {
		return nil, fmt.Errorf("activity '%s' not found", activityName)
	}

	participants, err := uc.participantRepo.GetAllParticipantsWithAccounts(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %v", err)
	}

	err = uc.competitionRepo.StartBotm(serverID, selectedActivity.ID, password)
	if err != nil {
		return nil, fmt.Errorf("failed to start activity: %v", err)
	}

	botm, err := uc.competitionRepo.GetBotm(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get started activity: %v", err)
	}

	for _, participant := range participants {
		for _, account := range participant.Accounts {
			startKC, err := uc.hiscoreRepo.GetBossKCCombined(account, selectedActivity.HiscoreNames)
			if err != nil {
				continue
			}

			err = uc.participantRepo.AddBotmParticipation(participant.DiscordID, botm.ID, startKC)
			if err != nil {
				continue
			}
		}
	}

	return botm, nil
}
