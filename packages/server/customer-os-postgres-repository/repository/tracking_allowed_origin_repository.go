package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TrackingAllowedOriginRepository interface {
	GetTenantForOrigin(origin string) (*string, error)
}

type trackingAllowedOriginRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewTrackingAllowedOriginRepository(gormDb *gorm.DB) TrackingAllowedOriginRepository {
	return &trackingAllowedOriginRepositoryImpl{gormDb: gormDb}
}

func (repo *trackingAllowedOriginRepositoryImpl) GetTenantForOrigin(origin string) (*string, error) {
	var result entity.TrackingAllowedOrigin
	err := repo.gormDb.Model(&entity.TrackingAllowedOrigin{}).Find(&result, "origin = ?", origin).Error
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
	return &result.Tenant, nil
}
