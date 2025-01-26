package postgres_entity

import (
	"time"

	"github.com/lib/pq"
)

type WebSession struct {
	ID                    string         `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant                string         `gorm:"column:tenant;type:varchar(255);index:idx_tenant" json:"tenant"`
	VisitorID             string         `gorm:"column:visitor_id;type:varchar(255);NOT NULL;index:idx_visitor_id" json:"visitorId"`
	IP                    string         `gorm:"column:ip;type:varchar(255);" json:"ip"`
	Hostname              string         `gorm:"column:hostname;type:varchar(255);" json:"hostname"`
	Domain                *string        `gorm:"column:domain;type:varchar(255);index:idx_domain" json:"domain"`
	Referrer              *string        `gorm:"column:referrer;type:varchar(255);" json:"referrer"`
	QueryParams           *string        `gorm:"column:query_params;type:varchar(255);" json:"queryParams"`
	UniquePageViews       pq.StringArray `gorm:"column:unique_page_views;type:text[];" json:"uniquePageViews"`
	StartTime             time.Time      `gorm:"column:start_time;type:timestamp;" json:"startTime"`
	LastEventType         string         `gorm:"column:last_event_type;type:varchar(255);" json:"lastEventType"`
	LastActivity          time.Time      `gorm:"column:last_activity;type:timestamp;" json:"lastActivity"`
	EndTime               *time.Time     `gorm:"column:end_time;type:timestamp;" json:"endTime"`
	IsActive              bool           `gorm:"column:is_active;type:boolean;default:true" json:"isActive"`
	DetectedExit          bool           `gorm:"column:detected_exit;type:boolean;default:false" json:"detectedExit"`
	PublishedEvent        bool           `gorm:"column:published_event;type:boolean;default:false" json:"publishedEvent"`
	SentSlackNotification *time.Time     `gorm:"column:sent_slack_notification;type:timestamp;" json:"sentSlackNotification"`
	IntentSignals         int8           `gorm:"column:intent_signals;type:smallint;default:0" json:"intentSignals"`
	CreatedAt             time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
}

func (WebSession) TableName() string {
	return "web_session"
}
