package scheduler

import "lostinsoba/ninhydrin/internal/monitoring/exporter"

type Metrics struct {
	statusesRefreshed func(float64)
}

func NewMetrics(exporter exporter.Exporter) *Metrics {
	return &Metrics{
		statusesRefreshed: exporter.RegisterCounter("statuses_refreshed"),
	}
}
