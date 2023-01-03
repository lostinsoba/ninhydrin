package controller

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
)

func (ctrl *Controller) RegisterTask(ctx context.Context, task *model.Task) error {
	return ctrl.storage.RegisterTask(ctx, task)
}

func (ctrl *Controller) DeregisterTask(ctx context.Context, taskID string) error {
	return ctrl.storage.DeregisterTask(ctx, taskID)
}

func (ctrl *Controller) ReadTask(ctx context.Context, taskID string) (*model.Task, bool, error) {
	task, err := ctrl.storage.ReadTask(ctx, taskID)
	switch err.(type) {
	case nil:
		return task, true, nil
	case model.ErrNotFound:
		return nil, false, nil
	default:
		return nil, false, err
	}
}

func (ctrl *Controller) ListTaskIDs(ctx context.Context) ([]string, error) {
	return ctrl.storage.ListTaskIDs(ctx)
}

func (ctrl *Controller) CaptureTaskIDs(ctx context.Context, limit int) ([]string, error) {
	return ctrl.storage.CaptureTaskIDs(ctx, limit)
}

func (ctrl *Controller) ReleaseTaskIDs(ctx context.Context, taskIDs []string, status string) error {
	return ctrl.storage.ReleaseTaskIDs(ctx, taskIDs, model.TaskStatus(status))
}

func (ctrl *Controller) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
	return ctrl.storage.RefreshTaskIDs(ctx)
}
