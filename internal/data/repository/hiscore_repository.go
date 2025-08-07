package repository

import (
	"misclicked-events/internal/data/datasource/api"
)

// HiscoreRepository provides business logic for hiscore operations
type HiscoreRepository struct {
	hiscoreDS api.HiscoreDataSource
}

// NewHiscoreRepository creates a new instance of HiscoreRepository
func NewHiscoreRepository(hiscoreDS api.HiscoreDataSource) *HiscoreRepository {
	return &HiscoreRepository{
		hiscoreDS: hiscoreDS,
	}
}

// CheckIfPlayerExists checks if a player exists in OSRS
func (r *HiscoreRepository) CheckIfPlayerExists(username string) (bool, error) {
	return r.hiscoreDS.CheckIfPlayerExists(username)
}

// FetchHiscore fetches the complete hiscore data for a player
func (r *HiscoreRepository) FetchHiscore(username string) ([]api.Skill, []api.Activity, error) {
	data, err := r.hiscoreDS.FetchHiscore(username)
	if err != nil {
		return nil, nil, err
	}
	return data.Skills, data.Activities, nil
}

// FindSkill searches for a skill by its name in the provided slice.
// It returns a pointer to the found Skill and true if found; otherwise, nil and false.
func FindSkill(skills []api.Skill, name string) (*api.Skill, bool) {
	for i, skill := range skills {
		if skill.Name == name {
			return &skills[i], true
		}
	}
	return nil, false
}

// FindActivity searches for an activity by its name in the provided slice.
// It returns a pointer to the found Activity and true if found; otherwise, nil and false.
func FindActivity(activities []api.Activity, name string) (*api.Activity, bool) {
	for i, activity := range activities {
		if activity.Name == name {
			return &activities[i], true
		}
	}
	return nil, false
}
