package postgres

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
)

func (s *Storage) RegisterNamespace(ctx context.Context, namespace *model.Namespace) error {
	var query = `insert into namespace (id) values ($1)`
	_, err := s.db.ExecContext(ctx, query, namespace.ID)
	if isAlreadyExist(err) {
		return model.ErrAlreadyExist{}
	}
	return err
}

func (s *Storage) DeregisterNamespace(ctx context.Context, namespaceID string) error {
	var query = `delete from namespace where id = $1`
	_, err := s.db.ExecContext(ctx, query, namespaceID)
	return err
}

func (s *Storage) CheckNamespaceExists(ctx context.Context, namespaceID string) (exists bool, err error) {
	var query = `select exists(select 1 from namespace where id = $1)`
	err = s.db.QueryRowContext(ctx, query, namespaceID).Scan(&exists)
	return
}

func (s *Storage) ReadNamespace(ctx context.Context, namespaceID string) (namespace *model.Namespace, err error) {
	var query = `select id from namespace where id = $1`
	var (
		id string
	)
	err = s.db.QueryRowContext(ctx, query, namespaceID).Scan(&id)
	if err != nil {
		if isNoRows(err) {
			return nil, model.ErrNotFound{}
		}
		return nil, err
	}
	return &model.Namespace{
		ID: id,
	}, nil
}

func (s *Storage) ListNamespaces(ctx context.Context) (namespaces []*model.Namespace, err error) {
	var query = `select id from namespace`
	rows, err := s.db.QueryContext(ctx, query)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}

	namespaces = make([]*model.Namespace, 0)
	for rows.Next() {
		var (
			id string
		)
		err = rows.Scan(&id)
		if err != nil {
			return
		}
		namespaces = append(namespaces, &model.Namespace{
			ID: id,
		})
	}
	return
}
