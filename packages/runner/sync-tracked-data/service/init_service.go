package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/repository"
	"github.com/sirupsen/logrus"
)

type InitService interface {
	Init()
}

type initService struct {
	repositories *repository.Repositories
}

func NewInitService(repositories *repository.Repositories) InitService {
	return &initService{
		repositories: repositories,
	}
}

func (s *initService) Init() {
	db := s.repositories.Drivers.GormDb

	err := db.AutoMigrate(&entity.SyncRun{})
	if err != nil {
		logrus.Fatal(err)
	}
}
