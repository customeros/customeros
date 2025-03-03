package postgres_entity

import (
	"time"
)

type WebSessionPageVisit struct {
	ID             string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	SessionID      string    `gorm:"column:session_id;type:uuid;" json:"sessionId"`
	Url            string    `gorm:"column:url;type:varchar(255);" json:"url"`
	Referer        string    `gorm:"column:referer;type:varchar(255);" json:"referer"`
	EntryTimestamp time.Time `gorm:"column:entry_timestamp;type:timestamp;" json:"entryTimestamp"`
	ExitTimestamp  time.Time `gorm:"column:exit_timestamp;type:timestamp;" json:"exitTimestamp"`
	Hostname       string    `gorm:"column:hostname;type:varchar(255);" json:"hostname"`
	Pathname       string    `gorm:"column:pathname;type:varchar(255);" json:"pathname"`
}

func (WebSessionPageVisit) TableName() string {
	return "web_session_page_visits"
}
