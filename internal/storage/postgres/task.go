package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/util"
)

func (s *Storage) RegisterTask(ctx context.Context, task *model.Task) error {
	var query = `insert into task (id, pool_id, timeout, retries_left, updated_at, status) values ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.ExecContext(ctx, query, task.ID, task.PoolID, task.Timeout, task.RetriesLeft, util.UnixEpoch(), string(task.Status))
	return err
}

func (s *Storage) DeregisterTask(ctx context.Context, taskID string) error {
	var query = `delete from task where id = $1`
	_, err := s.db.ExecContext(ctx, query, taskID)
	return err
}

func (s *Storage) ReadTask(ctx context.Context, taskID string) (task *model.Task, err error) {
	var query = `select id, pool_id, timeout, retries_left, updated_at, status from task where id = $1`
	var (
		id          string
		poolID      string
		timeout     int64
		retriesLeft int
		updatedAt   int64
		status      string
	)
	err = s.db.QueryRowContext(ctx, query, taskID).Scan(&id, &poolID, &timeout, &retriesLeft, &updatedAt, &status)
	switch err {
	case nil:
		return &model.Task{
			ID:          id,
			PoolID:      poolID,
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

func (s *Storage) CaptureTaskIDs(ctx context.Context, poolID string, limit int) (taskIDs []string, err error) {
	var query = `update task set status = $1, retries_left = retries_left-1, updated_at = $2
		where id in (
			select id
			from task
			where pool_id = any($3) and status = any($4) and retries_left > 0
			limit $5
		) returning id`
	taskCaptureStatuses := model.GetTaskCaptureStatuses()
	rows, err := s.db.QueryContext(ctx, query,
		model.TaskStatusInProgress,
		util.UnixEpoch(),
		poolID,
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

func (s *Storage) UpdateTaskStatus(ctx context.Context, taskID string, status model.TaskStatus) error {
	var query = `update task 
					set status = $1, 
					    retries_left = retries_left - (case when $2 then 1 else 0 end), 
					    updated_at = $3 
					where id = $4`
	_, err := s.db.ExecContext(ctx, query, string(status), status == model.TaskStatusFailed, util.UnixEpoch(), taskID)
	return err
}

func (s *Storage) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
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

func (s *Storage) ListTaskIDs(ctx context.Context, poolIDs ...string) (taskIDs []string, err error) {
	var query = `select id from task where $1 or pool_id <@ $2`
	rows, err := s.db.QueryContext(ctx, query, len(poolIDs) == 0, pq.Array(poolIDs))
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
