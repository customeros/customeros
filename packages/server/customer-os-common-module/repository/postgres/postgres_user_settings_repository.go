package repository

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	"gorm.io/gorm"
)

type UserSettingsRepository interface {
	GetByUserName(username string) helper.QueryResult
	Save(userSettings *entity.UserSettingsEntity) error
	Delete(username string) error
}

type userSettingsRepo struct {
	db *gorm.DB
}

func NewUserSettingsRepository(db *gorm.DB) UserSettingsRepository {
	return &userSettingsRepo{
		db: db,
	}
}

func (userSettingsRepo userSettingsRepo) GetByUserName(username string) helper.QueryResult {
	var userSettingsEntity entity.UserSettingsEntity

	err := userSettingsRepo.db.
		Where("tenant_name = ?", username).
		First(&userSettingsEntity).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return helper.QueryResult{Error: err}
	}
	if err == gorm.ErrRecordNotFound {
		return helper.QueryResult{Result: nil}
	}

	return helper.QueryResult{Result: userSettingsEntity}
}

func (userSettingsRepo userSettingsRepo) Save(userSettings *entity.UserSettingsEntity) error {

	result := userSettingsRepo.db.Save(userSettings)
	if result.Error != nil {
		return fmt.Errorf("SaveKeys: %w", result.Error)
	}

	return nil
}

func (userSettingsRepo userSettingsRepo) Delete(username string) error {
	//TODO implement me
	panic("implement me")
}
