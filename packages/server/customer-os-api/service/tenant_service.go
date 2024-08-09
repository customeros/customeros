package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TenantService interface {
	GetTenantForWorkspace(ctx context.Context, workspaceEntity entity.WorkspaceEntity) (*neo4jentity.TenantEntity, error)
	GetTenantForUserEmail(ctx context.Context, email string) (*neo4jentity.TenantEntity, error)
	Merge(ctx context.Context, tenantEntity neo4jentity.TenantEntity) (*neo4jentity.TenantEntity, error)
	GetTenantBillingProfiles(ctx context.Context) (*neo4jentity.TenantBillingProfileEntities, error)
	GetTenantBillingProfile(ctx context.Context, id string) (*neo4jentity.TenantBillingProfileEntity, error)
	GetDefaultTenantBillingProfile(ctx context.Context) (*neo4jentity.TenantBillingProfileEntity, error)
	CreateTenantBillingProfile(ctx context.Context, input model.TenantBillingProfileInput) (string, error)
	UpdateTenantBillingProfile(ctx context.Context, input model.TenantBillingProfileUpdateInput) error
	GetTenantSettings(ctx context.Context) (*neo4jentity.TenantSettingsEntity, error)
	UpdateTenantSettings(ctx context.Context, input *model.TenantSettingsInput) error

	HardDelete(ctx context.Context, tenant string) error
}

type tenantService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewTenantService(log logger.Logger, repository *repository.Repositories, grpcClients *grpc_client.Clients) TenantService {
	return &tenantService{
		log:          log,
		repositories: repository,
		grpcClients:  grpcClients,
	}
}

func (s *tenantService) Merge(ctx context.Context, tenantEntity neo4jentity.TenantEntity) (*neo4jentity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.Merge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "tenantEntity", tenantEntity)

	tenantName := strings.ReplaceAll(tenantEntity.Name, " ", "")
	if tenantName == "" {
		err := fmt.Errorf("tenant name is empty")
		tracing.TraceErr(span, err)
		return nil, err
	}

	for i := 0; i < 10; i++ {
		existNode, err := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantByName(ctx, tenantName)
		if err != nil {
			return nil, fmt.Errorf("merge: %w", err)
		}
		if existNode == nil {
			break
		}
		tenantName = fmt.Sprintf("%s%d", tenantName, rand.Intn(10))
	}
	span.LogFields(log.Object("tenantName", tenantName))
	tenantEntity.Name = tenantName
	tenant, err := s.repositories.Neo4jRepositories.TenantWriteRepository.CreateTenantIfNotExistAndReturn(ctx, tenantEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	// save tenant in postgres table
	_, err = s.repositories.PostgresRepositories.TenantRepository.Create(ctx, postgresentity.Tenant{
		Name: tenantName,
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}
	// create tenant specific api key
	err = s.repositories.PostgresRepositories.TenantWebhookApiKeyRepository.CreateApiKey(ctx, tenantName)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	err = s.repositories.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenantEntity.Name, "gmail", "gmail")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	err = s.repositories.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenantEntity.Name, "slack", "slack")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	err = s.repositories.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenantEntity.Name, "intercom", "intercom")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("merge: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantEntity(tenant), nil
}

func (s *tenantService) GetTenantForWorkspace(ctx context.Context, workspaceEntity entity.WorkspaceEntity) (*neo4jentity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantForWorkspace")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("workspace", workspaceEntity))

	tenant, err := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantForWorkspaceProvider(ctx, workspaceEntity.Name, workspaceEntity.Provider)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantForWorkspace: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantEntity(tenant), nil
}

func (s *tenantService) GetTenantForUserEmail(ctx context.Context, email string) (*neo4jentity.TenantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantForUserEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	tenant, err := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantForUserEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantForWorkspace: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantEntity(tenant), nil
}

func (s *tenantService) GetTenantBillingProfiles(ctx context.Context) (*neo4jentity.TenantBillingProfileEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantBillingProfiles")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbNodes, err := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantBillingProfiles(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantBillingProfiles: %w", err)
	}

	tenantBillingProfiles := neo4jentity.TenantBillingProfileEntities{}
	for _, dbNode := range dbNodes {
		tenantBillingProfiles = append(tenantBillingProfiles, *neo4jmapper.MapDbNodeToTenantBillingProfileEntity(dbNode))
	}

	return &tenantBillingProfiles, nil
}

func (s *tenantService) GetTenantBillingProfile(ctx context.Context, id string) (*neo4jentity.TenantBillingProfileEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("id", id))

	dbNode, err := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantBillingProfileById(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantBillingProfile: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantBillingProfileEntity(dbNode), nil
}

func (s *tenantService) CreateTenantBillingProfile(ctx context.Context, input model.TenantBillingProfileInput) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.CreateTenantBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	grpcRequest := tenantpb.AddBillingProfileRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
		Phone:                  utils.IfNotNilString(input.Phone),
		LegalName:              utils.IfNotNilString(input.LegalName),
		AddressLine1:           utils.IfNotNilString(input.AddressLine1),
		AddressLine2:           utils.IfNotNilString(input.AddressLine2),
		AddressLine3:           utils.IfNotNilString(input.AddressLine3),
		Locality:               utils.IfNotNilString(input.Locality),
		Country:                utils.IfNotNilString(input.Country),
		Region:                 utils.IfNotNilString(input.Region),
		Zip:                    utils.IfNotNilString(input.Zip),
		VatNumber:              utils.IfNotNilString(input.VatNumber),
		SendInvoicesFrom:       utils.IfNotNilString(input.SendInvoicesFrom),
		SendInvoicesBcc:        utils.IfNotNilString(input.SendInvoicesBcc),
		CanPayWithPigeon:       utils.IfNotNilBool(input.CanPayWithPigeon),
		CanPayWithBankTransfer: utils.IfNotNilBool(input.CanPayWithBankTransfer),
		Check:                  utils.IfNotNilBool(input.Check),
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*commonpb.IdResponse](func() (*commonpb.IdResponse, error) {
		return s.grpcClients.TenantClient.AddBillingProfile(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, model2.NodeLabelTenantBillingProfile, span)

	return response.Id, nil
}

func (s *tenantService) UpdateTenantBillingProfile(ctx context.Context, input model.TenantBillingProfileUpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.UpdateTenantBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	if input.ID == "" {
		err := fmt.Errorf("(TenantService.UpdateTenantBillingProfile) billing profile id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	billingProfileExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), input.ID, model2.NodeLabelTenantBillingProfile)
	if !billingProfileExists {
		err := fmt.Errorf("(TenantService.UpdateTenantBillingProfile) tenant billing profile with id {%s} not found", input.ID)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	var fieldsMask []tenantpb.TenantBillingProfileFieldMask
	updateRequest := tenantpb.UpdateBillingProfileRequest{
		Tenant:                 common.GetTenantFromContext(ctx),
		Id:                     input.ID,
		LoggedInUserId:         common.GetUserIdFromContext(ctx),
		AppSource:              constants.AppSourceCustomerOsApi,
		Phone:                  utils.IfNotNilString(input.Phone),
		LegalName:              utils.IfNotNilString(input.LegalName),
		AddressLine1:           utils.IfNotNilString(input.AddressLine1),
		AddressLine2:           utils.IfNotNilString(input.AddressLine2),
		AddressLine3:           utils.IfNotNilString(input.AddressLine3),
		Locality:               utils.IfNotNilString(input.Locality),
		Country:                utils.IfNotNilString(input.Country),
		Region:                 utils.IfNotNilString(input.Region),
		Zip:                    utils.IfNotNilString(input.Zip),
		VatNumber:              utils.IfNotNilString(input.VatNumber),
		SendInvoicesFrom:       utils.IfNotNilString(input.SendInvoicesFrom),
		SendInvoicesBcc:        utils.IfNotNilString(input.SendInvoicesBcc),
		CanPayWithPigeon:       utils.IfNotNilBool(input.CanPayWithPigeon),
		CanPayWithBankTransfer: utils.IfNotNilBool(input.CanPayWithBankTransfer),
		Check:                  utils.IfNotNilBool(input.Check),
	}

	if input.Patch != nil && *input.Patch {
		if input.Phone != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_PHONE)
		}
		if input.LegalName != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_LEGAL_NAME)
		}
		if input.AddressLine1 != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_1)
		}
		if input.AddressLine2 != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_2)
		}
		if input.AddressLine3 != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_3)
		}
		if input.Locality != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_LOCALITY)
		}
		if input.Country != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_COUNTRY)
		}
		if input.Region != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_REGION)
		}
		if input.Zip != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ZIP)
		}
		if input.VatNumber != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_VAT_NUMBER)
		}
		if input.SendInvoicesFrom != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_SEND_INVOICES_FROM)
		}
		if input.SendInvoicesBcc != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_SEND_INVOICES_BCC)
		}
		if input.CanPayWithPigeon != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_PIGEON)
		}
		if input.CanPayWithBankTransfer != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_BANK_TRANSFER)
		}
		if input.Check != nil {
			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CHECK)
		}
		if len(fieldsMask) == 0 {
			span.LogFields(log.String("result", "No fields to update"))
			return nil
		}
		updateRequest.FieldsMask = fieldsMask
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*commonpb.IdResponse](func() (*commonpb.IdResponse, error) {
		return s.grpcClients.TenantClient.UpdateBillingProfile(ctx, &updateRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *tenantService) GetTenantSettings(ctx context.Context) (*neo4jentity.TenantSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetTenantSettings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbNode, err := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode), nil
}

func (s *tenantService) UpdateTenantSettings(ctx context.Context, input *model.TenantSettingsInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.UpdateTenantSettings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	var baseCurrency string
	if input.BaseCurrency != nil {
		baseCurrency = input.BaseCurrency.String()
	}

	var fieldMask []tenantpb.TenantSettingsFieldMask
	updateRequest := tenantpb.UpdateTenantSettingsRequest{
		Tenant:               common.GetTenantFromContext(ctx),
		LoggedInUserId:       common.GetUserIdFromContext(ctx),
		AppSource:            constants.AppSourceCustomerOsApi,
		LogoRepositoryFileId: utils.IfNotNilString(input.LogoRepositoryFileID),
		WorkspaceLogo:        utils.IfNotNilString(input.WorkspaceLogo),
		WorkspaceName:        utils.IfNotNilString(input.WorkspaceName),
		BaseCurrency:         baseCurrency,
		InvoicingEnabled:     utils.IfNotNilBool(input.BillingEnabled),
	}

	if input.Patch != nil && *input.Patch {
		if input.LogoRepositoryFileID != nil {
			fieldMask = append(fieldMask, tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_LOGO_REPOSITORY_FILE_ID)
		}
		if input.BaseCurrency != nil {
			fieldMask = append(fieldMask, tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_BASE_CURRENCY)
		}
		if input.BillingEnabled != nil {
			fieldMask = append(fieldMask, tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_INVOICING_ENABLED)
		}
		if input.WorkspaceLogo != nil {
			fieldMask = append(fieldMask, tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_WORKSPACE_LOGO)
		}
		if input.WorkspaceName != nil {
			fieldMask = append(fieldMask, tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_WORKSPACE_NAME)
		}
		updateRequest.FieldsMask = fieldMask
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*emptypb.Empty](func() (*emptypb.Empty, error) {
		return s.grpcClients.TenantClient.UpdateTenantSettings(ctx, &updateRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *tenantService) HardDelete(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.HardDelete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "tenant", tenant)

	err := s.repositories.Neo4jRepositories.TenantWriteRepository.HardDeleteTenant(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *tenantService) GetDefaultTenantBillingProfile(ctx context.Context) (*neo4jentity.TenantBillingProfileEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.GetDefaultTenantBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenantBillingProfiles, err := s.GetTenantBillingProfiles(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetDefaultTenantBillingProfile: %w", err)
	}
	if tenantBillingProfiles == nil || len(*tenantBillingProfiles) == 0 {
		return nil, nil
	} else {
		return &(*tenantBillingProfiles)[0], nil
	}
}
