package scheduler

import (
	"context"
	"time"

	"lostinsoba/ninhydrin/controller"
	"lostinsoba/ninhydrin/internal/monitoring/logger"
)

type Scheduler struct {
	ctrl     *controller.Controller
	interval time.Duration
	cancel   context.CancelFunc

	log logger.Logger
}

func NewScheduler(ctrl *controller.Controller, interval time.Duration, log logger.Logger) *Scheduler {
	return &Scheduler{ctrl: ctrl, interval: interval, log: log}
}

func (s *Scheduler) Run(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		s.log.Debugf("started refreshing task statuses every %s", s.interval)
		for {
			select {
			case <-ticker.C:
				tasksUpdated, err := s.ctrl.RefreshTaskStatuses(ctx)
				if err != nil {
					s.log.Errorf("failed to refresh task statuses: %s", err)
				} else {
					s.log.Debugf("updated %d tasks statuses", tasksUpdated)
				}
			case <-ctx.Done():
				s.log.Debugf("context cancelled, stopping...")
				return
			}
		}
	}()
}
