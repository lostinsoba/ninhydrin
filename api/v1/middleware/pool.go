package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	urlVariablePoolID = "poolID"
)

func PoolID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		poolID := chi.URLParam(request, urlVariablePoolID)
		ctx := context.WithValue(request.Context(), urlVariablePoolID, poolID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetPoolID(request *http.Request) (string, error) {
	poolID, ok := request.Context().Value(urlVariablePoolID).(string)
	if !ok {
		return "", fmt.Errorf("%s not present", urlVariablePoolID)
	}
	return poolID, nil
}
