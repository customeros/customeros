package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
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
	LinkEmailToBillingProfile(ctx context.Context, organizationId, billingProfileId, emailId string, primary bool) error
	UnlinkEmailFromBillingProfile(ctx context.Context, organizationId, billingProfileId, emailId string) error
	LinkLocationToBillingProfile(ctx context.Context, organizationId, billingProfileId, locationId string) error
	UnlinkLocationFromBillingProfile(ctx context.Context, organizationId, billingProfileId, locationId string) error
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

	organizationExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), organizationId, neo4jutil.NodeLabelOrganization)
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
	response, err := utils.CallEventsPlatformGRPCWithRetry[*organizationpb.BillingProfileIdGrpcResponse](func() (*organizationpb.BillingProfileIdGrpcResponse, error) {
		return s.grpcClients.OrganizationClient.CreateBillingProfile(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, neo4jutil.NodeLabelBillingProfile, span)

	return response.Id, nil
}

func (s *billingProfileService) UpdateBillingProfile(ctx context.Context, organizationId, billingProfileId string, legalName, taxId *string, updatedAt *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileService.UpdateBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, billingProfileId)
	span.LogFields(log.Object("legalName", legalName), log.Object("taxId", legalName), log.String("organizationId", organizationId), log.Object("updatedAt", updatedAt))

	billingProfileExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsByIdLinkedFrom(ctx, common.GetTenantFromContext(ctx), billingProfileId, neo4jutil.NodeLabelBillingProfile, organizationId, neo4jutil.NodeLabelOrganization, "HAS_BILLING_PROFILE")
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
		AppSource:        constants.AppSourceCustomerOsApi,
	}
	fieldsMask := make([]organizationpb.BillingProfileFieldMask, 0)
	if legalName != nil {
		fieldsMask = append(fieldsMask, organizationpb.BillingProfileFieldMask_BILLING_PROFILE_PROPERTY_LEGAL_NAME)
	}
	if taxId != nil {
		fieldsMask = append(fieldsMask, organizationpb.BillingProfileFieldMask_BILLING_PROFILE_PROPERTY_TAX_ID)
	}
	if len(fieldsMask) == 0 {
		span.LogFields(log.String("result", "No fields to update"))
		return nil
	}
	grpcRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.BillingProfileIdGrpcResponse](func() (*organizationpb.BillingProfileIdGrpcResponse, error) {
		return s.grpcClients.OrganizationClient.UpdateBillingProfile(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *billingProfileService) LinkEmailToBillingProfile(ctx context.Context, organizationId, billingProfileId, emailId string, primary bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileService.LinkEmailToBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, billingProfileId)
	span.LogFields(log.String("organizationId", organizationId), log.String("emailId", emailId), log.Bool("primary", primary))

	billingProfileExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsByIdLinkedFrom(ctx, common.GetTenantFromContext(ctx), billingProfileId, neo4jutil.NodeLabelBillingProfile, organizationId, neo4jutil.NodeLabelOrganization, "HAS_BILLING_PROFILE")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !billingProfileExists {
		err = errors.New(fmt.Sprintf("Billing profile with id {%s} not found", billingProfileId))
		tracing.TraceErr(span, err)
		return err
	}

	emailExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), emailId, neo4jutil.NodeLabelEmail)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !emailExists {
		err = errors.New(fmt.Sprintf("Email with id {%s} not found", emailId))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := organizationpb.LinkEmailToBillingProfileGrpcRequest{
		Tenant:           common.GetTenantFromContext(ctx),
		OrganizationId:   organizationId,
		BillingProfileId: billingProfileId,
		EmailId:          emailId,
		Primary:          primary,
		LoggedInUserId:   common.GetUserIdFromContext(ctx),
		AppSource:        constants.AppSourceCustomerOsApi,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.BillingProfileIdGrpcResponse](func() (*organizationpb.BillingProfileIdGrpcResponse, error) {
		return s.grpcClients.OrganizationClient.LinkEmailToBillingProfile(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *billingProfileService) UnlinkEmailFromBillingProfile(ctx context.Context, organizationId, billingProfileId, emailId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileService.UnlinkEmailFromBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, billingProfileId)
	span.LogFields(log.String("organizationId", organizationId), log.String("emailId", emailId))

	billingProfileExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsByIdLinkedFrom(ctx, common.GetTenantFromContext(ctx), billingProfileId, neo4jutil.NodeLabelBillingProfile, organizationId, neo4jutil.NodeLabelOrganization, "HAS_BILLING_PROFILE")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !billingProfileExists {
		err = errors.New(fmt.Sprintf("Billing profile with id {%s} not found", billingProfileId))
		tracing.TraceErr(span, err)
		return err
	}

	emailExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), emailId, neo4jutil.NodeLabelEmail)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !emailExists {
		err = errors.New(fmt.Sprintf("Email with id {%s} not found", emailId))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := organizationpb.UnlinkEmailFromBillingProfileGrpcRequest{
		Tenant:           common.GetTenantFromContext(ctx),
		OrganizationId:   organizationId,
		BillingProfileId: billingProfileId,
		EmailId:          emailId,
		LoggedInUserId:   common.GetUserIdFromContext(ctx),
		AppSource:        constants.AppSourceCustomerOsApi,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.BillingProfileIdGrpcResponse](func() (*organizationpb.BillingProfileIdGrpcResponse, error) {
		return s.grpcClients.OrganizationClient.UnlinkEmailFromBillingProfile(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *billingProfileService) LinkLocationToBillingProfile(ctx context.Context, organizationId, billingProfileId, locationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileService.LinkLocationToBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, billingProfileId)
	span.LogFields(log.String("organizationId", organizationId), log.String("locationId", locationId))

	billingProfileExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsByIdLinkedFrom(ctx, common.GetTenantFromContext(ctx), billingProfileId, neo4jutil.NodeLabelBillingProfile, organizationId, neo4jutil.NodeLabelOrganization, "HAS_BILLING_PROFILE")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !billingProfileExists {
		err = errors.New(fmt.Sprintf("Billing profile with id {%s} not found", billingProfileId))
		tracing.TraceErr(span, err)
		return err
	}

	locationExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), locationId, neo4jutil.NodeLabelLocation)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !locationExists {
		err = errors.New(fmt.Sprintf("LocationlocationId with id {%s} not found", locationId))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := organizationpb.LinkLocationToBillingProfileGrpcRequest{
		Tenant:           common.GetTenantFromContext(ctx),
		OrganizationId:   organizationId,
		BillingProfileId: billingProfileId,
		LocationId:       locationId,
		LoggedInUserId:   common.GetUserIdFromContext(ctx),
		AppSource:        constants.AppSourceCustomerOsApi,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.BillingProfileIdGrpcResponse](func() (*organizationpb.BillingProfileIdGrpcResponse, error) {
		return s.grpcClients.OrganizationClient.LinkLocationToBillingProfile(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *billingProfileService) UnlinkLocationFromBillingProfile(ctx context.Context, organizationId, billingProfileId, locationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillingProfileService.UnlinkLocationFromBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, billingProfileId)
	span.LogFields(log.String("organizationId", organizationId), log.String("locationId", locationId))

	billingProfileExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsByIdLinkedFrom(ctx, common.GetTenantFromContext(ctx), billingProfileId, neo4jutil.NodeLabelBillingProfile, organizationId, neo4jutil.NodeLabelOrganization, "HAS_BILLING_PROFILE")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !billingProfileExists {
		err = errors.New(fmt.Sprintf("Billing profile with id {%s} not found", billingProfileId))
		tracing.TraceErr(span, err)
		return err
	}

	locationExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), locationId, neo4jutil.NodeLabelLocation)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !locationExists {
		err = errors.New(fmt.Sprintf("Location with id {%s} not found", locationId))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := organizationpb.UnlinkLocationFromBillingProfileGrpcRequest{
		Tenant:           common.GetTenantFromContext(ctx),
		OrganizationId:   organizationId,
		BillingProfileId: billingProfileId,
		LocationId:       locationId,
		LoggedInUserId:   common.GetUserIdFromContext(ctx),
		AppSource:        constants.AppSourceCustomerOsApi,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.BillingProfileIdGrpcResponse](func() (*organizationpb.BillingProfileIdGrpcResponse, error) {
		return s.grpcClients.OrganizationClient.UnlinkLocationFromBillingProfile(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}
