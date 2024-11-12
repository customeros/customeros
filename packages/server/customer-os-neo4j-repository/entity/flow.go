package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
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

	// Statistics

	Total        int64
	OnHold       int64
	Ready        int64
	Scheduled    int64
	InProgress   int64
	Completed    int64
	GoalAchieved int64
}

type FlowEntities []FlowEntity

type FlowActionEntity struct {
	DataLoaderKey
	Id         string
	ExternalId string
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Type string
	Data struct {
		//FLOW_START properties
		Entity      *string // CONTACT / ORGANIZATION / etc
		TriggerType *string

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

type FlowParticipantEntity struct {
	DataLoaderKey
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	EntityId   string
	EntityType model.EntityType

	Status          FlowParticipantStatus
	ScheduledAction *string
	ScheduledAt     *time.Time
}

type FlowParticipantEntities []FlowParticipantEntity

type FlowSenderEntity struct {
	DataLoaderKey
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	UserId *string
}

type FlowSenderEntities []FlowSenderEntity

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
	FlowActionTypeFlowEnd                   FlowActionType = "FLOW_END"
)

func GetFlowActionType(s string) FlowActionType {
	return FlowActionType(s)
}

type FlowParticipantStatus string

const (
	FlowParticipantStatusOnHold       FlowParticipantStatus = "ON_HOLD"
	FlowParticipantStatusReady        FlowParticipantStatus = "READY"
	FlowParticipantStatusScheduled    FlowParticipantStatus = "SCHEDULED"
	FlowParticipantStatusInProgress   FlowParticipantStatus = "IN_PROGRESS"
	FlowParticipantStatusCompleted    FlowParticipantStatus = "COMPLETED"
	FlowParticipantStatusGoalAchieved FlowParticipantStatus = "GOAL_ACHIEVED"
)

func GetFlowContactStatus(s string) FlowParticipantStatus {
	return FlowParticipantStatus(s)
}

type FlowExecutionSettingsEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	FlowId     string
	EntityId   string
	EntityType model.EntityType

	Mailbox *string
	UserId  *string
}

type FlowActionExecutionEntity struct {
	DataLoaderKey
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	FlowId        string
	ActionId      string
	ParticipantId string
	EntityId      string
	EntityType    model.EntityType

	// Scheduling Info
	ScheduledAt     time.Time
	ExecutedAt      *time.Time
	StatusUpdatedAt time.Time
	Status          FlowActionExecutionStatus

	//Config
	Mailbox   *string
	UserId    *string
	SocialUrl *string

	// Additional metadata
	Error *string // If execution fails, store the error message
}

type FlowActionExecutionEntities []FlowActionExecutionEntity

type FlowActionExecutionStatus string

const (
	FlowActionExecutionStatusScheduled     FlowActionExecutionStatus = "SCHEDULED"
	FlowActionExecutionStatusInProgress    FlowActionExecutionStatus = "IN_PROGRESS"
	FlowActionExecutionStatusSkipped       FlowActionExecutionStatus = "SKIPPED"
	FlowActionExecutionStatusSuccess       FlowActionExecutionStatus = "SUCCESS"
	FlowActionExecutionStatusTechError     FlowActionExecutionStatus = "TECH_ERROR"
	FlowActionExecutionStatusBusinessError FlowActionExecutionStatus = "BUSINESS_ERROR"
)

func GetFlowActionExecutionStatus(s string) FlowActionExecutionStatus {
	return FlowActionExecutionStatus(s)
}
