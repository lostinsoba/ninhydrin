package controller

import (
	"context"
	"fmt"
	"lostinsoba/ninhydrin/internal/util"
	"sync"
	"sync/atomic"

	"github.com/lostinsoba/chain"

	"lostinsoba/ninhydrin/internal/model"
)

const (
	defaultTaskStatus  = model.TaskStatusIdle
	defaultTaskRetries = 5
	defaultTaskTimeout = 360
)

func (ctrl *Controller) RegisterTask(ctx context.Context, task *model.Task) error {
	if task.Status == "" {
		task.Status = defaultTaskStatus
	} else {
		if !model.IsValidTaskStatus(task.Status) {
			return fmt.Errorf("invalid status")
		}
	}
	if task.RetriesLeft == 0 {
		task.RetriesLeft = defaultTaskRetries
	}
	if task.Timeout == 0 {
		task.Timeout = defaultTaskTimeout
	}
	task.UpdatedAt = util.UnixEpoch()
	exists, err := ctrl.storage.CheckNamespaceExists(ctx, task.NamespaceID)
	if err != nil {
		return fmt.Errorf("failed to check namespace %s existence: %w", task.NamespaceID, err)
	}
	if !exists {
		return fmt.Errorf("there's no such namespace: %s", task.NamespaceID)
	}
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

func (ctrl *Controller) ListTasks(ctx context.Context, namespaceID string) ([]*model.Task, error) {
	return ctrl.storage.ListTasks(ctx, namespaceID)
}

func (ctrl *Controller) CaptureTasks(ctx context.Context, namespaceID string, limit int) ([]*model.Task, error) {
	return ctrl.storage.CaptureTasks(ctx, namespaceID, limit)
}

func (ctrl *Controller) ReleaseTasks(ctx context.Context, namespaceID string, taskIDs []string, status model.TaskStatus) error {
	if !model.IsValidTaskStatus(status) {
		return fmt.Errorf("invalid status")
	}
	return ctrl.storage.ReleaseTasks(ctx, namespaceID, taskIDs, status)
}

func (ctrl *Controller) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
	namespaces, err := ctrl.ListNamespaces(ctx)
	if err != nil {
		return 0, err
	}

	errChan := make(chan error)
	doneChan := make(chan bool)

	var wg sync.WaitGroup
	var c chain.Chain

	c.SetStop(len(namespaces))
	c.SetStep(4)

	for c.Next() {
		left, right := c.Bounds()
		ns := namespaces[left:right]

		wg.Add(1)
		go func(wg *sync.WaitGroup, ctx context.Context, ns []*model.Namespace) {
			nsTaskUpdated, refreshErr := ctrl.refreshTaskStatuses(ctx, ns)
			if refreshErr != nil {
				errChan <- refreshErr
			} else {
				atomic.AddInt64(&tasksUpdated, nsTaskUpdated)
			}
			wg.Done()
		}(&wg, ctx, ns)
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

func (ctrl *Controller) refreshTaskStatuses(ctx context.Context, namespaces []*model.Namespace) (int64, error) {
	var tasksUpdatedTotal int64
	for _, namespace := range namespaces {
		tasksUpdated, err := ctrl.storage.RefreshTaskStatuses(ctx, namespace.ID)
		if err != nil {
			return 0, err
		}
		tasksUpdatedTotal += tasksUpdated
	}
	return tasksUpdatedTotal, nil
}
