package enum

type EnrichStatus string

const (
	EnrichNotStarted EnrichStatus = "NOT_STARTED"
	EnrichInProgress EnrichStatus = "IN_PROGRESS"
	EnrichCompleted  EnrichStatus = "COMPLETED"
	EnrichError      EnrichStatus = "ERROR"
)

func (s EnrichStatus) String() string {
	return string(s)
}
