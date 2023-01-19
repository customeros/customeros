package entity

import (
	"github.com/jackc/pgtype"
	"time"
)

type Email struct {
	Id                  string       `gorm:"column:id"`
	AirbyteAbId         string       `gorm:"column:_airbyte_ab_id"`
	AirbyteEmailsHashid string       `gorm:"column:_airbyte_engagements_emails_hashid"`
	CreateDate          time.Time    `gorm:"column:createdat"`
	UpdatedDate         time.Time    `gorm:"column:updatedat"`
	ContactsExternalIds pgtype.JSONB `gorm:"column:contacts;type:jsonb"`
}

type Emails []Email

func (Email) TableName() string {
	return "engagements_emails"
}
