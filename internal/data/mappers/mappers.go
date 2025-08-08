package mappers

import (
	"misclicked-events/internal/data/datasource/api"
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/domain"
)

// ServerMapper converts between Server domain entity and ServerModel
type ServerMapper struct{}

func NewServerMapper() *ServerMapper {
	return &ServerMapper{}
}

func (m *ServerMapper) ToDomain(model *sqlite.ServerModel) *domain.Server {
	if model == nil {
		return nil
	}
	return &domain.Server{
		ID:   model.ID,
		Name: model.Name,
	}
}

func (m *ServerMapper) ToModel(entity *domain.Server) *sqlite.ServerModel {
	if entity == nil {
		return nil
	}
	return &sqlite.ServerModel{
		ID:   entity.ID,
		Name: entity.Name,
	}
}

// ConfigMapper converts between Config domain entity and ConfigModel
type ConfigMapper struct{}

func NewConfigMapper() *ConfigMapper {
	return &ConfigMapper{}
}

func (m *ConfigMapper) ToDomain(model *sqlite.ConfigModel, serverID string) *domain.Config {
	if model == nil {
		return nil
	}
	return &domain.Config{
		RankingChannelID:  model.RankingChannelID,
		HiscoreChannelID:  model.HiscoreChannelID,
		CategoryChannelID: model.CategoryChannelID,
		RankingMessageID:  model.RankingMessageID,
		HiscoreMessageID:  model.HiscoreMessageID,
	}
}

func (m *ConfigMapper) ToModel(entity *domain.Config, serverID string) *sqlite.ConfigModel {
	if entity == nil {
		return nil
	}
	return &sqlite.ConfigModel{
		ServerID:          serverID,
		RankingChannelID:  entity.RankingChannelID,
		HiscoreChannelID:  entity.HiscoreChannelID,
		CategoryChannelID: entity.CategoryChannelID,
		RankingMessageID:  entity.RankingMessageID,
		HiscoreMessageID:  entity.HiscoreMessageID,
	}
}

// BotmMapper converts between Botm domain entity and BotmModel
type BotmMapper struct{}

func NewBotmMapper() *BotmMapper {
	return &BotmMapper{}
}

func (m *BotmMapper) ToDomain(model *sqlite.BotmModel) *domain.Botm {
	if model == nil {
		return nil
	}
	return &domain.Botm{
		ID:          model.ID,
		CurrentBoss: model.CurrentBoss,
		Password:    model.Password,
		Status:      model.Status,
	}
}

func (m *BotmMapper) ToModel(entity *domain.Botm, serverID string) *sqlite.BotmModel {
	if entity == nil {
		return nil
	}
	return &sqlite.BotmModel{
		ID:          entity.ID,
		ServerID:    serverID,
		CurrentBoss: entity.CurrentBoss,
		Password:    entity.Password,
		Status:      entity.Status,
	}
}

type KotsMapper struct{}

func NewKotsMapper() *KotsMapper {
	return &KotsMapper{}
}

func (m *KotsMapper) ToDomain(model *sqlite.KotsModel) *domain.Kots {
	if model == nil {
		return nil
	}
	return &domain.Kots{
		ID:                     model.ID,
		CurrentSkill:           model.CurrentSkill,
		CurrentKingParticipant: model.CurrentKingParticipant,
		Streak:                 model.Streak,
		StartDate:              model.StartDate,
		EndDate:                model.EndDate,
		Status:                 model.Status,
	}
}

func (m *KotsMapper) ToModel(entity *domain.Kots, serverID string) *sqlite.KotsModel {
	if entity == nil {
		return nil
	}
	return &sqlite.KotsModel{
		ID:                     entity.ID,
		ServerID:               serverID,
		CurrentSkill:           entity.CurrentSkill,
		CurrentKingParticipant: entity.CurrentKingParticipant,
		Streak:                 entity.Streak,
		StartDate:              entity.StartDate,
		EndDate:                entity.EndDate,
		Status:                 entity.Status,
	}
}

type HiscoreDataMapper struct{}

func NewHiscoreDataMapper() *HiscoreDataMapper {
	return &HiscoreDataMapper{}
}

func (m *HiscoreDataMapper) ToDomain(model *api.HiscoreDataModel) *domain.HiscoreData {
	if model == nil {
		return nil
	}

	skills := make([]domain.Skill, len(model.Skills))
	for i, skill := range model.Skills {
		skills[i] = domain.Skill{
			ID:    skill.ID,
			Name:  skill.Name,
			Rank:  skill.Rank,
			Level: skill.Level,
			XP:    skill.XP,
		}
	}

	activities := make([]domain.Activity, len(model.Activities))
	for i, activity := range model.Activities {
		activities[i] = domain.Activity{
			ID:    activity.ID,
			Name:  activity.Name,
			Rank:  activity.Rank,
			Score: activity.Score,
		}
	}

	return &domain.HiscoreData{
		Skills:     skills,
		Activities: activities,
	}
}
