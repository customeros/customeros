package entity

import "time"

type CacheEmailEnrow struct {
	ID            string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	RequestID     string    `gorm:"column:request_id;type:varchar(255);NOT NULL" json:"requestId"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamp;" json:"updatedAt"`
	Email         string    `gorm:"column:email;type:varchar(255);NOT NULL" json:"email"`
	Qualification string    `gorm:"column:qualification;type:varchar(255)" json:"result"`
	Data          string    `gorm:"column:data;type:text" json:"data"`
}

func (CacheEmailEnrow) TableName() string {
	return "cache_email_enrow"
}

type EnrowResponseBody struct {
	Email         string `json:"email"`
	Id            string `json:"id"`
	Qualification string `json:"qualification"`
}
