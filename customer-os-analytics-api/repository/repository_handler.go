package repository

import "gorm.io/gorm"

type RepositoryHandler struct {
	AppInfoRepo AppInfoRepository
}

func InitRepositories(db *gorm.DB) *RepositoryHandler {

	appInfoRepo := NewAppInfoRepo(db)

	return &RepositoryHandler{
		AppInfoRepo: appInfoRepo,
	}
}
