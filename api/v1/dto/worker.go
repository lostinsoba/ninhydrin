package dto

import (
	"fmt"
	"net/http"

	"lostinsoba/ninhydrin/internal/model"
)

type WorkerData struct {
	ID     string   `json:"id"`
	TagIDs []string `json:"tag_ids"`
}

func (*WorkerData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (workerData *WorkerData) Bind(r *http.Request) error {
	if workerData.ID == "" {
		return fmt.Errorf("id required")
	}
	var validTagIDsCount uint
	for _, tagIDs := range workerData.TagIDs {
		if tagIDs != "" {
			validTagIDsCount++
		}
	}
	if validTagIDsCount == 0 {
		return fmt.Errorf("tag_ids required")
	}
	return nil
}

func (workerData *WorkerData) ToModel() *model.Worker {
	return &model.Worker{
		ID:     workerData.ID,
		TagIDs: workerData.TagIDs,
	}
}

func ToWorkerData(worker *model.Worker) *WorkerData {
	return &WorkerData{
		ID:     worker.ID,
		TagIDs: worker.TagIDs,
	}
}

func ToWorkerIDListData(workerIDs []string) *WorkerIDListData {
	return &WorkerIDListData{List: workerIDs}
}

type WorkerIDListData struct {
	List []string `json:"list"`
}

func (*WorkerIDListData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
