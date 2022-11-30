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

func ToTagListData(tags []string) *TagListData {
	return &TagListData{List: tags}
}

type TagListData struct {
	List []string `json:"list"`
}

func (*TagListData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
