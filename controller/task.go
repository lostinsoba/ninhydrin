package controller

import (
	"context"
	"sync"

	"github.com/lostinsoba/chain"

	"lostinsoba/ninhydrin/internal/model"
)

func (ctrl *Controller) RegisterTask(ctx context.Context, task *model.Task) error {
	return ctrl.storage.RegisterTask(ctx, task)
}

func (ctrl *Controller) DeregisterTask(ctx context.Context, taskID string) error {
	return ctrl.storage.DeregisterTask(ctx, taskID)
}

func (ctrl *Controller) ReadTask(ctx context.Context, taskID string) (*model.Task, bool, error) {
	task, err := ctrl.storage.ReadTask(ctx, taskID)
	switch err.(type) {
	case nil:
		return task, true, nil
	case model.ErrNotFound:
		return nil, false, nil
	default:
		return nil, false, err
	}
}

func (ctrl *Controller) ListTaskIDs(ctx context.Context, poolIDs ...string) ([]string, error) {
	return ctrl.storage.ListTaskIDs(ctx, poolIDs...)
}

func (ctrl *Controller) CapturePoolTaskIDs(ctx context.Context, poolID string, limit int) ([]string, error) {
	return ctrl.storage.CapturePoolTaskIDs(ctx, poolID, limit)
}

func (ctrl *Controller) ReleasePoolTaskIDs(ctx context.Context, poolID string, taskIDs []string, status string) error {
	return ctrl.storage.ReleasePoolTaskIDs(ctx, poolID, taskIDs, model.TaskStatus(status))
}

func (ctrl *Controller) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
	poolIDs, err := ctrl.storage.ListPoolIDs(ctx)
	if err != nil {
		return
	}

	var c chain.Chain
	c.SetStop(len(poolIDs))
	c.SetStep(5)

	for c.Next() {
		start, end := c.Bounds()
		tasksUpdatedIncr, err := ctrl.refreshPoolTaskStatuses(ctx, poolIDs[start:end])
		if err != nil {
			return
		}
		tasksUpdated += tasksUpdatedIncr
	}

	return
}

func (ctrl *Controller) refreshPoolTaskStatuses(ctx context.Context, poolIDs []string) (tasksUpdated int64, err error) {
	errChan := make(chan error)
	doneChan := make(chan bool)

	var wg sync.WaitGroup
	for _, poolID := range poolIDs {
		wg.Add(1)
		go func(wg *sync.WaitGroup, poolID string) {
			tasksUpdatedIncr, err := ctrl.storage.RefreshPoolTaskIDs(ctx, poolID)
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
