package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type WhitelistDomainRepository interface {
	GetWildcardAllWhitelistDomain(tenant string) (*entity.WhitelistDomain, error)
	SaveWildcardAllWhitelistDomain(tenant, source, appSource string) error
	RemoveWildcardAllWhitelistDomain(tenant string) error

	SaveWhitelistDomain(whitelistDomain entity.WhitelistDomain) error
	GetWhitelistDomains(tenant string) ([]entity.WhitelistDomain, error)
}

type whitelistDomainRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewWhitelistDomainRepository(gormDb *gorm.DB) WhitelistDomainRepository {
	return &whitelistDomainRepositoryImpl{gormDb: gormDb}
}

func (repo *whitelistDomainRepositoryImpl) GetWildcardAllWhitelistDomain(tenant string) (*entity.WhitelistDomain, error) {
	var result entity.WhitelistDomain
	err := repo.gormDb.Model(&entity.WhitelistDomain{}).Find(&result, "tenant = ? AND domain = ?", tenant, "*").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logrus.Errorf("error while getting import allowed organization: %v", err)
		return nil, err
	}
	if result.Tenant == "" {
		return nil, nil
	}
	return &result, nil
}

func (repo *whitelistDomainRepositoryImpl) SaveWildcardAllWhitelistDomain(tenant, source, appSource string) error {
	existing, err := repo.GetWildcardAllWhitelistDomain(tenant)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	result := repo.gormDb.Create(&entity.WhitelistDomain{
		Tenant:    tenant,
		Name:      "*",
		Domain:    "*",
		Allowed:   true,
		CreatedAt: time.Now(),
	})
	if result.Error != nil {
		logrus.Errorf("error while saving import allowed organization: %v", err)
		return result.Error
	}
	return nil
}

func (repo *whitelistDomainRepositoryImpl) RemoveWildcardAllWhitelistDomain(tenant string) error {
	existing, err := repo.GetWildcardAllWhitelistDomain(tenant)
	if err != nil {
		return err
	}
	if existing == nil {
		return nil
	}

	result := repo.gormDb.Delete(existing)
	if result.Error != nil {
		logrus.Errorf("error while saving import allowed organization: %v", err)
		return result.Error
	}

	return nil
}

func (repo *whitelistDomainRepositoryImpl) SaveWhitelistDomain(whitelistDomain entity.WhitelistDomain) error {
	var existing entity.WhitelistDomain

	if whitelistDomain.ID != "" {
		err := repo.gormDb.Model(&entity.WhitelistDomain{}).First(&existing, "id = ?", whitelistDomain.ID).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				existing = entity.WhitelistDomain{}
			} else {
				logrus.Errorf("error while getting import allowed organization: %v", err)
				return err
			}
		}
	}

	if existing.ID == "" && whitelistDomain.Tenant != "" && whitelistDomain.Domain != "" && whitelistDomain.Name != "" {
		err := repo.gormDb.Model(&entity.WhitelistDomain{}).First(&existing, "tenant = ? AND domain = ? AND name = ?", whitelistDomain.Tenant, whitelistDomain.Domain, whitelistDomain.Name).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				existing = entity.WhitelistDomain{}
			} else {
				logrus.Errorf("error while getting import allowed organization: %v", err)
				return err
			}
		}

	}

	if existing.ID == "" {
		existing.Tenant = whitelistDomain.Tenant
		existing.Name = whitelistDomain.Name
		existing.Domain = whitelistDomain.Domain
	}

	existing.Allowed = whitelistDomain.Allowed

	err := repo.gormDb.Save(&existing).Error
	if err != nil {
		logrus.Errorf("error while saving import allowed organization: %v", err)
		return err
	}
	return nil
}

func (repo *whitelistDomainRepositoryImpl) GetWhitelistDomains(tenant string) ([]entity.WhitelistDomain, error) {
	var result []entity.WhitelistDomain
	err := repo.gormDb.Model(&entity.WhitelistDomain{}).Find(&result, "tenant = ? AND allowed = true", tenant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logrus.Errorf("error while getting whitelist domains: %v", err)
		return nil, err
	}
	return result, nil
}
