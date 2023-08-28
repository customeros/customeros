package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ImportAllowedOrganizationRepository interface {
	SaveOrganizationAllowedForImport(importAllowedOrganization entity.ImportAllowedOrganization) error
	GetOrganizationsAllowedForImport(tenant string) ([]entity.ImportAllowedOrganization, error)
}

type importAllowedOrganizationRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewImportAllowedOrganizationRepository(gormDb *gorm.DB) ImportAllowedOrganizationRepository {
	return &importAllowedOrganizationRepositoryImpl{gormDb: gormDb}
}

func (repo *importAllowedOrganizationRepositoryImpl) SaveOrganizationAllowedForImport(importAllowedOrganization entity.ImportAllowedOrganization) error {
	var existing entity.ImportAllowedOrganization

	if importAllowedOrganization.ID != "" {
		err := repo.gormDb.Model(&entity.ImportAllowedOrganization{}).First(&existing, "id = ?", importAllowedOrganization.ID).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				existing = entity.ImportAllowedOrganization{}
			} else {
				logrus.Errorf("error while getting import allowed organization: %v", err)
				return err
			}
		}
	}

	if existing.ID == "" && importAllowedOrganization.Tenant != "" && importAllowedOrganization.Domain != "" && importAllowedOrganization.Name != "" {
		err := repo.gormDb.Model(&entity.ImportAllowedOrganization{}).First(&existing, "tenant = ? AND domain = ? AND name = ?", importAllowedOrganization.Tenant, importAllowedOrganization.Domain, importAllowedOrganization.Name).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				existing = entity.ImportAllowedOrganization{}
			} else {
				logrus.Errorf("error while getting import allowed organization: %v", err)
				return err
			}
		}

	}

	if existing.ID == "" {
		existing.Tenant = importAllowedOrganization.Tenant
		existing.Source = importAllowedOrganization.Source
		existing.AppSource = importAllowedOrganization.AppSource
		existing.Name = importAllowedOrganization.Name
		existing.Domain = importAllowedOrganization.Domain
	}

	existing.Allowed = importAllowedOrganization.Allowed

	err := repo.gormDb.Save(&existing).Error
	if err != nil {
		logrus.Errorf("error while saving import allowed organization: %v", err)
		return err
	}
	return nil
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
