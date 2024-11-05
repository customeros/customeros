package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type ContractService interface {
	Create(ctx context.Context, contractDetails *ContractCreateData) (string, error)
	Update(ctx context.Context, input model.ContractUpdateInput) error
	SoftDeleteContract(ctx context.Context, contractId string) (bool, error)
	GetById(ctx context.Context, id string) (*neo4jentity.ContractEntity, error)
	GetContractsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.ContractEntities, error)
	GetContractsForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.ContractEntities, error)
	GetContractByServiceLineItem(ctx context.Context, serviceLineItemId string) (*neo4jentity.ContractEntity, error)
	ContractsExistForTenant(ctx context.Context) (bool, error)
	CountContracts(ctx context.Context, tenant string) (int64, error)
	RenewContract(ctx context.Context, contractId string, renewalDate *time.Time) error
	GetPaginatedContracts(ctx context.Context, page int, limit int) (*utils.Pagination, error)
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
	Input             model.ContractInput
	ExternalReference *neo4jentity.ExternalSystemEntity
	Source            neo4jentity.DataSource
	AppSource         string
}

func (s *contractService) Create(ctx context.Context, contractDetails *ContractCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "contractDetails", contractDetails)

	contractDataFields := data_fields.ContractSaveFields{
		OrganizationId:         &contractDetails.Input.OrganizationID,
		Name:                   contractDetails.Input.ContractName,
		ContractUrl:            contractDetails.Input.ContractURL,
		InvoicingEnabled:       contractDetails.Input.BillingEnabled,
		Source:                 utils.StringPtr(contractDetails.Source.String()),
		AppSource:              utils.StringPtr(contractDetails.AppSource),
		PayOnline:              utils.BoolPtr(true),
		PayAutomatically:       utils.BoolPtr(true),
		CanPayWithCard:         utils.BoolPtr(true),
		CanPayWithDirectDebit:  utils.BoolPtr(true),
		CanPayWithBankTransfer: utils.BoolPtr(true),
		Check:                  utils.BoolPtr(true),
		AutoRenew:              contractDetails.Input.AutoRenew,
		BillingCycleInMonths:   utils.Int64Ptr(1),
		DueDays:                contractDetails.Input.DueDays,
		Approved:               contractDetails.Input.Approved,
	}
	if common.GetUserIdFromContext(ctx) != "" {
		contractDataFields.CreatedByUserId = utils.StringPtr(common.GetUserIdFromContext(ctx))
	}
	if contractDetails.Input.ContractSigned != nil {
		contractDataFields.SignedAt = contractDetails.Input.ContractSigned
	} else if contractDetails.Input.SignedAt != nil {
		contractDataFields.SignedAt = contractDetails.Input.SignedAt
	}
	if contractDetails.Input.ServiceStarted != nil {
		contractDataFields.ServiceStartedAt = contractDetails.Input.ServiceStarted
	} else if contractDetails.Input.ServiceStartedAt != nil {
		contractDataFields.ServiceStartedAt = contractDetails.Input.ServiceStartedAt
	}
	if contractDetails.Input.InvoicingStartDate != nil {
		contractDataFields.InvoicingStartDate = contractDetails.Input.InvoicingStartDate
	}

	if contractDetails.Input.CommittedPeriodInMonths != nil {
		contractDataFields.LengthInMonths = contractDetails.Input.CommittedPeriodInMonths
	} else {
		renewalCycle := ""
		if contractDetails.Input.ContractRenewalCycle != nil {
			renewalCycle = contractDetails.Input.ContractRenewalCycle.String()
		} else if contractDetails.Input.RenewalCycle != nil {
			renewalCycle = contractDetails.Input.RenewalCycle.String()
		}
		switch renewalCycle {
		case model.ContractRenewalCycleMonthlyRenewal.String():
			contractDataFields.LengthInMonths = utils.Int64Ptr(1)
		case model.ContractRenewalCycleQuarterlyRenewal.String():
			contractDataFields.LengthInMonths = utils.Int64Ptr(3)
		case model.ContractRenewalCycleAnnualRenewal.String():
			contractDataFields.LengthInMonths = utils.Int64Ptr(12)
		default:
			contractDataFields.LengthInMonths = utils.Int64Ptr(0)
		}
		if *contractDataFields.LengthInMonths == 12 {
			if contractDetails.Input.CommittedPeriods != nil && *contractDetails.Input.CommittedPeriods > 1 {
				contractDataFields.LengthInMonths = utils.Int64Ptr(*contractDataFields.LengthInMonths * *contractDetails.Input.CommittedPeriods)
			} else if contractDetails.Input.RenewalPeriods != nil && *contractDetails.Input.RenewalPeriods > 1 {
				contractDataFields.LengthInMonths = utils.Int64Ptr(*contractDataFields.LengthInMonths * *contractDetails.Input.RenewalPeriods)
			}
		}
	}

	// set default fields
	// set currency
	if contractDetails.Input.Currency != nil && contractDetails.Input.Currency.String() != "" {
		contractDataFields.Currency = utils.ToPtr(neo4jenum.DecodeCurrency(contractDetails.Input.Currency.String()))
	} else {
		// if not provided, get default currency from tenant settings
		tenantSettingsEntity, err := s.services.CommonServices.TenantService.GetTenantSettings(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		if tenantSettingsEntity.BaseCurrency.String() != "" {
			contractDataFields.Currency = &tenantSettingsEntity.BaseCurrency
		}
	}

	tenantBillingProfileEntity, err := s.services.CommonServices.TenantService.GetDefaultTenantBillingProfile(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if tenantBillingProfileEntity != nil {
		// set country
		contractDataFields.Country = &tenantBillingProfileEntity.Country
	}

	// prepare external system fields
	if contractDetails.ExternalReference != nil && contractDetails.ExternalReference.ExternalSystemId != "" {
		contractDataFields.ExternalSystem = &neo4jmodel.ExternalSystem{
			ExternalSystemId: string(contractDetails.ExternalReference.ExternalSystemId),
			ExternalId:       contractDetails.ExternalReference.Relationship.ExternalId,
			ExternalUrl:      utils.IfNotNilString(contractDetails.ExternalReference.Relationship.ExternalUrl),
			ExternalSource:   utils.IfNotNilString(contractDetails.ExternalReference.Relationship.ExternalSource),
		}
	}

	contractId, err := s.services.CommonServices.ContractService.Save(ctx, nil, contractDataFields)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from create contract %s", err.Error())
		return "", err
	}

	span.LogFields(log.String("output - createdContractId", contractId))
	return contractId, nil
}

func (s *contractService) Update(ctx context.Context, input model.ContractUpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	if err := s.validateContractExists(ctx, input.ContractID, span); err != nil {
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
		AddressLine1:          utils.IfNotNilString(input.AddressLine1),
		AddressLine2:          utils.IfNotNilString(input.AddressLine2),
		Locality:              utils.IfNotNilString(input.Locality),
		Country:               utils.IfNotNilString(input.Country),
		Zip:                   utils.IfNotNilString(input.Zip),
		OrganizationLegalName: utils.IfNotNilString(input.OrganizationLegalName),
		InvoiceEmailTo:        utils.IfNotNilString(input.InvoiceEmail),
		InvoiceNote:           utils.IfNotNilString(input.InvoiceNote),
		InvoicingEnabled:      utils.IfNotNilBool(input.BillingEnabled),
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
		if input.BillingDetails.Check != nil {
			contractUpdateRequest.Check = *input.BillingDetails.Check
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
		if input.BillingDetails.Region != nil {
			contractUpdateRequest.Region = *input.BillingDetails.Region
		}
		if input.BillingDetails.PostalCode != nil {
			contractUpdateRequest.Zip = *input.BillingDetails.PostalCode
		}
		if input.BillingDetails.OrganizationLegalName != nil {
			contractUpdateRequest.OrganizationLegalName = *input.BillingDetails.OrganizationLegalName
		}
		if input.BillingDetails.BillingEmail != nil {
			contractUpdateRequest.InvoiceEmailTo = *input.BillingDetails.BillingEmail
		}
		if input.BillingDetails.BillingEmailCc != nil {
			contractUpdateRequest.InvoiceEmailCc = input.BillingDetails.BillingEmailCc
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL_CC)
		}
		if input.BillingDetails.BillingEmailBcc != nil {
			contractUpdateRequest.InvoiceEmailBcc = input.BillingDetails.BillingEmailBcc
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL_BCC)
		}
		if input.BillingDetails.InvoiceNote != nil {
			contractUpdateRequest.InvoiceNote = *input.BillingDetails.InvoiceNote
		}
		if input.BillingDetails.BillingCycleInMonths != nil {
			contractUpdateRequest.BillingCycleInMonths = *input.BillingDetails.BillingCycleInMonths
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_BILLING_CYCLE_IN_MONTHS)
		}
		contractUpdateRequest.PayOnline = utils.IfNotNilBool(input.BillingDetails.PayOnline)
		contractUpdateRequest.PayAutomatically = utils.IfNotNilBool(input.BillingDetails.PayAutomatically)
	}

	if input.Approved != nil {
		contractUpdateRequest.Approved = *input.Approved
		fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_APPROVED)
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

	if input.CommittedPeriodInMonths != nil {
		contractUpdateRequest.LengthInMonths = *input.CommittedPeriodInMonths
		fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_LENGTH_IN_MONTHS)
	} else {
		// prepare length in months from renewal cycle and periods
		renewalCycle := ""
		if input.ContractRenewalCycle != nil {
			renewalCycleEnum := *input.ContractRenewalCycle
			renewalCycle = renewalCycleEnum.String()
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_LENGTH_IN_MONTHS)
		} else if input.RenewalCycle != nil {
			renewalCycleEnum := *input.RenewalCycle
			renewalCycle = renewalCycleEnum.String()
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_LENGTH_IN_MONTHS)
		}
		switch renewalCycle {
		case model.ContractRenewalCycleMonthlyRenewal.String():
			contractUpdateRequest.LengthInMonths = 1
		case model.ContractRenewalCycleQuarterlyRenewal.String():
			contractUpdateRequest.LengthInMonths = 3
		case model.ContractRenewalCycleAnnualRenewal.String():
			contractUpdateRequest.LengthInMonths = 12
		default:
			contractUpdateRequest.LengthInMonths = 0
		}
		if contractUpdateRequest.LengthInMonths == 12 {
			if input.CommittedPeriods != nil && *input.CommittedPeriods > 1 {
				contractUpdateRequest.LengthInMonths *= *input.CommittedPeriods
			} else if input.RenewalPeriods != nil && *input.RenewalPeriods > 1 {
				contractUpdateRequest.LengthInMonths *= *input.RenewalPeriods
			}
		}
	}

	if input.BillingDetails == nil || input.BillingDetails.BillingCycleInMonths == nil {
		if input.BillingDetails != nil && input.BillingDetails.BillingCycle != nil {
			switch *input.BillingDetails.BillingCycle {
			case model.ContractBillingCycleMonthlyBilling:
				contractUpdateRequest.BillingCycleInMonths = 1
			case model.ContractBillingCycleQuarterlyBilling:
				contractUpdateRequest.BillingCycleInMonths = 3
			case model.ContractBillingCycleAnnualBilling:
				contractUpdateRequest.BillingCycleInMonths = 12
			default:
				contractUpdateRequest.BillingCycleInMonths = 0
			}
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_BILLING_CYCLE_IN_MONTHS)
		} else if input.BillingCycle != nil {
			switch *input.BillingCycle {
			case model.ContractBillingCycleMonthlyBilling:
				contractUpdateRequest.BillingCycleInMonths = 1
			case model.ContractBillingCycleQuarterlyBilling:
				contractUpdateRequest.BillingCycleInMonths = 3
			case model.ContractBillingCycleAnnualBilling:
				contractUpdateRequest.BillingCycleInMonths = 12
			default:
				contractUpdateRequest.BillingCycleInMonths = 0
			}
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_BILLING_CYCLE_IN_MONTHS)
		}
	}

	if input.Patch != nil && *input.Patch {
		if input.Name != nil || input.ContractName != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_NAME)
		}
		if input.ContractURL != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CONTRACT_URL)
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
		if input.BillingDetails != nil && input.BillingDetails.Region != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_REGION)
		}
		if input.Zip != nil || (input.BillingDetails != nil && input.BillingDetails.PostalCode != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ZIP)
		}
		if input.OrganizationLegalName != nil || (input.BillingDetails != nil && input.BillingDetails.OrganizationLegalName != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_ORGANIZATION_LEGAL_NAME)
		}
		if input.InvoiceEmail != nil || (input.BillingDetails != nil && input.BillingDetails.BillingEmail != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL_TO)
		}
		if input.InvoiceNote != nil || (input.BillingDetails != nil && input.BillingDetails.InvoiceNote != nil) {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_NOTE)
		}
		if input.BillingDetails != nil && input.BillingDetails.CanPayWithCard != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_CARD)
		}
		if input.BillingDetails != nil && input.BillingDetails.CanPayWithDirectDebit != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_DIRECT_DEBIT)
		}
		if input.BillingDetails != nil && input.BillingDetails.CanPayWithBankTransfer != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_BANK_TRANSFER)
		}
		if input.BillingDetails != nil && input.BillingDetails.Check != nil {
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_CHECK)
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
		if input.AutoRenew != nil {
			contractUpdateRequest.AutoRenew = *input.AutoRenew
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_AUTO_RENEW)
		}
		if input.BillingDetails != nil && input.BillingDetails.DueDays != nil {
			contractUpdateRequest.DueDays = *input.BillingDetails.DueDays
			fieldMask = append(fieldMask, contractpb.ContractFieldMask_CONTRACT_FIELD_DUE_DAYS)
		}
		contractUpdateRequest.FieldsMask = fieldMask
		if len(fieldMask) == 0 {
			span.LogFields(log.String("result", "No fields to update"))
			return nil
		}
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*contractpb.ContractIdGrpcResponse](func() (*contractpb.ContractIdGrpcResponse, error) {
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
	span.SetTag(tracing.SpanTagTenant, tenant)

	return s.repositories.Neo4jRepositories.ContractReadRepository.CountContracts(ctx, tenant)
}

func (s *contractService) SoftDeleteContract(ctx context.Context, contractId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.SoftDeleteContract")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	// check contract exists
	if err := s.validateContractExists(ctx, contractId, span); err != nil {
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
	_, err = utils.CallEventsPlatformGRPCWithRetry[*emptypb.Empty](func() (*emptypb.Empty, error) {
		return s.grpcClients.ContractClient.SoftDeleteContract(ctx, &deleteRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return false, err
	}

	// wait for contract to be deleted from graph db
	neo4jrepository.WaitForNodeDeletedFromNeo4j(ctx, s.repositories.Neo4jRepositories, contractId, model2.NodeLabelContract, span)

	return false, nil
}

func (s *contractService) RenewContract(ctx context.Context, contractId string, renewalDate *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.RenewContract")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, contractId)
	if renewalDate != nil {
		span.LogFields(log.Object("renewalDate", renewalDate.String()))
	}

	// check contract exists
	if err := s.validateContractExists(ctx, contractId, span); err != nil {
		return err
	}

	contractEntity, err := s.GetById(ctx, contractId)
	if err != nil {
		return err
	}

	// if contract is not renewable - return
	if contractEntity.LengthInMonths == 0 {
		span.LogFields(log.Bool("result.contractRenewable", false))
		return nil
	}

	opportunityDbNode, err := s.repositories.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, common.GetTenantFromContext(ctx), contractId)
	if err != nil {
		return err
	}
	// if no active renewal opportunity found create new
	if opportunityDbNode == nil {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err := utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
			return s.grpcClients.OpportunityClient.CreateRenewalOpportunity(ctx, &opportunitypb.CreateRenewalOpportunityGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				ContractId:     contractId,
				SourceFields: &commonpb.SourceFields{
					Source:    neo4jentity.DataSourceOpenline.String(),
					AppSource: constants.AppSourceCustomerOsApi,
				},
				RenewalApproved: true,
				RenewedAt:       utils.ConvertTimeToTimestampPtr(renewalDate),
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error from events processing: %s", err.Error())
			return err
		}
		return nil
	}
	opportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)

	// if renewal opportunity is not expired - approve next renewal
	if opportunityEntity.RenewalDetails.RenewedAt != nil && utils.Now().Before(*opportunityEntity.RenewalDetails.RenewedAt) {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err := utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
			grpcUpdateRequest := opportunitypb.UpdateRenewalOpportunityGrpcRequest{
				Id:              opportunityEntity.Id,
				Tenant:          common.GetTenantFromContext(ctx),
				LoggedInUserId:  common.GetUserIdFromContext(ctx),
				RenewalApproved: true,
				FieldsMask:      []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEW_APPROVED},
				SourceFields: &commonpb.SourceFields{
					Source:    neo4jentity.DataSourceOpenline.String(),
					AppSource: constants.AppSourceCustomerOsApi,
				},
			}
			if renewalDate != nil {
				grpcUpdateRequest.RenewedAt = utils.ConvertTimeToTimestampPtr(renewalDate)
				grpcUpdateRequest.FieldsMask = append(grpcUpdateRequest.FieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWED_AT)
			}
			return s.grpcClients.OpportunityClient.UpdateRenewalOpportunity(ctx, &grpcUpdateRequest)
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error from events processing: %s", err.Error())
			return err
		}
	} else {
		// if contract is draft - skip rollout renewal opportunity
		if contractEntity.ContractStatus == neo4jenum.ContractStatusDraft {
			return nil
		}
		// if renewal opportunity is expired - rollout renewal opportunity
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		if renewalDate != nil {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
				return s.grpcClients.OpportunityClient.UpdateRenewalOpportunity(ctx, &opportunitypb.UpdateRenewalOpportunityGrpcRequest{
					Id:             opportunityEntity.Id,
					Tenant:         common.GetTenantFromContext(ctx),
					LoggedInUserId: common.GetUserIdFromContext(ctx),
					RenewedAt:      utils.ConvertTimeToTimestampPtr(renewalDate),
					FieldsMask:     []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWED_AT},
					SourceFields: &commonpb.SourceFields{
						Source:    neo4jentity.DataSourceOpenline.String(),
						AppSource: constants.AppSourceCustomerOsApi,
					},
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error from events processing: %s", err.Error())
				return err
			}
			time.Sleep(500 * time.Millisecond)
		}
		_, err = utils.CallEventsPlatformGRPCWithRetry[*contractpb.ContractIdGrpcResponse](func() (*contractpb.ContractIdGrpcResponse, error) {
			return s.grpcClients.ContractClient.RolloutRenewalOpportunityOnExpiration(ctx, &contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest{
				Id:             contractId,
				Tenant:         common.GetTenantFromContext(ctx),
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      constants.AppSourceCustomerOsApi,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error from events processing: %s", err.Error())
			return err
		}
	}

	return nil
}

func (s *contractService) validateContractExists(ctx context.Context, contractId string, span opentracing.Span) error {
	if contractId == "" {
		err := fmt.Errorf("contract id is missing")
		tracing.TraceErr(span, err)
		s.log.Error(err.Error())
		return err
	}

	contractExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), contractId, model2.NodeLabelContract)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error(err.Error())
		return err
	}
	if !contractExists {
		err := fmt.Errorf("contract with id {%s} not found", contractId)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (s *contractService) GetContractByServiceLineItem(ctx context.Context, serviceLineItemId string) (*neo4jentity.ContractEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetContractByServiceLineItem")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))

	contract, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractByServiceLineItemId(ctx, common.GetTenantFromContext(ctx), serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract by service line item: %s", err.Error())
		return nil, err
	}
	if contract == nil {
		err = fmt.Errorf("Contract not found for service line item: %s", serviceLineItemId)
		tracing.TraceErr(span, err)
		return &neo4jentity.ContractEntity{}, err
	}
	return neo4jmapper.MapDbNodeToContractEntity(contract), nil
}

func (s *contractService) GetPaginatedContracts(ctx context.Context, page int, limit int) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetContractByServiceLineItem")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("page", page), log.Int("limit", limit))

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}

	dbNodesWithTotalCount, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetPaginatedContracts(ctx, common.GetContext(ctx).Tenant, paginatedResult.GetSkip(), paginatedResult.GetLimit())
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	contracts := neo4jentity.ContractEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		contracts = append(contracts, *neo4jmapper.MapDbNodeToContractEntity(v))
	}
	paginatedResult.SetRows(&contracts)
	return &paginatedResult, nil
}
