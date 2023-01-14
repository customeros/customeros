package entity

import (
	"time"
)

type EmailContent struct {
	Html    string   `json:"html"`
	Subject string   `json:"subject"`
	From    string   `json:"from"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
}

type EventType string

const (
	EMAIL EventType = "EMAIL"
)

type SenderType string

const (
	CONTACT SenderType = "CONTACT"
	USER    SenderType = "USER"
)

type Direction string

const (
	INBOUND  Direction = "INBOUND"
	OUTBOUND Direction = "OUTBOUND"
)

type ConversationEvent struct {
	ID             string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	EventUUID      string    `gorm:"type:uuid;default:gen_random_uuid()" json:"eventUuid"`
	TenantName     string    `gorm:"column:tenant_name;type:varchar(50);NOT NULL" json:"tenantId" binding:"required"`
	ConversationId string    `gorm:"column:conversation_id;type:varchar(50);NOT NULL" json:"conversationId" binding:"required"`
	Type           EventType `gorm:"column:type;type:varchar(50);NOT NULL;" json:"type" binding:"required"`
	//thread id
	Subtype           string     `gorm:"column:subtype;type:varchar(50);NOT NULL;" json:"subtype" binding:"required"`
	SenderId          string     `gorm:"column:sender_id;type:varchar(50);NOT NULL" json:"senderId" binding:"required"`
	SenderType        SenderType `gorm:"column:sender_type;type:varchar(50);NOT NULL" json:"senderType" binding:"required"`
	Source            string     `gorm:"column:source;type:varchar(50);NOT NULL;" json:"source" binding:"required"`
	ExternalId        string     `gorm:"column:external_id;type:varchar(100);" json:"externalId" binding:"required"`
	Direction         Direction  `gorm:"column:direction;type:varchar(10);NOT NULL;" json:"direction" binding:"required"`
	CreateDate        time.Time  `gorm:"column:created_at"`
	Content           string     `gorm:"column:content;type:text;NOT NULL;" json:"content" binding:"required"`
	InitiatorUsername string     `gorm:"column:initiator_username;type:varchar(50);NOT NULL" json:"initiatorUsername" binding:"required"`
	SenderUsername    string     `gorm:"column:sender_username;type:varchar(50);NOT NULL" json:"senderUsername" binding:"required"`
}

func (ConversationEvent) TableName() string {
	return "conversation_event"
}
