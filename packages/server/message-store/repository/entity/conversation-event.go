package entity

import "time"

type EventType string

const (
	WEB_CHAT EventType = "WEB_CHAT"
	EMAIL    EventType = "EMAIL"
	MESSAGE  EventType = "MESSAGE"
	VOICE    EventType = "VOICE"
)

type Source string

const (
	HUBSPOT Source = "HUBSPOT"
	ZENDESK Source = "ZENDESK"
	MANUAL  Source = "MANUAL"
	SYSTEM  Source = "SYSTEM"
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
	ID             string     `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	EventUUID      string     `gorm:"type:uuid;default:gen_random_uuid()" json:"eventUuid"`
	TenantId       string     `gorm:"column:tenant_id;type:varchar(50);NOT NULL" json:"tenantId" binding:"required"`
	ConversationId string     `gorm:"column:conversation_id;type:varchar(50);NOT NULL" json:"conversationId" binding:"required"`
	Type           EventType  `gorm:"column:type;type:varchar(50);NOT NULL;" json:"type" binding:"required"`
	SenderId       string     `gorm:"column:sender_id;type:varchar(50);NOT NULL" json:"senderId" binding:"required"`
	SenderType     SenderType `gorm:"column:sender_type;type:varchar(50);NOT NULL" json:"senderType" binding:"required"`
	Content        string     `gorm:"column:content;type:text;NOT NULL;" json:"content" binding:"required"`
	Source         Source     `gorm:"column:source;type:varchar(50);NOT NULL;" json:"source" binding:"required"`
	Direction      Direction  `gorm:"column:direction;type:varchar(10);NOT NULL;" json:"direction" binding:"required"`
	CreateDate     time.Time  `gorm:"column:created_at"`
}

func (ConversationEvent) TableName() string {
	return "conversation_event"
}
