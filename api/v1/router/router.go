package router

import (
	"github.com/go-chi/chi"

	"lostinsoba/ninhydrin/api/v1/middleware"
	"lostinsoba/ninhydrin/controller"
)

type Router struct {
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Router {
	return &Router{ctrl: ctrl}
}

func (r *Router) Route() func(router chi.Router) {
	return func(router chi.Router) {
		router.Use(middleware.Token)
		router.Route("/task", r.task)
	}
}
