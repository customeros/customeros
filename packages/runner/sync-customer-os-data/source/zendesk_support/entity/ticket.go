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
	Status               string       `gorm:"column:status"`
	Type                 string       `gorm:"column:type"`
	Tags                 pgtype.JSONB `gorm:"column:tags;type:jsonb"`
	CustomFieldsAsJson   pgtype.JSONB `gorm:"column:custom_fields;type:jsonb"`
	Priority             string       `gorm:"column:priority"`
	Description          string       `gorm:"column:description"`
	CollaboratorIds      pgtype.JSONB `gorm:"column:collaborator_ids;type:jsonb"`
	FollowerIds          pgtype.JSONB `gorm:"column:follower_ids;type:jsonb"`
	SubmitterId          int64        `gorm:"column:submitter_id"`
	RequesterId          int64        `gorm:"column:requester_id"`
	AssigneeId           int64        `gorm:"column:assignee_id"`
}

type Tickets []Ticket

func (Ticket) TableName() string {
	return "tickets"
}
