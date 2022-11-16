package controller

import (
	"context"
)

func (c *Controller) RegisterTag(ctx context.Context, tagID string) error {
	return c.storage.RegisterTag(ctx, tagID)
}

func (c *Controller) ListTags(ctx context.Context) ([]string, error) {
	return c.storage.ListTags(ctx)
}
