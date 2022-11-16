package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/dto"
)

func (r *Router) worker(router chi.Router) {
	router.Get("/", r.listWorkers)
	router.Post("/register", r.registerWorker)
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
