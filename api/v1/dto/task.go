package dto

import (
	"fmt"
	"net/http"

	"lostinsoba/ninhydrin/internal/model"
)

type TaskData struct {
	ID          string `json:"id"`
	NamespaceID string `json:"namespace_id"`
	Timeout     int64  `json:"timeout,omitempty"`
	RetriesLeft int    `json:"retries_left,omitempty"`
	UpdatedAt   int64  `json:"updated_at,omitempty"`
	Status      string `json:"status,omitempty"`
}

func (*TaskData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (taskData *TaskData) Bind(r *http.Request) error {
	if taskData.NamespaceID == "" {
		return fmt.Errorf("namespace_id required")
	}
	if taskData.ID == "" {
		return fmt.Errorf("id required")
	}
	return nil
}

func (taskData *TaskData) ToModel() *model.Task {
	return &model.Task{
		ID:          taskData.ID,
		NamespaceID: taskData.NamespaceID,
		Timeout:     taskData.Timeout,
		RetriesLeft: taskData.RetriesLeft,
		UpdatedAt:   taskData.UpdatedAt,
		Status:      ToTaskStatus(taskData.Status),
	}
}

func ToTaskData(task *model.Task) *TaskData {
	return &TaskData{
		ID:          task.ID,
		NamespaceID: task.NamespaceID,
		Timeout:     task.Timeout,
		RetriesLeft: task.RetriesLeft,
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
