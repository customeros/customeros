package agent_capability

type CapabilityParams struct {
	Completed          bool     `json:"completed"`
	Domain             string   `json:"domain"`
	ExecutionValidated bool     `json:"executionValidated"`
	Hostname           string   `json:"hostname"`
	IcpFitRationale    string   `json:"icpFitRationale"`
	IPAddress          string   `json:"ipAddress"`
	IsICPFit           string   `json:"isIcpFit"`
	IsNewCompanyVisit  bool     `json:"isNewCompanyVisit"`
	IsNewPersonVisit   bool     `json:"isNewPersonVisit"`
	LinkedInSlug       string   `json:"linkedinSlug"`
	Message            *string  `json:"message,omitempty"`
	OrganizationID     string   `json:"organizationId"`
	PageViews          []string `json:"pageViews"`
	Referrer           string   `json:"referrer"`
	SessionDuration    string   `json:"sessionDuration"`
	SessionID          string   `json:"sessionId"`
	SlackNotification  string   `json:"slackNotification"`
	VisitorID          string   `json:"visitorId"`
}
