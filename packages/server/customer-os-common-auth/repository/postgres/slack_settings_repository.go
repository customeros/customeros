package repository

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"gorm.io/gorm"
)

type SlackSettingsRepository interface {
	Get(tenant string) (*entity.SlackSettingsEntity, error)
	Save(slackSettings entity.SlackSettingsEntity) (*entity.SlackSettingsEntity, error)

	Delete(tenant string) error
}

type slackSettingsRepository struct {
	db *gorm.DB
}

func NewSlackSettingsRepository(db *gorm.DB) SlackSettingsRepository {
	return &slackSettingsRepository{
		db: db,
	}
}

func (repo *slackSettingsRepository) Get(tenant string) (*entity.SlackSettingsEntity, error) {
	var existing *entity.SlackSettingsEntity
	err := repo.db.Find(&existing, "tenant_name = ?", tenant).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	if existing == nil {
		return nil, nil
	} else if existing.TenantName == "" {
		return nil, nil
	}
	return existing, nil
}

func (repo *slackSettingsRepository) Save(slackSettings entity.SlackSettingsEntity) (*entity.SlackSettingsEntity, error) {
	result := repo.db.Save(&slackSettings)
	if result.Error != nil {
		return nil, fmt.Errorf("saving slack settings failed: %w", result.Error)
	}
	return &slackSettings, nil
}

func (repo *slackSettingsRepository) Delete(tenant string) error {
	existing, err := repo.Get(tenant)
	if err != nil {
		return err
	}

	err = repo.db.Delete(&existing).Error
	if err != nil {
		return err
	}

	return nil
}
