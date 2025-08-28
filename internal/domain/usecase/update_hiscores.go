package usecase

import (
	"fmt"
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/domain"
	"time"
)

type UpdateHiscoresUseCase struct {
	competitionRepo CompetitionRepository
	participantRepo ParticipantRepository
	configRepo      ConfigRepository
	hiscoreRepo     HiscoreRepository
}

// ConfigRepository interface for bot configuration operations
type ConfigRepository interface {
	FetchConfig(serverID string) (*domain.Config, error)
	EditHiscoreMessageID(serverID, messageID string) error
}

func NewUpdateHiscoresUseCase(
	competitionRepo CompetitionRepository,
	participantRepo ParticipantRepository,
	configRepo ConfigRepository,
	hiscoreRepo HiscoreRepository,
) *UpdateHiscoresUseCase {
	return &UpdateHiscoresUseCase{
		competitionRepo: competitionRepo,
		participantRepo: participantRepo,
		configRepo:      configRepo,
		hiscoreRepo:     hiscoreRepo,
	}
}

// Execute updates all tracked accounts and hiscores for a server, and returns the result
func (uc *UpdateHiscoresUseCase) Execute(serverID string) (*UpdateHiscoresResult, error) {
	// Check if there's an ongoing event
	currentBotm, err := uc.competitionRepo.GetBotm(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to check current activity: %w", err)
	}

	if currentBotm == nil || currentBotm.Status != "active" {
		// No ongoing event, return result indicating no event
		return &UpdateHiscoresResult{
			HasOngoingEvent: false,
			EventStatus:     "none",
		}, nil
	}

	// Get all participants with accounts for this server
	participants, err := uc.participantRepo.GetAllParticipantsWithAccounts(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	// Update KC for each participant's accounts
	for _, participant := range participants {
		err := uc.updateParticipantKC(participant, currentBotm)
		if err != nil {
			// Log error but continue with other participants
			fmt.Printf("Failed to update KC for participant %s: %v\n", participant.DiscordID, err)
			continue
		}
	}

	// Get all participants with their KC data for the leaderboard
	var allParticipantsKC []domain.ParticipantWithAccountKC
	allParticipants, err := uc.participantRepo.GetAllParticipantsWithAccounts(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants for leaderboard: %w", err)
	}

	for _, participant := range allParticipants {
		participantKC, err := uc.participantRepo.GetParticipantWithAccountKC(participant.DiscordID, currentBotm.ID)
		if err != nil {
			// Skip participants with errors
			continue
		}
		if participantKC != nil {
			allParticipantsKC = append(allParticipantsKC, *participantKC)
		}
	}

	// Get bot configuration
	config, err := uc.configRepo.FetchConfig(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bot configuration: %w", err)
	}

	// Build the result
	result := &UpdateHiscoresResult{
		HasOngoingEvent: true,
		EventStatus:     currentBotm.Status,
		EventName:       currentBotm.Activity.Name,
		BossNames:       currentBotm.Activity.HiscoreNames,
		Threshold:       currentBotm.Activity.Threshold,
		Participants:    allParticipantsKC,
		Config:          config,
		LastUpdated:     time.Now(),
	}

	return result, nil
}

// updateParticipantKC updates the kill count for a specific participant
func (uc *UpdateHiscoresUseCase) updateParticipantKC(participant sqlite.ParticipantWithAccounts, currentBotm *domain.BotmWithActivity) error {
	// Get current KC for each account
	accountStartingKC := make(map[string]int)
	for _, account := range participant.Accounts {
		if uc.hiscoreRepo == nil {
			return fmt.Errorf("hiscore repository is not initialized")
		}

		startKC, err := uc.hiscoreRepo.GetBossKCCombined(account, currentBotm.Activity.HiscoreNames)
		if err != nil {
			// Log error but continue with other accounts
			fmt.Printf("Failed to get KC for account %s: %v\n", account, err)
			continue
		}
		accountStartingKC[account] = startKC
	}

	// Update the participation with new KC values
	err := uc.participantRepo.AddBotmParticipation(participant.DiscordID, currentBotm.ID, accountStartingKC)
	if err != nil {
		return fmt.Errorf("failed to update participation: %w", err)
	}

	return nil
}

// UpdateHiscoreMessageID updates the hiscore message ID in the configuration
func (uc *UpdateHiscoresUseCase) UpdateHiscoreMessageID(serverID, messageID string) error {
	return uc.configRepo.EditHiscoreMessageID(serverID, messageID)
}

// UpdateHiscoresResult contains all the data needed by the service to update Discord messages
type UpdateHiscoresResult struct {
	HasOngoingEvent bool
	EventStatus     string
	EventName       string
	BossNames       []string
	Threshold       int
	Participants    []domain.ParticipantWithAccountKC
	Config          *domain.Config
	LastUpdated     time.Time
}
