package entity

type SequenceVariable struct {
	BaseEntity

	Name  string `gorm:"type:varchar(255);not null"`
	Value string `gorm:"type:varchar(255);not null"`
}

func (SequenceVariable) TableName() string {
	return "sequence_variable"
}

type Sequence struct {
	BaseEntity

	Tenant string `gorm:"not null" json:"-"`
	T      Tenant `gorm:"foreignKey:Tenant;references:name;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;" json:"-"`

	Name string `gorm:"type:varchar(255);not null"`

	Enabled bool `gorm:"not null;default:false"`
}

func (Sequence) TableName() string {
	return "sequence"
}

type SequenceStep struct {
	BaseEntity

	SequenceId string `gorm:"type:uuid;not null" json:"-"`

	Order int    `gorm:"not null"`
	Type  string `gorm:"type:varchar(255);not null"`
	Name  string `gorm:"type:varchar(255);not null"`

	Text     *string `gorm:"type:varchar(255)"`
	Template *string `gorm:"type:varchar(255)"`
}

func (SequenceStep) TableName() string {
	return "sequence_step"
}

type SequenceContact struct {
	BaseEntity

	SequenceId string `gorm:"type:uuid;not null" json:"-"`

	FirstName   *string
	LastName    *string
	Email       string `gorm:"not null"`
	LinkedinUrl *string
}

func (SequenceContact) TableName() string {
	return "sequence_contact"
}

type SequenceMailbox struct {
	BaseEntity

	SequenceId string `gorm:"type:uuid;not null" json:"-"`

	MailboxId string `gorm:"not null"`
}

func (SequenceMailbox) TableName() string {
	return "sequence_mailbox"
}
