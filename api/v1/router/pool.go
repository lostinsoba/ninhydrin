package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/dto"
)

func (r *Router) pool(router chi.Router) {
	router.Get("/", r.listPools)
	router.Post("/register", r.registerPool)
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
