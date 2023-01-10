package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	urlVariableNamespaceID = "namespaceID"
)

func NamespaceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		namespaceID := chi.URLParam(request, urlVariableNamespaceID)
		ctx := context.WithValue(request.Context(), urlVariableNamespaceID, namespaceID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetNamespaceID(request *http.Request) (string, error) {
	namespaceID, ok := request.Context().Value(urlVariableNamespaceID).(string)
	if !ok {
		return "", fmt.Errorf("%s not present", urlVariableNamespaceID)
	}
	return namespaceID, nil
}

const (
	queryParamNamespaceID = "namespace_id"
)

func QueryGetNamespaceID(request *http.Request) (string, error) {
	namespaceID := request.URL.Query().Get(queryParamNamespaceID)
	if namespaceID == "" {
		return "", fmt.Errorf("%s not present", queryParamNamespaceID)
	}
	return namespaceID, nil
}
