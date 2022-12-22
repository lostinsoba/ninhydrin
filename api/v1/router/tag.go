package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/dto"
	"lostinsoba/ninhydrin/api/v1/middleware"
)

func (r *Router) tag(router chi.Router) {
	router.Get("/", r.listTagIDs)
	router.Post("/", r.registerTag)
	router.Route("/{tagID}", func(router chi.Router) {
		router.Use(middleware.TagID)
		router.Get("/", r.readTag)
		router.Delete("/", r.deregisterTag)
	})
}

func (r *Router) listTagIDs(writer http.ResponseWriter, request *http.Request) {
	list, err := r.ctrl.ListTagIDs(request.Context())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	response := dto.ToTagIDListData(list)
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

func (r *Router) readTag(writer http.ResponseWriter, request *http.Request) {
	tagID, err := middleware.GetTagID(request)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	tag, ok, err := r.ctrl.ReadTag(request.Context(), tagID)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	if !ok {
		render.Status(request, http.StatusNoContent)
		return
	}

	response := dto.ToTagData(tag)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
}

func (r *Router) deregisterTag(writer http.ResponseWriter, request *http.Request) {
	tagID, err := middleware.GetTagID(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.DeregisterTag(request.Context(), tagID)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusOK)
}
