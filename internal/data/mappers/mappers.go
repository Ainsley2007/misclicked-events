package mappers

import (
	"misclicked-events/internal/data/datasource/api"
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/domain"
)

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

type AccountMapper struct{}

func NewAccountMapper() *AccountMapper {
	return &AccountMapper{}
}

func (m *AccountMapper) ToDomain(participantModel *sqlite.ParticipantModel, accountModels []*sqlite.AccountModel, botmParticipationModels []*sqlite.BotmParticipationModel, kotsParticipationModels []*sqlite.KotsParticipationModel) *domain.Account {
	if participantModel == nil {
		return nil
	}

	osrsAccounts := make([]domain.OSRSAccount, len(accountModels))
	for i, accountModel := range accountModels {
		osrsAccounts[i] = domain.OSRSAccount{
			ID:        accountModel.ID,
			Name:      accountModel.Name,
			CreatedAt: accountModel.CreatedAt,
			UpdatedAt: accountModel.UpdatedAt,
		}
	}

	var botmParticipation *domain.BotmParticipation
	if len(botmParticipationModels) > 0 {
		latestBotm := botmParticipationModels[0]
		for _, participation := range botmParticipationModels {
			if participation.CreatedAt.After(latestBotm.CreatedAt) {
				latestBotm = participation
			}
		}
		botmParticipation = &domain.BotmParticipation{
			ID:            latestBotm.ID,
			BotmID:        latestBotm.BotmID,
			StartAmount:   latestBotm.StartAmount,
			CurrentAmount: latestBotm.CurrentAmount,
			CreatedAt:     latestBotm.CreatedAt,
			UpdatedAt:     latestBotm.UpdatedAt,
		}
	}

	var kotsParticipation *domain.KotsParticipation
	if len(kotsParticipationModels) > 0 {
		latestKots := kotsParticipationModels[0]
		for _, participation := range kotsParticipationModels {
			if participation.CreatedAt.After(latestKots.CreatedAt) {
				latestKots = participation
			}
		}
		kotsParticipation = &domain.KotsParticipation{
			ID:            latestKots.ID,
			KotsID:        latestKots.KotsID,
			StartAmount:   latestKots.StartAmount,
			CurrentAmount: latestKots.CurrentAmount,
			CreatedAt:     latestKots.CreatedAt,
			UpdatedAt:     latestKots.UpdatedAt,
		}
	}

	return &domain.Account{
		ID:                participantModel.ID,
		DiscordID:         participantModel.DiscordID,
		Points:            participantModel.Points,
		OSRSAccounts:      osrsAccounts,
		BotmParticipation: botmParticipation,
		KotsParticipation: kotsParticipation,
		BotmEnabled:       participantModel.BotmEnabled,
		KotsEnabled:       participantModel.KotsEnabled,
		CreatedAt:         participantModel.CreatedAt,
		UpdatedAt:         participantModel.UpdatedAt,
	}
}

func (m *AccountMapper) ToModels(entity *domain.Account, serverID string) (*sqlite.ParticipantModel, []*sqlite.AccountModel, []*sqlite.BotmParticipationModel, []*sqlite.KotsParticipationModel) {
	if entity == nil {
		return nil, nil, nil, nil
	}

	participantModel := &sqlite.ParticipantModel{
		ID:          entity.ID,
		ServerID:    serverID,
		DiscordID:   entity.DiscordID,
		Points:      entity.Points,
		BotmEnabled: entity.BotmEnabled,
		KotsEnabled: entity.KotsEnabled,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}

	accountModels := make([]*sqlite.AccountModel, len(entity.OSRSAccounts))
	for i, osrsAccount := range entity.OSRSAccounts {
		accountModels[i] = &sqlite.AccountModel{
			ID:            osrsAccount.ID,
			ParticipantID: entity.ID,
			Name:          osrsAccount.Name,
			CreatedAt:     osrsAccount.CreatedAt,
			UpdatedAt:     osrsAccount.UpdatedAt,
		}
	}

	var botmParticipationModels []*sqlite.BotmParticipationModel
	if entity.BotmParticipation != nil {
		botmParticipationModels = []*sqlite.BotmParticipationModel{
			{
				ID:            entity.BotmParticipation.ID,
				ParticipantID: entity.ID,
				BotmID:        entity.BotmParticipation.BotmID,
				StartAmount:   entity.BotmParticipation.StartAmount,
				CurrentAmount: entity.BotmParticipation.CurrentAmount,
				CreatedAt:     entity.BotmParticipation.CreatedAt,
				UpdatedAt:     entity.BotmParticipation.UpdatedAt,
			},
		}
	}

	var kotsParticipationModels []*sqlite.KotsParticipationModel
	if entity.KotsParticipation != nil {
		kotsParticipationModels = []*sqlite.KotsParticipationModel{
			{
				ID:            entity.KotsParticipation.ID,
				ParticipantID: entity.ID,
				KotsID:        entity.KotsParticipation.KotsID,
				StartAmount:   entity.KotsParticipation.StartAmount,
				CurrentAmount: entity.KotsParticipation.CurrentAmount,
				CreatedAt:     entity.KotsParticipation.CreatedAt,
				UpdatedAt:     entity.KotsParticipation.UpdatedAt,
			},
		}
	}

	return participantModel, accountModels, botmParticipationModels, kotsParticipationModels
}
