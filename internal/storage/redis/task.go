package redis

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/util"
)

func (s *Storage) RegisterTask(ctx context.Context, task *model.Task) error {
	pipe := s.client.TxPipeline()
	pipe.SAdd(ctx, namespaceTaskKey(task.NamespaceID), task.ID)
	pipe.SAdd(ctx, taskStatusKey(task.Status), task.ID)
	pipe.Set(ctx, taskTimeoutKey(task.ID), task.Timeout, -1)
	pipe.Set(ctx, taskRetriesLeftKey(task.ID), task.RetriesLeft, -1)
	pipe.Set(ctx, taskUpdatedAtKey(task.ID), util.UnixEpoch(), -1)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *Storage) DeregisterTask(ctx context.Context, taskID string) error {
	return nil
}

func (s *Storage) ReadTask(ctx context.Context, taskID string) (task *model.Task, err error) {
	return nil, err
}

func (s *Storage) CaptureTasks(ctx context.Context, namespaceID string, limit int) (tasks []*model.Task, err error) {
	return nil, err
}

func (s *Storage) ReleaseTasks(ctx context.Context, taskIDs []string, status model.TaskStatus) error {
	return nil
}

func (s *Storage) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
	return 0, nil
}

func (s *Storage) ListTasks(ctx context.Context, namespaceID string) (tasks []*model.Task, err error) {
	return nil, nil
}

const (
	namespaceTaskPrefix = "namespace-task"
	taskTimeoutPrefix   = "task-timeout"
	taskRetriesLeft     = "task-retries-left"
	taskUpdatedAt       = "task-updated-at"
	taskStatusPrefix    = "task-status"
)

func namespaceTaskKey(namespaceID string) string {
	return fmt.Sprintf("%s:%s", namespaceTaskPrefix, namespaceID)
}

func taskStatusKey(status model.TaskStatus) string {
	return fmt.Sprintf("%s:%s", taskStatusPrefix, status)
}

func taskTimeoutKey(taskID string) string {
	return fmt.Sprintf("%s:%s", taskTimeoutPrefix, taskID)
}

func taskRetriesLeftKey(taskID string) string {
	return fmt.Sprintf("%s:%s", taskRetriesLeft, taskID)
}

func taskUpdatedAtKey(taskID string) string {
	return fmt.Sprintf("%s:%s", taskUpdatedAt, taskID)
}
