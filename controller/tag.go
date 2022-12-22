package controller

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
)

func (c *Controller) RegisterTag(ctx context.Context, tagID string) error {
	return c.storage.RegisterTag(ctx, tagID)
}

func (c *Controller) ListTagIDs(ctx context.Context) ([]string, error) {
	return c.storage.ListTagIDs(ctx)
}

func (c *Controller) ReadTag(ctx context.Context, tagID string) (string, bool, error) {
	tag, err := c.storage.ReadTag(ctx, tagID)
	switch err.(type) {
	case nil:
		return tag, true, nil
	case model.ErrNotFound:
		return tag, false, nil
	default:
		return tag, false, err
	}
}

func (c *Controller) DeregisterTag(ctx context.Context, tagID string) error {
	isInUse, err := c.isTagInUse(ctx, tagID)
	if err != nil {
		return fmt.Errorf("can't check whether tag is in use: %w", err)
	}
	if isInUse {
		return fmt.Errorf("tag is being used")
	}
	return c.storage.DeregisterTag(ctx, tagID)
}

func (c *Controller) isTagInUse(ctx context.Context, tagID string) (bool, error) {
	poolIDs, err := c.storage.ListPoolIDs(ctx, tagID)
	if err != nil {
		return false, err
	}
	taskIDs, err := c.storage.ListTaskIDs(ctx, poolIDs...)
	if err != nil {
		return false, err
	}
	return len(poolIDs)+len(taskIDs) > 0, nil
}
