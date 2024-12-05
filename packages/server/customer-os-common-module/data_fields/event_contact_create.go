package data_fields

type ContactCreateEvent struct {
	Email       string
	LinkedinURL string
}

func (c ContactCreateEvent) Type() string {
	return "ContactCreateEvent"
}
