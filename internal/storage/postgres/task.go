package postgres

import (
	"context"

	"github.com/lib/pq"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/util"
)

func (s *Storage) RegisterTask(ctx context.Context, task *model.Task) error {
	var query = `insert into task (id, namespace_id, retries_left, timeout, updated_at, status) values ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.ExecContext(ctx, query, task.ID, task.NamespaceID, task.RetriesLeft, task.Timeout, util.UnixEpoch(), string(task.Status))
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
	var query = `select id, namespace_id, retries_left, timeout, updated_at, status from task where id = $1`
	var (
		id          string
		namespaceID string
		retriesLeft int
		timeout     int64
		updatedAt   int64
		status      string
	)
	err = s.db.QueryRowContext(ctx, query, taskID).Scan(&id, &namespaceID, &retriesLeft, &timeout, &updatedAt, &status)
	if err != nil {
		if isNoRows(err) {
			return nil, model.ErrNotFound{}
		}
		return nil, err
	}
	return &model.Task{
		ID:          id,
		NamespaceID: namespaceID,
		RetriesLeft: retriesLeft,
		Timeout:     timeout,
		UpdatedAt:   updatedAt,
		Status:      model.TaskStatus(status),
	}, nil
}

func (s *Storage) CaptureTasks(ctx context.Context, namespaceID string, limit int) (tasks []*model.Task, err error) {
	var query = `update task set status = $1, retries_left = retries_left-1, updated_at = $2
		where id in (
			select id
			from task
			where namespace_id = $3 and status = any($4) and retries_left > 0
			limit $5
		) returning id, namespace_id, retries_left, timeout, updated_at, status`
	taskCaptureStatuses := model.GetTaskCaptureStatuses()
	rows, err := s.db.QueryContext(ctx, query,
		model.TaskStatusInProgress,
		util.UnixEpoch(),
		namespaceID,
		pq.Array(&taskCaptureStatuses),
		limit,
	)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	tasks = make([]*model.Task, 0)
	for rows.Next() {
		var (
			id          string
			nsID        string
			retriesLeft int
			timeout     int64
			updatedAt   int64
			status      string
		)
		err = rows.Scan(&id, &nsID, &retriesLeft, &timeout, &updatedAt, &status)
		if err != nil {
			return
		}
		tasks = append(tasks, &model.Task{
			ID:          id,
			NamespaceID: nsID,
			RetriesLeft: retriesLeft,
			Timeout:     timeout,
			UpdatedAt:   updatedAt,
			Status:      model.TaskStatus(status),
		})
	}
	return
}

func (s *Storage) ReleaseTasks(ctx context.Context, namespaceID string, taskIDs []string, status model.TaskStatus) error {
	var query = `update task set status = $1, updated_at = $2 where id = any($3)`
	_, err := s.db.ExecContext(ctx, query, status, util.UnixEpoch(), pq.Array(taskIDs))
	return err
}

func (s *Storage) RefreshTaskStatuses(ctx context.Context, namespaceID string) (tasksUpdated int64, err error) {
	var query = `
		update task set status = $1, retries_left = 0, updated_at = $2
		where namespace_id = $3 and status = $4 and $2 - updated_at > timeout`
	result, err := s.db.ExecContext(ctx, query, model.TaskStatusTimeout, util.UnixEpoch(), namespaceID, model.TaskStatusInProgress)
	if err != nil {
		return
	}
	tasksUpdated, err = result.RowsAffected()
	return
}

func (s *Storage) ListTasks(ctx context.Context, namespaceID string) (tasks []*model.Task, err error) {
	var query = `select id, namespace_id, retries_left, timeout, updated_at, status from task where namespace_id = $1`
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
			id          string
			nsID        string
			retriesLeft int
			timeout     int64
			updatedAt   int64
			status      string
		)
		err = rows.Scan(&id, &nsID, &retriesLeft, &timeout, &updatedAt, &status)
		if err != nil {
			return
		}
		tasks = append(tasks, &model.Task{
			ID:          id,
			NamespaceID: nsID,
			RetriesLeft: retriesLeft,
			Timeout:     timeout,
			UpdatedAt:   updatedAt,
			Status:      model.TaskStatus(status),
		})
	}
	return
}
