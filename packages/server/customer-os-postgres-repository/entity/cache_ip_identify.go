package postgres_entity

import "time"

type CacheIPIdentify struct {
	ID           string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	IPAddress    string    `gorm:"column:ip_address;type:varchar(255)" json:"ipAddress"`
	Domain       string    `gorm:"column:domain;type:varchar(255)" json:"domain"`
	LinkedinSlug string    `gorm:"column:linkedin_slug;type:varchar(255)" json:"linkedinSlug"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;" json:"updatedAt"`
	SnitcherData string    `gorm:"column:snitcher_data;type:text" json:"snitcherData"`
	SourceEmail  string    `gorm:"column:source_email;type:text" json:"sourceEmail"`
}

func (CacheIPIdentify) TableName() string {
	return "cache_ip_identify"
}
