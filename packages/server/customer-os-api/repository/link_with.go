package repository

type LinkedWith string
type LinkedNature string

const (
	LINKED_WITH_INTERACTION_SESSION LinkedWith = "InteractionSession"
	LINKED_WITH_INTERACTION_EVENT   LinkedWith = "InteractionEvent"
	LINKED_WITH_MEETING             LinkedWith = "Meeting"
	LINKED_WITH_NOTE                LinkedWith = "Note"

	LINKED_NATURE_RECORDING LinkedNature = "Recording"
)
