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

	if len(selectedActivity.HiscoreNames) == 0 {
		return nil, fmt.Errorf("activity '%s' has no hiscore names configured", activityName)
	}

	participants, err := uc.participantRepo.GetAllParticipantsWithAccounts(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %v", err)
	}

	fmt.Printf("DEBUG: Found %d participants for server %s\n", len(participants), serverID)
	for i, participant := range participants {
		fmt.Printf("DEBUG: Participant %d: DiscordID=%s, Accounts=%v\n", i, participant.DiscordID, participant.Accounts)
	}

	err = uc.competitionRepo.StartBotm(serverID, selectedActivity.ID, password)
	if err != nil {
		return nil, fmt.Errorf("failed to start activity: %v", err)
	}

	botm, err := uc.competitionRepo.GetBotm(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get started activity: %v", err)
	}

	fmt.Printf("DEBUG: Processing %d participants for BOTM participation\n", len(participants))
	for _, participant := range participants {
		fmt.Printf("DEBUG: Processing participant %s with %d accounts\n", participant.DiscordID, len(participant.Accounts))

		// Combine KC from all accounts for this participant
		totalStartKC := 0
		for _, account := range participant.Accounts {
			fmt.Printf("DEBUG: Processing account: %s\n", account)
			if uc.hiscoreRepo == nil {
				return nil, fmt.Errorf("hiscore repository is not initialized")
			}
			startKC, err := uc.hiscoreRepo.GetBossKCCombined(account, selectedActivity.HiscoreNames)
			if err != nil {
				fmt.Printf("DEBUG: Failed to get KC for account %s: %v\n", account, err)
				continue
			}
			fmt.Printf("DEBUG: Got KC %d for account %s\n", startKC, account)
			totalStartKC += startKC
		}

		fmt.Printf("DEBUG: Total combined KC for participant %s: %d\n", participant.DiscordID, totalStartKC)

		// Add participation with the combined KC
		err = uc.participantRepo.AddBotmParticipation(participant.DiscordID, botm.ID, totalStartKC)
		if err != nil {
			fmt.Printf("DEBUG: Failed to add participation for participant %s: %v\n", participant.DiscordID, err)
			continue
		}
		fmt.Printf("DEBUG: Successfully added participation for participant %s with total KC %d\n", participant.DiscordID, totalStartKC)
	}

	return botm, nil
}
