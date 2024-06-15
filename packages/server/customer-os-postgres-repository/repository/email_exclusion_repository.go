package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmailExclusionRepository interface {
	GetExclusionList() ([]entity.EmailExclusion, error)
}

type emailExclusionRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewEmailExclusionRepository(gormDb *gorm.DB) EmailExclusionRepository {
	return &emailExclusionRepositoryImpl{gormDb: gormDb}
}

func (repo *emailExclusionRepositoryImpl) GetExclusionList() ([]entity.EmailExclusion, error) {
	result := []entity.EmailExclusion{}
	err := repo.gormDb.Find(&result).Limit(5000).Error

	if err != nil {
		logrus.Errorf("error while getting personal email provider list: %v", err)
		return nil, err
	}

	return result, nil
}
