package dto

import (
	"fmt"
	"net/http"

	"lostinsoba/ninhydrin/internal/model"
)

type PoolData struct {
	ID          string   `json:"id"`
	Description string   `json:"description,omitempty"`
	TagIDs      []string `json:"tag_ids"`
}

func (*PoolData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (poolData *PoolData) Bind(r *http.Request) error {
	if poolData.ID == "" {
		return fmt.Errorf("id required")
	}
	if len(poolData.TagIDs) == 0 {
		return fmt.Errorf("tag_ids required")
	}
	return nil
}

func (poolData *PoolData) ToModel() *model.Pool {
	return &model.Pool{
		ID:          poolData.ID,
		Description: poolData.Description,
		TagIDs:      poolData.TagIDs,
	}
}

func ToPoolData(pool *model.Pool) *PoolData {
	return &PoolData{
		ID:          pool.ID,
		Description: pool.Description,
		TagIDs:      pool.TagIDs,
	}
}

func ToPoolIDListData(pools []string) *PoolIDListData {
	return &PoolIDListData{List: pools}
}

type PoolIDListData struct {
	List []string `json:"list"`
}

func (*PoolIDListData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ReleaseData struct {
	TaskIDs []string `json:"task_ids"`
	Status  string   `json:"status"`
}

func (releaseData *ReleaseData) Bind(r *http.Request) error {
	if len(releaseData.TaskIDs) == 0 {
		return fmt.Errorf("task_ids required")
	}
	if releaseData.Status == "" {
		return fmt.Errorf("status required")
	}
	return nil
}
