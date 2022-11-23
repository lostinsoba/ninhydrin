package exporter

import (
	"fmt"

	"lostinsoba/ninhydrin/internal/monitoring/exporter/prometheus"
)

type Exporter interface {
	RegisterCounter(name string) func(float64)
	RegisterGauge(name string) func(float64)
	Start()
}

func NewExporter(kind string, settings map[string]string, labels map[string]string) (Exporter, error) {
	switch kind {
	case prometheus.Kind:
		return prometheus.NewExporter(settings, labels)
	default:
		return nil, fmt.Errorf("unknown logger kind: %s", kind)
	}
}
