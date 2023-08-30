package entity

import "time"

type PersonalIntegration struct {
	ID         uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	TenantName string    `gorm:"column:tenant_name;type:varchar(255);NOT NULL" json:"tenantName" binding:"required"`
	Name       string    `gorm:"column:name;type:varchar(255);NOT NULL" json:"name" binding:"required"`
	Email      string    `gorm:"column:email;type:varchar(255);NOT NULL" json:"email" binding:"required"`
	Secret     string    `gorm:"column:key;type:varchar(255);NOT NULL;index:idx_key,unique" json:"key" binding:"required"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`
}

func (PersonalIntegration) TableName() string {
	return "personal_integrations"
}

func (PersonalIntegration) UniqueIndex() [][]string {
	return [][]string{
		{"TenantName", "Name", "Email"},
	}
}
