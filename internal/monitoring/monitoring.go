package monitoring

import "lostinsoba/ninhydrin/internal/monitoring/logger"

const (
	labelService   = "service"
	labelVersion   = "version"
	labelGitCommit = "git_commit"
)

type Monitoring struct {
	labels map[string]string
}

func NewMonitoring(service, version, gitCommit string) *Monitoring {
	return &Monitoring{labels: map[string]string{
		labelService:   service,
		labelVersion:   version,
		labelGitCommit: gitCommit,
	}}
}

func (m *Monitoring) NewLogger(level string) logger.Logger {
	return logger.NewLogger(level, m.labels)
}
