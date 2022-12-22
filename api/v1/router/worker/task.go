package worker

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/dto"
	"lostinsoba/ninhydrin/api/v1/middleware"
)

func (r *Router) task(router chi.Router) {
	router.Get("/capture", r.captureTasks)
	router.Route("/{taskID}", func(router chi.Router) {
		router.Use(middleware.TaskID)
		router.Put("/status", r.updateTaskStatus)
	})
}

func (r *Router) captureTasks(writer http.ResponseWriter, request *http.Request) {
	poolID, err := middleware.GetPoolID(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	limit, err := middleware.GetTaskCaptureLimit(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	list, err := r.ctrl.CaptureTaskIDs(request.Context(), poolID, limit)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	response := dto.ToTaskIDListData(list)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
}

func (r *Router) updateTaskStatus(writer http.ResponseWriter, request *http.Request) {
	taskID, err := middleware.GetTaskID(request)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	statusUpdateData := dto.TaskStatusUpdateData{}
	err = render.Bind(request, &statusUpdateData)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.UpdateTaskStatus(request.Context(), taskID, statusUpdateData.ToModel())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusAccepted)
}
