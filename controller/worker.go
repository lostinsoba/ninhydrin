package controller

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
)

func (c *Controller) RegisterWorker(ctx context.Context, worker *model.Worker) error {
	return c.storage.RegisterWorker(ctx, worker)
}

func (c *Controller) ReadWorker(ctx context.Context, workerID string) (*model.Worker, bool, error) {
	worker, err := c.storage.ReadWorker(ctx, workerID)
	switch err.(type) {
	case nil:
		return worker, true, nil
	case model.ErrNotFound:
		return nil, false, nil
	default:
		return nil, false, err
	}
}

func (c *Controller) ListWorkerIDs(ctx context.Context) ([]string, error) {
	return c.storage.ListWorkerIDs(ctx)
}

func (c *Controller) DeregisterWorker(ctx context.Context, workerID string) error {
	return c.storage.DeregisterWorker(ctx, workerID)
}
