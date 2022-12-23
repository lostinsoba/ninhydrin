package controller

import (
	"context"
	"sync"

	"lostinsoba/ninhydrin/internal/model"
)

func (c *Controller) RegisterTask(ctx context.Context, task *model.Task) error {
	return c.storage.RegisterTask(ctx, task)
}

func (c *Controller) DeregisterTask(ctx context.Context, taskID string) error {
	return c.storage.DeregisterTask(ctx, taskID)
}

func (c *Controller) ReadTask(ctx context.Context, taskID string) (*model.Task, bool, error) {
	task, err := c.storage.ReadTask(ctx, taskID)
	switch err.(type) {
	case nil:
		return task, true, nil
	case model.ErrNotFound:
		return nil, false, nil
	default:
		return nil, false, err
	}
}

func (c *Controller) ListTaskIDs(ctx context.Context, poolIDs ...string) ([]string, error) {
	return c.storage.ListTaskIDs(ctx, poolIDs...)
}

func (c *Controller) CapturePoolTaskIDs(ctx context.Context, poolID string, limit int) ([]string, error) {
	return c.storage.CapturePoolTaskIDs(ctx, poolID, limit)
}

func (c *Controller) ReleasePoolTaskIDs(ctx context.Context, poolID string, taskIDs []string, status string) error {
	return c.storage.ReleasePoolTaskIDs(ctx, poolID, taskIDs, model.TaskStatus(status))
}

func (c *Controller) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
	poolIDs, err := c.storage.ListPoolIDs(ctx)
	if err != nil {
		return
	}

	totalCount := len(poolIDs)
	if totalCount == 0 {
		return
	}

	const maxBatchSize = 5

	for start := 0; start < totalCount; start += maxBatchSize {
		end := start + maxBatchSize
		if end > totalCount {
			end = totalCount
		}

		tasksUpdatedIncr, err := c.refreshPoolTaskStatuses(ctx, poolIDs[start:end])
		if err != nil {
			return
		}
		tasksUpdated += tasksUpdatedIncr
	}

	return
}

func (c *Controller) refreshPoolTaskStatuses(ctx context.Context, poolIDs []string) (tasksUpdated int64, err error) {
	errChan := make(chan error)
	doneChan := make(chan bool)

	var wg sync.WaitGroup
	for _, poolID := range poolIDs {
		wg.Add(1)
		go func(wg *sync.WaitGroup, poolID string) {
			tasksUpdatedIncr, err := c.storage.RefreshPoolTaskIDs(ctx, poolID)
			if err != nil {
				errChan <- err
			} else {
				tasksUpdated += tasksUpdatedIncr
			}
			wg.Done()
		}(&wg, poolID)
	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case <-doneChan:
		break
	case err = <-errChan:
		close(errChan)
	}

	return
}
