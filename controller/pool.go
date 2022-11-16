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
