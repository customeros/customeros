package model

type IssueDataFields struct {
	GroupId                   *string
	Subject                   string
	Description               string
	Status                    string
	Priority                  string
	ReportedByOrganizationId  *string
	SubmittedByOrganizationId *string
	SubmittedByUserId         *string
}
