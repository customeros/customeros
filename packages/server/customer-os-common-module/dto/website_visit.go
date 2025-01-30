package dto

type WebsiteVisit struct {
	SessionID string `json:"sessionId"`
	Tenant    string `json:"tenant"`
	IPAddress string `json:"ipAddress"`
	VisitorID string `json:"visitorId"`
	Hostname  string `json:"hostname"`
}

func (f WebsiteVisit) Type() string {
	return "WebsiteVisit"
}
