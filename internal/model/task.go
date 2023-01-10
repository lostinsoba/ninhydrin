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

func IsValidTaskStatus(status TaskStatus) bool {
	validStatuses := getAllTaskStatuses()
	for _, validStatus := range validStatuses {
		if validStatus == status {
			return true
		}
	}
	return false
}

func getAllTaskStatuses() []TaskStatus {
	return []TaskStatus{
		TaskStatusIdle,
		TaskStatusInProgress,
		TaskStatusTimeout,
		TaskStatusFailed,
		TaskStatusDone,
	}
}

type TaskStatus string

type Task struct {
	ID          string
	NamespaceID string
	Timeout     int64
	RetriesLeft int
	UpdatedAt   int64
	Status      TaskStatus
}
