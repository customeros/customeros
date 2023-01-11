package entity

import "time"

type EventType string

const (
	WEB_CHAT EventType = "WEB_CHAT"
	EMAIL    EventType = "EMAIL"
	MESSAGE  EventType = "MESSAGE"
	VOICE    EventType = "VOICE"
)

type EventSubtype string

const (
	TEXT EventSubtype = "TEXT"
	FILE EventSubtype = "FILE"
)

type Source string

const (
	HUBSPOT  Source = "hubspot"
	ZENDESK  Source = "zendesk"
	MANUAL   Source = "manual"
	OPENLINE Source = "openline"
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
	ID             string `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	EventUUID      string `gorm:"type:uuid;default:gen_random_uuid()" json:"eventUuid"`
	TenantName     string `gorm:"column:tenant_name;type:varchar(50);NOT NULL" json:"tenantId" binding:"required"`
	ConversationId string `gorm:"column:conversation_id;type:varchar(50);NOT NULL" json:"conversationId" binding:"required"`
	//email
	Type EventType `gorm:"column:type;type:varchar(50);NOT NULL;" json:"type" binding:"required"`
	//thread id
	Subtype EventSubtype `gorm:"column:subtype;type:varchar(50);NOT NULL;" json:"subtype" binding:"required"`

	//used in websockets for web chat messages
	InitiatorUsername string `gorm:"column:initiator_username;type:varchar(50);NOT NULL" json:"initiatorUsername" binding:"required"`

	SenderId       string     `gorm:"column:sender_id;type:varchar(50);NOT NULL" json:"senderId" binding:"required"`
	SenderType     SenderType `gorm:"column:sender_type;type:varchar(50);NOT NULL" json:"senderType" binding:"required"`
	SenderUsername string     `gorm:"column:sender_username;type:varchar(50);NOT NULL" json:"senderUsername" binding:"required"`

	//vezi mai jos json
	Content    string `gorm:"column:content;type:text;NOT NULL;" json:"content" binding:"required"`
	Source     Source `gorm:"column:source;type:varchar(50);NOT NULL;" json:"source" binding:"required"`
	ExternalId string `gorm:"column:external_id;type:varchar(100);" json:"externalId" binding:"required"`

	Direction  Direction `gorm:"column:direction;type:varchar(10);NOT NULL;" json:"direction" binding:"required"`
	CreateDate time.Time `gorm:"column:created_at"`

	OriginalJson string `gorm:"column:direction;type:text;NOT NULL;" json:"direction" binding:"required"`
	//
	//Content: {
	//	"message": "Hello World",
	//	"fileId": "1234"
	//}
	//
	//Content: {
	//	to: "",
	//cc: "",
	//bcc: "",
	//	from: "",
	//	body: "",
	//	subject: "",
	//}
	//
	//Content: {
	//	callerId: "",
	//	calledId: "",
	//	callerPhoneNumber: "",
	//	calledPhoneNumber: "",
	//}
}

func (ConversationEvent) TableName() string {
	return "conversation_event"
}
