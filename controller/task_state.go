package controller

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/util"
)

func (ctrl *Controller) CaptureTasks(ctx context.Context, namespaceID string, limit int) ([]*model.TaskState, error) {
	return ctrl.storage.CaptureTasks(ctx, namespaceID, limit)
}

func (ctrl *Controller) ReadTaskState(ctx context.Context, taskID string) (*model.TaskState, error) {
	return ctrl.storage.ReadTaskState(ctx, taskID)
}

func (ctrl *Controller) UpdateTaskState(ctx context.Context, taskState *model.TaskState) error {
	if taskState.Status == "" {
		taskState.Status = defaultTaskStatus
	} else {
		if !model.IsValidTaskStatus(taskState.Status) {
			return fmt.Errorf("invalid status")
		}
	}
	if taskState.RetriesLeft == 0 {
		taskState.RetriesLeft = defaultTaskRetries
	}
	taskState.UpdatedAt = util.UnixEpoch()

	return ctrl.storage.UpdateTaskState(ctx, taskState)
}
