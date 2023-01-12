package exporter

import (
	"fmt"

	"lostinsoba/ninhydrin/internal/model"
	"lostinsoba/ninhydrin/internal/monitoring/exporter/prometheus"
)

type Exporter interface {
	RegisterCounter(name string) (incr func(float64))
	RegisterGauge(name string) (set func(float64))
	Start()
}

func NewExporter(kind string, settings model.Settings, labels map[string]string) (Exporter, error) {
	switch kind {
	case prometheus.Kind:
		return prometheus.NewExporter(settings, labels)
	default:
		return nil, fmt.Errorf("unknown logger kind: %s", kind)
	}
}
