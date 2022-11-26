package controller

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
)

func (c *Controller) RegisterPool(ctx context.Context, pool *model.Pool) error {
	return c.storage.RegisterPool(ctx, pool)
}

func (c *Controller) ListPools(ctx context.Context) ([]*model.Pool, error) {
	return c.storage.ListPools(ctx)
}

func (c *Controller) DeregisterPool(ctx context.Context, poolID string) error {
	return c.storage.DeregisterPool(ctx, poolID)
}

func (c *Controller) UpdatePool(ctx context.Context, pool *model.Pool) error {
	return c.storage.UpdatePool(ctx, pool)
}
