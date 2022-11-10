package repository

import (
	"gorm.io/gorm"
)

type RepositoryContainer struct {
	AppInfoRepo  AppInfoRepository
	SessionsRepo SessionsRepository
	PageViewRepo PageViewRepository
}

func InitRepositories(db *gorm.DB) *RepositoryContainer {
	return &RepositoryContainer{
		AppInfoRepo:  NewAppInfoRepo(db),
		SessionsRepo: NewSessionsRepo(db),
		PageViewRepo: NewPageViewRepo(db),
	}
}
