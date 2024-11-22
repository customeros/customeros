package service

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
)

type AuthService interface {
	ValidateAPIKey(ctx context.Context, c *gin.Context) error
}

type authService struct {
	services *Services
}

func NewAuthService(services *Services) AuthService {
	return &authService{
		services: services,
	}
}

func (a *authService) ValidateAPIKey(ctx context.Context, c *gin.Context) error {
	apiKey := c.Query(security.ApiKeyHeader)
	if apiKey == "" {
		return errors.New("missing API key")
	}

	appKey, err := a.services.CommonServices.PostgresRepositories.AppKeyRepository.FindByKey(
		ctx,
		string(security.CUSTOMER_OS_WEBHOOKS),
		apiKey,
	)
	if err != nil || appKey == nil {
		return errors.New("invalid API key")
	}
	return nil
}
