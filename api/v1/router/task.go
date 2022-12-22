package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/dto"
	"lostinsoba/ninhydrin/api/v1/middleware"
)

func (r *Router) task(router chi.Router) {
	router.Get("/", r.listTasks)
	router.Post("/", r.registerTask)
	router.Route("/{taskID}", func(router chi.Router) {
		router.Use(middleware.TaskID)
		router.Get("/", r.readTask)
		router.Delete("/", r.deregisterTask)
	})
}

func (r *Router) listTasks(writer http.ResponseWriter, request *http.Request) {
	list, err := r.ctrl.ListTaskIDs(request.Context())
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

func (r *Router) readTask(writer http.ResponseWriter, request *http.Request) {
	taskID, err := middleware.GetTaskID(request)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	task, ok, err := r.ctrl.ReadTask(request.Context(), taskID)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	if !ok {
		render.Status(request, http.StatusNoContent)
		return
	}

	response := dto.ToTaskData(task)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
}

func (r *Router) deregisterTask(writer http.ResponseWriter, request *http.Request) {
	taskID, err := middleware.GetTaskID(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.DeregisterTask(request.Context(), taskID)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusOK)
}
