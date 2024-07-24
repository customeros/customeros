package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

type EmailingService interface {
	GenerateSpyPixelUrl(ctx context.Context, tenant string) (string, error)
	GenerateLinkUrl(ctx context.Context, tenant string, socialId string) (string, string, error)
}

type emailingService struct {
	log      logger.Logger
	services *Services
}

func NewEmailingService(log logger.Logger, services *Services) EmailingService {
	return &emailingService{
		log:      log,
		services: services,
	}
}

func (e emailingService) GenerateSpyPixelUrl(ctx context.Context, tenant string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (e emailingService) GenerateLinkUrl(ctx context.Context, tenant string, socialId string) (string, string, error) {
	//TODO implement me
	panic("implement me")
}
