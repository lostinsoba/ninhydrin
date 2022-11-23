package scheduler

import "lostinsoba/ninhydrin/internal/monitoring/exporter"

type metrics struct {
	statusesRefreshed func(float64)
}

func newMetrics(exporter exporter.Exporter) *metrics {
	return &metrics{
		statusesRefreshed: exporter.RegisterCounter("statuses_refreshed"),
	}
}
