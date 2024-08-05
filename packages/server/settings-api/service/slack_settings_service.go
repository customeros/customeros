package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/model"
	"golang.org/x/net/context"
)

type SlackSettingsService interface {
	GetSlackSettings(ctx context.Context, tenant string) (*model.SlackSettingsResponse, error)
}

type slackSettingsService struct {
	services *Services
	log      logger.Logger
}

func NewSlackSettingsService(services *Services, log logger.Logger) SlackSettingsService {
	return &slackSettingsService{
		services: services,
		log:      log,
	}
}

func (u slackSettingsService) GetSlackSettings(ctx context.Context, tenant string) (*model.SlackSettingsResponse, error) {
	slackSettings, err := u.services.CommonServices.PostgresRepositories.SlackSettingsRepository.Get(ctx, tenant)
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
