package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
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
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type ContractService interface {
	Create(ctx context.Context, contract *ContractCreateData) (string, error)
	Update(ctx context.Context, input model.ContractUpdateInput) error
	SoftDeleteContract(ctx context.Context, contractId string) (bool, error)
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
	ExternalReference *neo4jentity.ExternalSystemEntity
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
			AppSource: contractDetails.AppSource,
		},
		RenewalPeriods:   contractDetails.ContractEntity.RenewalPeriods,
		PayOnline:        true,
		PayAutomatically: true,
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
		// if not provided, get default currency from tenant settings
		dbNode, err := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, common.GetTenantFromContext(ctx))
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode)
		if tenantSettingsEntity.BaseCurrency.String() != "" {
			createContractRequest.Currency = tenantSettingsEntity.BaseCurrency.String()
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
	response, err := CallEventsPlatformGRPCWithRetry[*contractpb.ContractIdGrpcResponse](func() (*contractpb.ContractIdGrpcResponse, error) {
		return s.grpcClients.ContractClient.CreateContract(ctx, &createContractRequest)
	})

	WaitForNodeCreatedInNeo4j(ctx, s.repositories, response.Id, neo4jutil.NodeLabelContact, span)
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

	var fieldMask []contractpb.ContractFieldMask
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
	if input.BillingDetails != nil {
		if input.BillingDetails.CanPayWithCard != nil {
			contractUpdateRequest.CanPayWithCard = *input.BillingDetails.CanPayWithCard
		}
		if input.BillingDetails.CanPayWithDirectDebit != nil {
			contractUpdateRequest.CanPayWithDirectDebit = *input.BillingDetails.CanPayWithDirectDebit
		}
		if input.BillingDetails.CanPayWithBankTransfer != nil {
			contractUpdateRequest.CanPayWithBankTransfer = *input.BillingDetails.CanPayWithBankTransfer
		}
		if input.BillingDetails.AddressLine1 != nil {
			contractUpdateRequest.AddressLine1 = *input.BillingDetails.AddressLine1
		}
		if input.BillingDetails.AddressLine2 != nil {
			contractUpdateRequest.AddressLine2 = *input.BillingDetails.AddressLine2
		}
		if input.BillingDetails.Locality != nil {
			contractUpdateRequest.Locality = *input.BillingDetails.Locality
		}
		if input.BillingDetails.Country != nil {
			contractUpdateRequest.Country = *input.BillingDetails.Country
		}
		if input.BillingDetails.PostalCode != nil {
			contractUpdateRequest.Zip = *input.BillingDetails.PostalCode
		}
		if input.BillingDetails.OrganizationLegalName != nil {
			contractUpdateRequest.OrganizationLegalName = *input.BillingDetails.OrganizationLegalName
		}
		if input.BillingDetails.BillingEmail != nil {
			contractUpdateRequest.InvoiceEmail = *input.BillingDetails.BillingEmail
		}
		if input.BillingDetails.InvoiceNote != nil {
			contractUpdateRequest.InvoiceNote = *input.BillingDetails.InvoiceNote
		}
		contractUpdateRequest.PayOnline = utils.IfNotNilBool(input.BillingDetails.PayOnline)
		contractUpdateRequest.PayAutomatically = utils.IfNotNilBool(input.BillingDetails.PayAutomatically)
	}
	if input.CommittedPeriods != nil {
		contractUpdateRequest.RenewalPeriods = input.CommittedPeriods
	}
	if input.ContractName != nil {
		contractUpdateRequest.Name = *input.ContractName
	}

	nullTime := time.Time{}

	if input.SignedAt != nil {
		if *input.SignedAt != nullTime {
			contractUpdateRequest.SignedAt = utils.ConvertTimeToTimestampPtr(input.SignedAt)
		} else {
			contractUpdateRequest.SignedAt = nil
		}
	}
	if input.ContractSigned != nil {
		if *input.ContractSigned != nullTime {
			contractUpdateRequest.SignedAt = utils.ConvertTimeToTimestampPtr(input.ContractSigned)
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
	if input.ServiceStarted != nil {
		if *input.ServiceStarted != nullTime {
			contractUpdateRequest.ServiceStartedAt = utils.ConvertTimeToTimestampPtr(input.ServiceStarted)
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
	if input.ContractEnded != nil {
		if *input.ContractEnded != nullTime {
			contractUpdateRequest.EndedAt = utils.ConvertTimeToTimestampPtr(input.ContractEnded)
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
	if input.BillingDetails != nil && input.BillingDetails.InvoicingStarted != nil {
		if *input.BillingDetails.InvoicingStarted != nullTime {
			contractUpdateRequest.InvoicingStartDate = utils.ConvertTimeToTimestampPtr(input.BillingDetails.InvoicingStarted)
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
	if input.ContractRenewalCycle != nil {
		switch *input.ContractRenewalCycle {
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
	if input.BillingDetails != nil && input.BillingDetails.BillingCycle != nil {
		switch *input.BillingDetails.BillingCycle {
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
		if input.Name != nil || input.ContractName != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_NAME)
		}
		if input.ContractURL != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CONTRACT_URL)
		}
		if input.RenewalCycle != nil || input.ContractRenewalCycle != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_RENEWAL_CYCLE)
		}
		if input.RenewalPeriods != nil || input.CommittedPeriods != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_RENEWAL_PERIODS)
		}
		if input.ServiceStartedAt != nil || input.ServiceStarted != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_SERVICE_STARTED_AT)
		}
		if input.SignedAt != nil || input.ContractSigned != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_SIGNED_AT)
		}
		if input.EndedAt != nil || input.ContractEnded != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ENDED_AT)
		}
		if input.InvoicingStartDate != nil || (input.BillingDetails != nil && input.BillingDetails.InvoicingStarted != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICING_START_DATE)
		}
		if input.Currency != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CURRENCY)
		}
		if input.BillingCycle != nil || (input.BillingDetails != nil && input.BillingDetails.BillingCycle != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_BILLING_CYCLE)
		}
		if input.AddressLine1 != nil || (input.BillingDetails != nil && input.BillingDetails.AddressLine1 != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ADDRESS_LINE_1)
		}
		if input.AddressLine2 != nil || (input.BillingDetails != nil && input.BillingDetails.AddressLine2 != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ADDRESS_LINE_2)
		}
		if input.Locality != nil || (input.BillingDetails != nil && input.BillingDetails.Locality != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_LOCALITY)
		}
		if input.Country != nil || (input.BillingDetails != nil && input.BillingDetails.Country != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_COUNTRY)
		}
		if input.Zip != nil || (input.BillingDetails != nil && input.BillingDetails.PostalCode != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ZIP)
		}
		if input.OrganizationLegalName != nil || (input.BillingDetails != nil && input.BillingDetails.OrganizationLegalName != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ORGANIZATION_LEGAL_NAME)
		}
		if input.InvoiceEmail != nil || (input.BillingDetails != nil && input.BillingDetails.BillingEmail != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL)
		}
		if input.InvoiceNote != nil || (input.BillingDetails != nil && input.BillingDetails.InvoiceNote != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_NOTE)
		}
		if input.CanPayWithCard != nil || (input.BillingDetails != nil && input.BillingDetails.CanPayWithCard != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_CARD)
		}
		if input.CanPayWithDirectDebit != nil || (input.BillingDetails != nil && input.BillingDetails.CanPayWithDirectDebit != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_DIRECT_DEBIT)
		}
		if input.CanPayWithBankTransfer != nil || (input.BillingDetails != nil && input.BillingDetails.CanPayWithBankTransfer != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_BANK_TRANSFER)
		}
		if input.BillingDetails != nil && input.BillingDetails.PayOnline != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_PAY_ONLINE)
		}
		if input.BillingDetails != nil && input.BillingDetails.PayAutomatically != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_PAY_AUTOMATICALLY)
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
	_, err := CallEventsPlatformGRPCWithRetry[*contractpb.ContractIdGrpcResponse](func() (*contractpb.ContractIdGrpcResponse, error) {
		return s.grpcClients.ContractClient.UpdateContract(ctx, &contractUpdateRequest)
	})
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

func (s *contractService) SoftDeleteContract(ctx context.Context, contractId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.SoftDeleteContract")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	// check contract exists
	contractExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), contractId, neo4jutil.NodeLabelContract)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error on checking if contract exists: %s", err.Error())
		return false, err
	}
	if !contractExists {
		err := fmt.Errorf("contract with id {%s} not found", contractId)
		tracing.TraceErr(span, err)
		s.log.Errorf(err.Error())
		return false, err
	}

	// check contract has no invoices
	countInvoices, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.CountNonDryRunInvoicesForContract(ctx, common.GetTenantFromContext(ctx), contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error on counting invoices for contract: %s", err.Error())
		return false, err
	}
	if countInvoices > 0 {
		err := fmt.Errorf("contract with id {%s} has invoices", contractId)
		tracing.TraceErr(span, err)
		s.log.Errorf(err.Error())
		return false, err
	}

	deleteRequest := contractpb.SoftDeleteContractGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             contractId,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		AppSource:      constants.AppSourceCustomerOsApi,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = CallEventsPlatformGRPCWithRetry[*emptypb.Empty](func() (*emptypb.Empty, error) {
		return s.grpcClients.ContractClient.SoftDeleteContract(ctx, &deleteRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return false, err
	}

	// wait for contract to be deleted from graph db
	WaitForNodeDeletedFromNeo4j(ctx, s.repositories, contractId, neo4jutil.NodeLabelContract, span)

	return false, nil
}
