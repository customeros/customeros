package entity

import (
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type FlowSequenceStepTemplateVariable struct {
	BaseEntity

	Name  string `gorm:"type:varchar(255);not null"`
	Value string `gorm:"type:varchar(255);not null"`
}

func (FlowSequenceStepTemplateVariable) TableName() string {
	return "flow_sequence_step_template_variable"
}

type Flow struct {
	BaseEntity
	Tenant string `gorm:"not null" json:"-"`

	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text;" json:"description"`
	Active      bool   `gorm:"not null;default:false" json:"active"`

	//Schedule
	ActiveDaysString string `gorm:"type:varchar(255)" json:"-"`

	ActiveTimeWindowStart    string `gorm:"type:varchar(255)" json:"activeTimeWindowStart"` //09:00:00
	ActiveTimeWindowEnd      string `gorm:"type:varchar(255)" json:"activeTimeWindowEnd"`   //09:00:00
	PauseOnHolidays          bool   `json:"pauseOnHolidays"`
	RespectRecipientTimezone bool   `json:"respectRecipientTimezone"`

	MinutesDelayBetweenEmails int `json:"minutesDelayBetweenEmails"`

	EmailsPerMailboxPerHour int `json:"emailsPerMailboxPerHour"`
	EmailsPerMailboxPerDay  int `json:"emailsPerMailboxPerDay"`
}

func (s Flow) MarshalJSON() ([]byte, error) {
	type Alias Flow
	return json.Marshal(&struct {
		*Alias
		ActiveDays []string `json:"activeDays"`
	}{
		Alias:      (*Alias)(&s),
		ActiveDays: utils.StringToSlice(s.ActiveDaysString),
	})
}

type FlowSequence struct {
	BaseEntity

	FlowId string `gorm:"type:uuid;not null" json:"-"`

	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text;not null" json:"description"`
	Active      bool   `gorm:"not null;default:false" json:"active"`

	PersonasString string `gorm:"type:text" json:"-"` //used to store days in DB
}

func (s FlowSequence) MarshalJSON() ([]byte, error) {
	type Alias FlowSequence
	return json.Marshal(&struct {
		*Alias
		Personas []string `json:"personas"`
	}{
		Alias:    (*Alias)(&s),
		Personas: utils.StringToSlice(s.PersonasString),
	})
}

func (FlowSequence) TableName() string {
	return "flow_sequence"
}

type FlowSequenceStep struct {
	BaseEntity

	SequenceId string `gorm:"type:uuid;not null" json:"-"`

	Active bool `gorm:"not null;default:false" json:"active"`

	Order int    `gorm:"not null" json:"order"`
	Type  string `gorm:"type:varchar(255);not null" json:"type"`
	Name  string `gorm:"type:varchar(255);not null" json:"name"`

	Text     *string `gorm:"type:varchar(255)" json:"text"`
	Template *string `gorm:"type:varchar(255)" json:"template"`
}

func (FlowSequenceStep) TableName() string {
	return "flow_sequence_step"
}

type FlowSequenceContact struct {
	BaseEntity

	SequenceId string `gorm:"type:uuid;not null" json:"-"`

	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	Email       string  `gorm:"not null" json:"email"`
	LinkedinUrl *string `json:"linkedinUrl"`
}

func (FlowSequenceContact) TableName() string {
	return "flow_sequence_contact"
}

type FlowSequenceSender struct {
	BaseEntity

	SequenceId string `gorm:"type:uuid;not null" json:"-"`

	MailboxId string `gorm:"not null"`
}

func (FlowSequenceSender) TableName() string {
	return "flow_sequence_sender"
}
