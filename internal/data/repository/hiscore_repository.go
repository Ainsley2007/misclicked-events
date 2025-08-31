package repository

import (
	"fmt"
	"misclicked-events/internal/data/datasource/api"
	"misclicked-events/internal/data/mappers"
	"misclicked-events/internal/domain"
	"misclicked-events/internal/utils"
)

type HiscoreRepository struct {
	hiscoreDS api.HiscoreDataSource
	mapper    *mappers.HiscoreDataMapper
}

func NewHiscoreRepository(hiscoreDS api.HiscoreDataSource) *HiscoreRepository {
	return &HiscoreRepository{
		hiscoreDS: hiscoreDS,
		mapper:    mappers.NewHiscoreDataMapper(),
	}
}

func (r *HiscoreRepository) CheckIfPlayerExists(username string) (bool, error) {
	if username == "" {
		utils.Error("CheckIfPlayerExists called with empty username")
		return false, fmt.Errorf("username cannot be empty")
	}

	exists, err := r.hiscoreDS.CheckIfPlayerExists(username)
	if err != nil {
		utils.Error("Failed to check if player %s exists: %v", username, err)
		return false, fmt.Errorf("failed to check player")
	}

	return exists, nil
}

// FetchHiscore fetches the complete hiscore data for a player
func (r *HiscoreRepository) FetchHiscore(username string) (*domain.HiscoreData, error) {
	if username == "" {
		utils.Error("FetchHiscore called with empty username")
		return nil, fmt.Errorf("username cannot be empty")
	}

	data, err := r.hiscoreDS.FetchHiscore(username)
	if err != nil {
		utils.Error("Failed to fetch hiscore data for player %s: %v", username, err)
		return nil, fmt.Errorf("failed to fetch hiscore")
	}

	domainData := r.mapper.ToDomain(data)
	return domainData, nil
}

func (r *HiscoreRepository) GetBossKC(username, bossName string) (int, error) {
	if username == "" {
		utils.Error("GetBossKC called with empty username")
		return 0, fmt.Errorf("username cannot be empty")
	}

	if bossName == "" {
		utils.Error("GetBossKC called with empty boss name")
		return 0, fmt.Errorf("boss name cannot be empty")
	}

	hiscoreData, err := r.FetchHiscore(username)
	if err != nil {
		utils.Error("Failed to fetch hiscore data for boss KC check for player %s on boss %s: %v", username, bossName, err)
		return 0, fmt.Errorf("failed to get boss KC")
	}

	activity, found := FindActivity(hiscoreData.Activities, bossName)
	if !found {
		return 0, nil
	}

	return activity.Score, nil
}

func (r *HiscoreRepository) GetBossKCCombined(username string, bossNames []string) (int, error) {
	if username == "" {
		utils.Error("GetBossKCCombined called with empty username")
		return 0, fmt.Errorf("username cannot be empty")
	}

	if len(bossNames) == 0 {
		utils.Error("GetBossKCCombined called with empty boss names")
		return 0, fmt.Errorf("boss names cannot be empty")
	}

	hiscoreData, err := r.FetchHiscore(username)
	if err != nil {
		utils.Error("Failed to fetch hiscore data for combined boss KC check for player %s: %v", username, err)
		return 0, fmt.Errorf("failed to get boss KC")
	}

	totalKC := 0
	for _, bossName := range bossNames {
		activity, found := FindActivity(hiscoreData.Activities, bossName)
		if found {
			totalKC += activity.Score
		}
	}

	return totalKC, nil
}

func FindSkill(skills []domain.Skill, name string) (*domain.Skill, bool) {
	for i, skill := range skills {
		if skill.Name == name {
			return &skills[i], true
		}
	}
	return nil, false
}

func FindActivity(activities []domain.Activity, name string) (*domain.Activity, bool) {
	for i, activity := range activities {
		if activity.Name == name {
			return &activities[i], true
		}
	}
	return nil, false
}
