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
		Timeout:     taskData.Timeout,
	}
}

func ToTaskData(task *model.Task) *TaskData {
	return &TaskData{
		ID:          task.ID,
		NamespaceID: task.NamespaceID,
		Timeout:     task.Timeout,
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
