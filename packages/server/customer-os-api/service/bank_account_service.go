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
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
)

type BankAccountService interface {
	GetTenantBankAccounts(ctx context.Context) (*neo4jentity.BankAccountEntities, error)
	//GetTenantBankAccount(ctx context.Context, id string) (*neo4jentity.BankAccountEntity, error)
	//CreateTenantBillingProfile(ctx context.Context, input model.TenantBillingProfileInput) (string, error)
	//UpdateTenantBillingProfile(ctx context.Context, input model.TenantBillingProfileUpdateInput) error
}

type bankAccountService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewBankAccountService(log logger.Logger, repository *repository.Repositories, grpcClients *grpc_client.Clients) BankAccountService {
	return &bankAccountService{
		log:          log,
		repositories: repository,
		grpcClients:  grpcClients,
	}
}

func (s *bankAccountService) GetTenantBankAccounts(ctx context.Context) (*neo4jentity.BankAccountEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountService.GetTenantBankAccounts")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	dbNodes, err := s.repositories.Neo4jRepositories.BankAccountReadRepository.GetBankAccounts(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantBankAccounts: %w", err)
	}

	tenantBankAccounts := neo4jentity.BankAccountEntities{}
	for _, dbNode := range dbNodes {
		tenantBankAccounts = append(tenantBankAccounts, *neo4jmapper.MapDbNodeToBankAccountEntity(dbNode))
	}

	return &tenantBankAccounts, nil
}

//func (s *bankAccountService) GetTenantBillingProfile(ctx context.Context, id string) (*neo4jentity.TenantBillingProfileEntity, error) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountService.GetTenantBillingProfile")
//	defer span.Finish()
//	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
//	span.LogFields(log.String("id", id))
//
//	dbNode, err := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantBillingProfileById(ctx, common.GetTenantFromContext(ctx), id)
//	if err != nil {
//		tracing.TraceErr(span, err)
//		return nil, fmt.Errorf("GetTenantBillingProfile: %w", err)
//	}
//
//	return neo4jmapper.MapDbNodeToTenantBillingProfileEntity(dbNode), nil
//}

//func (s *bankAccountService) CreateTenantBillingProfile(ctx context.Context, input model.TenantBillingProfileInput) (string, error) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountService.CreateTenantBillingProfile")
//	defer span.Finish()
//	tracing.SetDefaultServiceSpanTags(ctx, span)
//	tracing.LogObjectAsJson(span, "input", input)
//
//	grpcRequest := tenantpb.AddBillingProfileRequest{
//		Tenant:         common.GetTenantFromContext(ctx),
//		LoggedInUserId: common.GetUserIdFromContext(ctx),
//		SourceFields: &commonpb.SourceFields{
//			Source:    neo4jentity.DataSourceOpenline.String(),
//			AppSource: constants.AppSourceCustomerOsApi,
//		},
//		Phone:                         utils.IfNotNilString(input.Phone),
//		LegalName:                     utils.IfNotNilString(input.LegalName),
//		AddressLine1:                  utils.IfNotNilString(input.AddressLine1),
//		AddressLine2:                  utils.IfNotNilString(input.AddressLine2),
//		AddressLine3:                  utils.IfNotNilString(input.AddressLine3),
//		Locality:                      utils.IfNotNilString(input.Locality),
//		Country:                       utils.IfNotNilString(input.Country),
//		Zip:                           utils.IfNotNilString(input.Zip),
//		DomesticPaymentsBankInfo:      utils.IfNotNilString(input.DomesticPaymentsBankInfo),
//		InternationalPaymentsBankInfo: utils.IfNotNilString(input.InternationalPaymentsBankInfo),
//		VatNumber:                     utils.IfNotNilString(input.VatNumber),
//		SendInvoicesFrom:              utils.IfNotNilString(input.SendInvoicesFrom),
//		SendInvoicesBcc:               utils.IfNotNilString(input.SendInvoicesBcc),
//		CanPayWithCard:                utils.IfNotNilBool(input.CanPayWithCard),
//		CanPayWithDirectDebitSEPA:     utils.IfNotNilBool(input.CanPayWithDirectDebitSepa),
//		CanPayWithDirectDebitACH:      utils.IfNotNilBool(input.CanPayWithDirectDebitAch),
//		CanPayWithDirectDebitBacs:     utils.IfNotNilBool(input.CanPayWithDirectDebitBacs),
//		CanPayWithPigeon:              utils.IfNotNilBool(input.CanPayWithPigeon),
//	}
//
//	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
//	response, err := CallEventsPlatformGRPCWithRetry[*commonpb.IdResponse](func() (*commonpb.IdResponse, error) {
//		return s.grpcClients.TenantClient.AddBillingProfile(ctx, &grpcRequest)
//	})
//	if err != nil {
//		tracing.TraceErr(span, err)
//		s.log.Errorf("Error from events processing: %s", err.Error())
//		return "", err
//	}
//
//	WaitForNodeCreatedInNeo4j(ctx, s.repositories, response.Id, neo4jutil.NodeLabelTenantBillingProfile, span)
//
//	return response.Id, nil
//}

//func (s *bankAccountService) UpdateTenantBillingProfile(ctx context.Context, input model.TenantBillingProfileUpdateInput) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountService.UpdateTenantBillingProfile")
//	defer span.Finish()
//	tracing.SetDefaultServiceSpanTags(ctx, span)
//	tracing.LogObjectAsJson(span, "input", input)
//
//	if input.ID == "" {
//		err := fmt.Errorf("(BankAccountService.UpdateTenantBillingProfile) billing profile id is missing")
//		s.log.Error(err.Error())
//		tracing.TraceErr(span, err)
//		return err
//	}
//
//	billingProfileExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), input.ID, neo4jutil.NodeLabelTenantBillingProfile)
//	if !billingProfileExists {
//		err := fmt.Errorf("(BankAccountService.UpdateTenantBillingProfile) tenant billing profile with id {%s} not found", input.ID)
//		s.log.Error(err.Error())
//		tracing.TraceErr(span, err)
//		return err
//	}
//
//	var fieldsMask []tenantpb.TenantBillingProfileFieldMask
//	updateRequest := tenantpb.UpdateBillingProfileRequest{
//		Tenant:                        common.GetTenantFromContext(ctx),
//		Id:                            input.ID,
//		LoggedInUserId:                common.GetUserIdFromContext(ctx),
//		AppSource:                     constants.AppSourceCustomerOsApi,
//		Phone:                         utils.IfNotNilString(input.Phone),
//		LegalName:                     utils.IfNotNilString(input.LegalName),
//		AddressLine1:                  utils.IfNotNilString(input.AddressLine1),
//		AddressLine2:                  utils.IfNotNilString(input.AddressLine2),
//		AddressLine3:                  utils.IfNotNilString(input.AddressLine3),
//		Locality:                      utils.IfNotNilString(input.Locality),
//		Country:                       utils.IfNotNilString(input.Country),
//		Zip:                           utils.IfNotNilString(input.Zip),
//		DomesticPaymentsBankInfo:      utils.IfNotNilString(input.DomesticPaymentsBankInfo),
//		InternationalPaymentsBankInfo: utils.IfNotNilString(input.InternationalPaymentsBankInfo),
//		VatNumber:                     utils.IfNotNilString(input.VatNumber),
//		SendInvoicesFrom:              utils.IfNotNilString(input.SendInvoicesFrom),
//		SendInvoicesBcc:               utils.IfNotNilString(input.SendInvoicesBcc),
//		CanPayWithCard:                utils.IfNotNilBool(input.CanPayWithCard),
//		CanPayWithDirectDebitSEPA:     utils.IfNotNilBool(input.CanPayWithDirectDebitSepa),
//		CanPayWithDirectDebitACH:      utils.IfNotNilBool(input.CanPayWithDirectDebitAch),
//		CanPayWithDirectDebitBacs:     utils.IfNotNilBool(input.CanPayWithDirectDebitBacs),
//		CanPayWithPigeon:              utils.IfNotNilBool(input.CanPayWithPigeon),
//	}
//
//	if input.Patch != nil && *input.Patch {
//		if input.Phone != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_PHONE)
//		}
//		if input.LegalName != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_LEGAL_NAME)
//		}
//		if input.AddressLine1 != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_1)
//		}
//		if input.AddressLine2 != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_2)
//		}
//		if input.AddressLine3 != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_3)
//		}
//		if input.Locality != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_LOCALITY)
//		}
//		if input.Country != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_COUNTRY)
//		}
//		if input.Zip != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ZIP)
//		}
//		if input.DomesticPaymentsBankInfo != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_DOMESTIC_PAYMENTS_BANK_INFO)
//		}
//		if input.InternationalPaymentsBankInfo != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_INTERNATIONAL_PAYMENTS_BANK_INFO)
//		}
//		if input.VatNumber != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_VAT_NUMBER)
//		}
//		if input.SendInvoicesFrom != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_SEND_INVOICES_FROM)
//		}
//		if input.SendInvoicesBcc != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_SEND_INVOICES_BCC)
//		}
//		if input.CanPayWithCard != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_CARD)
//		}
//		if input.CanPayWithDirectDebitSepa != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_DIRECT_DEBIT_SEPA)
//		}
//		if input.CanPayWithDirectDebitAch != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_DIRECT_DEBIT_ACH)
//		}
//		if input.CanPayWithDirectDebitBacs != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_DIRECT_DEBIT_BACS)
//		}
//		if input.CanPayWithPigeon != nil {
//			fieldsMask = append(fieldsMask, tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_PIGEON)
//		}
//		if len(fieldsMask) == 0 {
//			span.LogFields(log.String("result", "No fields to update"))
//			return nil
//		}
//		updateRequest.FieldsMask = fieldsMask
//	}
//
//	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
//	_, err := CallEventsPlatformGRPCWithRetry[*commonpb.IdResponse](func() (*commonpb.IdResponse, error) {
//		return s.grpcClients.TenantClient.UpdateBillingProfile(ctx, &updateRequest)
//	})
//	if err != nil {
//		tracing.TraceErr(span, err)
//		s.log.Errorf("Error from events processing: %s", err.Error())
//		return err
//	}
//
//	return nil
//}
