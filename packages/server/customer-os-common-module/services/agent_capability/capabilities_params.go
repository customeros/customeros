package agent_capability

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type CapabilityParams struct {
	Completed                bool        `json:"completed"`
	DisqualificationCriteria string      `json:"disqualificationCriteria"`
	Domain                   string      `json:"domain"`
	ExecutionValidated       bool        `json:"executionValidated"`
	Hostname                 string      `json:"hostname"`
	IcpFit                   enum.IcpFit `json:"isIcpFit"`
	IcpFitRationale          string      `json:"icpFitRationale"`
	IPAddress                string      `json:"ipAddress"`
	IsNewCompanyVisit        bool        `json:"isNewCompanyVisit"`
	IsNewPersonVisit         bool        `json:"isNewPersonVisit"`
	LinkedInSlug             string      `json:"linkedinSlug"`
	MarkdownEventID          string      `json:"markdownEventId"`
	Message                  string      `json:"message"`
	OrganizationID           string      `json:"organizationId"`
	PageViews                []string    `json:"pageViews"`
	QualificationCriteria    string      `json:"qualificationCriteria"`
	Referrer                 string      `json:"referrer"`
	SessionDuration          string      `json:"sessionDuration"`
	SessionID                string      `json:"sessionId"`
	SlackNotification        string      `json:"slackNotification"`
	VisitorID                string      `json:"visitorId"`

	CompanyName         string `json:"companyName"`
	PrimaryDomain       string `json:"primaryDomain"`
	CompanyDescriptions string `json:"companyDescriptions"`
	IndustryNAICSName   string `json:"industryName"`
	YearCompanyFounded  string `json:"yearCompanyFounded"`
	EmployeeCount       int64  `json:"employeeCount"`
	CompanyCity         string `json:"companyCity"`
	CompanyRegion       string `json:"companyRegion"`
	CompanyCountryA2    string `json:"companyCountry"`
}
