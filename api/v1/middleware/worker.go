package middleware

import (
	"context"
	"fmt"
	"net/http"
)

const (
	reqHeaderWorkerID = "X-NINHYDRIN-WORKER-ID"
)

func WorkerID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userLogin := request.Header.Get(reqHeaderWorkerID)
		ctx := context.WithValue(request.Context(), reqHeaderWorkerID, userLogin)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetWorkerID(request *http.Request) (string, error) {
	workerID, ok := request.Context().Value(reqHeaderWorkerID).(string)
	if !ok {
		return "", fmt.Errorf("%s not present", reqHeaderWorkerID)
	}
	return workerID, nil
}
