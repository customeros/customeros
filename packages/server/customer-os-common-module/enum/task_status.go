package enum

type TaskStatus string

const (
	Todo       TaskStatus = "TODO"
	InProgress TaskStatus = "IN_PROGRESS"
	Done       TaskStatus = "DONE"
)

var AllTaskStatuses = []TaskStatus{
	Todo,
	InProgress,
	Done,
}

func DecodeTaskStatus(s string) TaskStatus {
	if IsValidTaskStatus(s) {
		return TaskStatus(s)
	}
	return Todo // Default to TODO if invalid
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
