package postgrest

import (
	commonEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"gorm.io/gorm"
)

func InsertTenantAPIKey(db *gorm.DB, tenantName, key, value string) error {
	newRecord := &commonEntity.GoogleServiceAccountKey{
		TenantName: tenantName,
		Key:        key,
		Value:      value,
	}

	result := db.Create(newRecord)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
