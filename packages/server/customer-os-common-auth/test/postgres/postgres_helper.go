package postgrest

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"gorm.io/gorm"
)

func InsertTenantAPIKey(db *gorm.DB, tenantName, key, value string) error {
	newRecord := &entity.TenantAPIKey{
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
