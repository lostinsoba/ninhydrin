package controller

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
)

func (c *Controller) RegisterWorker(ctx context.Context, worker *model.Worker) error {
	return c.storage.RegisterWorker(ctx, worker)
}

func (c *Controller) ReadWorker(ctx context.Context, workerID string) (*model.Worker, error) {
	return c.storage.ReadWorker(ctx, workerID)
}

func (c *Controller) ListWorkerIDs(ctx context.Context) ([]string, error) {
	return c.storage.ListWorkerIDs(ctx)
}

func (c *Controller) DeregisterWorker(ctx context.Context, workerID string) error {
	return c.storage.DeregisterWorker(ctx, workerID)
}
