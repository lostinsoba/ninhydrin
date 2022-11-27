package storage

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/storage/postgres"
)

type Storage interface {
	RegisterTag(ctx context.Context, tagID string) error
	DeregisterTag(ctx context.Context, tagID string) error
	ListTags(ctx context.Context) (tagIDs []string, err error)

	RegisterPool(ctx context.Context, pool *model.Pool) error
	DeregisterPool(ctx context.Context, poolID string) error
	ListPools(ctx context.Context, tagIDs ...string) (pools []*model.Pool, err error)
	ListPoolIDsByTagIDs(ctx context.Context, tagIDs ...string) (poolIDs []string, err error)
	UpdatePool(ctx context.Context, pool *model.Pool) error

	RegisterTask(ctx context.Context, task *model.Task) error
	CaptureTasks(ctx context.Context, poolIDs []string, limit int) (tasks []*model.Task, err error)
	UpdateTaskStatus(ctx context.Context, taskID string, status model.TaskStatus) error
	ListCurrentTasks(ctx context.Context) (tasks []*model.Task, err error)
	ListTaskIDsByPoolIDs(ctx context.Context, poolIDs ...string) (taskIDs []string, err error)
	RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error)

	RegisterWorker(ctx context.Context, worker *model.Worker) error
	DeregisterWorker(ctx context.Context, workerID string) error
	ListWorkers(ctx context.Context) (workers []*model.Worker, err error)
	ListWorkerTagIDs(ctx context.Context, workerID string) (tagIDs []string, err error)
	ListWorkerIDsByTagIDs(ctx context.Context, tagIDs ...string) (workerIDs []string, err error)
}

func NewStorage(kind string, settings map[string]string) (Storage, error) {
	switch kind {
	case postgres.Kind:
		return postgres.NewPostgres(settings)
	default:
		return nil, fmt.Errorf("unknown storage kind: %s", kind)
	}
}
