package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/dto"
)

func (r *Router) tag(router chi.Router) {
	router.Get("/", r.listTags)
	router.Post("/register", r.registerTag)
}

func (r *Router) listTags(writer http.ResponseWriter, request *http.Request) {
	list, err := r.ctrl.ListTags(request.Context())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	response := dto.ToTagListData(list)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
}

func (r *Router) registerTag(writer http.ResponseWriter, request *http.Request) {
	tagData := dto.TagData{}
	err := render.Bind(request, &tagData)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.RegisterTag(request.Context(), tagData.ID)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusCreated)
}
