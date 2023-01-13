package redis

import (
	"context"
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
)

func (s *Storage) RegisterNamespace(ctx context.Context, namespace *model.Namespace) error {
	return s.client.SAdd(ctx, namespaceKey(namespace.ID), namespace.ID).Err()
}

func (s *Storage) DeregisterNamespace(ctx context.Context, namespaceID string) error {
	return s.client.SRem(ctx, namespaceKey(namespaceID), namespaceID).Err()
}

func (s *Storage) CheckNamespaceExists(ctx context.Context, namespaceID string) (exists bool, err error) {
	res, err := s.client.Exists(ctx, namespaceKey(namespaceID)).Result()
	exists = res == 1
	return
}

func (s *Storage) ReadNamespace(ctx context.Context, namespaceID string) (namespace *model.Namespace, err error) {
	res := s.client.Get(ctx, namespaceKey(namespaceID))
	if res.Err() != nil {
		return nil, res.Err()
	}
	return &model.Namespace{ID: res.String()}, err
}

func (s *Storage) ListNamespaces(ctx context.Context) (namespaces []*model.Namespace, err error) {
	res := s.client.SMembers(ctx, namespacePrefix)
	if res.Err() != nil {
		return nil, res.Err()
	}
	namespaces = make([]*model.Namespace, 0, len(res.Val()))
	for _, val := range res.Val() {
		namespaces = append(namespaces, &model.Namespace{ID: val})
	}
	return
}

const (
	namespacePrefix = "namespace"
)

func namespaceKey(namespaceID string) string {
	return fmt.Sprintf("%s:%s", namespacePrefix, namespaceID)
}
