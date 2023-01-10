package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

const (
	urlVariableTaskID = "taskID"
)

func TaskID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		taskID := chi.URLParam(request, urlVariableTaskID)
		ctx := context.WithValue(request.Context(), urlVariableTaskID, taskID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetTaskID(request *http.Request) (string, error) {
	taskID, ok := request.Context().Value(urlVariableTaskID).(string)
	if !ok {
		return "", fmt.Errorf("%s not present", urlVariableTaskID)
	}
	return taskID, nil
}

const (
	queryParamTaskCaptureLimit = "limit"
)

func QueryGetTaskCaptureLimit(request *http.Request) (int, error) {
	taskCaptureLimitStr := request.URL.Query().Get(queryParamTaskCaptureLimit)
	if taskCaptureLimitStr == "" {
		return 0, fmt.Errorf("%s not present", queryParamTaskCaptureLimit)
	}
	taskCaptureLimit, err := strconv.Atoi(taskCaptureLimitStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %s", queryParamTaskCaptureLimit, taskCaptureLimitStr)
	}
	return taskCaptureLimit, nil
}
