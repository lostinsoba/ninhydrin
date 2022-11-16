package router

import (
	"github.com/go-chi/chi"

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
		router.Route("/tag", r.tag)
		router.Route("/pool", r.pool)
		router.Route("/task", r.task)
		router.Route("/worker", r.worker)
	}
}
