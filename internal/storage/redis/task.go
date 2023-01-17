package redis

import (
	"context"
	"fmt"
	"strconv"

	r "github.com/go-redis/redis/v9"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/util"
)

const (
	TaskScoreTimeout = iota
	TaskScoreFailed
	TaskScoreIdle
	TaskScoreInProgress
	TaskScoreDone
)

func toScore(status model.TaskStatus) float64 {
	scoreMap := map[model.TaskStatus]int{
		model.TaskStatusTimeout:    TaskScoreTimeout,
		model.TaskStatusFailed:     TaskScoreFailed,
		model.TaskStatusIdle:       TaskScoreIdle,
		model.TaskStatusInProgress: TaskScoreInProgress,
		model.TaskStatusDone:       TaskScoreDone,
	}
	return float64(scoreMap[status])
}

func fromScore(score float64) model.TaskStatus {
	statuses := []model.TaskStatus{
		model.TaskStatusTimeout,
		model.TaskStatusFailed,
		model.TaskStatusIdle,
		model.TaskStatusInProgress,
		model.TaskStatusDone,
	}
	ind := int(score)
	return statuses[ind]
}

const (
	sortedSetNamespaceTask = "namespace-task"
)

func namespaceTaskKey(namespaceID string) string {
	return fmt.Sprintf("%s:%s", sortedSetNamespaceTask, namespaceID)
}

func (s *Storage) RegisterTask(ctx context.Context, task *model.Task) error {
	exists, err := s.checkTaskExists(ctx, task.ID)
	if err != nil {
		return err
	}
	if exists {
		return model.ErrAlreadyExist{}
	}

	pipe := s.client.TxPipeline()
	pipe.HSet(ctx, "task", task.ID, true)
	pipe.HSet(ctx, "task-namespace", task.ID, task.NamespaceID)
	pipe.HSet(ctx, "task-retries-left", task.ID, task.RetriesLeft)
	pipe.HSet(ctx, "task-timeout", task.ID, task.Timeout)
	pipe.HSet(ctx, "task-updated-at", task.ID, task.UpdatedAt)
	pipe.ZAddNX(ctx, namespaceTaskKey(task.NamespaceID), r.Z{
		Member: task.ID,
		Score:  toScore(task.Status),
	})
	_, err = pipe.Exec(ctx)
	return err
}

func (s *Storage) checkTaskExists(ctx context.Context, taskID string) (bool, error) {
	cmd := s.client.HExists(ctx, "task", taskID)
	exists, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) DeregisterTask(ctx context.Context, taskID string) error {
	cmd := s.client.HGet(ctx, "task-namespace", taskID)
	namespaceID, err := cmd.Result()
	if err != nil {
		return err
	}

	pipe := s.client.TxPipeline()
	pipe.HSet(ctx, "task", taskID)
	pipe.HSet(ctx, "task-namespace", taskID)
	pipe.HSet(ctx, "task-retries-left", taskID)
	pipe.HSet(ctx, "task-timeout", taskID)
	pipe.HSet(ctx, "task-updated-at", taskID)
	pipe.ZRem(ctx, namespaceTaskKey(namespaceID))
	_, err = pipe.Exec(ctx)
	return err
}

func (s *Storage) ReadTask(ctx context.Context, taskID string) (*model.Task, error) {
	cmd := s.client.HGet(ctx, "task-namespace", taskID)
	namespaceID, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	pipe := s.client.TxPipeline()
	var (
		retriesLeftCmd = pipe.HGet(ctx, "task-retries-left", taskID)
		timeoutCmd     = pipe.HGet(ctx, "task-timeout", taskID)
		updatedAtCmd   = pipe.HGet(ctx, "task-updated-at", taskID)
		scoreCmd       = pipe.ZScore(ctx, namespaceTaskKey(namespaceID), taskID)
	)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	var (
		retriesLeft int
		timeout     int64
		updatedAt   int64
		score       float64
	)
	retriesLeft, err = retriesLeftCmd.Int()
	if err != nil {
		return nil, err
	}
	timeout, err = timeoutCmd.Int64()
	if err != nil {
		return nil, err
	}
	updatedAt, err = updatedAtCmd.Int64()
	if err != nil {
		return nil, err
	}
	score, err = scoreCmd.Result()
	if err != nil {
		return nil, err
	}

	return &model.Task{
		ID:          taskID,
		NamespaceID: namespaceID,
		RetriesLeft: retriesLeft,
		Timeout:     timeout,
		UpdatedAt:   updatedAt,
		Status:      fromScore(score),
	}, nil
}

func (s *Storage) CaptureTasks(ctx context.Context, namespaceID string, limit int) ([]*model.Task, error) {
	cmd := s.client.ZRangeArgs(ctx, r.ZRangeArgs{
		Key:     namespaceTaskKey(namespaceID),
		Start:   TaskScoreFailed,
		Stop:    TaskScoreIdle,
		ByScore: true,
		Count:   int64(limit),
	})
	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	captureTaskCmds := make(map[string]*r.IntCmd)

	capturePipe := s.client.TxPipeline()
	for _, taskID := range res {
		newScore := r.ZAddArgs{
			XX: true,
			Ch: true,
			Members: []r.Z{
				{
					Score:  TaskScoreInProgress,
					Member: taskID,
				},
			},
		}
		captureTaskCmd := capturePipe.ZAddArgs(ctx, namespaceTaskKey(namespaceID), newScore)
		captureTaskCmds[taskID] = captureTaskCmd
	}
	_, err = capturePipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	taskIDs := make([]string, 0, len(captureTaskCmds))
	for taskID, captureTaskCmd := range captureTaskCmds {
		if isCaptured(captureTaskCmd) {
			taskIDs = append(taskIDs, taskID)
		}
	}

	ts := util.UnixEpoch()

	fields, err := s.readTaskFields(ctx, taskIDs...)
	for retriesLeftInd := range fields.retriesLeft {
		fields.retriesLeft[retriesLeftInd] = fields.retriesLeft[retriesLeftInd] - 1
	}
	for updatedAtInd := range fields.updatedAt {
		fields.updatedAt[updatedAtInd] = ts
	}

	err = s.writeTaskFields(ctx, taskIDs, fields)
	if err != nil {
		return nil, err
	}

	tasks := make([]*model.Task, 0, len(taskIDs))
	for taskIDInd := range taskIDs {
		tasks = append(tasks, &model.Task{
			ID:          taskIDs[taskIDInd],
			NamespaceID: namespaceID,
			RetriesLeft: fields.retriesLeft[taskIDInd],
			Timeout:     fields.timeouts[taskIDInd],
			UpdatedAt:   ts,
			Status:      model.TaskStatusInProgress,
		})
	}
	return tasks, nil
}

func isCaptured(captureTaskCmd *r.IntCmd) bool {
	affected, err := captureTaskCmd.Result()
	if err != nil {
		return false
	}
	return affected > 0
}

func (s *Storage) ReleaseTasks(ctx context.Context, namespaceID string, taskIDs []string, status model.TaskStatus) error {
	members := make([]r.Z, 0, len(taskIDs))
	for _, taskID := range taskIDs {
		members = append(members, r.Z{
			Score:  toScore(status),
			Member: taskID,
		})
	}
	newScore := r.ZAddArgs{
		XX:      true,
		Members: members,
	}
	cmd := s.client.ZAddArgs(ctx, namespaceTaskKey(namespaceID), newScore)
	_, err := cmd.Result()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) RefreshTaskStatuses(ctx context.Context, namespaceID string) (int64, error) {
	cmd := s.client.ZRangeArgs(ctx, r.ZRangeArgs{
		Key:     namespaceTaskKey(namespaceID),
		Start:   TaskScoreInProgress,
		Stop:    TaskScoreInProgress,
		ByScore: true,
	})
	taskIDs, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	fields, err := s.readTaskFields(ctx, taskIDs...)
	if err != nil {
		return 0, err
	}

	taskIDsToRefresh := make([]string, 0)
	retriesLeftToRefresh := make([]int, 0)
	updatedAtToRefresh := make([]int64, 0)

	ts := util.UnixEpoch()
	for taskIDInd := range taskIDs {
		var (
			timeoutVal   = fields.timeouts[taskIDInd]
			updatedAtVal = fields.updatedAt[taskIDInd]
		)
		if ts > updatedAtVal+timeoutVal {
			taskIDsToRefresh = append(taskIDsToRefresh, taskIDs[taskIDInd])
			retriesLeftToRefresh = append(retriesLeftToRefresh, fields.retriesLeft[taskIDInd]-1)
			updatedAtToRefresh = append(updatedAtToRefresh, ts)
		}
	}

	members := make([]r.Z, 0, len(taskIDs))
	for _, taskID := range taskIDsToRefresh {
		members = append(members, r.Z{
			Score:  TaskScoreTimeout,
			Member: taskID,
		})
	}

	updateCmd := s.client.ZAddArgs(ctx, namespaceTaskKey(namespaceID), r.ZAddArgs{
		XX:      true,
		Ch:      true,
		Members: members,
	})
	tasksUpdated, err := updateCmd.Result()
	if err != nil {
		return 0, err
	}

	newFields := &taskFields{
		retriesLeft: retriesLeftToRefresh,
		updatedAt:   updatedAtToRefresh,
	}
	err = s.writeTaskFields(ctx, taskIDsToRefresh, newFields)
	if err != nil {
		return 0, err
	}

	return tasksUpdated, nil
}

func (s *Storage) ListTasks(ctx context.Context, namespaceID string) ([]*model.Task, error) {
	cmd := s.client.ZRangeWithScores(ctx, namespaceTaskKey(namespaceID), 0, -1)
	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	taskIDs := make([]string, 0, len(res))
	taskScores := make([]float64, 0, len(res))
	for _, item := range res {
		taskID, ok := item.Member.(string)
		if ok {
			taskIDs = append(taskIDs, taskID)
			taskScores = append(taskScores, item.Score)
		}
	}

	fields, err := s.readTaskFields(ctx, taskIDs...)
	if err != nil {
		return nil, err
	}

	tasks := make([]*model.Task, 0, len(taskIDs))
	for taskIDInd := range taskIDs {
		tasks = append(tasks, &model.Task{
			ID:          taskIDs[taskIDInd],
			NamespaceID: namespaceID,
			RetriesLeft: fields.retriesLeft[taskIDInd],
			Timeout:     fields.timeouts[taskIDInd],
			UpdatedAt:   fields.updatedAt[taskIDInd],
			Status:      fromScore(taskScores[taskIDInd]),
		})
	}

	return tasks, nil
}

type taskFields struct {
	retriesLeft []int
	timeouts    []int64
	updatedAt   []int64
}

func (s *Storage) readTaskFields(ctx context.Context, taskIDs ...string) (*taskFields, error) {
	pipe := s.client.TxPipeline()
	var (
		retriesLeftCmd = pipe.HMGet(ctx, "task-retries-left", taskIDs...)
		timeoutCmd     = pipe.HMGet(ctx, "task-timeout", taskIDs...)
		updatedAtCmd   = pipe.HMGet(ctx, "task-updated-at", taskIDs...)
	)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	retriesLeft := make([]int, 0, len(taskIDs))
	retriesLeftRaw, err := retriesLeftCmd.Result()
	if err != nil {
		return nil, err
	}
	for _, item := range retriesLeftRaw {
		var (
			strval string
			value  int
		)
		strval, _ = item.(string)
		value, err = strconv.Atoi(strval)
		if err != nil {
			return nil, err
		}
		retriesLeft = append(retriesLeft, value)
	}

	timeouts := make([]int64, 0, len(taskIDs))
	timeoutsRaw, err := timeoutCmd.Result()
	if err != nil {
		return nil, err
	}
	for _, item := range timeoutsRaw {
		var (
			strval string
			value  int64
		)
		strval, _ = item.(string)
		value, err = strconv.ParseInt(strval, 10, 64)
		if err != nil {
			return nil, err
		}
		timeouts = append(timeouts, value)
	}

	updatedAt := make([]int64, 0, len(taskIDs))
	updatedAtRaw, err := updatedAtCmd.Result()
	if err != nil {
		return nil, err
	}
	for _, item := range updatedAtRaw {
		var (
			strval string
			value  int64
		)
		strval, _ = item.(string)
		value, err = strconv.ParseInt(strval, 10, 64)
		if err != nil {
			return nil, err
		}
		updatedAt = append(updatedAt, value)
	}

	return &taskFields{
		retriesLeft: retriesLeft,
		timeouts:    timeouts,
		updatedAt:   updatedAt,
	}, nil
}

func (s *Storage) writeTaskFields(ctx context.Context, taskIDs []string, fields *taskFields) error {
	var (
		eqRetriesLen   = len(taskIDs) == len(fields.retriesLeft)
		eqTimeoutsLen  = len(taskIDs) == len(fields.timeouts)
		eqUpdatedAtLen = len(taskIDs) == len(fields.updatedAt)
	)

	if !eqRetriesLen && !eqTimeoutsLen && !eqUpdatedAtLen {
		return nil
	}

	pipe := s.client.TxPipeline()

	if eqRetriesLen {
		newRetriesLeft := make([]interface{}, 0, len(taskIDs))
		for taskInd := range taskIDs {
			newRetriesLeft = append(newRetriesLeft, taskIDs[taskInd])
			newRetriesLeft = append(newRetriesLeft, fields.retriesLeft[taskInd])
		}
		pipe.HMSet(ctx, "task-retries-left", newRetriesLeft)
	}

	if eqTimeoutsLen {
		newTimeouts := make([]interface{}, 0, len(taskIDs))
		for taskInd := range taskIDs {
			newTimeouts = append(newTimeouts, taskIDs[taskInd])
			newTimeouts = append(newTimeouts, fields.timeouts[taskInd])
		}
		pipe.HMSet(ctx, "task-timeout", newTimeouts)
	}

	if eqUpdatedAtLen {
		newUpdatedAt := make([]interface{}, 0, len(taskIDs))
		for taskInd := range taskIDs {
			newUpdatedAt = append(newUpdatedAt, taskIDs[taskInd])
			newUpdatedAt = append(newUpdatedAt, fields.updatedAt[taskInd])
		}
		pipe.HMSet(ctx, "task-updated-at", newUpdatedAt)
	}

	_, err := pipe.Exec(ctx)
	return err
}
