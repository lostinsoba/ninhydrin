package prometheus

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const Kind = "prometheus"

const (
	settingNamespace = "namespace"
	settingListen    = "listen"
)

type Exporter struct {
	addr     string
	registry *Registry
}

func NewExporter(settings map[string]string, labels map[string]string) (*Exporter, error) {
	namespaceStr, ok := settings[settingNamespace]
	if !ok {
		return nil, fmt.Errorf("%s setting not present", settingNamespace)
	}
	listenStr, ok := settings[settingListen]
	if !ok {
		return nil, fmt.Errorf("%s setting not present", settingListen)
	}
	return &Exporter{
		addr:     listenStr,
		registry: NewRegistry(namespaceStr, labels),
	}, nil
}

func (exporter *Exporter) RegisterCounter(name string) func(float64) {
	return exporter.registry.RegisterCounter(name)
}

func (exporter *Exporter) RegisterGauge(name string) func(float64) {
	return exporter.registry.RegisterGauge(name)
}

func (exporter *Exporter) Start() {
	opts := promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	}
	handler := promhttp.HandlerFor(exporter.registry.internal, opts)
	mux := http.NewServeMux()
	mux.Handle("/metrics", handler)
	httpServer := &http.Server{Addr: exporter.addr, Handler: mux}
	go func() {
		_ = httpServer.ListenAndServe()
	}()
}
