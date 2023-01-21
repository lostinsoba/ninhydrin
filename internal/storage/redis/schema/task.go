package schema

import "lostinsoba/ninhydrin/internal/model"

type TaskData struct {
	ID          string
	RetriesLeft int
	Timeout     int64
	UpdatedAt   int64
}

func (d *TaskData) ToTask(namespaceID string, score float64) *model.Task {
	return &model.Task{
		ID:          d.ID,
		NamespaceID: namespaceID,
		RetriesLeft: d.RetriesLeft,
		Timeout:     d.Timeout,
		UpdatedAt:   d.UpdatedAt,
		Status:      StatusFromScore(score),
	}
}

func ToTaskData(task *model.Task) *TaskData {
	return &TaskData{
		ID:          task.ID,
		RetriesLeft: task.RetriesLeft,
		Timeout:     task.Timeout,
		UpdatedAt:   task.UpdatedAt,
	}
}
