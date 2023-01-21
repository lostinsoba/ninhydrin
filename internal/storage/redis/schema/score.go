package schema

import "lostinsoba/ninhydrin/internal/model"

const (
	TaskScoreTimeout = iota
	TaskScoreFailed
	TaskScoreIdle
	TaskScoreInProgress
	TaskScoreDone
)

func StatusToScore(status model.TaskStatus) float64 {
	scoreMap := map[model.TaskStatus]int{
		model.TaskStatusTimeout:    TaskScoreTimeout,
		model.TaskStatusFailed:     TaskScoreFailed,
		model.TaskStatusIdle:       TaskScoreIdle,
		model.TaskStatusInProgress: TaskScoreInProgress,
		model.TaskStatusDone:       TaskScoreDone,
	}
	return float64(scoreMap[status])
}

func StatusFromScore(score float64) model.TaskStatus {
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
