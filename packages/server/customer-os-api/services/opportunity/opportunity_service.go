package api_opportunity

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	model2 "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
)

type opportunityService struct {
	log          logger.Logger
	repositories *repository.Repositories
	opportunity  interfaces.OpportunityService
	organization interfaces.OrganizationService
}

func NewOpportunityService(log logger.Logger, repositories *repository.Repositories, opportunity interfaces.OpportunityService, org interfaces.OrganizationService) cosapi_interfaces.OpportunityService {
	return &opportunityService{
		log:          log,
		repositories: repositories,
		opportunity:  opportunity,
		organization: org,
	}
}

func (s *opportunityService) UpdateRenewal(ctx context.Context, opportunityId string, renewalLikelihood neo4jenum.RenewalLikelihood, amount *float64, comments *string, ownerUserId *string, adjustedRate *int64, appSource string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "OpportunityService.UpdateRenewal")
	defer spans.Finish()
	spans.LogKV("opportunityId", opportunityId,
		"renewalLikelihood", renewalLikelihood,
		"amount", amount,
		"comments", comments,
		"appSource", appSource)

	if opportunityId == "" {
		err := fmt.Errorf("(OpportunityService.UpdateRenewal) opportunity id is missing")
		s.log.Error(err.Error())
		spans.TraceError(err)
		return err
	}

	opportunityExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), opportunityId, model2.NodeLabelOpportunity)
	if !opportunityExists {
		err := fmt.Errorf("(OpportunityService.UpdateRenewal) opportunity with id {%s} not found", opportunityId)
		s.log.Error(err.Error())
		spans.TraceError(err)
		return err
	}

	dataFields := data_fields.OpportunityFields{
		OwnerId:             ownerUserId,
		Amount:              amount,
		Comments:            comments,
		RenewalAdjustedRate: adjustedRate,
		RenewalLikelihood:   utils.ToPtr(renewalLikelihood),
	}

	_, err := s.opportunity.Save(ctx, nil, &opportunityId, &dataFields)
	if err != nil {
		spans.TraceError(err)
		s.log.Error("error while updating renewal opportunity: ", err.Error())
		return err
	}

	return nil
}

func (s *opportunityService) UpdateRenewalsForOrganization(ctx context.Context, organizationId string, renewalLikelihood neo4jenum.RenewalLikelihood, renewalAdjustedRate *int64) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "OpportunityService.UpdateRenewalsForOrganization")
	defer span.Finish()
	span.LogKV("organizationId", organizationId,
		"renewalLikelihood", renewalLikelihood.String(),
		"renewalAdjustedRate", renewalAdjustedRate)

	tenant := common.GetTenantFromContext(ctx)

	_, err := s.organization.GetById(ctx, tenant, organizationId)
	if err != nil {
		span.TraceError(err)
		return err
	}

	opportunityDbNodes, err := s.repositories.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunitiesForOrganization(ctx, common.GetTenantFromContext(ctx), organizationId, true)
	if err != nil {
		span.TraceError(err)
		return err
	}

	for _, opportunityDbNode := range opportunityDbNodes {
		opportunity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)
		if err := s.UpdateRenewal(ctx, opportunity.Id, renewalLikelihood, nil, nil, nil, renewalAdjustedRate, constants.AppSourceCustomerOsApi); err != nil {
			span.TraceError(err)
			return err
		}
	}

	return nil
}

func (s *opportunityService) GetPaginatedOrganizationOpportunities(ctx context.Context, page int, limit int) (*utils.Pagination, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "OpportunityService.GetPaginatedOrganizationOpportunities")
	defer span.Finish()
	span.LogKV("page", page, "limit", limit)

	paginatedResult := utils.Pagination{
		Limit: limit,
		Page:  page,
	}

	dbNodesWithTotalCount, err := s.repositories.Neo4jRepositories.OpportunityReadRepository.GetPaginatedOpportunitiesLinkedToAnOrganization(ctx, common.GetContext(ctx).Tenant, paginatedResult.GetSkip(), paginatedResult.GetLimit())
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	opportunities := neo4jentity.OpportunityEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		opportunities = append(opportunities, *neo4jmapper.MapDbNodeToOpportunityEntity(v))
	}
	paginatedResult.SetRows(&opportunities)
	return &paginatedResult, nil
}
