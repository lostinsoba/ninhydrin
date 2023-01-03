package model

const (
	TaskStatusIdle       TaskStatus = "idle"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusTimeout    TaskStatus = "timeout"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusDone       TaskStatus = "done"
)

func GetTaskCaptureStatuses() []string {
	return []string{
		string(TaskStatusIdle),
		string(TaskStatusFailed),
	}
}

type TaskStatus string

type Task struct {
	ID          string
	Timeout     int64
	RetriesLeft int
	UpdatedAt   int64
	Status      TaskStatus
}
