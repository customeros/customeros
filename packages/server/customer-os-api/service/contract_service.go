package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type ContractService interface {
	Create(ctx context.Context, contract *ContractCreateData) (string, error)
	Update(ctx context.Context, input model.ContractUpdateInput) error
	GetById(ctx context.Context, id string) (*neo4jentity.ContractEntity, error)
	GetContractsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.ContractEntities, error)
	GetContractsForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.ContractEntities, error)
	ContractsExistForTenant(ctx context.Context) (bool, error)
	CountContracts(ctx context.Context, tenant string) (int64, error)
}
type contractService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewContractService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) ContractService {
	return &contractService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

type ContractCreateData struct {
	ContractEntity    *neo4jentity.ContractEntity
	OrganizationId    string
	ExternalReference *entity.ExternalSystemEntity
	Source            neo4jentity.DataSource
	AppSource         string
}

func (s *contractService) Create(ctx context.Context, contractDetails *ContractCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contractDetails", contractDetails))

	if contractDetails.ContractEntity == nil {
		err := fmt.Errorf("contract entity is nil")
		tracing.TraceErr(span, err)
		return "", err
	}

	contractId, err := s.createContractWithEvents(ctx, contractDetails)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	span.LogFields(log.String("output - createdContractId", contractId))
	return contractId, nil
}

func (s *contractService) createContractWithEvents(ctx context.Context, contractDetails *ContractCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.createContractWithEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	createContractRequest := contractpb.CreateContractGrpcRequest{
		Tenant:             common.GetTenantFromContext(ctx),
		OrganizationId:     contractDetails.OrganizationId,
		Name:               contractDetails.ContractEntity.Name,
		ContractUrl:        contractDetails.ContractEntity.ContractUrl,
		SignedAt:           utils.ConvertTimeToTimestampPtr(contractDetails.ContractEntity.SignedAt),
		ServiceStartedAt:   utils.ConvertTimeToTimestampPtr(contractDetails.ContractEntity.ServiceStartedAt),
		InvoicingStartDate: utils.ConvertTimeToTimestampPtr(contractDetails.ContractEntity.InvoicingStartDate),
		InvoicingEnabled:   contractDetails.ContractEntity.InvoicingEnabled,
		LoggedInUserId:     common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    string(contractDetails.Source),
			AppSource: utils.StringFirstNonEmpty(contractDetails.AppSource, constants.AppSourceCustomerOsApi),
		},
		RenewalPeriods: contractDetails.ContractEntity.RenewalPeriods,
	}

	// prepare renewal cycle
	switch contractDetails.ContractEntity.RenewalCycle {
	case neo4jenum.RenewalCycleMonthlyRenewal:
		createContractRequest.RenewalCycle = contractpb.RenewalCycle_MONTHLY_RENEWAL
	case neo4jenum.RenewalCycleQuarterlyRenewal:
		createContractRequest.RenewalCycle = contractpb.RenewalCycle_QUARTERLY_RENEWAL
	case neo4jenum.RenewalCycleAnnualRenewal:
		createContractRequest.RenewalCycle = contractpb.RenewalCycle_ANNUALLY_RENEWAL
	default:
		createContractRequest.RenewalCycle = contractpb.RenewalCycle_NONE
	}

	// prepare billing cycle
	switch contractDetails.ContractEntity.BillingCycle {
	case neo4jenum.BillingCycleMonthlyBilling:
		createContractRequest.BillingCycle = commonpb.BillingCycle_MONTHLY_BILLING
	case neo4jenum.BillingCycleQuarterlyBilling:
		createContractRequest.BillingCycle = commonpb.BillingCycle_QUARTERLY_BILLING
	case neo4jenum.BillingCycleAnnuallyBilling:
		createContractRequest.BillingCycle = commonpb.BillingCycle_ANNUALLY_BILLING
	default:
		createContractRequest.BillingCycle = commonpb.BillingCycle_NONE_BILLING
	}

	// prepare currency
	if contractDetails.ContractEntity.Currency.String() != "" {
		createContractRequest.Currency = contractDetails.ContractEntity.Currency.String()
	} else {
		// if not privided, get default currency from tenant settings
		dbNode, err := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, common.GetTenantFromContext(ctx))
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode)
		if tenantSettingsEntity.DefaultCurrency.String() != "" {
			createContractRequest.Currency = tenantSettingsEntity.DefaultCurrency.String()
		}
	}

	// prepare external system fields
	if contractDetails.ExternalReference != nil && contractDetails.ExternalReference.ExternalSystemId != "" {
		createContractRequest.ExternalSystemFields = &commonpb.ExternalSystemFields{
			ExternalSystemId: string(contractDetails.ExternalReference.ExternalSystemId),
			ExternalId:       contractDetails.ExternalReference.Relationship.ExternalId,
			ExternalUrl:      utils.IfNotNilString(contractDetails.ExternalReference.Relationship.ExternalUrl),
			ExternalSource:   utils.IfNotNilString(contractDetails.ExternalReference.Relationship.ExternalSource),
		}
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := s.grpcClients.ContractClient.CreateContract(ctx, &createContractRequest)

	WaitForObjectCreationAndLogSpan(ctx, s.repositories, response.Id, neo4jutil.NodeLabelContact, span)
	return response.Id, err
}

func (s *contractService) Update(ctx context.Context, input model.ContractUpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	if input.ContractID == "" {
		err := fmt.Errorf("(ContractService.Update) contract id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	contractExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), input.ContractID, neo4jutil.NodeLabelContract)
	if !contractExists {
		err := fmt.Errorf("(ContractService.Update) contract with id {%s} not found", input.ContractID)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	fieldMask := []contractpb.ContractFieldMask{}
	contractUpdateRequest := contractpb.UpdateContractGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             input.ContractID,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		Name:           utils.IfNotNilString(input.Name),
		ContractUrl:    utils.IfNotNilString(input.ContractURL),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: utils.StringFirstNonEmpty(utils.IfNotNilString(input.AppSource), constants.AppSourceCustomerOsApi),
		},
		RenewalPeriods:         input.RenewalPeriods,
		AddressLine1:           utils.IfNotNilString(input.AddressLine1),
		AddressLine2:           utils.IfNotNilString(input.AddressLine2),
		Locality:               utils.IfNotNilString(input.Locality),
		Country:                utils.IfNotNilString(input.Country),
		Zip:                    utils.IfNotNilString(input.Zip),
		OrganizationLegalName:  utils.IfNotNilString(input.OrganizationLegalName),
		InvoiceEmail:           utils.IfNotNilString(input.InvoiceEmail),
		InvoiceNote:            utils.IfNotNilString(input.InvoiceNote),
		CanPayWithCard:         utils.IfNotNilBool(input.CanPayWithCard),
		CanPayWithDirectDebit:  utils.IfNotNilBool(input.CanPayWithDirectDebit),
		CanPayWithBankTransfer: utils.IfNotNilBool(input.CanPayWithBankTransfer),
		InvoicingEnabled:       utils.IfNotNilBool(input.BillingEnabled),
	}
	if input.Currency != nil {
		contractUpdateRequest.Currency = mapper.MapCurrencyFromModel(*input.Currency).String()
	}

	nullTime := time.Time{}

	if input.SignedAt != nil {
		if *input.SignedAt != nullTime {
			contractUpdateRequest.SignedAt = utils.ConvertTimeToTimestampPtr(input.SignedAt)
		} else {
			contractUpdateRequest.SignedAt = nil
		}
	}

	if input.ServiceStartedAt != nil {
		if *input.ServiceStartedAt != nullTime {
			contractUpdateRequest.ServiceStartedAt = utils.ConvertTimeToTimestampPtr(input.ServiceStartedAt)
		} else {
			contractUpdateRequest.ServiceStartedAt = nil
		}
	}

	if input.EndedAt != nil {
		if *input.EndedAt != nullTime {
			contractUpdateRequest.EndedAt = utils.ConvertTimeToTimestampPtr(input.EndedAt)
		} else {
			contractUpdateRequest.EndedAt = nil
		}
	}

	if input.InvoicingStartDate != nil {
		if *input.InvoicingStartDate != nullTime {
			contractUpdateRequest.InvoicingStartDate = utils.ConvertTimeToTimestampPtr(input.InvoicingStartDate)
		} else {
			contractUpdateRequest.InvoicingStartDate = nil
		}
	}

	// prepare renewal cycle
	if input.RenewalCycle != nil {
		switch *input.RenewalCycle {
		case model.ContractRenewalCycleMonthlyRenewal:
			contractUpdateRequest.RenewalCycle = contractpb.RenewalCycle_MONTHLY_RENEWAL
		case model.ContractRenewalCycleQuarterlyRenewal:
			contractUpdateRequest.RenewalCycle = contractpb.RenewalCycle_QUARTERLY_RENEWAL
		case model.ContractRenewalCycleAnnualRenewal:
			contractUpdateRequest.RenewalCycle = contractpb.RenewalCycle_ANNUALLY_RENEWAL
		default:
			contractUpdateRequest.RenewalCycle = contractpb.RenewalCycle_NONE
		}
	}

	// prepare billing cycle
	if input.BillingCycle != nil {
		switch *input.BillingCycle {
		case model.ContractBillingCycleMonthlyBilling:
			contractUpdateRequest.BillingCycle = commonpb.BillingCycle_MONTHLY_BILLING
		case model.ContractBillingCycleQuarterlyBilling:
			contractUpdateRequest.BillingCycle = commonpb.BillingCycle_QUARTERLY_BILLING
		case model.ContractBillingCycleAnnualBilling:
			contractUpdateRequest.BillingCycle = commonpb.BillingCycle_ANNUALLY_BILLING
		default:
			contractUpdateRequest.BillingCycle = commonpb.BillingCycle_NONE_BILLING
		}
	}

	if input.Patch != nil && *input.Patch {
		if input.Name != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_NAME)
		}
		if input.ContractURL != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CONTRACT_URL)
		}
		if input.RenewalCycle != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_RENEWAL_CYCLE)
		}
		if input.RenewalPeriods != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_RENEWAL_PERIODS)
		}
		if input.ServiceStartedAt != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_SERVICE_STARTED_AT)
		}
		if input.SignedAt != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_SIGNED_AT)
		}
		if input.EndedAt != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ENDED_AT)
		}
		if input.InvoicingStartDate != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICING_START_DATE)
		}
		if input.Currency != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CURRENCY)
		}
		if input.BillingCycle != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_BILLING_CYCLE)
		}
		if input.AddressLine1 != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ADDRESS_LINE_1)
		}
		if input.AddressLine2 != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ADDRESS_LINE_2)
		}
		if input.Locality != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_LOCALITY)
		}
		if input.Country != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_COUNTRY)
		}
		if input.Zip != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ZIP)
		}
		if input.OrganizationLegalName != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ORGANIZATION_LEGAL_NAME)
		}
		if input.InvoiceEmail != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL)
		}
		if input.InvoiceNote != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_NOTE)
		}
		if input.CanPayWithCard != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_CARD)
		}
		if input.CanPayWithDirectDebit != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_DIRECT_DEBIT)
		}
		if input.CanPayWithBankTransfer != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_BANK_TRANSFER)
		}
		if input.BillingEnabled != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICING_ENABLED)
		}
		contractUpdateRequest.FieldsMask = fieldMask
		if len(fieldMask) == 0 {
			span.LogFields(log.String("result", "No fields to update"))
			return nil
		}
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := s.grpcClients.ContractClient.UpdateContract(ctx, &contractUpdateRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *contractService) GetById(ctx context.Context, contractId string) (*neo4jentity.ContractEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contractId", contractId))

	if contractDbNode, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, common.GetContext(ctx).Tenant, contractId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Contract with id {%s} not found", contractId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToContractEntity(contractDbNode), nil
	}
}

func (s *contractService) GetContractsForOrganizations(ctx context.Context, organizationIDs []string) (*neo4jentity.ContractEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetContractsForOrganizations")
	defer span.Finish()
	span.LogFields(log.Object("organizationIDs", organizationIDs))

	contracts, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractsForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
	if err != nil {
		return nil, err
	}
	contractEntities := make(neo4jentity.ContractEntities, 0, len(contracts))
	for _, v := range contracts {
		contractEntity := neo4jmapper.MapDbNodeToContractEntity(v.Node)
		contractEntity.DataloaderKey = v.LinkedNodeId
		contractEntities = append(contractEntities, *contractEntity)
	}
	return &contractEntities, nil
}

func (s *contractService) GetContractsForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.ContractEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetContractsForInvoices")
	defer span.Finish()
	span.LogFields(log.Object("invoiceIds", invoiceIds))

	contracts, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractsForInvoices(ctx, common.GetTenantFromContext(ctx), invoiceIds)
	if err != nil {
		return nil, err
	}
	contractEntities := make(neo4jentity.ContractEntities, 0, len(contracts))
	for _, v := range contracts {
		contractEntity := neo4jmapper.MapDbNodeToContractEntity(v.Node)
		contractEntity.DataloaderKey = v.LinkedNodeId
		contractEntities = append(contractEntities, *contractEntity)
	}
	return &contractEntities, nil
}

func (s *contractService) ContractsExistForTenant(ctx context.Context) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.ContractsExistForTenant")
	defer span.Finish()

	contractsExistForTenant, err := s.repositories.Neo4jRepositories.ContractReadRepository.TenantsHasAtLeastOneContract(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		return false, err
	}
	return contractsExistForTenant, nil
}

func (s *contractService) CountContracts(ctx context.Context, tenant string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.CountContracts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("tenant", tenant))

	return s.repositories.Neo4jRepositories.ContractReadRepository.CountContracts(ctx, tenant)
}
