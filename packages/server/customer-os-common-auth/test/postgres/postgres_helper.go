package postgrest

import (
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

func InsertTenantAPIKey(db *gorm.DB, tenantName, key, value string) error {
	newRecord := &postgresEntity.GoogleServiceAccountKey{
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
