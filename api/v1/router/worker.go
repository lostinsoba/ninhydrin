package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/dto"
	"lostinsoba/ninhydrin/api/v1/middleware"
)

func (r *Router) worker(router chi.Router) {
	router.Get("/", r.listWorkers)
	router.Post("/", r.registerWorker)
	router.Route("/{workerID}", func(router chi.Router) {
		router.Use(middleware.WorkerID)
		router.Get("/", r.readWorker)
		router.Delete("/", r.deregisterWorker)
	})
}

func (r *Router) listWorkers(writer http.ResponseWriter, request *http.Request) {
	list, err := r.ctrl.ListWorkers(request.Context())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}

	response := dto.ToWorkerListData(list)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
}

func (r *Router) registerWorker(writer http.ResponseWriter, request *http.Request) {
	workerData := dto.WorkerData{}
	err := render.Bind(request, &workerData)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.RegisterWorker(request.Context(), workerData.ToModel())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusCreated)
}

func (r *Router) readWorker(writer http.ResponseWriter, request *http.Request) {
	workerID, err := middleware.GetWorkerID(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	worker, err := r.ctrl.ReadWorker(request.Context(), workerID)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	response := dto.ToWorkerData(worker)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
}

func (r *Router) deregisterWorker(writer http.ResponseWriter, request *http.Request) {
	workerID, err := middleware.GetWorkerID(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.DeregisterWorker(request.Context(), workerID)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusOK)
}
