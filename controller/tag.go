package controller

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
)

func (ctrl *Controller) RegisterTag(ctx context.Context, tagID string) error {
	return ctrl.storage.RegisterTag(ctx, tagID)
}

func (ctrl *Controller) ListTagIDs(ctx context.Context) ([]string, error) {
	return ctrl.storage.ListTagIDs(ctx)
}

func (ctrl *Controller) ReadTag(ctx context.Context, tagID string) (string, bool, error) {
	tag, err := ctrl.storage.ReadTag(ctx, tagID)
	switch err.(type) {
	case nil:
		return tag, true, nil
	case model.ErrNotFound:
		return tag, false, nil
	default:
		return tag, false, err
	}
}

func (ctrl *Controller) DeregisterTag(ctx context.Context, tagID string) error {
	isInUse, err := ctrl.isTagInUse(ctx, tagID)
	if err != nil {
		return fmt.Errorf("can't check whether tag is in use: %w", err)
	}
	if isInUse {
		return fmt.Errorf("tag is being used")
	}
	return ctrl.storage.DeregisterTag(ctx, tagID)
}

func (ctrl *Controller) isTagInUse(ctx context.Context, tagID string) (bool, error) {
	poolIDs, err := ctrl.storage.ListPoolIDs(ctx, tagID)
	if err != nil {
		return false, err
	}
	return len(poolIDs) > 0, nil
}
