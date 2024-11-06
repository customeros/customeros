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

	contractDataFields := data_fields.ContractSaveFields{}
	contractDataFields.Name = input.Name
	contractDataFields.ContractUrl = input.ContractURL
	contractDataFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
	contractDataFields.AppSource = utils.StringPtr(utils.StringFirstNonEmpty(utils.IfNotNilString(input.AppSource), constants.AppSourceCustomerOsApi))
	contractDataFields.AddressLine1 = input.AddressLine1
	contractDataFields.AddressLine2 = input.AddressLine2
	contractDataFields.Locality = input.Locality
	contractDataFields.Country = input.Country
	contractDataFields.Zip = input.Zip
	contractDataFields.OrganizationLegalName = input.OrganizationLegalName
	contractDataFields.InvoiceEmail = input.InvoiceEmail
	contractDataFields.InvoiceNote = input.InvoiceNote
	contractDataFields.InvoicingEnabled = input.BillingEnabled
	if input.Currency != nil {
		contractDataFields.Currency = utils.ToPtr(mapper.MapCurrencyFromModel(*input.Currency))
	}
	if input.BillingDetails != nil {
		contractDataFields.CanPayWithCard = input.BillingDetails.CanPayWithCard
		contractDataFields.CanPayWithDirectDebit = input.BillingDetails.CanPayWithDirectDebit
		contractDataFields.CanPayWithBankTransfer = input.BillingDetails.CanPayWithBankTransfer
		contractDataFields.Check = input.BillingDetails.Check
		contractDataFields.AddressLine1 = input.BillingDetails.AddressLine1
		contractDataFields.AddressLine2 = input.BillingDetails.AddressLine2
		contractDataFields.Locality = input.BillingDetails.Locality
		contractDataFields.Country = input.BillingDetails.Country
		contractDataFields.Region = input.BillingDetails.Region
		contractDataFields.Zip = input.BillingDetails.PostalCode
		contractDataFields.OrganizationLegalName = input.BillingDetails.OrganizationLegalName
		contractDataFields.InvoiceEmail = input.BillingDetails.BillingEmail
		if input.BillingDetails.BillingEmailCc != nil {
			contractDataFields.InvoiceEmailCC = utils.ToPtr(input.BillingDetails.BillingEmailCc)
		}
		if input.BillingDetails.BillingEmailBcc != nil {
			contractDataFields.InvoiceEmailBCC = utils.ToPtr(input.BillingDetails.BillingEmailBcc)
		}
		contractDataFields.InvoiceNote = input.BillingDetails.InvoiceNote
		contractDataFields.BillingCycleInMonths = input.BillingDetails.BillingCycleInMonths
		contractDataFields.PayOnline = input.BillingDetails.PayOnline
		contractDataFields.PayAutomatically = input.BillingDetails.PayAutomatically
		contractDataFields.DueDays = input.BillingDetails.DueDays
	}

	contractDataFields.Approved = input.Approved

	if input.ContractName != nil {
		contractDataFields.Name = input.ContractName
	}

	zeroTime := time.Time{}

	if input.SignedAt != nil && *input.SignedAt != zeroTime {
		contractDataFields.SignedAt = input.SignedAt
	}
	if input.ContractSigned != nil && *input.ContractSigned != zeroTime {
		contractDataFields.SignedAt = input.ContractSigned
	}
	if input.ServiceStartedAt != nil && *input.ServiceStartedAt != zeroTime {
		contractDataFields.SignedAt = input.ServiceStartedAt
	}
	if input.ServiceStarted != nil && *input.ServiceStarted != zeroTime {
		contractDataFields.SignedAt = input.ServiceStarted
	}
	if input.EndedAt != nil && *input.EndedAt != zeroTime {
		contractDataFields.SignedAt = input.EndedAt
	}
	if input.ContractEnded != nil && *input.ContractEnded != zeroTime {
		contractDataFields.SignedAt = input.ContractEnded
	}
	if input.InvoicingStartDate != nil && *input.InvoicingStartDate != zeroTime {
		contractDataFields.SignedAt = input.InvoicingStartDate
	}
	if input.BillingDetails != nil && input.BillingDetails.InvoicingStarted != nil {
		if *input.BillingDetails.InvoicingStarted != zeroTime {
			contractDataFields.InvoicingStartDate = input.BillingDetails.InvoicingStarted
		}
	}

	if input.CommittedPeriodInMonths != nil {
		contractDataFields.LengthInMonths = input.CommittedPeriodInMonths
	} else {
		// prepare length in months from renewal cycle and periods
		renewalCycle := ""
		if input.ContractRenewalCycle != nil {
			renewalCycleEnum := *input.ContractRenewalCycle
			renewalCycle = renewalCycleEnum.String()
		} else if input.RenewalCycle != nil {
			renewalCycleEnum := *input.RenewalCycle
			renewalCycle = renewalCycleEnum.String()
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
			if input.CommittedPeriods != nil && *input.CommittedPeriods > 1 {
				contractDataFields.LengthInMonths = utils.Int64Ptr(*contractDataFields.LengthInMonths * *input.CommittedPeriods)
			} else if input.RenewalPeriods != nil && *input.RenewalPeriods > 1 {
				contractDataFields.LengthInMonths = utils.Int64Ptr(*contractDataFields.LengthInMonths * *input.RenewalPeriods)
			}
		}
	}

	if input.BillingDetails == nil || input.BillingDetails.BillingCycleInMonths == nil {
		if input.BillingDetails != nil && input.BillingDetails.BillingCycle != nil {
			switch *input.BillingDetails.BillingCycle {
			case model.ContractBillingCycleMonthlyBilling:
				contractDataFields.BillingCycleInMonths = utils.Int64Ptr(1)
			case model.ContractBillingCycleQuarterlyBilling:
				contractDataFields.BillingCycleInMonths = utils.Int64Ptr(3)
			case model.ContractBillingCycleAnnualBilling:
				contractDataFields.BillingCycleInMonths = utils.Int64Ptr(12)
			default:
				contractDataFields.BillingCycleInMonths = utils.Int64Ptr(0)
			}
		} else if input.BillingCycle != nil {
			switch *input.BillingCycle {
			case model.ContractBillingCycleMonthlyBilling:
				contractDataFields.BillingCycleInMonths = utils.Int64Ptr(1)
			case model.ContractBillingCycleQuarterlyBilling:
				contractDataFields.BillingCycleInMonths = utils.Int64Ptr(3)
			case model.ContractBillingCycleAnnualBilling:
				contractDataFields.BillingCycleInMonths = utils.Int64Ptr(12)
			default:
				contractDataFields.BillingCycleInMonths = utils.Int64Ptr(0)
			}
		}
	}
	contractDataFields.AutoRenew = input.AutoRenew

	_, err := s.services.CommonServices.ContractService.Save(ctx, &input.ContractID, contractDataFields)
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
