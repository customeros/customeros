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

	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text;"`
	Active      bool   `gorm:"not null;default:false"`
}

type FlowSequence struct {
	BaseEntity

	FlowId string `gorm:"type:uuid;not null" json:"-"`

	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text;not null"`
	Active      bool   `gorm:"not null;default:false"`
}

func (FlowSequence) TableName() string {
	return "flow_sequence"
}

type FlowSequenceStep struct {
	BaseEntity

	SequenceId string `gorm:"type:uuid;not null" json:"-"`

	Active bool `gorm:"not null;default:false"`

	Order int    `gorm:"not null"`
	Type  string `gorm:"type:varchar(255);not null"`
	Name  string `gorm:"type:varchar(255);not null"`

	Text     *string `gorm:"type:varchar(255)"`
	Template *string `gorm:"type:varchar(255)"`
}

func (FlowSequenceStep) TableName() string {
	return "flow_sequence_step"
}

type FlowSequenceContact struct {
	BaseEntity

	SequenceId string `gorm:"type:uuid;not null" json:"-"`

	FirstName   *string
	LastName    *string
	Email       string `gorm:"not null"`
	LinkedinUrl *string
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
