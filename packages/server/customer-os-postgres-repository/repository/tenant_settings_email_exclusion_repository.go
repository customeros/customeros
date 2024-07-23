package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TenantSettingsEmailExclusionRepository interface {
	GetExclusionList() ([]entity.TenantSettingsEmailExclusion, error)
}

type tenantSettingsEmailExclusionRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewEmailExclusionRepository(gormDb *gorm.DB) TenantSettingsEmailExclusionRepository {
	return &tenantSettingsEmailExclusionRepositoryImpl{gormDb: gormDb}
}

func (repo *tenantSettingsEmailExclusionRepositoryImpl) GetExclusionList() ([]entity.TenantSettingsEmailExclusion, error) {
	result := []entity.TenantSettingsEmailExclusion{}
	err := repo.gormDb.Find(&result).Limit(5000).Error

	if err != nil {
		logrus.Errorf("error while getting personal email provider list: %v", err)
		return nil, err
	}

	return result, nil
}
