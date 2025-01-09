package entity

import "time"

type WebSession struct {
	ID             string     `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant         string     `gorm:"column:tenant;type:varchar(255);" json:"tenant"`
	VisitorID      string     `gorm:"column:visitor_id;type:varchar(255);NOT NULL;" json:"visitorId"`
	IP             string     `gorm:"column:ip;type:varchar(255);" json:"ip"`
	Domain         *string    `gorm:"column:domain;type:varchar(255);" json:"domain"`
	StartTime      time.Time  `gorm:"column:start_time;type:timestamp;" json:"startTime"`
	LastActivity   time.Time  `gorm:"column:last_activity;type:timestamp;" json:"lastActivity"`
	EndTime        *time.Time `gorm:"column:end_time;type:timestamp;" json:"endTime"`
	IsActive       bool       `gorm:"column:is_active;type:boolean;default:true" json:"isActive"`
	ExitDetected   bool       `gorm:"column:exit_detected;type:boolean;default:false" json:"exitDetected"`
	EventPublished bool       `gorm:"column:event_published;type:boolean;default:false" json:"eventPublished"`
	CreatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
}

func (WebSession) TableName() string {
	return "web_session"
}
