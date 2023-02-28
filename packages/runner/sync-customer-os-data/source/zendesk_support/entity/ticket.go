package entity

import (
	"github.com/jackc/pgtype"
	"time"
)

type Ticket struct {
	Id                   int64        `gorm:"column:id"`
	AirbyteAbId          string       `gorm:"column:_airbyte_ab_id"`
	AirbyteTicketsHashid string       `gorm:"column:_airbyte_tickets_hashid"`
	CreateDate           time.Time    `gorm:"column:created_at"`
	UpdatedDate          time.Time    `gorm:"column:updated_at"`
	Url                  string       `gorm:"column:url"`
	Subject              string       `gorm:"column:subject"`
	CollaboratorIds      pgtype.JSONB `gorm:"column:collaborator_ids;type:jsonb"`
}

type Tickets []Ticket

func (Ticket) TableName() string {
	return "tickets"
}
