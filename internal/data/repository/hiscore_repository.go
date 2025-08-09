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
	utils.Debug("Creating new HiscoreRepository")
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

	utils.Debug("Checking if player %s exists in OSRS", username)
	exists, err := r.hiscoreDS.CheckIfPlayerExists(username)
	if err != nil {
		utils.Error("Failed to check if player %s exists: %v", username, err)
		return false, fmt.Errorf("failed to check player")
	}

	if exists {
		utils.Debug("Player %s exists in OSRS", username)
	} else {
		utils.Debug("Player %s does not exist in OSRS", username)
	}

	return exists, nil
}

// FetchHiscore fetches the complete hiscore data for a player
func (r *HiscoreRepository) FetchHiscore(username string) (*domain.HiscoreData, error) {
	if username == "" {
		utils.Error("FetchHiscore called with empty username")
		return nil, fmt.Errorf("username cannot be empty")
	}

	utils.Debug("Fetching hiscore data for player %s", username)
	data, err := r.hiscoreDS.FetchHiscore(username)
	if err != nil {
		utils.Error("Failed to fetch hiscore data for player %s: %v", username, err)
		return nil, fmt.Errorf("failed to fetch hiscore")
	}

	utils.Debug("Mapping hiscore data for player %s", username)
	domainData := r.mapper.ToDomain(data)
	utils.Debug("Successfully mapped hiscore data for player %s: %d skills, %d activities", username, len(domainData.Skills), len(domainData.Activities))

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

	utils.Debug("Getting boss KC for player %s on boss %s", username, bossName)
	hiscoreData, err := r.FetchHiscore(username)
	if err != nil {
		utils.Error("Failed to fetch hiscore data for boss KC check for player %s on boss %s: %v", username, bossName, err)
		return 0, fmt.Errorf("failed to get boss KC")
	}

	utils.Debug("Searching for boss %s in activities for player %s", bossName, username)
	activity, found := FindActivity(hiscoreData.Activities, bossName)
	if !found {
		utils.Debug("Boss %s not found in activities for player %s", bossName, username)
		return 0, nil
	}

	utils.Debug("Found boss %s for player %s with KC %d", bossName, username, activity.Score)
	return activity.Score, nil
}

func FindSkill(skills []domain.Skill, name string) (*domain.Skill, bool) {
	utils.Debug("Searching for skill %s in %d skills", name, len(skills))
	for i, skill := range skills {
		if skill.Name == name {
			utils.Debug("Found skill %s at index %d", name, i)
			return &skills[i], true
		}
	}
	utils.Debug("Skill %s not found", name)
	return nil, false
}

func FindActivity(activities []domain.Activity, name string) (*domain.Activity, bool) {
	utils.Debug("Searching for activity %s in %d activities", name, len(activities))
	for i, activity := range activities {
		if activity.Name == name {
			utils.Debug("Found activity %s at index %d", name, i)
			return &activities[i], true
		}
	}
	utils.Debug("Activity %s not found", name)
	return nil, false
}
