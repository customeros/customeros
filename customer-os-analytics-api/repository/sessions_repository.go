package repository

import (
	"github.com.openline-ai.customer-os-analytics-api/repository/entity"
	"gorm.io/gorm"
)

type SessionsRepo struct {
	db *gorm.DB
}

type SessionsRepository interface {
	FindAllByApplication(appIdentifier entity.ApplicationUniqueIdentifier) RepositoryResult
}

func NewSessionsRepo(db *gorm.DB) *SessionsRepo {
	return &SessionsRepo{db: db}
}

func (r *SessionsRepo) FindAllByApplication(appIdentifier entity.ApplicationUniqueIdentifier) RepositoryResult {
	var sessions entity.SessionEntities

	err := r.db.Where(&entity.SessionEntity{Tenant: appIdentifier.Tenant, AppId: appIdentifier.AppId, TrackerName: appIdentifier.TrackerName}).Find(&sessions).Error

	if err != nil {
		return RepositoryResult{Error: err}
	}

	return RepositoryResult{Result: &sessions}
}
