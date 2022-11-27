package controller

import (
	"context"
	"fmt"
)

func (c *Controller) RegisterTag(ctx context.Context, tagID string) error {
	return c.storage.RegisterTag(ctx, tagID)
}

func (c *Controller) ListTags(ctx context.Context) ([]string, error) {
	return c.storage.ListTags(ctx)
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
	poolIDs, err := c.storage.ListPoolIDsByTagIDs(ctx, tagID)
	if err != nil {
		return false, err
	}

	taskIDs, err := c.storage.ListTaskIDsByPoolIDs(ctx, poolIDs...)
	if err != nil {
		return false, err
	}

	workerIDs, err := c.storage.ListWorkerIDsByTagIDs(ctx, tagID)
	if err != nil {
		return false, err
	}

	return len(poolIDs)+len(taskIDs)+len(workerIDs) > 0, nil
}
