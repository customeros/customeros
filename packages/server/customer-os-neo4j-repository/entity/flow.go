package entity

import (
	"time"
)

type FlowEntity struct {
	DataLoaderKey
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Description string
	Status      FlowStatus
}

type FlowEntities []FlowEntity

type FlowActionEntity struct {
	DataLoaderKey
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Name  string
	Index int64

	Status FlowActionStatus

	ActionType FlowActionType
	ActionData FlowActionData
}

type FlowActionEntities []FlowActionEntity

type FlowActionData struct {
	Wait                      *FlowActionDataWait
	EmailNew                  *FlowActionDataEmail
	EmailReply                *FlowActionDataEmail
	LinkedinConnectionRequest *FlowActionDataLinkedinConnectionRequest
	LinkedinMessage           *FlowActionDataLinkedinMessage
}

type FlowActionDataEmail struct {
	ReplyToId    *string
	Subject      string
	BodyTemplate string
}

type FlowActionDataLinkedinConnectionRequest struct {
	MessageTemplate string
}

type FlowActionDataLinkedinMessage struct {
	MessageTemplate string
}

type FlowActionDataWait struct {
	Minutes int64
}

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

type FlowSequenceStatus string

const (
	FlowSequenceStatusInactive FlowSequenceStatus = "INACTIVE"
	FlowSequenceStatusActive   FlowSequenceStatus = "ACTIVE"
	FlowSequenceStatusPaused   FlowSequenceStatus = "PAUSED"
	FlowSequenceStatusArchived FlowSequenceStatus = "ARCHIVED"
)

func GetFlowSequenceStatus(s string) FlowSequenceStatus {
	return FlowSequenceStatus(s)
}

type FlowActionStatus string

const (
	FlowActionStatusInactive FlowActionStatus = "INACTIVE"
	FlowActionStatusActive   FlowActionStatus = "ACTIVE"
	FlowActionStatusPaused   FlowActionStatus = "PAUSED"
	FlowActionStatusArchived FlowActionStatus = "ARCHIVED"
)

func GetFlowActionStatus(s string) FlowActionStatus {
	return FlowActionStatus(s)
}

type FlowActionType string

const (
	FlowActionTypeWait                      FlowActionType = "WAIT"
	FlowActionTypeEmailNew                  FlowActionType = "EMAIL_NEW"
	FlowActionTypeEmailReply                FlowActionType = "EMAIL_REPLY"
	FlowActionTypeLinkedinConnectionRequest FlowActionType = "LINKEDIN_CONNECTION_REQUEST"
	FlowActionTypeLinkedinMessage           FlowActionType = "LINKEDIN_MESSAGE"
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
