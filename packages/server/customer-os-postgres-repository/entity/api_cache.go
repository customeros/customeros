package entity

import (
	"time"
)

type ApiCache struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL;default:CURRENT_TIMESTAMP"`
	Tenant    string    `gorm:"column:tenant;type:varchar(100);NOT NULL;" json:"tenant" binding:"required"`
	Type      string    `gorm:"column:type;type:varchar(255);NOT NULL;" json:"type" binding:"required"`
	Data      string    `gorm:"column:data;type:text;NOT NULL;" json:"data" binding:"required"`
}

func (ApiCache) TableName() string {
	return "api_cache"
}
