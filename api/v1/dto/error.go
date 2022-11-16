package dto

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponseData struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *ErrorResponseData) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func InternalServerError(err error) *ErrorResponseData {
	return &ErrorResponseData{
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
	}
}

func InvalidRequestError(err error) *ErrorResponseData {
	return &ErrorResponseData{
		Message:    err.Error(),
		StatusCode: http.StatusBadRequest,
	}
}
