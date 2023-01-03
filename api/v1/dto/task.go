package dto

import (
	"fmt"
	"net/http"

	"lostinsoba/ninhydrin/internal/model"
)

type TaskData struct {
	ID          string `json:"id"`
	Timeout     int64  `json:"timeout,omitempty"`
	RetriesLeft int    `json:"retries_left,omitempty"`
	UpdatedAt   int64  `json:"updated_at,omitempty"`
	Status      string `json:"status,omitempty"`
}

func (*TaskData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (taskData *TaskData) Bind(r *http.Request) error {
	if taskData.ID == "" {
		return fmt.Errorf("id required")
	}
	return nil
}

func (taskData *TaskData) ToModel() *model.Task {
	return &model.Task{
		ID:          taskData.ID,
		Timeout:     taskData.Timeout,
		RetriesLeft: taskData.RetriesLeft,
		UpdatedAt:   taskData.UpdatedAt,
		Status:      ToTaskStatus(taskData.Status),
	}
}

func ToTaskData(task *model.Task) *TaskData {
	return &TaskData{
		ID:          task.ID,
		Timeout:     task.Timeout,
		RetriesLeft: task.RetriesLeft,
		UpdatedAt:   task.UpdatedAt,
		Status:      string(task.Status),
	}
}

func ToTaskIDListData(taskIDs []string) *TaskIDListData {
	return &TaskIDListData{List: taskIDs}
}

type TaskIDListData struct {
	List []string `json:"list"`
}

func (*TaskIDListData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ReleaseData struct {
	Status  string   `json:"status"`
	TaskIDs []string `json:"task_ids"`
}

func (releaseData *ReleaseData) Bind(r *http.Request) error {
	if releaseData.Status == "" {
		return fmt.Errorf("status required")
	}
	if len(releaseData.TaskIDs) == 0 {
		return fmt.Errorf("task_ids required")
	}
	return nil
}

func ToTaskStatus(status string) model.TaskStatus {
	return model.TaskStatus(status)
}
