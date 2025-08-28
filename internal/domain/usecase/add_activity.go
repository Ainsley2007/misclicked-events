package usecase

import (
	"misclicked-events/internal/data/repository"
	"misclicked-events/internal/domain"
	"misclicked-events/internal/utils"
)

type AddActivityUseCase struct {
	activityRepo repository.ActivityRepository
}

func NewAddActivityUseCase(activityRepo repository.ActivityRepository) *AddActivityUseCase {
	return &AddActivityUseCase{
		activityRepo: activityRepo,
	}
}

func (uc *AddActivityUseCase) Execute(name, activityType string, hiscoreNames []string, threshold int) error {
	utils.Debug("AddActivityUseCase: Starting execution for activity %s of type %s", name, activityType)

	if len(hiscoreNames) == 0 {
		hiscoreNames = []string{name}
		utils.Debug("AddActivityUseCase: No hiscore names provided, using activity name: %s", name)
	}

	activity := &domain.ActivityEntity{
		Name:         name,
		Type:         activityType,
		HiscoreNames: hiscoreNames,
		Threshold:    threshold,
	}

	err := uc.activityRepo.AddActivity(activity)
	if err != nil {
		utils.Error("AddActivityUseCase: Failed to add activity %s: %v", name, err)
		return err
	}

	utils.Info("AddActivityUseCase: Successfully added activity %s with %d hiscore names", name, len(hiscoreNames))
	return nil
}
