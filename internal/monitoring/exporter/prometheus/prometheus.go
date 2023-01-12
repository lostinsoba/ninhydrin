package prometheus

import (
	"net/http"

	"lostinsoba/ninhydrin/internal/model"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const Kind = "prometheus"

const (
	settingNamespace = "namespace"
	settingListen    = "listen"
)

type Exporter struct {
	addr     string
	registry *registry
}

func NewExporter(settings model.Settings, labels map[string]string) (*Exporter, error) {
	namespace, err := settings.ReadStr(settingNamespace)
	if err != nil {
		return nil, err
	}
	listen, err := settings.ReadStr(settingListen)
	if err != nil {
		return nil, err
	}
	return &Exporter{
		addr:     listen,
		registry: newRegistry(namespace, labels),
	}, nil
}

func (exporter *Exporter) RegisterCounter(name string) func(float64) {
	return exporter.registry.registerCounter(name)
}

func (exporter *Exporter) RegisterGauge(name string) func(float64) {
	return exporter.registry.registerGauge(name)
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
