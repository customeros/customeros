package data_fields

// Nil fields wil be skipped from update
type WebsiteVisitEvent struct {
	ID             string  `json:"id"`
	Tenant         string  `json:"tenant"`
	Domain         string  `json:"domain"`
	Referrer       string  `json:"referrer"`
	OrganizationID *string `json:"organizationId"`
}

func (f WebsiteVisitEvent) Type() string {
	return "WebsiteVisitEvent"
}
