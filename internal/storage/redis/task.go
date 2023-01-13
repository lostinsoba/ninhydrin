package redis

import (
	"context"
	"fmt"
	r "github.com/go-redis/redis/v9"

	"lostinsoba/ninhydrin/internal/model"
)

const (
	TaskScoreTimeout = iota
	TaskScoreFailed
	TaskScoreIdle
	TaskScoreInProgress
	TaskScoreDone
)

func toScore(status model.TaskStatus) float64 {
	scoreMap := map[model.TaskStatus]int{
		model.TaskStatusTimeout:    TaskScoreTimeout,
		model.TaskStatusFailed:     TaskScoreFailed,
		model.TaskStatusIdle:       TaskScoreIdle,
		model.TaskStatusInProgress: TaskScoreInProgress,
		model.TaskStatusDone:       TaskScoreDone,
	}
	return float64(scoreMap[status])
}

func (s *Storage) RegisterTask(ctx context.Context, task *model.Task) error {
	data, err := encode(task)
	if err != nil {
		return err
	}
	pipe := s.client.TxPipeline()
	pipe.Set(ctx, taskKey(task.ID), data, -1)
	pipe.ZAddNX(ctx, namespaceTaskKey(task.NamespaceID), r.Z{
		Member: task.ID,
		Score:  toScore(task.Status),
	})
	_, err = pipe.Exec(ctx)
	return err
}

func (s *Storage) DeregisterTask(ctx context.Context, taskID string) error {
	return nil
}

func (s *Storage) ReadTask(ctx context.Context, taskID string) (task *model.Task, err error) {
	cmd := s.client.Get(ctx, taskKey(taskID))
	if cmd.Err() != nil {
		return nil, err
	}
	data, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}
	err = decode(data, &task)
	return
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
	taskPrefix          = "task"
	namespaceTaskPrefix = "namespace-task"
)

func taskKey(taskID string) string {
	return fmt.Sprintf("%s:%s", taskPrefix, taskID)
}

func namespaceTaskKey(namespaceID string) string {
	return fmt.Sprintf("%s:%s", namespaceTaskPrefix, namespaceID)
}
