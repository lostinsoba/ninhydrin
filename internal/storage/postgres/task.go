package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/util"
)

func (s *Storage) RegisterTask(ctx context.Context, task *model.Task) error {
	var query = `insert into task (id, timeout, retries_left, updated_at, status) values ($1, $2, $3, $4, $5)`
	_, err := s.db.ExecContext(ctx, query, task.ID, task.Timeout, task.RetriesLeft, util.UnixEpoch(), string(task.Status))
	return err
}

func (s *Storage) DeregisterTask(ctx context.Context, taskID string) error {
	var query = `delete from task where id = $1`
	_, err := s.db.ExecContext(ctx, query, taskID)
	return err
}

func (s *Storage) ReadTask(ctx context.Context, taskID string) (task *model.Task, err error) {
	var query = `select id, timeout, retries_left, updated_at, status from task where id = $1`
	var (
		id          string
		timeout     int64
		retriesLeft int
		updatedAt   int64
		status      string
	)
	err = s.db.QueryRowContext(ctx, query, taskID).Scan(&id, &timeout, &retriesLeft, &updatedAt, &status)
	switch err {
	case nil:
		return &model.Task{
			ID:          id,
			Timeout:     timeout,
			RetriesLeft: retriesLeft,
			UpdatedAt:   updatedAt,
			Status:      model.TaskStatus(status),
		}, nil
	case sql.ErrNoRows:
		return nil, model.ErrNotFound{}
	default:
		return nil, err
	}
}

func (s *Storage) CaptureTaskIDs(ctx context.Context, limit int) (taskIDs []string, err error) {
	var query = `update task set status = $1, retries_left = retries_left-1, updated_at = $2
		where id in (
			select id
			from task
			where status = any($3) and retries_left > 0
			limit $4
		) returning id`
	taskCaptureStatuses := model.GetTaskCaptureStatuses()
	rows, err := s.db.QueryContext(ctx, query,
		model.TaskStatusInProgress,
		util.UnixEpoch(),
		pq.Array(&taskCaptureStatuses),
		limit,
	)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	taskIDs = make([]string, 0)
	for rows.Next() {
		var (
			id string
		)
		err = rows.Scan(&id)
		if err != nil {
			return
		}
		taskIDs = append(taskIDs, id)
	}
	return
}

func (s *Storage) ReleaseTaskIDs(ctx context.Context, taskIDs []string, status model.TaskStatus) error {
	var query = `update task set status = $1, updated_at = $2 where id = any($3)`
	_, err := s.db.ExecContext(ctx, query, status, util.UnixEpoch(), pq.Array(taskIDs))
	return err
}

func (s *Storage) RefreshTaskIDs(ctx context.Context) (tasksUpdated int64, err error) {
	var query = `
		update task set status = $1, retries_left = 0, updated_at = $2
		where status = $3 and $2 - updated_at > timeout`
	result, err := s.db.ExecContext(ctx, query, model.TaskStatusTimeout, util.UnixEpoch(), model.TaskStatusInProgress)
	if err != nil {
		return
	}
	tasksUpdated, err = result.RowsAffected()
	return
}

func (s *Storage) ListTaskIDs(ctx context.Context) (taskIDs []string, err error) {
	var query = `select id from task`
	rows, err := s.db.QueryContext(ctx, query)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	taskIDs = make([]string, 0)
	for rows.Next() {
		var (
			id string
		)
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		taskIDs = append(taskIDs, id)
	}
	return taskIDs, nil
}
