package entity

import "time"

type TenantSettingsOpportunityStage struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`

	Visible bool   `gorm:"column:visible;value:boolean;NOT NULL" json:"visible"`
	Value   string `gorm:"column:name;value:varchar(255);NOT NULL" json:"value"`
	Order   int    `gorm:"column:idx;type:int;NOT NULL" json:"order"`
	Label   string `gorm:"column:label;type:varchar(255);NOT NULL" json:"label"`
}

func (TenantSettingsOpportunityStage) TableName() string {
	return "tenant_settings_opportunity_stage"
}
