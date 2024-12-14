package enum

type Status string

const (
	StatusError          Status = "error"
	StatusPartialSuccess Status = "partial success"
	StatusProcessing     Status = "processing"
	StatusSuccess        Status = "success"
	StatusWarning        Status = "warning"
)
