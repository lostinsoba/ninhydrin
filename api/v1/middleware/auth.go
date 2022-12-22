package middleware

import (
	"context"
	"fmt"
	"net/http"
)

const (
	reqHeaderToken = "X-Ninhydrin-Token"
)

func Token(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token := request.Header.Get(reqHeaderToken)
		ctx := context.WithValue(request.Context(), reqHeaderToken, token)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetToken(request *http.Request) (string, error) {
	token, ok := request.Context().Value(reqHeaderToken).(string)
	if !ok {
		return "", fmt.Errorf("%s not present", reqHeaderToken)
	}
	return token, nil
}
