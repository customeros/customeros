package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PersonalEmailProviderRepository interface {
	GetPersonalEmailProviders() ([]entity.PersonalEmailProvider, error)
}

type personalEmailProviderRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewPersonalEmailProviderRepository(gormDb *gorm.DB) PersonalEmailProviderRepository {
	return &personalEmailProviderRepositoryImpl{gormDb: gormDb}
}

func (repo *personalEmailProviderRepositoryImpl) GetPersonalEmailProviders() ([]entity.PersonalEmailProvider, error) {
	result := []entity.PersonalEmailProvider{}
	err := repo.gormDb.Find(&result).Limit(1000).Error

	if err != nil {
		logrus.Errorf("error while getting personal email provider list: %v", err)
		return nil, err
	}

	return result, nil
}
