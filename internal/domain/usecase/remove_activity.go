package usecase

import (
	"misclicked-events/internal/data/repository"
	"misclicked-events/internal/utils"
)

type RemoveActivityUseCase struct {
	activityRepo repository.ActivityRepository
}

func NewRemoveActivityUseCase(activityRepo repository.ActivityRepository) *RemoveActivityUseCase {
	return &RemoveActivityUseCase{
		activityRepo: activityRepo,
	}
}

func (uc *RemoveActivityUseCase) Execute(name string) error {
	utils.Debug("RemoveActivityUseCase: Starting execution for activity %s", name)

	err := uc.activityRepo.RemoveActivity(name)
	if err != nil {
		utils.Error("RemoveActivityUseCase: Failed to remove activity %s: %v", name, err)
		return err
	}

	utils.Info("RemoveActivityUseCase: Successfully removed activity %s", name)
	return nil
}
