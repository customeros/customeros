package entity

import (
	"github.com/google/uuid"
	"time"
)

type SlackChannelNotificationWorkflow string

const (
	SlackChannelNotificationWorkflowWebsiteVisit   SlackChannelNotificationWorkflow = "REVEAL-AI-WEBSITE-VISIT"
	SlackChannelNotificationWorkflowMailstackReply SlackChannelNotificationWorkflow = "MAILSTACK-REPLY"
)

// send notifications to slack channels
// tenants can configure this for particular workflows
type SlackChannelNotification struct {
	ID        uuid.UUID                        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time                        `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp"`
	UpdatedAt time.Time                        `gorm:"column:updated_at;type:timestamp"`
	Tenant    string                           `gorm:"column:tenant;type:varchar(255);NOT NULL"`
	ChannelId string                           `gorm:"column:channel_id;type:varchar(255);NOT NULL"`
	Workflow  SlackChannelNotificationWorkflow `gorm:"column:workflow;type:varchar(255);"`
}

func (SlackChannelNotification) TableName() string {
	return "slack_channel_notification"
}

func (SlackChannelNotification) UniqueIndex() [][]string {
	return [][]string{
		{"Tenant", "ChannelId", "Workflow"},
	}
}
