package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type WorkflowService interface {
	ExecuteWorkflow(ctx context.Context, tenant string, workflowId uint64) error
}

type workflowService struct {
	log      logger.Logger
	services *Services
}

func NewWorkflowService(log logger.Logger, services *Services) WorkflowService {
	return &workflowService{
		log:      log,
		services: services,
	}
}

func (s *workflowService) ExecuteWorkflow(ctx context.Context, tenant string, workflowId uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.ExecuteWorkflow")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.Uint64("workflowId", workflowId))

	// get workflow by id
	workflow, err := s.services.PostgresRepositories.WorkflowRepository.GetWorkflowByTenantAndId(ctx, tenant, workflowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !workflow.Live {
		span.LogFields(log.String("result", "skipping, workflow not live"))
		return nil
	}

	// evaluation condition
	var organizationIds []string

	switch workflow.WorkflowType {
	case postgresentity.WorkflowTypeIdealCustomerProfile:
		organizationIds, err = s.findOrganizationIds(ctx, tenant, workflow)
	}

	// execute actions
	for _, organizationId := range organizationIds {
		_ = s.executeOrganizationAction(ctx, tenant, organizationId, workflow.WorkflowType)
	}

	return nil
}

func (s *workflowService) findOrganizationIds(ctx context.Context, tenant string, workflow postgresentity.Workflow) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.findOrganizationIds")
	defer span.Finish()

	var organizationIds []string

	// unmarshal condition into filter
	filter, err := model.UnmarshalFilter(workflow.Condition)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal filter"))
		return nil, err
	}

	switch workflow.WorkflowType {
	case postgresentity.WorkflowTypeIdealCustomerProfile:
		// add condition to filter that stage is Lead
		filter.And = append(filter.And, &model.Filter{
			Filter: &model.FilterItem{
				Property:  neo4jrepo.SearchParamStage,
				Operation: model.ComparisonOperatorEq,
				Value:     model.AnyTypeValue{Str: utils.StringPtr(neo4jenum.Lead.String())},
				JsonValue: neo4jenum.Lead.String(),
			},
		})
		organizationIds, err = s.services.Neo4jRepositories.OrganizationWithFiltersReadRepository.GetFilteredOrganizationIds(ctx, tenant, filter)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	return organizationIds, nil

}

func (s *workflowService) executeOrganizationAction(ctx context.Context, tenant string, organizationId string, workflowType postgresentity.WorkflowType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.executeOrganizationAction")
	defer span.Finish()

	switch workflowType {
	case postgresentity.WorkflowTypeIdealCustomerProfile:
		_, err := utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			request := organizationpb.UpdateOrganizationGrpcRequest{
				Tenant:         tenant,
				OrganizationId: organizationId,
				SourceFields: &commonpb.SourceFields{
					AppSource: string(workflowType),
					Source:    neo4jentity.DataSourceOpenline.String(),
				},
				Stage:      neo4jenum.Target.String(),
				FieldsMask: []organizationpb.OrganizationMaskField{organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_STAGE},
			}
			return s.services.GrpcClients.OrganizationClient.UpdateOrganization(ctx, &request)
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error updating organization %s: %s", organizationId, err.Error())
			return err
		}
	}
	return nil
}
