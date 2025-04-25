package enum

type OutboxStatus string

const (
	OutboxPending    OutboxStatus = "pending"
	OutboxProcessing OutboxStatus = "processing"
	OutboxCompleted  OutboxStatus = "completed"
	OutboxFailed     OutboxStatus = "failed"
)
