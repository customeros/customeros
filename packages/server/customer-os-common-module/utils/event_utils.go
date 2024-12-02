package utils

type EventCompletedDetails struct {
	Create bool
	Update bool
	Delete bool
}

func NewEventCompletedDetails() *EventCompletedDetails {
	return &EventCompletedDetails{}
}

func (ecd *EventCompletedDetails) WithCreate() *EventCompletedDetails {
	ecd.Create = true
	return ecd
}

func (ecd *EventCompletedDetails) WithUpdate() *EventCompletedDetails {
	ecd.Update = true
	return ecd
}

func (ecd *EventCompletedDetails) WithDelete() *EventCompletedDetails {
	ecd.Delete = true
	return ecd
}
