package controller

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
)

func (c *Controller) RegisterPool(ctx context.Context, pool *model.Pool) error {
	return c.storage.RegisterPool(ctx, pool)
}

func (c *Controller) ReadPool(ctx context.Context, poolID string) (*model.Pool, bool, error) {
	pool, err := c.storage.ReadPool(ctx, poolID)
	switch err.(type) {
	case nil:
		return pool, true, nil
	case model.ErrNotFound:
		return nil, false, nil
	default:
		return nil, false, err
	}
}

func (c *Controller) ListPoolIDs(ctx context.Context) ([]string, error) {
	return c.storage.ListPoolIDs(ctx)
}

func (c *Controller) DeregisterPool(ctx context.Context, poolID string) error {
	isInUse, err := c.isPoolInUse(ctx, poolID)
	if err != nil {
		return fmt.Errorf("can't check whether pool is in use: %w", err)
	}
	if isInUse {
		return fmt.Errorf("pool is being used")
	}
	return c.storage.DeregisterPool(ctx, poolID)
}

func (c *Controller) UpdatePool(ctx context.Context, pool *model.Pool) error {
	return c.storage.UpdatePool(ctx, pool)
}

func (c *Controller) isPoolInUse(ctx context.Context, poolID string) (bool, error) {
	taskIDs, err := c.storage.ListTaskIDs(ctx, poolID)
	if err != nil {
		return false, err
	}
	return len(taskIDs) > 0, nil
}
