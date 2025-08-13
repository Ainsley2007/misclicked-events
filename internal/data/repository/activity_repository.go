package repository

import (
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/data/mappers"
	"misclicked-events/internal/domain"
	"misclicked-events/internal/utils"
)

type ActivityRepository interface {
	AddActivity(activity *domain.ActivityEntity) error
	RemoveActivity(name string) error
	GetActivityListByType(activityType string) ([]*domain.ActivityEntity, error)
	GetAllActivities() ([]*domain.ActivityEntity, error)
	GetActivityByName(name string) (*domain.ActivityEntity, error)
}

type activityRepository struct {
	datasource sqlite.ActivityDataSource
	mapper     *mappers.ActivityMapper
}

func NewActivityRepository(datasource sqlite.ActivityDataSource) ActivityRepository {
	return &activityRepository{
		datasource: datasource,
		mapper:     mappers.NewActivityMapper(),
	}
}

func (r *activityRepository) AddActivity(activity *domain.ActivityEntity) error {
	utils.Info("Adding activity: %s", activity.Name)
	model := r.mapper.ToModel(activity)
	err := r.datasource.AddActivity(model)
	if err != nil {
		utils.Error("Failed to add activity %s: %v", activity.Name, err)
	} else {
		utils.Info("Successfully added activity: %s", activity.Name)
	}
	return err
}

func (r *activityRepository) RemoveActivity(name string) error {
	utils.Info("Removing activity: %s", name)
	err := r.datasource.RemoveActivity(name)
	if err != nil {
		utils.Error("Failed to remove activity %s: %v", name, err)
	} else {
		utils.Info("Successfully removed activity: %s", name)
	}
	return err
}

func (r *activityRepository) GetActivityListByType(activityType string) ([]*domain.ActivityEntity, error) {
	utils.Info("Getting activity list for type: %s", activityType)
	models, err := r.datasource.GetActivityListByType(activityType)
	if err != nil {
		utils.Error("Failed to get activity list for type %s: %v", activityType, err)
		return nil, err
	}

	entities := make([]*domain.ActivityEntity, len(models))
	for i, model := range models {
		entities[i] = r.mapper.ToDomain(model)
	}

	utils.Info("Retrieved %d activities for type: %s", len(entities), activityType)
	return entities, nil
}

func (r *activityRepository) GetAllActivities() ([]*domain.ActivityEntity, error) {
	utils.Info("Getting all activities")
	models, err := r.datasource.GetAllActivities()
	if err != nil {
		utils.Error("Failed to get all activities: %v", err)
		return nil, err
	}

	entities := make([]*domain.ActivityEntity, len(models))
	for i, model := range models {
		entities[i] = r.mapper.ToDomain(model)
	}

	utils.Info("Retrieved %d total activities", len(entities))
	return entities, nil
}

func (r *activityRepository) GetActivityByName(name string) (*domain.ActivityEntity, error) {
	utils.Info("Getting activity by name: %s", name)
	model, err := r.datasource.GetActivityByName(name)
	if err != nil {
		utils.Error("Failed to get activity by name %s: %v", name, err)
		return nil, err
	}
	if model == nil {
		utils.Info("No activity found for name: %s", name)
		return nil, nil
	}

	entity := r.mapper.ToDomain(model)
	utils.Info("Successfully retrieved activity: %s", entity.Name)
	return entity, nil
}
