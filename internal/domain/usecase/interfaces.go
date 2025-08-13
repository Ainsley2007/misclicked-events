package usecase

import (
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/domain"
)

type CompetitionRepository interface {
	GetBotm(serverID string) (*domain.BotmWithActivity, error)
	StartBotm(serverID string, activityID int64, password string) error
}

type ParticipantRepository interface {
	GetAllParticipantsWithAccounts(serverID string) ([]sqlite.ParticipantWithAccounts, error)
	AddBotmParticipation(participantID string, botmID int64, startingKC int) error
}

type ActivityRepository interface {
	GetActivityByName(name string) (*domain.ActivityEntity, error)
}

type HiscoreRepository interface {
	GetBossKCCombined(username string, bossNames []string) (int, error)
}
