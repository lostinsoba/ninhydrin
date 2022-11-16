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
	var validTagIDsCount uint
	for _, tagIDs := range poolData.TagIDs {
		if tagIDs != "" {
			validTagIDsCount++
		}
	}
	if validTagIDsCount == 0 {
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

func ToPoolListData(pools []*model.Pool) *PoolListData {
	list := make([]*PoolData, 0, len(pools))
	for _, pool := range pools {
		list = append(list, ToPoolData(pool))
	}
	return &PoolListData{List: list}
}

type PoolListData struct {
	List []*PoolData `json:"list"`
}

func (*PoolListData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
