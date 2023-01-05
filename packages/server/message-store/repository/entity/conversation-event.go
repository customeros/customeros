package entity

import "time"

type eventType string

const (
	EMAIL      eventType = "EMAIL"
	MESSAGE    eventType = "MESSAGE"
	PHONE_CALL eventType = "PHONE_CALL"
)

type source string

const (
	HUBSPOT source = "HUBSPOT"
	ZENDESK source = "ZENDESK"
	MANUAL  source = "MANUAL"
	SYSTEM  source = "SYSTEM"
)

type senderType string

const (
	CONTACT senderType = "CONTACT"
	USER    senderType = "USER"
)

type direction string

const (
	INBOUND  direction = "INBOUND"
	OUTBOUND direction = "OUTBOUND"
)

type ConversationEvent struct {
	ID             string     `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	EventUUID      string     `gorm:"type:uuid;default:gen_random_uuid()" json:"eventUuid"`
	TenantId       string     `gorm:"column:tenant_id;type:varchar(50);NOT NULL" json:"tenantId" binding:"required"`
	ConversationId string     `gorm:"column:conversation_id;type:varchar(50);NOT NULL" json:"conversationId" binding:"required"`
	Type           eventType  `gorm:"column:type;type:varchar(50);NOT NULL;" json:"type" binding:"required"`
	SenderId       string     `gorm:"column:sender_id;type:varchar(50);NOT NULL" json:"senderId" binding:"required"`
	SenderType     senderType `gorm:"column:sender_type;type:varchar(50);NOT NULL" json:"senderType" binding:"required"`
	Content        string     `gorm:"column:content;type:text;NOT NULL;" json:"content" binding:"required"`
	Source         source     `gorm:"column:source;type:varchar(50);NOT NULL;" json:"source" binding:"required"`
	Direction      direction  `gorm:"column:direction;type:varchar(10);NOT NULL;" json:"direction" binding:"required"`
	CreateDate     time.Time  `gorm:"column:created_at"`
}

func (ConversationEvent) TableName() string {
	return "conversation_event"
}
