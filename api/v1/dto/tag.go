package dto

import (
	"fmt"
	"net/http"
)

type TagData struct {
	ID string `json:"id"`
}

func (*TagData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (tagData *TagData) Bind(r *http.Request) error {
	if tagData.ID == "" {
		return fmt.Errorf("id required")
	}
	return nil
}

func ToTagData(tag string) *TagData {
	return &TagData{ID: tag}
}

func ToTagIDListData(tags []string) *TagIDListData {
	return &TagIDListData{List: tags}
}

type TagIDListData struct {
	List []string `json:"list"`
}

func (*TagIDListData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
