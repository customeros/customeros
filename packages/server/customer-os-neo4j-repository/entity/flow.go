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

type FlowSequenceEntity struct {
	DataLoaderKey
	Id          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Description string
	Status      FlowSequenceStatus

	////Schedule
	//ActiveDaysString string `gorm:"type:varchar(255)" json:"-"`
	//
	//ActiveTimeWindowStart    string `gorm:"type:varchar(255)" json:"activeTimeWindowStart"` //09:00:00
	//ActiveTimeWindowEnd      string `gorm:"type:varchar(255)" json:"activeTimeWindowEnd"`   //09:00:00
	//PauseOnHolidays          bool   `json:"pauseOnHolidays"`
	//RespectRecipientTimezone bool   `json:"respectRecipientTimezone"`
	//
	//MinutesDelayBetweenEmails int `json:"minutesDelayBetweenEmails"`
	//
	//EmailsPerMailboxPerHour int `json:"emailsPerMailboxPerHour"`
	//EmailsPerMailboxPerDay  int `json:"emailsPerMailboxPerDay"`

}

type FlowSequenceEntities []FlowSequenceEntity

type FlowSequenceStepEntity struct {
	DataLoaderKey
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string

	Status FlowSequenceStepStatus

	Type    FlowSequenceStepType
	Subtype *FlowSequenceStepSubtype
	Body    string
}

type FlowSequenceStepEntities []FlowSequenceStepEntity

type FlowSequenceContactEntity struct {
	DataLoaderKey
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	ContactId string
	EmailId   string
}

type FlowSequenceContactEntities []FlowSequenceContactEntity

type FlowSequenceSenderEntity struct {
	DataLoaderKey
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Mailbox string
}

type FlowSequenceSenderEntities []FlowSequenceSenderEntity

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

type FlowSequenceStepStatus string

const (
	FlowSequenceStepStatusInactive FlowSequenceStepStatus = "INACTIVE"
	FlowSequenceStepStatusActive   FlowSequenceStepStatus = "ACTIVE"
	FlowSequenceStepStatusPaused   FlowSequenceStepStatus = "PAUSED"
	FlowSequenceStepStatusArchived FlowSequenceStepStatus = "ARCHIVED"
)

type FlowSequenceStepType string

const (
	FlowSequenceStepTypeEmail    FlowSequenceStepType = "EMAIL"
	FlowSequenceStepTypeLinkedin FlowSequenceStepType = "LINKEDIN"
)

type FlowSequenceStepSubtype string

const (
	FlowSequenceStepSubtypeLinkedinConnectionRequest FlowSequenceStepSubtype = "LINKEDIN_CONNECTION_REQUEST"
	FlowSequenceStepSubtypeLinkedinMessage           FlowSequenceStepSubtype = "LINKEDIN_MESSAGE"
)
