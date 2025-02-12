package agent_capability

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type CapabilityParams struct {
	ActionItems                    []string    `json:"actionItems"`
	ContactEmails                  []string    `json:"contactEmails"`
	ContactIDs                     []string    `json:"contactIds"`
	CompanyName                    string      `json:"companyName"`
	CompanyCity                    string      `json:"companyCity"`
	CompanyCountryA2               string      `json:"companyCountry"`
	CompanyDescriptions            string      `json:"companyDescriptions"`
	CompanyNeedsHelp               bool        `json:"companyNeedsHelp"`
	CompanyRegion                  string      `json:"companyRegion"`
	DisqualificationCriteria       string      `json:"disqualificationCriteria"`
	Domain                         string      `json:"domain"`
	EmployeeCount                  int64       `json:"employeeCount"`
	ExecutionValidated             bool        `json:"executionValidated"`
	HelpNeeded                     []string    `json:"helpNeeded"`
	Hostname                       string      `json:"hostname"`
	IcpFit                         enum.IcpFit `json:"isIcpFit"`
	IcpFitRationale                string      `json:"icpFitRationale"`
	IndustryNAICSName              string      `json:"industryName"`
	IPAddress                      string      `json:"ipAddress"`
	IsNewCompanyVisit              bool        `json:"isNewCompanyVisit"`
	IsNewPersonVisit               bool        `json:"isNewPersonVisit"`
	LinkedInSlug                   string      `json:"linkedinSlug"`
	MarkdownEventID                string      `json:"markdownEventId"`
	MeetingContent                 string      `json:"meetingContent"`
	MeetingParticipantEmails       []string    `json:"meetingParticipantEmails"`
	MeetingParticipantEmailsTenant []string    `json:"meetingParticipantEmailsTenant"`
	MeetingRecordingUrl            string      `json:"meetingRecordingUrl"`
	MeetingSource                  string      `json:"meetingSource"`
	MeetingSummary                 string      `json:"meetingSummary"`
	MeetingTimestamp               time.Time   `json:"meetingTimestamp"`
	MeetingTitle                   string      `json:"MeetingTitle"`
	Message                        string      `json:"message"`
	OrganizationID                 string      `json:"organizationId"`
	OrganizationIDs                []string    `json:"organizationIds"`
	PageViews                      []string    `json:"pageViews"`
	PrimaryDomain                  string      `json:"primaryDomain"`
	QualificationCriteria          string      `json:"qualificationCriteria"`
	Referrer                       string      `json:"referrer"`
	SessionDuration                string      `json:"sessionDuration"`
	SlackNotification              string      `json:"slackNotification"`
	UniquePageViews                []string    `json:"uniquePageViews"`
	VisitorID                      string      `json:"visitorId"`
	WebSessionID                   string      `json:"webSessionId"`
	YearCompanyFounded             string      `json:"yearCompanyFounded"`
}
