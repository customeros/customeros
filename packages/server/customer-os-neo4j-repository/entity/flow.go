package entity

import (
	"time"
)

type FlowEntity struct {
	DataLoaderKey
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string

	Nodes string
	Edges string

	Status FlowStatus
}

type FlowEntities []FlowEntity

type FlowActionEntity struct {
	DataLoaderKey
	Id         string
	ExternalId string
	Type       string
	Data       struct {
		Entity *string // CONTACT / ORGANIZATION / etc

		WaitBefore int64 // in minutes

		Action FlowActionType

		//ActionData fields below

		//Email
		Subject      *string
		BodyTemplate *string

		//Linkedin
		MessageTemplate *string
	} `json:"data"`
	Json string
}

type FlowActionEntities []FlowActionEntity

type FlowContactEntity struct {
	DataLoaderKey
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	ContactId string
}

type FlowContactEntities []FlowContactEntity

type FlowActionSenderEntity struct {
	DataLoaderKey
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Mailbox *string
	UserId  *string
}

type FlowActionSenderEntities []FlowActionSenderEntity

type FlowStatus string

const (
	FlowStatusInactive FlowStatus = "INACTIVE"
	FlowStatusActive   FlowStatus = "ACTIVE"
	FlowStatusPaused   FlowStatus = "PAUSED"
	FlowStatusArchived FlowStatus = "ARCHIVED"
)

func GetFlowStatus(s string) FlowStatus {
	return FlowStatus(s)
}

type FlowActionType string

const (
	FlowActionTypeFlowStart                 FlowActionType = "FLOW_START"
	FlowActionTypeEmailNew                  FlowActionType = "EMAIL_NEW"
	FlowActionTypeEmailReply                FlowActionType = "EMAIL_REPLY"
	FlowActionTypeLinkedinConnectionRequest FlowActionType = "LINKEDIN_CONNECTION_REQUEST"
	FlowActionTypeLinkedinMessage           FlowActionType = "LINKEDIN_MESSAGE"
	FlowActionTypeWait                      FlowActionType = "WAIT"
	FlowActionTypeFlowEnd                   FlowActionType = "FLOW_END"
)

func GetFlowActionType(s string) FlowActionType {
	return FlowActionType(s)
}

type FlowExecutionSettingsEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	FlowId   string
	EntityId string

	Mailbox *string
}

type FlowActionExecutionEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	FlowId    string
	ActionId  string
	ContactId string

	// Scheduling Info
	ScheduledAt time.Time
	ExecutedAt  *time.Time
	Status      FlowActionExecutionStatus

	// Execution details for email
	Subject *string
	Body    *string
	From    *string
	To      []string
	Cc      []string
	Bcc     []string
	Mailbox *string

	// Additional metadata
	Error *string // If execution fails, store the error message
}

type FlowActionExecutionStatus string

const (
	FlowActionExecutionStatusPending FlowActionExecutionStatus = "PENDING"
	FlowActionExecutionStatusSuccess FlowActionExecutionStatus = "SUCCESS"
	FlowActionExecutionStatusError   FlowActionExecutionStatus = "ERROR"
)

func GetFlowActionExecutionStatus(s string) FlowActionExecutionStatus {
	return FlowActionExecutionStatus(s)
}
