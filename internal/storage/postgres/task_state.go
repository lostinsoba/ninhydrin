package postgres

import (
	"context"

	"github.com/lib/pq"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/util"
)

func (s *Storage) CaptureTasks(ctx context.Context, namespaceID string, limit int) (taskStates []*model.TaskState, err error) {
	var query = `update task_state set status = $1, retries_left = retries_left-1, updated_at = $2
		where task_id in (
			select task_id
			from task_state
			where namespace_id = $3 and status = any($4) and retries_left > 0
			limit $5
		) returning task_id, namespace_id, retries_left, updated_at, status`
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
	taskStates = make([]*model.TaskState, 0)
	for rows.Next() {
		var (
			taskID      string
			retriesLeft int
			updatedAt   int64
			status      string
		)
		err = rows.Scan(&taskID, &retriesLeft, &updatedAt, &status)
		if err != nil {
			return
		}
		taskStates = append(taskStates, &model.TaskState{
			TaskID:      taskID,
			RetriesLeft: retriesLeft,
			UpdatedAt:   updatedAt,
			Status:      model.TaskStatus(status),
		})
	}
	return
}

func (s *Storage) UpdateTaskState(ctx context.Context, taskState *model.TaskState) error {
	var query = `update task_state set status = $1, retries_left = $2, updated_at = $2 where task_id = any($3)`
	_, err := s.db.ExecContext(ctx, query, taskState.Status, taskState.RetriesLeft, taskState.UpdatedAt, taskState.TaskID)
	return err
}

func (s *Storage) ReadTaskState(ctx context.Context, taskID string) (taskState *model.TaskState, err error) {
	var query = `select task_id, status, retries_left, updated_at from task_state where task_id = $1`
	var (
		id          string
		status      string
		retriesLeft int
		updatedAt   int64
	)
	err = s.db.QueryRowContext(ctx, query, taskID).Scan(&id, &status, &retriesLeft, &updatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, model.ErrNotFound{}
		}
		return nil, err
	}
	return &model.TaskState{
		TaskID:      id,
		RetriesLeft: retriesLeft,
		UpdatedAt:   updatedAt,
		Status:      model.TaskStatus(status),
	}, nil
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
