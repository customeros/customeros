package repository

import (
	"gorm.io/gorm"
)

type RepositoryHandler struct {
	AppInfoRepo  AppInfoRepository
	SessionsRepo SessionsRepository
	PageViewRepo PageViewRepository
}

func InitRepositories(db *gorm.DB) *RepositoryHandler {
	return &RepositoryHandler{
		AppInfoRepo:  NewAppInfoRepo(db),
		SessionsRepo: NewSessionsRepo(db),
		PageViewRepo: NewPageViewRepo(db),
	}
}
