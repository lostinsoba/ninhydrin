package controller

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
)

const (
	defaultTaskStatus  = model.TaskStatusIdle
	defaultTaskRetries = 5
	defaultTaskTimeout = 360
)

func (ctrl *Controller) RegisterTask(ctx context.Context, task *model.Task) error {
	if task.Status == "" {
		task.Status = defaultTaskStatus
	} else {
		if !model.IsValidTaskStatus(task.Status) {
			return fmt.Errorf("invalid status")
		}
	}
	if task.RetriesLeft == 0 {
		task.RetriesLeft = defaultTaskRetries
	}
	if task.Timeout == 0 {
		task.Timeout = defaultTaskTimeout
	}
	exists, err := ctrl.storage.CheckNamespaceExists(ctx, task.NamespaceID)
	if err != nil {
		return fmt.Errorf("failed to check namespace %s existence: %w", task.NamespaceID, err)
	}
	if !exists {
		return fmt.Errorf("there's no such namespace: %s", task.NamespaceID)
	}
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

func (ctrl *Controller) ListTasks(ctx context.Context, namespaceID string) ([]*model.Task, error) {
	return ctrl.storage.ListTasks(ctx, namespaceID)
}

func (ctrl *Controller) CaptureTasks(ctx context.Context, namespaceID string, limit int) ([]*model.Task, error) {
	return ctrl.storage.CaptureTasks(ctx, namespaceID, limit)
}

func (ctrl *Controller) ReleaseTaskIDs(ctx context.Context, taskIDs []string, status model.TaskStatus) error {
	if !model.IsValidTaskStatus(status) {
		return fmt.Errorf("invalid status")
	}
	return ctrl.storage.ReleaseTaskIDs(ctx, taskIDs, status)
}

func (ctrl *Controller) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
	return ctrl.storage.RefreshTaskIDs(ctx)
}
