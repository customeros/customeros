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
	CreateBillingProfile(ctx context.Context, organizationId, legalName, taxId string, createdAt *time.Time) (string, error)
	UpdateBillingProfile(ctx context.Context, organizationId, billingProfileId string, legalName, taxId *string, updatedAt *time.Time) error
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

func (s *billingProfileService) CreateBillingProfile(ctx context.Context, organizationId, legalName, taxId string, createdAt *time.Time) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileService.CreateBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("legalName", legalName), log.String("taxId", taxId), log.String("organizationId", organizationId), log.Object("createdAt", createdAt))

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
		LegalName:      legalName,
		TaxId:          taxId,
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

func (s *billingProfileService) UpdateBillingProfile(ctx context.Context, organizationId, billingProfileId string, legalName, taxId *string, updatedAt *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileService.UpdateBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, billingProfileId)
	span.LogFields(log.Object("legalName", legalName), log.Object("taxId", legalName), log.String("organizationId", organizationId), log.Object("updatedAt", updatedAt))

	billingProfileExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsByIdLinkedFrom(ctx, common.GetTenantFromContext(ctx), billingProfileId, neo4jentity.NodeLabelBillingProfile, organizationId, neo4jentity.NodeLabelOrganization, "HAS_BILLING_PROFILE")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !billingProfileExists {
		err = errors.New(fmt.Sprintf("Billing profile with id {%s} not found", billingProfileId))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := organizationpb.UpdateBillingProfileGrpcRequest{
		Tenant:           common.GetTenantFromContext(ctx),
		OrganizationId:   organizationId,
		BillingProfileId: billingProfileId,
		LegalName:        utils.IfNotNilString(legalName),
		TaxId:            utils.IfNotNilString(taxId),
		LoggedInUserId:   common.GetUserIdFromContext(ctx),
		UpdatedAt:        utils.ConvertTimeToTimestampPtr(updatedAt),
	}
	fieldsMask := make([]organizationpb.BillingProfileFieldMask, 0)
	if legalName != nil {
		fieldsMask = append(fieldsMask, organizationpb.BillingProfileFieldMask_BILLING_PROFILE_PROPERTY_LEGAL_NAME)
	}
	if taxId != nil {
		fieldsMask = append(fieldsMask, organizationpb.BillingProfileFieldMask_BILLING_PROFILE_PROPERTY_TAX_ID)
	}
	grpcRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = s.grpcClients.OrganizationClient.UpdateBillingProfile(ctx, &grpcRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}
