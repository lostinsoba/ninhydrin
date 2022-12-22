package controller

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
)

func (c *Controller) RegisterTask(ctx context.Context, task *model.Task) error {
	return c.storage.RegisterTask(ctx, task)
}

func (c *Controller) DeregisterTask(ctx context.Context, taskID string) error {
	return c.storage.DeregisterTask(ctx, taskID)
}

func (c *Controller) ReadTask(ctx context.Context, taskID string) (*model.Task, bool, error) {
	task, err := c.storage.ReadTask(ctx, taskID)
	switch err.(type) {
	case nil:
		return task, true, nil
	case model.ErrNotFound:
		return nil, false, nil
	default:
		return nil, false, err
	}
}

func (c *Controller) ListTaskIDs(ctx context.Context, poolIDs ...string) ([]string, error) {
	return c.storage.ListTaskIDs(ctx, poolIDs...)
}

func (c *Controller) CaptureTaskIDs(ctx context.Context, poolID string, limit int) ([]string, error) {
	return c.storage.CaptureTaskIDs(ctx, poolID, limit)
}

func (c *Controller) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
	return c.storage.RefreshTaskStatuses(ctx)
}
