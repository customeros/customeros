package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TrackingRepository interface {
	Store(tracking entity.Tracking) (string, error)
}

type trackingRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewTrackingRepository(gormDb *gorm.DB) TrackingRepository {
	return &trackingRepositoryImpl{gormDb: gormDb}
}

func (repo *trackingRepositoryImpl) Store(tracking entity.Tracking) (string, error) {
	err := repo.gormDb.Save(&tracking).Error

	if err != nil {
		logrus.Errorf("Failed storing tracking: %v", tracking)
		return "", err
	}

	return tracking.ID, nil
}
