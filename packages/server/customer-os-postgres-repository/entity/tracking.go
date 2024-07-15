package entity

import "time"

type Tracking struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	Tenant string `gorm:"column:tenant;type:varchar(255);" json:"tenant"`

	UserId    string `gorm:"column:user_id;type:varchar(255);NOT NULL;" json:"user_id"`
	IP        string `gorm:"column:ip;type:varchar(255);" json:"ip" `
	EventType string `gorm:"column:event_type;type:varchar(255);" json:"event_type"`
	EventData string `gorm:"column:event_data;type:text;" json:"event_data"`
	Timestamp int    `gorm:"column:timestamp;type:integer;" json:"timestamp"`

	Href             string `gorm:"column:href;type:varchar(255);" json:"href"`
	Origin           string `gorm:"column:origin;type:varchar(255);" json:"origin"`
	Search           string `gorm:"column:search;type:varchar(255);" json:"search"`
	Hostname         string `gorm:"column:hostname;type:varchar(255);" json:"hostname"`
	Pathname         string `gorm:"column:pathname;type:varchar(255);" json:"pathname"`
	Referrer         string `gorm:"column:referrer;type:varchar(255);" json:"referrer"`
	UserAgent        string `gorm:"column:user_agent;type:text;" json:"user_agent"`
	Language         string `gorm:"column:language;type:varchar(255);" json:"language"`
	CookiesEnabled   bool   `gorm:"column:cookies_enabled;type:boolean;" json:"cookies_enabled"`
	ScreenResolution string `gorm:"column:screen_resolution;type:varchar(255);" json:"screen_resolution"`
}

func (Tracking) TableName() string {
	return "tracking"
}
