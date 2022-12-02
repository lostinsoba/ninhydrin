package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	urlVariableWorkerID = "workerID"
)

func WorkerID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		workerID := chi.URLParam(request, urlVariableWorkerID)
		ctx := context.WithValue(request.Context(), urlVariableWorkerID, workerID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetWorkerID(request *http.Request) (string, error) {
	workerID, ok := request.Context().Value(urlVariableWorkerID).(string)
	if !ok {
		return "", fmt.Errorf("%s not present", urlVariableWorkerID)
	}
	return workerID, nil
}
