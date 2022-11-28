package controller

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
)

func (c *Controller) RegisterTask(ctx context.Context, task *model.Task) error {
	return c.storage.RegisterTask(ctx, task)
}

func (c *Controller) ReadTask(ctx context.Context, taskID string) (*model.Task, error) {
	return c.storage.ReadTask(ctx, taskID)
}

func (c *Controller) ListCurrentTasks(ctx context.Context) ([]*model.Task, error) {
	return c.storage.ListCurrentTasks(ctx)
}

func (c *Controller) CaptureTasks(ctx context.Context, workerID string, limit int) ([]*model.Task, error) {
	tagIDs, err := c.storage.ListWorkerTagIDs(ctx, workerID)
	if err != nil {
		return nil, err
	}
	if len(tagIDs) == 0 {
		return nil, nil
	}
	poolIDs, err := c.storage.ListPoolIDs(ctx, tagIDs...)
	if err != nil {
		return nil, err
	}
	if len(poolIDs) == 0 {
		return nil, nil
	}
	return c.storage.CaptureTasks(ctx, poolIDs, limit)
}

func (c *Controller) UpdateTaskStatus(ctx context.Context, taskID string, status model.TaskStatus) error {
	return c.storage.UpdateTaskStatus(ctx, taskID, status)
}

func (c *Controller) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
	return c.storage.RefreshTaskStatuses(ctx)
}
