package enum

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "TODO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusDone       TaskStatus = "DONE"
)

var AllTaskStatuses = []TaskStatus{
	TaskStatusTodo,
	TaskStatusInProgress,
	TaskStatusDone,
}

func DecodeTaskStatus(s string) TaskStatus {
	if IsValidTaskStatus(s) {
		return TaskStatus(s)
	}
	return TaskStatusTodo // Default to TODO if invalid
}

func IsValidTaskStatus(s string) bool {
	for _, status := range AllTaskStatuses {
		if status == TaskStatus(s) {
			return true
		}
	}
	return false
}

func (t TaskStatus) String() string {
	return string(t)
}
