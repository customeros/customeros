package entity

import "github.com/google/uuid"

type TenantAPIKey struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantName string    `gorm:"size:255;not null;index:idx_tenant_api_keys"`
	Key        string    `gorm:"size:255;not null;index:idx_tenant_api_keys"`
	Value      string    `gorm:"type:text"`
}

func (TenantAPIKey) TableName() string {
	return "tenant_api_keys"
}
