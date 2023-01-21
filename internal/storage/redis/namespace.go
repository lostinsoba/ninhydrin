package redis

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/storage/redis/schema"
)

func (s *Storage) RegisterNamespace(ctx context.Context, namespace *model.Namespace) error {
	res, err := s.client.SAdd(ctx, schema.NamespaceKey(), namespace.ID).Result()
	if err != nil {
		return err
	}
	if res == 0 {
		return model.ErrAlreadyExist{}
	}
	return nil
}

func (s *Storage) DeregisterNamespace(ctx context.Context, namespaceID string) error {
	// todo: try check rows affected in postgres and on conflict do nothing
	res, err := s.client.SRem(ctx, schema.NamespaceKey(), namespaceID).Result()
	if err != nil {
		return err
	}
	if res == 0 {
		return model.ErrNotFound{}
	}
	return nil
}

func (s *Storage) CheckNamespaceExists(ctx context.Context, namespaceID string) (exists bool, err error) {
	return s.checkNamespaceExists(ctx, namespaceID)
}

func (s *Storage) ReadNamespace(ctx context.Context, namespaceID string) (namespace *model.Namespace, err error) {
	exists, err := s.checkNamespaceExists(ctx, namespaceID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, model.ErrNotFound{}
	}
	return &model.Namespace{ID: namespaceID}, err
}

func (s *Storage) checkNamespaceExists(ctx context.Context, namespaceID string) (bool, error) {
	return s.client.SIsMember(ctx, schema.NamespaceKey(), namespaceID).Result()
}

func (s *Storage) ListNamespaces(ctx context.Context) (namespaces []*model.Namespace, err error) {
	cmd := s.client.SMembers(ctx, schema.NamespaceKey())
	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	namespaces = make([]*model.Namespace, 0, len(res))
	for _, val := range res {
		namespaces = append(namespaces, &model.Namespace{ID: val})
	}
	return
}
