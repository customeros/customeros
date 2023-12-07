package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OAuthTokenService interface {
	Save(ctx context.Context, tokenEntity entity.OAuthTokenEntity) (*entity.OAuthTokenEntity, error)
	GetByPlayerIdAndProvider(ctx context.Context, playerId string, provider string) (*entity.OAuthTokenEntity, error)
	DeleteByPlayerIdAndProvider(ctx context.Context, playerId string, provider string) error
}

type oAuthTokenService struct {
	repositories *repository.Repositories
}

func NewOAuthTokenService(repositories *repository.Repositories) OAuthTokenService {
	return &oAuthTokenService{
		repositories: repositories,
	}
}

func (o oAuthTokenService) Save(ctx context.Context, tokenEntity entity.OAuthTokenEntity) (*entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenService.Save")
	defer span.Finish()

	result, err := o.repositories.OAuthTokenRepository.Save(ctx, tokenEntity)
	return result, err
}

func (o oAuthTokenService) GetByPlayerIdAndProvider(ctx context.Context, playerId, provider string) (*entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenService.GetByPlayerIdAndProvider")
	defer span.Finish()
	span.LogFields(log.String("playerId", playerId), log.String("provider", provider))

	authTokenEntity, err := o.repositories.OAuthTokenRepository.GetByPlayerIdAndProvider(ctx, playerId, provider)

	if err != nil {
		return nil, err
	}

	return authTokenEntity, nil
}

func (o oAuthTokenService) DeleteByPlayerIdAndProvider(ctx context.Context, playerId, provider string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenService.DeleteByPlayerIdAndProvider")
	defer span.Finish()
	span.LogFields(log.String("playerId", playerId), log.String("provider", provider))

	return o.repositories.OAuthTokenRepository.DeleteByPlayerIdAndProvider(ctx, playerId, provider)
}
