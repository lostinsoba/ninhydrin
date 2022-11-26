package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	urlVariableTagID = "tagID"
)

func TagID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tagID := chi.URLParam(request, urlVariableTagID)
		ctx := context.WithValue(request.Context(), urlVariableTagID, tagID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetTagID(request *http.Request) (string, error) {
	tagID, ok := request.Context().Value(urlVariableTagID).(string)
	if !ok {
		return "", fmt.Errorf("%s not present", urlVariableTagID)
	}
	return tagID, nil
}
