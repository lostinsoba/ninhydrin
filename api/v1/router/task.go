package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/dto"
	"lostinsoba/ninhydrin/api/v1/middleware"
)

func (r *Router) task(router chi.Router) {
	router.Get("/", r.listCurrentTasks)
	router.Post("/register", r.registerTask)
	router.With(middleware.WorkerID).Get("/capture", r.captureTasks)
	router.With(middleware.TaskID).Put("/{id}/status", r.updateTaskStatus)
}

func (r *Router) listCurrentTasks(writer http.ResponseWriter, request *http.Request) {
	list, err := r.ctrl.ListCurrentTasks(request.Context())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	response := dto.ToTaskListData(list)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
}

func (r *Router) registerTask(writer http.ResponseWriter, request *http.Request) {
	taskData := dto.TaskData{}
	err := render.Bind(request, &taskData)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.RegisterTask(request.Context(), taskData.ToModel())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusCreated)
}

func (r *Router) captureTasks(writer http.ResponseWriter, request *http.Request) {
	workerID, err := middleware.GetWorkerID(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	limit, err := middleware.GetTaskCaptureLimit(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	list, err := r.ctrl.CaptureTasks(request.Context(), workerID, limit)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	response := dto.ToTaskListData(list)
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