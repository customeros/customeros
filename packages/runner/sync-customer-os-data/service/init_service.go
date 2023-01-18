package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/sirupsen/logrus"
)

type InitService interface {
	Init()
}

type initService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewInitService(repositories *repository.Repositories, services *Services) InitService {
	return &initService{
		repositories: repositories,
		services:     services,
	}
}

func (s *initService) Init() {
	db := s.repositories.Dbs.ControlDb

	err := db.AutoMigrate(&entity.TenantSyncSettings{})
	if err != nil {
		logrus.Fatal(err)
	}

	err = db.AutoMigrate(&entity.SyncRun{})
	if err != nil {
		logrus.Fatal(err)
	}

	err = db.AutoMigrate(&entity.ConversationEvent{})
	if err != nil {
		logrus.Fatal(err)
	}
}
