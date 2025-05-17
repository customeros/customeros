package api_dashboard

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
)

type dashboardService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewDashboardService(log logger.Logger, repositories *repository.Repositories) cosapi_interfaces.DashboardService {
	return &dashboardService{
		log:          log,
		repositories: repositories,
	}
}

func (s *dashboardService) GetDashboardViewOrganizationsData(ctx context.Context, requestDetails cosapi_interfaces.DashboardViewOrganizationsRequest) (*utils.Pagination, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "DashboardService.GetDashboardViewOrganizationsData")
	defer spans.Finish()
	spans.LogKV("page", requestDetails.Page)
	spans.LogKV("limit", requestDetails.Limit)
	spans.LogObjectAsJson("where", requestDetails.Where)
	spans.LogObjectAsJson("sort", requestDetails.Sort)

	paginatedResult := utils.Pagination{
		Limit: requestDetails.Limit,
		Page:  requestDetails.Page,
	}

	dbNodes, err := s.repositories.DashboardRepository.GetDashboardViewOrganizationData(ctx, common.GetContext(ctx).Tenant, paginatedResult.GetSkip(), paginatedResult.GetLimit(), requestDetails.Where, requestDetails.Sort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodes.Count)

	organizationEntities := neo4jentity.OrganizationEntities{}

	for _, v := range dbNodes.Nodes {
		organizationEntities = append(organizationEntities, *neo4jmapper.MapDbNodeToOrganizationEntity(v))
	}

	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}

func (s *dashboardService) GetDashboardViewRenewalsData(ctx context.Context, requestDetails cosapi_interfaces.DashboardViewRenewalsRequest) (*utils.Pagination, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "DashboardService.GetDashboardViewRenewalsData")
	defer spans.Finish()
	spans.LogKV("page", requestDetails.Page)
	spans.LogKV("limit", requestDetails.Limit)
	spans.LogObjectAsJson("filter", requestDetails.Where)
	spans.LogObjectAsJson("sort", requestDetails.Sort)

	paginatedResult := utils.Pagination{
		Limit: requestDetails.Limit,
		Page:  requestDetails.Page,
	}

	dbRecords, err := s.repositories.DashboardRepository.GetDashboardViewRenewalData(ctx, common.GetContext(ctx).Tenant, paginatedResult.GetSkip(), paginatedResult.GetLimit(), requestDetails.Where, requestDetails.Sort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbRecords.Count)

	renewalRecordEntities := entity.RenewalsRecordEntities{}

	for _, v := range dbRecords.Records {
		renewalRecordEntity := entity.RenewalsRecordEntity{}
		if v.Values[0] != nil {
			renewalRecordEntity.Organization = *neo4jmapper.MapDbNodeToOrganizationEntity(utils.ToPtr(v.Values[0].(dbtype.Node)))
		}
		if v.Values[1] != nil {
			renewalRecordEntity.Contract = *neo4jmapper.MapDbNodeToContractEntity(utils.NodePtr(v.Values[1].(dbtype.Node)))
		}
		if v.Values[2] != nil {
			renewalRecordEntity.Opportunity = *neo4jmapper.MapDbNodeToOpportunityEntity(utils.NodePtr(v.Values[2].(dbtype.Node)))
		}
		renewalRecordEntities = append(renewalRecordEntities, renewalRecordEntity)
	}

	paginatedResult.SetRows(&renewalRecordEntities)
	return &paginatedResult, nil
}
