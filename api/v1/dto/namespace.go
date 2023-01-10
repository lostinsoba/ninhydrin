package dto

import (
	"fmt"
	"net/http"

	"lostinsoba/ninhydrin/internal/model"
)

type NamespaceData struct {
	ID string `json:"id"`
}

func (*NamespaceData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (namespaceData *NamespaceData) Bind(r *http.Request) error {
	if namespaceData.ID == "" {
		return fmt.Errorf("id required")
	}
	return nil
}

func (namespaceData *NamespaceData) ToModel() *model.Namespace {
	return &model.Namespace{
		ID: namespaceData.ID,
	}
}

func ToNamespaceData(namespace *model.Namespace) *NamespaceData {
	return &NamespaceData{
		ID: namespace.ID,
	}
}

func ToNamespaceListData(namespaces []*model.Namespace) *NamespaceListData {
	list := make([]*NamespaceData, 0, len(namespaces))
	for _, namespace := range namespaces {
		list = append(list, ToNamespaceData(namespace))
	}
	return &NamespaceListData{List: list}
}

type NamespaceListData struct {
	List []*NamespaceData `json:"list"`
}

func (*NamespaceListData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
