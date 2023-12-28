package entity

type IssueData struct {
	BaseData
	Subject              string                  `json:"subject,omitempty"`
	Status               string                  `json:"status,omitempty"`
	Priority             string                  `json:"priority,omitempty"`
	Description          string                  `json:"description,omitempty"`
	Collaborators        []ReferencedParticipant `json:"collaborators,omitempty"`
	Followers            []ReferencedParticipant `json:"followers,omitempty"`
	Assignee             ReferencedUser          `json:"userAssignee,omitempty"`
	Reporter             ReferencedParticipant   `json:"reporter,omitempty"`
	Submitter            ReferencedParticipant   `json:"submitter,omitempty"`
	OrganizationRequired bool                    `json:"organizationRequired,omitempty"`
}
