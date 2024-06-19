package entity

import "time"

type ExternalAppKeys struct {
	ID         uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	App        string    `gorm:"column:app;type:varchar(255);NOT NULL;index:idx_external_app_key_unique,unique" json:"app"`
	AppKey     string    `gorm:"column:app_key;type:varchar(255);NOT NULL;index:idx_external_app_key_unique,unique" json:"appKey"`
	Group1     string    `gorm:"column:group1;type:varchar(255);index:idx_external_app_key_unique,unique" json:"group1"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	UsageCount int       `gorm:"column:usage_count;type:int;NOT NULL" json:"usageCount"`
}

func (ExternalAppKeys) TableName() string {
	return "external_app_keys"
}
