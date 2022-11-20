package postgres

import (
	"context"

	"github.com/lib/pq"

	"lostinsoba/ninhydrin/internal/model"
)

func (s *Storage) RegisterTask(ctx context.Context, task *model.Task) error {
	var query = `insert into task (id, pool_id, timeout, retries_left, status) values ($1, $2, $3, $4, $5)`
	_, err := s.db.ExecContext(ctx, query, task.ID, task.PoolID, task.Timeout, task.RetriesLeft, string(task.Status))
	return err
}

func (s *Storage) CaptureTasks(ctx context.Context, poolIDs []string, limit int) (tasks []*model.Task, err error) {
	var query = `update task set status = $1
		where id in (
			select id
			from task
			where pool_id = any($2) and status = any($3) and retries_left > 0
			limit $4
		) returning id, pool_id, timeout, retries_left, updated_at, status`
	taskCaptureStatuses := model.GetTaskCaptureStatuses()
	rows, err := s.db.QueryContext(ctx, query,
		model.TaskStatusInProgress,
		pq.Array(poolIDs),
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
			poolID      string
			timeout     int64
			retriesLeft int
			createdAt   int64
			updatedAt   int64
			status      string
		)
		err = rows.Scan(&id, &poolID, &timeout, &retriesLeft, &createdAt, &updatedAt, &status)
		if err != nil {
			return
		}
		tasks = append(tasks, &model.Task{
			ID:          id,
			PoolID:      poolID,
			Timeout:     timeout,
			RetriesLeft: retriesLeft,
			UpdatedAt:   updatedAt,
			Status:      model.TaskStatus(status),
		})
	}
	return
}

func (s *Storage) UpdateTaskStatus(ctx context.Context, taskID string, status model.TaskStatus) error {
	var query = `update task 
					set status = $1, 
					    retries_left = retries_left - (case when $2 then 1 else 0 end), 
					    updated_at = now() at time zone 'utc' 
					where id = $3`
	_, err := s.db.ExecContext(ctx, query, string(status), status == model.TaskStatusFailed, taskID)
	return err
}

func (s *Storage) ListCurrentTasks(ctx context.Context) (tasks []*model.Task, err error) {
	var query = `select id, pool_id, timeout, retries_left, updated_at, status 
					from task where status <> $1 order by updated_at desc`
	rows, err := s.db.QueryContext(ctx, query, model.TaskStatusDone)
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
			poolID      string
			timeout     int64
			retriesLeft int
			createdAt   int64
			updatedAt   int64
			status      string
		)
		err = rows.Scan(&id, &poolID, &timeout, &retriesLeft, &createdAt, &updatedAt, &status)
		if err != nil {
			return
		}
		tasks = append(tasks, &model.Task{
			ID:          id,
			PoolID:      poolID,
			Timeout:     timeout,
			RetriesLeft: retriesLeft,
			UpdatedAt:   updatedAt,
			Status:      model.TaskStatus(status),
		})
	}
	return
}

func (s *Storage) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
	var query = `
		update task set status = $1, retries_left = 0, updated_at = now() at time zone 'utc'
		where status = $2 and now() at time zone 'utc' - updated_at > timeout`
	result, err := s.db.ExecContext(ctx, query, model.TaskStatusTimeout, model.TaskStatusInProgress)
	if err != nil {
		return
	}
	tasksUpdated, err = result.RowsAffected()
	return
}