package api

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"lostinsoba/ninhydrin/api/v1/router/user"
	"lostinsoba/ninhydrin/api/v1/router/worker"
	"lostinsoba/ninhydrin/controller"
	"lostinsoba/ninhydrin/internal/monitoring/logger"
)

const (
	defaultCtxTimeout = 5 * time.Second
)

type Service struct {
	httpServer *http.Server
	log        logger.Logger
}

func New(addr string, ctrl *controller.Controller, log logger.Logger) *Service {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))

	var (
		userRouter   = user.New(ctrl)
		workerRouter = worker.New(ctrl)
	)
	router.Route("/v1", func(router chi.Router) {
		router.Route("/user", userRouter.Route())
		router.Route("/worker", workerRouter.Route())
	})

	return &Service{
		httpServer: &http.Server{
			Handler: router,
			Addr:    addr,
		},
		log: log,
	}
}

func (s *Service) Start() error {
	s.log.Infof("starting api service on %s", s.httpServer.Addr)
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}
	go func() {
		_ = s.httpServer.Serve(listener)
	}()
	return nil
}

func (s *Service) Stop() error {
	s.log.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), defaultCtxTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
