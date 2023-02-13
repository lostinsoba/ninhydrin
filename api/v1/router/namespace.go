package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/dto"
	"lostinsoba/ninhydrin/api/v1/middleware"
)

func (r *Router) namespace(router chi.Router) {
	router.Get("/", r.listNamespaces)
	router.Post("/", r.registerNamespace)
	router.Route("/{namespaceID}", func(router chi.Router) {
		router.Use(middleware.NamespaceID)
		router.Get("/", r.readNamespace)
		router.Delete("/", r.deregisterNamespace)
		router.Get("/capture", r.captureTasks)
	})
}

func (r *Router) listNamespaces(writer http.ResponseWriter, request *http.Request) {
	list, err := r.ctrl.ListNamespaces(request.Context())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	response := dto.ToNamespaceListData(list)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
}

func (r *Router) registerNamespace(writer http.ResponseWriter, request *http.Request) {
	namespaceData := dto.NamespaceData{}
	err := render.Bind(request, &namespaceData)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.RegisterNamespace(request.Context(), namespaceData.ToModel())
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusCreated)
}

func (r *Router) readNamespace(writer http.ResponseWriter, request *http.Request) {
	namespaceID, err := middleware.GetNamespaceID(request)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	namespace, ok, err := r.ctrl.ReadNamespace(request.Context(), namespaceID)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	if !ok {
		render.NoContent(writer, request)
		return
	}

	response := dto.ToNamespaceData(namespace)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
}

func (r *Router) deregisterNamespace(writer http.ResponseWriter, request *http.Request) {
	namespaceID, err := middleware.GetNamespaceID(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	err = r.ctrl.DeregisterNamespace(request.Context(), namespaceID)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	render.Status(request, http.StatusOK)
}

func (r *Router) captureTasks(writer http.ResponseWriter, request *http.Request) {
	namespaceID, err := middleware.GetNamespaceID(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	limit, err := middleware.QueryGetTaskCaptureLimit(request)
	if err != nil {
		render.Render(writer, request, dto.InvalidRequestError(err))
		return
	}
	list, err := r.ctrl.CaptureTasks(request.Context(), namespaceID, limit)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
	response := dto.ToTaskStateListData(list)
	err = render.Render(writer, request, response)
	if err != nil {
		render.Render(writer, request, dto.InternalServerError(err))
		return
	}
}
