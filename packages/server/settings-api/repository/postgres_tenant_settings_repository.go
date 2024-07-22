package repository

import (
	"fmt"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type TenantSettingsRepository interface {
	FindForTenantName(tenantName string) helper.QueryResult
	Save(tenantSettings *postgresentity.TenantSettings) helper.QueryResult
	CheckKeysExist(tenantName string, keyName []string) (bool, error)
}

type tenantSettingsRepo struct {
	db *gorm.DB
}

func NewTenantSettingsRepository(db *gorm.DB) TenantSettingsRepository {
	return &tenantSettingsRepo{
		db: db,
	}
}

func (r *tenantSettingsRepo) FindForTenantName(tenantName string) helper.QueryResult {
	var tenantSettings postgresentity.TenantSettings

	err := r.db.
		Where("tenant_name = ?", tenantName).
		First(&tenantSettings).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return helper.QueryResult{Error: err}
	}
	if err == gorm.ErrRecordNotFound {
		return helper.QueryResult{Result: nil}
	}

	return helper.QueryResult{Result: tenantSettings}
}

func (r *tenantSettingsRepo) CheckKeysExist(tenantName string, keyName []string) (bool, error) {
	var rows int64
	exists := true
	for _, key := range keyName {
		log.Printf("CheckKeysExist: %s, %s", tenantName, key)
		err := r.db.Model(&postgresentity.GoogleServiceAccountKey{}).
			Where(&postgresentity.GoogleServiceAccountKey{TenantName: tenantName, Key: key}, "tenant_name", "key").Count(&rows).Error

		if err != nil {
			return false, fmt.Errorf("CheckKeysExist: %w", err)
		}
		if rows == 0 {
			exists = false
		}

	}
	return exists, nil
}

func (r *tenantSettingsRepo) SaveKeys(keys []postgresentity.GoogleServiceAccountKey) error {

	for _, key := range keys {
		result := r.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tenant_name"}, {Name: "key"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"value": key.Value}),
		}).Save(&key)
		if result.Error != nil {
			return fmt.Errorf("SaveKeys: %w", result.Error)
		}
	}
	return nil
}

func (r *tenantSettingsRepo) DeleteKeys(keys []postgresentity.GoogleServiceAccountKey) error {

	var deletedItem postgresentity.GoogleServiceAccountKey
	for _, key := range keys {
		log.Printf("DeleteKeys: %s, %s", key.TenantName, key.Key)
		err := r.db.
			Where(&key, "tenant_name", "key").Delete(&deletedItem).Error

		if err != nil {
			return fmt.Errorf("DeleteKeys: %w", err)
		}
	}
	return nil
}

func (r *tenantSettingsRepo) Save(tenantSettings *postgresentity.TenantSettings) helper.QueryResult {

	result := r.db.Save(tenantSettings)

	if result.Error != nil {
		return helper.QueryResult{Error: result.Error}
	}

	return helper.QueryResult{Result: tenantSettings}
}
