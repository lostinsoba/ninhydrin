package monitoring

import (
	"lostinsoba/ninhydrin/internal/monitoring/exporter"
	"lostinsoba/ninhydrin/internal/monitoring/logger"
)

const (
	labelService   = "service"
	labelVersion   = "version"
	labelGitCommit = "git_commit"
)

type Monitoring struct {
	labels map[string]string
}

func NewMonitoring(service, version, gitCommit string) *Monitoring {
	return &Monitoring{
		labels: map[string]string{
			labelService:   service,
			labelVersion:   version,
			labelGitCommit: gitCommit,
		},
	}
}

func (m *Monitoring) NewLogger(kind string, settings map[string]string) logger.Logger {
	return logger.NewLogger(kind, settings, m.labels)
}

func (m *Monitoring) NewExporter(kind string, settings map[string]string) (exporter.Exporter, error) {
	return exporter.NewExporter(kind, settings, m.labels)
}
