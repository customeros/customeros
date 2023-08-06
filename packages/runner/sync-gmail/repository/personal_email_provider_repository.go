package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PersonalEmailProviderRepository interface {
	GetPersonalEmailProviderList() ([]entity.PersonalEmailProvider, error)
}

type personalEmailProviderRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewPersonalEmailProviderRepository(gormDb *gorm.DB) PersonalEmailProviderRepository {
	return &personalEmailProviderRepositoryImpl{gormDb: gormDb}
}

func (repo *personalEmailProviderRepositoryImpl) GetPersonalEmailProviderList() ([]entity.PersonalEmailProvider, error) {
	result := []entity.PersonalEmailProvider{}
	err := repo.gormDb.Find(&result).Limit(1000).Error

	if err != nil {
		logrus.Errorf("error while getting personal email provider list: %v", err)
		return nil, err
	}

	return result, nil
}
