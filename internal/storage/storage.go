package storage

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/storage/postgres"
	"lostinsoba/ninhydrin/internal/storage/redis"
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

	CaptureTasks(ctx context.Context, namespaceID string, limit int) (taskStates []*model.TaskState, err error)

	ReadTaskState(ctx context.Context, taskID string) (taskState *model.TaskState, err error)
	UpdateTaskState(ctx context.Context, taskState *model.TaskState) error
}

func NewStorage(kind string, settings model.Settings) (Storage, error) {
	switch kind {
	case postgres.Kind:
		return postgres.NewPostgres(settings)
	case redis.Kind:
		return redis.NewRedis(settings)
	case redis.KindSentinel:
		return redis.NewRedisSentinel(settings)
	default:
		return nil, fmt.Errorf("unknown storage kind: %s", kind)
	}
}
