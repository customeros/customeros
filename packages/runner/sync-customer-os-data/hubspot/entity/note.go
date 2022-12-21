package entity

import (
	"github.com/jackc/pgtype"
	"time"
)

type Note struct {
	Id                  string       `gorm:"column:id"`
	AirbyteAbId         string       `gorm:"column:_airbyte_ab_id"`
	AirbyteNotesHashid  string       `gorm:"column:_airbyte_engagements_notes_hashid"`
	CreateDate          time.Time    `gorm:"column:createdat"`
	UpdatedDate         time.Time    `gorm:"column:updatedat"`
	ContactsExternalIds pgtype.JSONB `gorm:"column:contacts;type:jsonb"`
}

type Notes []Note

func (Note) TableName() string {
	return "engagements_notes"
}
