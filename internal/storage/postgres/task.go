package postgres

import (
	"context"

	"lostinsoba/ninhydrin/internal/model"
)

func (s *Storage) RegisterTask(ctx context.Context, task *model.Task) error {
	var query = `insert into task (id, namespace_id, timeout) values ($1, $2, $3)`
	_, err := s.db.ExecContext(ctx, query, task.ID, task.NamespaceID, task.Timeout)
	if isAlreadyExist(err) {
		return model.ErrAlreadyExist{}
	}
	return err
}

func (s *Storage) DeregisterTask(ctx context.Context, taskID string) error {
	var query = `delete from task where id = $1`
	_, err := s.db.ExecContext(ctx, query, taskID)
	return err
}

func (s *Storage) ReadTask(ctx context.Context, taskID string) (task *model.Task, err error) {
	var query = `select id, namespace_id, timeout from task where id = $1`
	var (
		id          string
		namespaceID string
		timeout     int64
	)
	err = s.db.QueryRowContext(ctx, query, taskID).Scan(&id, &namespaceID, &timeout)
	if err != nil {
		if isNoRows(err) {
			return nil, model.ErrNotFound{}
		}
		return nil, err
	}
	return &model.Task{
		ID:          id,
		NamespaceID: namespaceID,
		Timeout:     timeout,
	}, nil
}

func (s *Storage) ListTasks(ctx context.Context, namespaceID string) (tasks []*model.Task, err error) {
	var query = `select id, namespace_id, timeout from task where namespace_id = $1`
	rows, err := s.db.QueryContext(ctx, query, namespaceID)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	tasks = make([]*model.Task, 0)
	for rows.Next() {
		var (
			id      string
			nsID    string
			timeout int64
		)
		err = rows.Scan(&id, &nsID, &timeout)
		if err != nil {
			return
		}
		tasks = append(tasks, &model.Task{
			ID:          id,
			NamespaceID: nsID,
			Timeout:     timeout,
		})
	}
	return
}
