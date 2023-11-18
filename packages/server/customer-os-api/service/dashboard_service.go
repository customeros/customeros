package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	entityDashboard "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity/dashboard"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"math/rand"
)

type DashboardViewOrganizationsRequest struct {
	Where *model.Filter
	Sort  *model.SortBy
	Page  int
	Limit int
}

type DashboardService interface {
	GetDashboardViewOrganizationsData(ctx context.Context, requestDetails DashboardViewOrganizationsRequest) (*utils.Pagination, error)

	GetDashboardNewCustomersData(ctx context.Context, year int) (*entityDashboard.DashboardNewCustomersData, error)
	GetDashboardRetentionRateData(ctx context.Context, year int) (*entityDashboard.DashboardRetentionRateData, error)
	GetDashboardRevenueAtRiskData(ctx context.Context, year int) (*entityDashboard.DashboardRevenueAtRiskData, error)
	GetDashboardARRBreakdownData(ctx context.Context, year int) (*entityDashboard.DashboardARRBreakdownData, error)
}

type dashboardService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewDashboardService(log logger.Logger, repositories *repository.Repositories, services *Services) DashboardService {
	return &dashboardService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *dashboardService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *dashboardService) GetDashboardViewOrganizationsData(ctx context.Context, requestDetails DashboardViewOrganizationsRequest) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardViewOrganizationsData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("page", requestDetails.Page), log.Int("limit", requestDetails.Limit))
	if requestDetails.Where != nil {
		span.LogFields(log.Object("filter", *requestDetails.Where))
	}
	if requestDetails.Sort != nil {
		span.LogFields(log.Object("sort", *requestDetails.Sort))
	}

	var paginatedResult = utils.Pagination{
		Limit: requestDetails.Limit,
		Page:  requestDetails.Page,
	}

	dbNodes, err := s.repositories.QueryRepository.GetDashboardViewOrganizationData(ctx, common.GetContext(ctx).Tenant, paginatedResult.GetSkip(), paginatedResult.GetLimit(), requestDetails.Where, requestDetails.Sort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodes.Count)

	organizationEntities := entity.OrganizationEntities{}

	for _, v := range dbNodes.Nodes {
		organizationEntities = append(organizationEntities, *s.services.OrganizationService.mapDbNodeToOrganizationEntity(*v))
	}

	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}

func (s *dashboardService) GetDashboardNewCustomersData(ctx context.Context, year int) (*entityDashboard.DashboardNewCustomersData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardNewCustomersData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("year", year))

	response := entityDashboard.DashboardNewCustomersData{}

	response.ThisMonthCount = 127
	response.ThisMonthIncreasePercentage = 3.1

	for i := 1; i <= 12; i++ {
		response.Months = append(response.Months, &entityDashboard.DashboardNewCustomerMonthData{
			Month: i,
			Count: i*10 + 7,
		})
	}

	return &response, nil
}

func (s *dashboardService) GetDashboardRetentionRateData(ctx context.Context, year int) (*entityDashboard.DashboardRetentionRateData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardNewCustomersData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("year", year))

	response := entityDashboard.DashboardRetentionRateData{}

	response.RetentionRate = 86
	response.IncreasePercentage = 2.6

	min := 2
	max := 5
	for i := 1; i <= 12; i++ {
		response.Months = append(response.Months, &entityDashboard.DashboardRetentionRatePerMonthData{
			Month:      i,
			RenewCount: i*(rand.Intn(max-min)+min) + 7,
			ChurnCount: i + 2,
		})
	}

	return &response, nil
}

func (s *dashboardService) GetDashboardRevenueAtRiskData(ctx context.Context, year int) (*entityDashboard.DashboardRevenueAtRiskData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardRevenueAtRiskData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("year", year))

	response := entityDashboard.DashboardRevenueAtRiskData{}

	response.HighConfidence = 1504990
	response.AtRisk = 355300

	return &response, nil
}

func (s *dashboardService) GetDashboardARRBreakdownData(ctx context.Context, year int) (*entityDashboard.DashboardARRBreakdownData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardARRBreakdownData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("year", year))

	response := entityDashboard.DashboardARRBreakdownData{}

	response.ArrBreakdown = 1830990
	response.IncreasePercentage = 2.3

	min := 1
	max := 50
	for i := 1; i <= 12; i++ {
		response.Months = append(response.Months, &entityDashboard.DashboardARRBreakdownPerMonthData{
			Month:           i,
			NewlyContracted: rand.Intn(max-min) + min,
			Renewals:        rand.Intn(max-min) + min,
			Upsells:         rand.Intn(max-min) + min,
			Downgrades:      rand.Intn(max-min) + min,
			Cancellations:   rand.Intn(max-min) + min,
			Churned:         rand.Intn(max-min) + min,
		})
	}

	return &response, nil
}
