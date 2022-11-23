package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Registry struct {
	internal  *prometheus.Registry
	labels    prometheus.Labels
	namespace string
}

func NewRegistry(namespace string, labels map[string]string) *Registry {
	internalRegistry := prometheus.NewRegistry()
	registry := &Registry{
		internal:  internalRegistry,
		labels:    prometheus.Labels(labels),
		namespace: namespace,
	}
	registry.RegisterGauge("build_info")(1)
	registry.registerProcessCollector()
	return registry
}

func (registry *Registry) RegisterGauge(name string) func(float64) {
	opts := prometheus.GaugeOpts{
		Name:        name,
		ConstLabels: registry.labels,
		Namespace:   registry.namespace,
	}
	gauge := prometheus.NewGauge(opts)
	registry.internal.MustRegister(gauge)
	return gauge.Set
}

func (registry *Registry) RegisterCounter(name string) func(float64) {
	opts := prometheus.CounterOpts{
		Name:        name,
		ConstLabels: registry.labels,
		Namespace:   registry.namespace,
	}
	counter := prometheus.NewCounter(opts)
	registry.internal.MustRegister(counter)
	return counter.Add
}

func (registry *Registry) registerProcessCollector() {
	opts := collectors.ProcessCollectorOpts{
		Namespace:    registry.namespace,
		ReportErrors: false,
	}
	collector := collectors.NewProcessCollector(opts)
	registry.internal.MustRegister(collector)
}
