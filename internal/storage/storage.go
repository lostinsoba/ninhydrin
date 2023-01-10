package storage

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/storage/postgres"
)

type Storage interface {
	RegisterNamespace(ctx context.Context, namespace *model.Namespace) error
	DeregisterNamespace(ctx context.Context, namespaceID string) error
	CheckNamespaceExists(ctx context.Context, namespaceID string) (exists bool, err error)
	ReadNamespace(ctx context.Context, namespaceID string) (namespace *model.Namespace, err error)
	ListNamespaces(ctx context.Context) (namespaces []*model.Namespace, err error)

	RegisterTask(ctx context.Context, task *model.Task) error
	DeregisterTask(ctx context.Context, taskID string) error
	ReadTask(ctx context.Context, taskID string) (task *model.Task, err error)
	ListTasks(ctx context.Context, namespaceID string) (tasks []*model.Task, err error)

	CaptureTasks(ctx context.Context, namespaceID string, limit int) (tasks []*model.Task, err error)
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
