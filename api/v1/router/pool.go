package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/dto"
	"lostinsoba/ninhydrin/api/v1/middleware"
)

func (r *Router) pool(router chi.Router) {
	router.Get("/", r.listPools)
	router.Post("/register", r.registerPool)
	router.With(middleware.PoolID).Delete("/{poolID}", r.deregisterPool)
	router.With(middleware.PoolID).Put("/{poolID}", r.updatePool)
}

func (r *Router) listPools(writer http.ResponseWriter, request *http.Request) {
	list, err := r.ctrl.ListPools(request.Context())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}

	response := dto.ToPoolListData(list)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
}

func (r *Router) registerPool(writer http.ResponseWriter, request *http.Request) {
	poolData := dto.PoolData{}
	err := render.Bind(request, &poolData)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.RegisterPool(request.Context(), poolData.ToModel())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusCreated)
}

func (r *Router) deregisterPool(writer http.ResponseWriter, request *http.Request) {
	poolID, err := middleware.GetPoolID(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.DeregisterPool(request.Context(), poolID)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusOK)
}

func (r *Router) updatePool(writer http.ResponseWriter, request *http.Request) {
	poolData := dto.PoolData{}
	err := render.Bind(request, &poolData)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.UpdatePool(request.Context(), poolData.ToModel())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusOK)
}
