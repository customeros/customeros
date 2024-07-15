package entity

import "time"

type Tracking struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	UserId    string `gorm:"column:user_id;type:varchar(255);NOT NULL;" json:"user_id" binding:"required"`
	IP        string `gorm:"column:ip;type:varchar(255);" json:"ip" binding:"required"`
	EventType string `gorm:"column:event_type;type:varchar(255);" json:"event_type" binding:"required"`
	EventData string `gorm:"column:event_data;type:text;" json:"event_data" binding:"required"`
	Timestamp int    `gorm:"column:timestamp;type:text;" json:"timestamp" binding:"required"`

	Href             string `gorm:"column:href;type:varchar(255);" json:"href" binding:"href"`
	Origin           string `gorm:"column:origin;type:varchar(255);" json:"origin" binding:"origin"`
	Search           string `gorm:"column:search;type:varchar(255);" json:"search" binding:"search"`
	Hostname         string `gorm:"column:hostname;type:varchar(255);" json:"hostname" binding:"hostname"`
	Pathname         string `gorm:"column:pathname;type:varchar(255);" json:"pathname" binding:"pathname"`
	Referrer         string `gorm:"column:referrer;type:varchar(255);" json:"referrer" binding:"referrer"`
	UserAgent        string `gorm:"column:user_agent;type:varchar(255);" json:"user_agent" binding:"user_agent"`
	Language         string `gorm:"column:language;type:varchar(255);" json:"language" binding:"language"`
	CookiesEnabled   bool   `gorm:"column:cookies_enabled;type:boolean;" json:"cookies_enabled" binding:"cookies_enabled"`
	ScreenResolution string `gorm:"column:screen_resolution;type:varchar(255);" json:"screen_resolution" binding:"screen_resolution"`
}

func (Tracking) TableName() string {
	return "tracking"
}
