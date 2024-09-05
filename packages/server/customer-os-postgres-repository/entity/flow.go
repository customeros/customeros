package entity

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

	Name        string     `gorm:"type:varchar(255);" json:"name"`
	Description string     `gorm:"type:text;" json:"description"`
	Status      FlowStatus `gorm:"type:varchar(100);not null;" json:"active"`
}

type FlowSequence struct {
	BaseEntity
	Tenant string `gorm:"not null" json:"-"`

	FlowId string `gorm:"type:uuid;not null" json:"-"`

	Name        string             `gorm:"type:varchar(255);" json:"name"`
	Description string             `gorm:"type:text;" json:"description"`
	Status      FlowSequenceStatus `gorm:"type:varchar(100);not null;" json:"status"`

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

//func (s Flow) MarshalJSON() ([]byte, error) {
//	type Alias Flow
//	return json.Marshal(&struct {
//		*Alias
//		ActiveDays []string `json:"activeDays"`
//	}{
//		Alias:      (*Alias)(&s),
//		ActiveDays: utils.StringToSlice(s.ActiveDaysString),
//	})
//}
//
//func (s FlowSequence) MarshalJSON() ([]byte, error) {
//	type Alias FlowSequence
//	return json.Marshal(&struct {
//		*Alias
//		Personas []string `json:"personas"`
//	}{
//		Alias:    (*Alias)(&s),
//		Personas: utils.StringToSlice(s.PersonasString),
//	})
//}

type FlowSequenceStep struct {
	BaseEntity
	Tenant string `gorm:"not null" json:"-"`

	SequenceId string `gorm:"type:uuid;not null" json:"-"`

	Name string `gorm:"type:varchar(255);" json:"name"`

	Status FlowSequenceStepStatus `gorm:"type:varchar(100);not null;" json:"status"`

	Type    FlowSequenceStepType     `gorm:"type:varchar(100);not null" json:"type"`
	Subtype *FlowSequenceStepSubtype `gorm:"type:varchar(100);" json:"subtype"`
	Body    string                   `gorm:"type:text;not null" json:"body"`
}

type FlowSequenceContact struct {
	BaseEntity
	Tenant string `gorm:"not null" json:"-"`

	SequenceId string `gorm:"type:uuid;not null" json:"-"`

	ContactId   string  `gorm:"not null" json:"contactId"`
	EmailId     string  `gorm:"not null" json:"emailId"`
	LinkedinUrl *string `json:"linkedinUrl"`
}

type FlowSequenceSender struct {
	BaseEntity
	Tenant string `gorm:"not null" json:"-"`

	SequenceId string `gorm:"type:uuid;not null" json:"-"`

	MailboxId string `gorm:"not null"`
}

func (Flow) TableName() string {
	return "flow"
}
func (FlowSequence) TableName() string {
	return "flow_sequence"
}
func (FlowSequenceStep) TableName() string {
	return "flow_sequence_step"
}
func (FlowSequenceContact) TableName() string {
	return "flow_sequence_contact"
}
func (FlowSequenceSender) TableName() string {
	return "flow_sequence_sender"
}

type FlowStatus string

const (
	FlowStatusInactive FlowStatus = "INACTIVE"
	FlowStatusActive   FlowStatus = "ACTIVE"
	FlowStatusPaused   FlowStatus = "PAUSED"
	FlowStatusArchived FlowStatus = "ARCHIVED"
)

type FlowSequenceStatus string

const (
	FlowSequenceStatusInactive FlowSequenceStatus = "INACTIVE"
	FlowSequenceStatusActive   FlowSequenceStatus = "ACTIVE"
	FlowSequenceStatusPaused   FlowSequenceStatus = "PAUSED"
	FlowSequenceStatusArchived FlowSequenceStatus = "ARCHIVED"
)

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
