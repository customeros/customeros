package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type BillingProfileService interface {
	CreateBillingProfile(ctx context.Context, organizationId, name string, createdAt *time.Time) (string, error)
}
type billingProfileService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewBillingProfileService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) BillingProfileService {
	return &billingProfileService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *billingProfileService) CreateBillingProfile(ctx context.Context, organizationId, name string, createdAt *time.Time) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileService.CreateBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("name", name), log.String("organizationId", organizationId), log.Object("createdAt", createdAt))

	organizationExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), organizationId, neo4jentity.NodeLabelOrganization)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if !organizationExists {
		err = errors.New(fmt.Sprintf("Organization with id {%s} not found", organizationId))
		tracing.TraceErr(span, err)
		return "", err
	}

	grpcRequest := organizationpb.CreateBillingProfileGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		OrganizationId: organizationId,
		Name:           name,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		CreatedAt:      utils.ConvertTimeToTimestampPtr(createdAt),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := s.grpcClients.OrganizationClient.CreateBillingProfile(ctx, &grpcRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	WaitForObjectCreationAndLogSpan(ctx, s.repositories, response.Id, neo4jentity.NodeLabelBillingProfile, span)

	return response.Id, nil
}
