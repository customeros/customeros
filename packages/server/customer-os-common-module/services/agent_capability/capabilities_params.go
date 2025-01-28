package agent_capability

type CapabilityParams struct {
	ExecutionValidated bool     `json:"executionValidated"`
	Completed          bool     `json:"completed"`
	SessionID          string   `json:"sessionId"`
	IPAddress          string   `json:"ipAddress"`
	VisitorID          string   `json:"visitorId"`
	Hostname           string   `json:"hostname"`
	Domain             string   `json:"domain"`
	LinkedInSlug       string   `json:"linkedinSlug"`
	OrganizationID     string   `json:"organizationId"`
	IsNewCompanyVisit  bool     `json:"isNewCompanyVisit"`
	IsNewPersonVisit   bool     `json:"isNewPersonVisit"`
	PageViews          []string `json:"pageViews"`
	SessionDuration    string   `json:"sessionDuration"`
	SlackNotification  string   `json:"slackNotification"`
	Referrer           string   `json:"referrer"`
	Message            *string  `json:"message,omitempty"`
}
