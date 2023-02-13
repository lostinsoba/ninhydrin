package redis

import (
	"context"

	r "github.com/go-redis/redis/v9"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/storage/redis/schema"
	"lostinsoba/ninhydrin/internal/util"
)

func (s *Storage) RegisterTask(ctx context.Context, task *model.Task) error {
	exists, err := s.checkTaskExists(ctx, task.ID)
	if err != nil {
		return err
	}
	if exists {
		return model.ErrAlreadyExist{}
	}
	return s.registerTask(ctx, task)
}

func (s *Storage) registerTask(ctx context.Context, task *model.Task) error {
	taskData := schema.ToTaskData(task)
	data, err := schema.Encode(taskData)
	if err != nil {
		return err
	}

	pipe := s.client.TxPipeline()
	pipe.HSet(ctx, "task", task.ID, data)
	pipe.ZAddNX(ctx, schema.NamespaceTaskKey(task.NamespaceID), r.Z{
		Member: task.ID,
		Score:  schema.StatusToScore(task.Status),
	})
	_, err = pipe.Exec(ctx)
	return err
}

func (s *Storage) DeregisterTask(ctx context.Context, taskID string) error {
	namespaceID, err := s.getTaskNamespaceID(ctx, taskID)
	if err != nil {
		return err
	}
	return s.deregisterTask(ctx, namespaceID, taskID)
}

func (s *Storage) deregisterTask(ctx context.Context, namespaceID, taskID string) error {
	pipe := s.client.TxPipeline()
	pipe.HDel(ctx, "task", taskID)
	pipe.ZRem(ctx, schema.NamespaceTaskKey(namespaceID))
	_, err := pipe.Exec(ctx)
	return err
}

func (s *Storage) ReadTask(ctx context.Context, taskID string) (*model.Task, error) {
	namespaceID, err := s.getTaskNamespaceID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return s.readTask(ctx, namespaceID, taskID)
}

func (s *Storage) checkTaskExists(ctx context.Context, taskID string) (bool, error) {
	cmd := s.client.HExists(ctx, "task", taskID)
	exists, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) readTask(ctx context.Context, namespaceID, taskID string) (*model.Task, error) {
	pipe := s.client.TxPipeline()
	var (
		taskCmd  = pipe.HGet(ctx, "task", taskID)
		scoreCmd = pipe.ZScore(ctx, schema.NamespaceTaskKey(namespaceID), taskID)
	)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	data, err := taskCmd.Bytes()
	if err != nil {
		return nil, err
	}
	var taskData *schema.TaskData
	err = schema.Decode(data, &taskData)
	if err != nil {
		return nil, err
	}

	score, err := scoreCmd.Result()
	if err != nil {
		return nil, err
	}
	return taskData.ToTask(namespaceID, score), nil
}

func (s *Storage) CaptureTasks(ctx context.Context, namespaceID string, limit int) ([]*model.Task, error) {
	tasks, err := s.listTasks(ctx, namespaceID, schema.TaskScoreFailed, schema.TaskScoreIdle)
	if err != nil {
		return nil, err
	}

	capturedIDs := make(map[string]bool)
	for _, task := range tasks {
		task.Status = model.TaskStatusInProgress
		capturedIDs[task.ID] = false
	}

	capturedTaskIDs, err := s.updateTasks(ctx, tasks)
	if err != nil {
		return nil, err
	}

	for _, capturedID := range capturedTaskIDs {
		capturedIDs[capturedID] = true
	}

	capturedTasks := make([]*model.Task, 0, len(capturedTaskIDs))
	for _, task := range tasks {
		_, isCaptured := capturedIDs[task.ID]
		if isCaptured {
			capturedTasks = append(capturedTasks, task)
		}
	}

	return capturedTasks, nil
}

func (s *Storage) ReleaseTasks(ctx context.Context, namespaceID string, taskIDs []string, status model.TaskStatus) error {

	return nil
}

func (s *Storage) RefreshTaskStatuses(ctx context.Context, namespaceID string) (int64, error) {
	tasks, err := s.listTasks(ctx, namespaceID, schema.TaskScoreInProgress, schema.TaskScoreInProgress)
	if err != nil {
		return 0, err
	}

	ts := util.UnixEpoch()
	for _, task := range tasks {
		if ts > task.UpdatedAt+task.Timeout {
			task.RetriesLeft = task.RetriesLeft - 1
			task.Status = model.TaskStatusTimeout
		}
	}

	tasksUpdated, err := s.updateTasks(ctx, tasks)
	if err != nil {
		return 0, err
	}
	return int64(len(tasksUpdated)), nil
}

func (s *Storage) ListTasks(ctx context.Context, namespaceID string) ([]*model.Task, error) {
	return s.listTasks(ctx, namespaceID, schema.TaskScoreTimeout, schema.TaskScoreDone)
}

func (s *Storage) listTasks(ctx context.Context, namespaceID string, fromScore, toScore int) ([]*model.Task, error) {
	taskIDsWithScoresCmd := s.client.ZRangeArgsWithScores(ctx, r.ZRangeArgs{
		Key:     schema.NamespaceTaskKey(namespaceID),
		Start:   fromScore,
		Stop:    toScore,
		ByScore: true,
	})
	taskIDsWithScores, err := taskIDsWithScoresCmd.Result()
	if err != nil {
		return nil, err
	}

	taskCmds := make([]*r.StringCmd, 0, len(taskIDsWithScores))

	pipe := s.client.TxPipeline()
	for _, taskIDWithScore := range taskIDsWithScores {
		taskID, _ := taskIDWithScore.Member.(string)
		taskCmds = append(taskCmds, pipe.HGet(ctx, "task", taskID))
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	tasks := make([]*model.Task, 0, len(taskIDsWithScores))

	for taskCmdInd := range taskCmds {
		data, err := taskCmds[taskCmdInd].Bytes()
		if err != nil {
			return nil, err
		}
		var taskData *schema.TaskData
		err = schema.Decode(data, &taskData)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, taskData.ToTask(namespaceID, taskIDsWithScores[taskCmdInd].Score))
	}

	return tasks, nil
}

func (s *Storage) updateTasks(ctx context.Context, tasks []*model.Task) (tasksUpdated []string, err error) {
	scoreMap := make(map[string]*r.IntCmd)

	pipe := s.client.TxPipeline()

	for _, task := range tasks {
		taskData := schema.ToTaskData(task)
		data, err := schema.Encode(taskData)
		if err != nil {
			return
		}
		pipe.HSet(ctx, "task", task.ID, data)
		newScore := r.ZAddArgs{
			XX: true,
			Ch: true,
			Members: []r.Z{
				{
					Score:  schema.StatusToScore(task.Status),
					Member: task.ID,
				},
			},
		}
		scoreCmd := pipe.ZAddArgs(ctx, schema.NamespaceTaskKey(task.ID), newScore)
		scoreMap[task.ID] = scoreCmd
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return
	}

	for taskID, scoreCmd := range scoreMap {
		if isAffected(scoreCmd) {
			tasksUpdated = append(tasksUpdated, taskID)
		}
	}
	return
}

func isAffected(scoreCmd *r.IntCmd) bool {
	affected, err := scoreCmd.Result()
	if err != nil {
		return false
	}
	return affected > 0
}

func (s *Storage) getTaskNamespaceID(ctx context.Context, taskID string) (string, error) {
	return s.client.HGet(ctx, "task-namespace", taskID).Result()
}
