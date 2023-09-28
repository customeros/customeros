package entity

import "github.com/google/uuid"

const GSUITE_SERVICE_PRIVATE_KEY = "GSUITE_SERVICE_PRIVATE_KEY"
const GSUITE_SERVICE_EMAIL_ADDRESS = "GSUITE_SERVICE_EMAIL_ADDRESS"

type TenantAPIKey struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantName string    `gorm:"size:255;not null;index:idx_tenant_api_keys"`
	Key        string    `gorm:"size:255;not null;index:idx_tenant_api_keys"`
	Value      string    `gorm:"type:text"`
}

func (TenantAPIKey) TableName() string {
	return "tenant_api_keys"
}
