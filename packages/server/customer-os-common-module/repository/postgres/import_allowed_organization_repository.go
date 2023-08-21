package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ImportAllowedOrganizationRepository interface {
	GetOrganizationsAllowedForImport(tenant string) ([]entity.ImportAllowedOrganization, error)
}

type importAllowedOrganizationRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewImportAllowedOrganizationRepository(gormDb *gorm.DB) ImportAllowedOrganizationRepository {
	return &importAllowedOrganizationRepositoryImpl{gormDb: gormDb}
}

func (repo *importAllowedOrganizationRepositoryImpl) GetOrganizationsAllowedForImport(tenant string) ([]entity.ImportAllowedOrganization, error) {
	var result []entity.ImportAllowedOrganization
	err := repo.gormDb.Model(&entity.ImportAllowedOrganization{}).Find(&result, "tenant = ? AND allowed = true", tenant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logrus.Errorf("error while getting import allowed organization: %v", err)
		return nil, err
	}
	return result, nil
}
