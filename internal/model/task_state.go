package model

const (
	TaskStatusTimeout    TaskStatus = "timeout"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusIdle       TaskStatus = "idle"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

func GetTaskCaptureStatuses() []string {
	return []string{
		string(TaskStatusFailed),
		string(TaskStatusIdle),
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
		TaskStatusTimeout,
		TaskStatusFailed,
		TaskStatusIdle,
		TaskStatusInProgress,
		TaskStatusDone,
	}
}

type TaskStatus string

type TaskState struct {
	TaskID      string
	RetriesLeft int
	UpdatedAt   int64
	Status      TaskStatus
}