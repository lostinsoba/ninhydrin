package middleware

import (
	"context"
	"fmt"
	"net/http"
)

const (
	reqHeaderWorkerToken = "X-Ninhydrin-Worker-Token"
)

func WorkerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token := request.Header.Get(reqHeaderWorkerToken)
		ctx := context.WithValue(request.Context(), reqHeaderWorkerToken, token)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetWorkerToken(request *http.Request) (string, error) {
	token, ok := request.Context().Value(reqHeaderWorkerToken).(string)
	if !ok {
		return "", fmt.Errorf("%s not present", reqHeaderWorkerToken)
	}
	return token, nil
}
