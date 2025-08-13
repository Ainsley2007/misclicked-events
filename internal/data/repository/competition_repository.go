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

func (r *CompetitionRepository) StartBotm(serverID string, activityID int64, password string) error {
	return r.botmDS.Start(serverID, activityID, password)
}

func (r *CompetitionRepository) StopBotm(serverID string) error {
	return r.botmDS.Stop(serverID)
}

func (r *CompetitionRepository) GetBotm(serverID string) (*domain.BotmWithActivity, error) {
	botmModel, err := r.botmDS.GetCurrentBotmWithActivity(serverID)
	if err != nil {
		return nil, err
	}
	if botmModel == nil {
		return nil, nil
	}

	activityEntity := r.botmMapper.ToDomainActivity(botmModel.Activity)
	return &domain.BotmWithActivity{
		ID:         botmModel.ID,
		ServerID:   botmModel.ServerID,
		ActivityID: botmModel.ActivityID,
		Password:   botmModel.Password,
		Status:     botmModel.Status,
		Activity:   activityEntity,
	}, nil
}
