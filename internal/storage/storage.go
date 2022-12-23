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
	ReadTag(ctx context.Context, tagID string) (tag string, err error)
	ListTagIDs(ctx context.Context) (tagIDs []string, err error)

	RegisterPool(ctx context.Context, pool *model.Pool) error
	DeregisterPool(ctx context.Context, poolID string) error
	ReadPool(ctx context.Context, poolID string) (pool *model.Pool, err error)
	UpdatePool(ctx context.Context, pool *model.Pool) error
	ListPoolIDs(ctx context.Context, tagIDs ...string) (poolIDs []string, err error)

	RegisterTask(ctx context.Context, task *model.Task) error
	DeregisterTask(ctx context.Context, taskID string) error
	ReadTask(ctx context.Context, taskID string) (task *model.Task, err error)
	ListTaskIDs(ctx context.Context, poolIDs ...string) (taskIDs []string, err error)

	CapturePoolTaskIDs(ctx context.Context, poolID string, limit int) (taskIDs []string, err error)
	ReleasePoolTaskIDs(ctx context.Context, poolID string, taskIDs []string, status model.TaskStatus) error
	RefreshPoolTaskIDs(ctx context.Context, poolID string) (tasksUpdated int64, err error)
}

func NewStorage(kind string, settings map[string]string) (Storage, error) {
	switch kind {
	case postgres.Kind:
		return postgres.NewPostgres(settings)
	default:
		return nil, fmt.Errorf("unknown storage kind: %s", kind)
	}
}
