package repository

import (
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/data/mappers"
	"misclicked-events/internal/domain"
)

type CompetitionRepository struct {
	botmDS     sqlite.BotmDataSource
	kotsDS     sqlite.KotsDataSource
	botmMapper *mappers.BotmMapper
	kotsMapper *mappers.KotsMapper
}

func NewCompetitionRepository(botmDS sqlite.BotmDataSource, kotsDS sqlite.KotsDataSource) *CompetitionRepository {
	return &CompetitionRepository{
		botmDS:     botmDS,
		kotsDS:     kotsDS,
		botmMapper: mappers.NewBotmMapper(),
		kotsMapper: mappers.NewKotsMapper(),
	}
}

func (r *CompetitionRepository) HasRunningBotmCompetition(serverID string) (bool, error) {
	competition, err := r.botmDS.GetCurrentBotm(serverID)
	if err != nil {
		return false, err
	}
	return competition != nil, nil
}

func (r *CompetitionRepository) GetBotm(serverID string) (*domain.Botm, error) {
	botmModel, err := r.botmDS.GetCurrentBotm(serverID)
	if err != nil {
		return nil, err
	}
	return r.botmMapper.ToDomain(botmModel), nil
}

func (r *CompetitionRepository) StartBotm(serverID, currentBoss, password string) error {
	return r.botmDS.Start(serverID, currentBoss, password)
}

func (r *CompetitionRepository) StopBotm(serverID string) error {
	return r.botmDS.Stop(serverID)
}
