package repository

import (
	"misclicked-events/internal/data/datasource/sqlite"
)

type CompetitionRepository struct {
	botmDS sqlite.BotmDataSource
	kotsDS sqlite.KotsDataSource
}

func NewCompetitionRepository(botmDS sqlite.BotmDataSource, kotsDS sqlite.KotsDataSource) *CompetitionRepository {
	return &CompetitionRepository{botmDS: botmDS, kotsDS: kotsDS}
}

func (r *CompetitionRepository) HasRunningBotmCompetition(serverID string) (bool, error) {
	competition, err := r.botmDS.GetCurrentBotm(serverID)
	if err != nil {
		return false, err
	}
	return competition != nil, nil
}

func (r *CompetitionRepository) GetBotm(serverID string) (*sqlite.Botm, error) {
	return r.botmDS.GetCurrentBotm(serverID)
}

func (r *CompetitionRepository) StartBotm(serverID, currentBoss, password string) error {
	return r.botmDS.Start(serverID, currentBoss, password)
}

func (r *CompetitionRepository) StopBotm(serverID string) error {
	return r.botmDS.Stop(serverID)
}
