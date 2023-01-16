package dto

import (
	"fmt"
	"net/http"

	"lostinsoba/ninhydrin/internal/model"
)

type TaskData struct {
	ID          string `json:"id"`
	NamespaceID string `json:"namespace_id"`
	RetriesLeft int    `json:"retries_left,omitempty"`
	Timeout     int64  `json:"timeout,omitempty"`
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
	if taskData.NamespaceID == "" {
		return fmt.Errorf("namespace_id required")
	}
	return nil
}

func (taskData *TaskData) ToModel() *model.Task {
	return &model.Task{
		ID:          taskData.ID,
		NamespaceID: taskData.NamespaceID,
		RetriesLeft: taskData.RetriesLeft,
		Timeout:     taskData.Timeout,
		UpdatedAt:   taskData.UpdatedAt,
		Status:      ToTaskStatus(taskData.Status),
	}
}

func ToTaskData(task *model.Task) *TaskData {
	return &TaskData{
		ID:          task.ID,
		NamespaceID: task.NamespaceID,
		RetriesLeft: task.RetriesLeft,
		Timeout:     task.Timeout,
		UpdatedAt:   task.UpdatedAt,
		Status:      string(task.Status),
	}
}

func ToTaskListData(tasks []*model.Task) *TaskListData {
	list := make([]*TaskData, 0, len(tasks))
	for _, task := range tasks {
		list = append(list, ToTaskData(task))
	}
	return &TaskListData{List: list}
}

type TaskListData struct {
	List []*TaskData `json:"list"`
}

func (*TaskListData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ReleaseData struct {
	NamespaceID string   `json:"namespace_id"`
	TaskIDs     []string `json:"task_ids"`
	Status      string   `json:"status"`
}

func (releaseData *ReleaseData) Bind(r *http.Request) error {
	if releaseData.NamespaceID == "" {
		return fmt.Errorf("namespace_id required")
	}
	if len(releaseData.TaskIDs) == 0 {
		return fmt.Errorf("task_ids required")
	}
	if releaseData.Status == "" {
		return fmt.Errorf("status required")
	}
	return nil
}

func ToTaskStatus(status string) model.TaskStatus {
	return model.TaskStatus(status)
}
