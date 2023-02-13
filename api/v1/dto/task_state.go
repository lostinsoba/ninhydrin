package dto

import (
	"fmt"
	"net/http"

	"lostinsoba/ninhydrin/internal/model"
)

type TaskStateData struct {
	TaskID      string `json:"task_id"`
	RetriesLeft int    `json:"retries_left"`
	UpdatedAt   int64  `json:"updated_at"`
	Status      string `json:"status"`
}

func (*TaskStateData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (taskStateData *TaskStateData) Bind(r *http.Request) error {
	if taskStateData.TaskID == "" {
		return fmt.Errorf("task id required")
	}
	if taskStateData.Status == "" {
		return fmt.Errorf("status required")
	}
	return nil
}

func (taskStateData *TaskStateData) ToModel() *model.TaskState {
	return &model.TaskState{
		TaskID:      taskStateData.TaskID,
		RetriesLeft: taskStateData.RetriesLeft,
		UpdatedAt:   taskStateData.UpdatedAt,
		Status:      model.TaskStatus(taskStateData.Status),
	}
}

func ToTaskStateData(taskState *model.TaskState) *TaskStateData {
	return &TaskStateData{
		TaskID:      taskState.TaskID,
		RetriesLeft: taskState.RetriesLeft,
		UpdatedAt:   taskState.UpdatedAt,
		Status:      string(taskState.Status),
	}
}

func ToTaskStateListData(taskStates []*model.TaskState) *TaskStateListData {
	list := make([]*TaskStateData, 0, len(taskStates))
	for _, taskState := range taskStates {
		list = append(list, ToTaskStateData(taskState))
	}
	return &TaskStateListData{List: list}
}

type TaskStateListData struct {
	List []*TaskStateData `json:"list"`
}

func (*TaskStateListData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
