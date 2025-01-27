package dto

type SupportEventSource string

const (
	SupportWebVisit SupportEventSource = "support_web_visit"
)

type SupportEvent struct {
	Tenant       string
	Domain       string
	WebSessionID string
	Source       SupportEventSource
}

func (d SupportEvent) Type() string {
	return "SupportEvent"
}
