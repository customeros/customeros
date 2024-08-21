package entity

import "time"

type CacheEmailScrubby struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	CheckedAt time.Time `gorm:"column:checked_at;type:timestamp" json:"checkedAt"`
	Email     string    `gorm:"column:email;type:varchar(255);NOT NULL" json:"domain"`
	Status    string    `gorm:"column:status;type:varchar(255)" json:"status"`
}

func (CacheEmailScrubby) TableName() string {
	return "cache_email_scrubby"
}
