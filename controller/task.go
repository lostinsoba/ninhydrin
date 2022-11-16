package controller

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
)

func (c *Controller) RegisterTask(ctx context.Context, task *model.Task) error {
	return c.storage.RegisterTask(ctx, task)
}

func (c *Controller) ListCurrentTasks(ctx context.Context) ([]*model.Task, error) {
	return c.storage.ListCurrentTasks(ctx)
}

func (c *Controller) CaptureTasks(ctx context.Context, workerID string, limit int) ([]*model.Task, error) {
	tagIDs, err := c.storage.ListWorkerTagIDs(ctx, workerID)
	if err != nil {
		return nil, err
	}
	pools, err := c.storage.ListPools(ctx, tagIDs...)
	if err != nil {
		return nil, err
	}
	poolIDs := make([]string, 0, len(pools))
	for _, pool := range pools {
		poolIDs = append(poolIDs, pool.ID)
	}
	return c.storage.CaptureTasks(ctx, poolIDs, limit)
}

func (c *Controller) UpdateTaskStatus(ctx context.Context, taskID string, status model.TaskStatus) error {
	return c.storage.UpdateTaskStatus(ctx, taskID, status)
}

func (c *Controller) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
	return c.storage.RefreshTaskStatuses(ctx)
}
