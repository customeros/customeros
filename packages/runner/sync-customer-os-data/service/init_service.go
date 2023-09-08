package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
)

type InitService interface {
	Init()
}

type initService struct {
	repositories *repository.Repositories
	services     *Services
	log          logger.Logger
}

func NewInitService(repositories *repository.Repositories, services *Services, log logger.Logger) InitService {
	return &initService{
		repositories: repositories,
		services:     services,
		log:          log,
	}
}

func (s *initService) Init() {
	db := s.repositories.Dbs.GormDB

	err := db.AutoMigrate(&entity.TenantSyncSettings{})
	if err != nil {
		s.log.Fatal(err)
	}

	err = db.AutoMigrate(&entity.SyncRun{})
	if err != nil {
		s.log.Fatal(err)
	}

}
