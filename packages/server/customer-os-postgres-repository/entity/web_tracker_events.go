package postgres_entity

import "time"

type WebTrackerEvents struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);" json:"tenant"`
	SessionID string    `gorm:"column:session_id;type:uuid;" json:"sessionId"`
	VisitorID string    `gorm:"column:visitor_id;type:varchar(255);NOT NULL;" json:"visitorId"`
	IP        string    `gorm:"column:ip;type:varchar(255);" json:"ip" `
	EventType string    `gorm:"column:event_type;type:varchar(255);" json:"eventType"`
	EventData string    `gorm:"column:event_data;type:text;" json:"eventData"`
	Timestamp time.Time `gorm:"column:timestamp;type:timestamp;" json:"timestamp"`

	Href             string `gorm:"column:href;type:varchar(1000);" json:"href"`
	Origin           string `gorm:"column:origin;type:varchar(255);" json:"origin"`
	Search           string `gorm:"column:search;type:varchar(1000);" json:"search"`
	Hostname         string `gorm:"column:hostname;type:varchar(255);" json:"hostname"`
	Pathname         string `gorm:"column:pathname;type:varchar(255);" json:"pathname"`
	Referrer         string `gorm:"column:referrer;type:varchar(2000);" json:"referrer"`
	UserAgent        string `gorm:"column:user_agent;type:text;" json:"userAgent"`
	Language         string `gorm:"column:language;type:varchar(255);" json:"language"`
	CookiesEnabled   bool   `gorm:"column:cookies_enabled;type:boolean;" json:"cookiesEnabled"`
	ScreenResolution string `gorm:"column:screen_resolution;type:varchar(255);" json:"screenResolution"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (WebTrackerEvents) TableName() string {
	return "web_tracker_events"
}
