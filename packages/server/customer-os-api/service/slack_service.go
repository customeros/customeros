package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SlackService interface {
	GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) (*utils.Pagination, error)
}

type slackService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewSlackService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) SlackService {
	return &slackService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *slackService) GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.GetSlackChannels")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("page", page))
	span.LogFields(log.Object("limit", limit))

	channels, totalCount, err := s.services.CommonServices.SlackService.GetPaginatedSlackChannels(ctx, tenant, page, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var paginatedResult = utils.Pagination{
		Limit:     page,
		Page:      limit,
		TotalRows: totalCount,
		Rows:      channels,
	}

	return &paginatedResult, nil
}
