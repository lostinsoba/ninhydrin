package controller

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
)

func (c *Controller) RegisterWorker(ctx context.Context, worker *model.Worker) error {
	return c.storage.RegisterWorker(ctx, worker)
}

func (c *Controller) ListWorkers(ctx context.Context) ([]*model.Worker, error) {
	return c.storage.ListWorkers(ctx)
}

func (c *Controller) DeregisterWorker(ctx context.Context, workerID string) error {
	return c.storage.DeregisterWorker(ctx, workerID)
}
