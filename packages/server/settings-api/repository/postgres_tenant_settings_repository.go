package repository

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type TenantSettingsRepository interface {
	FindForTenantName(tenantName string) helper.QueryResult
	Save(tenantSettings *entity.TenantSettings) helper.QueryResult
	SaveKeys(keys []entity.TenantAPIKey) error
	CheckKeysExist(tenantName string, keyName []string) (bool, error)
	DeleteKeys(keys []entity.TenantAPIKey) error
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
	var tenantSettings entity.TenantSettings

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
		err := r.db.Model(&entity.TenantAPIKey{}).
			Where(&entity.TenantAPIKey{TenantName: tenantName, Key: key}, "tenant_name", "key").Count(&rows).Error

		if err != nil {
			return false, fmt.Errorf("CheckKeysExist: %w", err)
		}
		if rows == 0 {
			exists = false
		}

	}
	return exists, nil
}

func (r *tenantSettingsRepo) SaveKeys(keys []entity.TenantAPIKey) error {

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

func (r *tenantSettingsRepo) DeleteKeys(keys []entity.TenantAPIKey) error {

	var deletedItem entity.TenantAPIKey
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

func (r *tenantSettingsRepo) Save(tenantSettings *entity.TenantSettings) helper.QueryResult {

	result := r.db.Save(tenantSettings)

	if result.Error != nil {
		return helper.QueryResult{Error: result.Error}
	}

	return helper.QueryResult{Result: tenantSettings}
}
