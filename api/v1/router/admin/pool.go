package admin

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/dto"
	"lostinsoba/ninhydrin/api/v1/middleware"
)

func (r *Router) pool(router chi.Router) {
	router.Get("/", r.listPoolIDs)
	router.Post("/", r.registerPool)
	router.Route("/{poolID}", func(router chi.Router) {
		router.Use(middleware.PoolID)
		router.Get("/", r.readPool)
		router.Put("/", r.updatePool)
		router.Delete("/", r.deregisterPool)
	})
}

func (r *Router) listPoolIDs(writer http.ResponseWriter, request *http.Request) {
	list, err := r.ctrl.ListPoolIDs(request.Context())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}

	response := dto.ToPoolIDListData(list)
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

func (r *Router) readPool(writer http.ResponseWriter, request *http.Request) {
	poolID, err := middleware.GetPoolID(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	pool, ok, err := r.ctrl.ReadPool(request.Context(), poolID)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	if !ok {
		render.Status(request, http.StatusNoContent)
		return
	}

	response := dto.ToPoolData(pool)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
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
