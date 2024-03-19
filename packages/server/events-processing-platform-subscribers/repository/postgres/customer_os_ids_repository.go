package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository/postgres/entity"
	"gorm.io/gorm"
)

type CustomerOsIdsRepository interface {
	Reserve(customerOsIds entity.CustomerOsIds) error
}

type customerOsIdsRepository struct {
	gormDb *gorm.DB
}

func NewCustomerOsIdsRepository(gormDb *gorm.DB) CustomerOsIdsRepository {
	return &customerOsIdsRepository{gormDb: gormDb}
}

func (repo *customerOsIdsRepository) Reserve(customerOsIds entity.CustomerOsIds) error {
	err := repo.gormDb.Save(&customerOsIds).Error
	return err
}
