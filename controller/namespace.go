package controller

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
)

func (ctrl *Controller) RegisterNamespace(ctx context.Context, namespace *model.Namespace) error {
	return ctrl.storage.RegisterNamespace(ctx, namespace)
}

func (ctrl *Controller) DeregisterNamespace(ctx context.Context, namespaceID string) error {
	return ctrl.storage.DeregisterNamespace(ctx, namespaceID)
}

func (ctrl *Controller) ReadNamespace(ctx context.Context, namespaceID string) (*model.Namespace, bool, error) {
	namespace, err := ctrl.storage.ReadNamespace(ctx, namespaceID)
	switch err.(type) {
	case nil:
		return namespace, true, nil
	case model.ErrNotFound:
		return nil, false, nil
	default:
		return nil, false, err
	}
}

func (ctrl *Controller) ListNamespaces(ctx context.Context) ([]*model.Namespace, error) {
	return ctrl.storage.ListNamespaces(ctx)
}
