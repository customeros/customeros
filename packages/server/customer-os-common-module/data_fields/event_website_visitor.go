package data_fields

type WebsiteVisitEvent struct {
	SessionID string `json:"sessionId"`
	Tenant    string `json:"tenant"`
	IPAddress string `json:"ipAddress"`
	VisitorID string `json:"visitorId"`
}

func (f WebsiteVisitEvent) Type() string {
	return "WebsiteVisitEvent"
}
