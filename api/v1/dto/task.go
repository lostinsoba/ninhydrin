package dto

import (
	"fmt"
	"net/http"

	"lostinsoba/ninhydrin/internal/model"
)

type TaskData struct {
	ID          string `json:"id"`
	PoolID      string `json:"pool_id"`
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
	if taskData.PoolID == "" {
		return fmt.Errorf("pool_id required")
	}
	return nil
}

func (taskData *TaskData) ToModel() *model.Task {
	return &model.Task{
		ID:          taskData.ID,
		PoolID:      taskData.PoolID,
		Timeout:     taskData.Timeout,
		RetriesLeft: taskData.RetriesLeft,
		UpdatedAt:   taskData.UpdatedAt,
		Status:      model.TaskStatus(taskData.Status),
	}
}

func ToTaskData(task *model.Task) *TaskData {
	return &TaskData{
		ID:          task.ID,
		PoolID:      task.PoolID,
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

type TaskStatusUpdateData struct {
	Status string `json:"status"`
}

func (statusUpdateData *TaskStatusUpdateData) Bind(r *http.Request) error {
	if statusUpdateData.Status == "" {
		return fmt.Errorf("status required")
	}
	return nil
}

func (statusUpdateData *TaskStatusUpdateData) ToModel() model.TaskStatus {
	return model.TaskStatus(statusUpdateData.Status)
}
