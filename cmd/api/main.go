package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"lostinsoba/ninhydrin/api"
	"lostinsoba/ninhydrin/controller"
	"lostinsoba/ninhydrin/internal/config"
	"lostinsoba/ninhydrin/internal/monitoring"
	"lostinsoba/ninhydrin/internal/storage"
)

var (
	name      = "api"
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
	log := serviceMonitoring.NewLogger(cfg.Monitoring.Logger.Level)

	serviceStorage, err := storage.NewStorage(
		cfg.Storage.Kind,
		cfg.Storage.Settings,
	)
	if err != nil {
		log.Fatalf("failed to open storage: %s", err)
	}

	serviceController := controller.New(serviceStorage)

	service := api.New(cfg.API.Listen, serviceController, log)
	err = service.Start()
	if err != nil {
		log.Fatalf("failed to start api service: %s", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	err = service.Stop()
	if err != nil {
		log.Fatalf("failed to gracefully stop api: %s", err)
	}
}
