package entity

import "time"

type TechLimit struct {
	ID         uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	Key        string    `gorm:"column:key;type:varchar(255);NOT NULL;index:idx_key_unique,unique" json:"key"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	UsageCount int       `gorm:"column:usage_count;type:int;NOT NULL" json:"limit"`
}

func (TechLimit) TableName() string {
	return "tech_limit"
}
