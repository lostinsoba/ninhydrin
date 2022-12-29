package controller

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
)

func (ctrl *Controller) RegisterPool(ctx context.Context, pool *model.Pool) error {
	return ctrl.storage.RegisterPool(ctx, pool)
}

func (ctrl *Controller) ReadPool(ctx context.Context, poolID string) (*model.Pool, bool, error) {
	pool, err := ctrl.storage.ReadPool(ctx, poolID)
	switch err.(type) {
	case nil:
		return pool, true, nil
	case model.ErrNotFound:
		return nil, false, nil
	default:
		return nil, false, err
	}
}

func (ctrl *Controller) ListPoolIDs(ctx context.Context) ([]string, error) {
	return ctrl.storage.ListPoolIDs(ctx)
}

func (ctrl *Controller) DeregisterPool(ctx context.Context, poolID string) error {
	isInUse, err := ctrl.isPoolInUse(ctx, poolID)
	if err != nil {
		return fmt.Errorf("can't check whether pool is in use: %w", err)
	}
	if isInUse {
		return fmt.Errorf("pool is being used")
	}
	return ctrl.storage.DeregisterPool(ctx, poolID)
}

func (ctrl *Controller) UpdatePool(ctx context.Context, pool *model.Pool) error {
	return ctrl.storage.UpdatePool(ctx, pool)
}

func (ctrl *Controller) isPoolInUse(ctx context.Context, poolID string) (bool, error) {
	taskIDs, err := ctrl.storage.ListTaskIDs(ctx, poolID)
	if err != nil {
		return false, err
	}
	return len(taskIDs) > 0, nil
}
