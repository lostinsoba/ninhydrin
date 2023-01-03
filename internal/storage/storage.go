package storage

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/storage/postgres"
)

type Storage interface {
	RegisterTask(ctx context.Context, task *model.Task) error
	DeregisterTask(ctx context.Context, taskID string) error
	ReadTask(ctx context.Context, taskID string) (task *model.Task, err error)
	ListTaskIDs(ctx context.Context) (taskIDs []string, err error)

	CaptureTaskIDs(ctx context.Context, limit int) (taskIDs []string, err error)
	ReleaseTaskIDs(ctx context.Context, taskIDs []string, status model.TaskStatus) error
	RefreshTaskIDs(ctx context.Context) (tasksUpdated int64, err error)
}

func NewStorage(kind string, settings map[string]string) (Storage, error) {
	switch kind {
	case postgres.Kind:
		return postgres.NewPostgres(settings)
	default:
		return nil, fmt.Errorf("unknown storage kind: %s", kind)
	}
}
