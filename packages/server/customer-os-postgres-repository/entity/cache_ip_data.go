package entity

import "time"

type CacheIpData struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Ip        string    `gorm:"column:ip;type:varchar(255);NOT NULL;index:idx_cache_ip_data_ip,unique" json:"ip"`
}

func (CacheIpData) TableName() string {
	return "cache_ip_data"
}
