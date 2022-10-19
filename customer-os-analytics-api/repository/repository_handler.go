package repository

import "gorm.io/gorm"

type RepositoryHandler struct {
	AppInfoRepo  AppInfoRepository
	SessionsRepo SessionsRepository
}

func InitRepositories(db *gorm.DB) *RepositoryHandler {

	appInfoRepo := NewAppInfoRepo(db)
	sessionsRepo := NewSessionsRepo(db)

	return &RepositoryHandler{
		AppInfoRepo:  appInfoRepo,
		SessionsRepo: sessionsRepo,
	}
}
