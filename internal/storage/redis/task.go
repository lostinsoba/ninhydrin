package redis

import (
	"context"
	"fmt"
	r "github.com/go-redis/redis/v9"
	"lostinsoba/ninhydrin/internal/util"

	"lostinsoba/ninhydrin/internal/model"
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

const (
	byteTask = "task"
)

func taskKey(taskID string) string {
	return fmt.Sprintf("%s:%s", byteTask, taskID)
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

	data, err := encode(task)
	if err != nil {
		return err
	}

	pipe := s.client.TxPipeline()
	pipe.Set(ctx, taskKey(task.ID), data, -1)
	pipe.ZAddNX(ctx, namespaceTaskKey(task.NamespaceID), r.Z{
		Member: task.ID,
		Score:  toScore(task.Status),
	})
	_, err = pipe.Exec(ctx)
	return err
}

func (s *Storage) checkTaskExists(ctx context.Context, taskID string) (bool, error) {
	cmd := s.client.Exists(ctx, taskKey(taskID))
	res, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return res > 0, err
}

func (s *Storage) DeregisterTask(ctx context.Context, taskID string) error {
	task, err := s.ReadTask(ctx, taskID)
	if err != nil {
		return err
	}

	pipe := s.client.TxPipeline()
	pipe.ZRem(ctx, namespaceTaskKey(task.NamespaceID))
	pipe.Del(ctx, taskKey(taskID))
	_, err = pipe.Exec(ctx)
	return err
}

func (s *Storage) ReadTask(ctx context.Context, taskID string) (task *model.Task, err error) {
	cmd := s.client.Get(ctx, taskKey(taskID))
	if cmd.Err() != nil {
		if isNil(cmd.Err()) {
			return nil, model.ErrNotFound{}
		}
		return nil, err
	}
	data, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}
	err = decode(data, &task)
	return
}

func (s *Storage) CaptureTasks(ctx context.Context, namespaceID string, limit int) (tasks []*model.Task, err error) {
	cmd := s.client.ZRangeArgs(ctx, r.ZRangeArgs{
		Key:     namespaceTaskKey(namespaceID),
		Start:   TaskScoreFailed,
		Stop:    TaskScoreIdle,
		ByScore: true,
		Count:   int64(limit),
	})
	res, err := cmd.Result()
	if err != nil {
		return
	}

	captureTaskCmds := make(map[string]*r.IntCmd)

	pipe := s.client.TxPipeline()
	for _, taskID := range res {
		newScore := r.Z{
			Score:  TaskScoreInProgress,
			Member: taskID,
		}
		captureTaskCmd := pipe.ZAddXX(ctx, namespaceTaskKey(namespaceID), newScore)
		captureTaskCmds[taskID] = captureTaskCmd
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return
	}

	taskIDs := make([]string, 0, len(captureTaskCmds))
	for taskID, captureTaskCmd := range captureTaskCmds {
		if isCaptured(captureTaskCmd) {
			taskIDs = append(taskIDs, taskID)
		}
	}

	tasks, err = s.listTasks(ctx, taskIDs...)
	if err != nil {
		return
	}

	for _, task := range tasks {
		task.Status = model.TaskStatusInProgress
		task.RetriesLeft = task.RetriesLeft - 1
		task.UpdatedAt = util.UnixEpoch()
	}
	err = s.updateTasks(ctx, tasks...)
	if err != nil {
		return
	}

	return nil, err
}

func (s *Storage) updateTasks(ctx context.Context, tasks ...*model.Task) error {
	encoded := make(map[string][]byte)
	for _, task := range tasks {
		data, err := encode(task)
		if err != nil {
			return err
		}
		encoded[task.ID] = data
	}

	pipe := s.client.TxPipeline()

	for _, task := range tasks {
		data := encoded[task.ID]
		pipe.Set(ctx, taskKey(task.ID), data, -1)
		pipe.ZAddXX(ctx, namespaceTaskKey(task.NamespaceID), r.Z{
			Member: task.ID,
			Score:  toScore(task.Status),
		})
	}
	_, err := pipe.Exec(ctx)
	return err
}

func isCaptured(captureTaskCmd *r.IntCmd) bool {
	affected, err := captureTaskCmd.Result()
	if err != nil {
		return false
	}
	return affected > 0
}

func (s *Storage) ReleaseTasks(ctx context.Context, taskIDs []string, status model.TaskStatus) error {
	return nil
}

func (s *Storage) RefreshTaskStatuses(ctx context.Context) (tasksUpdated int64, err error) {
	return 0, nil
}

func (s *Storage) ListTasks(ctx context.Context, namespaceID string) (tasks []*model.Task, err error) {
	cmd := s.client.ZRange(ctx, namespaceTaskKey(namespaceID), 0, -1)
	taskIDs, err := cmd.Result()
	if err != nil {
		return
	}
	return s.listTasks(ctx, taskIDs...)
}

func (s *Storage) listTasks(ctx context.Context, taskIDs ...string) (tasks []*model.Task, err error) {
	tasks = make([]*model.Task, 0, len(taskIDs))

	pipe := s.client.TxPipeline()
	getCmds := make([]*r.StringCmd, 0, len(taskIDs))
	for _, taskID := range taskIDs {
		getCmd := pipe.Get(ctx, taskKey(taskID))
		getCmds = append(getCmds, getCmd)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return
	}

	for _, getCmd := range getCmds {
		data, getCmdErr := getCmd.Bytes()
		if getCmdErr != nil {
			if isNil(getCmdErr) {
				continue
			}
			err = getCmdErr
			return
		}
		var task *model.Task
		err = decode(data, &task)
		if err != nil {
			return
		}
		tasks = append(tasks, task)
	}
	return
}
