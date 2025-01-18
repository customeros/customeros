package api_slack

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
)

type slackService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	slack        interfaces.SlackService
}

func NewSlackService(
	log logger.Logger,
	repositories *repository.Repositories,
	grpcClients *grpc_client.Clients,
	slack interfaces.SlackService,
) cosapi_interfaces.SlackService {
	return &slackService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		slack:        slack,
	}
}

func (s *slackService) GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.GetSlackChannels")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("page", page))
	span.LogFields(log.Object("limit", limit))

	channels, totalCount, err := s.slack.GetPaginatedSlackChannels(ctx, tenant, page, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	paginatedResult := utils.Pagination{
		Limit:     page,
		Page:      limit,
		TotalRows: totalCount,
		Rows:      channels,
	}

	return &paginatedResult, nil
}

func (s *slackService) GetSlackSettings(ctx context.Context, tenant string) (*cosapi_interfaces.SlackSettingsResponse, error) {
	slackSettings, err := s.repositories.PostgresRepositories.SlackSettingsRepository.Get(ctx, tenant)
	if err != nil {
		return nil, err
	}

	if slackSettings == nil {
		return &cosapi_interfaces.SlackSettingsResponse{
			SlackEnabled: false,
		}, nil
	}

	slackSettingsResponse := cosapi_interfaces.SlackSettingsResponse{
		SlackEnabled: true,
	}

	return &slackSettingsResponse, nil
}
