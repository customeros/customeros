package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
)

type SlackSettingsService interface {
	GetSlackSettings(tenant string) (*model.SlackSettingsResponse, error)
}

type slackSettingsService struct {
	repositories *repository.PostgresRepositories
	log          logger.Logger
}

func NewSlackSettingsService(repositories *repository.PostgresRepositories, log logger.Logger) SlackSettingsService {
	return &slackSettingsService{
		repositories: repositories,
		log:          log,
	}
}

func (u slackSettingsService) GetSlackSettings(tenant string) (*model.SlackSettingsResponse, error) {
	slackSettings, err := u.repositories.AuthRepositories.SlackSettingsRepository.Get(tenant)
	if err != nil {
		return nil, err
	}

	if slackSettings == nil {
		return &model.SlackSettingsResponse{
			SlackEnabled: false,
		}, nil
	}

	var slackSettingsResponse = model.SlackSettingsResponse{
		SlackEnabled: true,
	}

	return &slackSettingsResponse, nil
}
