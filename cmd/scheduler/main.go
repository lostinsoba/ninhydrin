package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"lostinsoba/ninhydrin/controller"
	"lostinsoba/ninhydrin/internal/config"
	"lostinsoba/ninhydrin/internal/monitoring"
	"lostinsoba/ninhydrin/internal/storage"
	"lostinsoba/ninhydrin/scheduler"
)

var (
	name      = "scheduler"
	version   = "unknown"
	gitCommit = "unknown"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	serviceMonitoring := monitoring.NewMonitoring(name, version, gitCommit)
	log := serviceMonitoring.NewLogger(
		cfg.Monitoring.Logger.Kind,
		cfg.Monitoring.Logger.Settings,
	)

	exporter, err := serviceMonitoring.NewExporter(
		cfg.Monitoring.Exporter.Kind,
		cfg.Monitoring.Exporter.Settings,
	)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %s", err)
	}
	exporter.Start()

	serviceStorage, err := storage.NewStorage(
		cfg.Storage.Kind,
		cfg.Storage.Settings,
	)
	if err != nil {
		log.Fatalf("failed to open storage: %s", err)
	}

	serviceController := controller.New(serviceStorage)

	ctx := context.Background()

	metrics := scheduler.NewMetrics(exporter)
	service := scheduler.NewScheduler(serviceController, cfg.Scheduler.Interval, metrics, log)
	service.Run(ctx)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Info("shutting down")
}
