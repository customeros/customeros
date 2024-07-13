package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
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
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BankAccountService interface {
	CreateTenantBankAccount(ctx context.Context, input *model.BankAccountCreateInput) (string, error)
	UpdateTenantBankAccount(ctx context.Context, input *model.BankAccountUpdateInput) error
	GetTenantBankAccounts(ctx context.Context) (*neo4jentity.BankAccountEntities, error)
	GetTenantBankAccount(ctx context.Context, id string) (*neo4jentity.BankAccountEntity, error)
	DeleteTenantBankAccount(ctx context.Context, id string) (bool, error)
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
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbNodes, err := s.repositories.Neo4jRepositories.BankAccountReadRepository.GetBankAccounts(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantBankAccounts: %s", err.Error())
	}

	tenantBankAccounts := neo4jentity.BankAccountEntities{}
	for _, dbNode := range dbNodes {
		tenantBankAccounts = append(tenantBankAccounts, *neo4jmapper.MapDbNodeToBankAccountEntity(dbNode))
	}

	return &tenantBankAccounts, nil
}

func (s *bankAccountService) GetTenantBankAccount(ctx context.Context, id string) (*neo4jentity.BankAccountEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountService.GetTenantBankAccount")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("bankAccountId", id))

	dbNode, err := s.repositories.Neo4jRepositories.BankAccountReadRepository.GetBankAccountById(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantBankAccount: %s", err.Error())
	}

	return neo4jmapper.MapDbNodeToBankAccountEntity(dbNode), nil
}

func (s *bankAccountService) CreateTenantBankAccount(ctx context.Context, input *model.BankAccountCreateInput) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountService.CreateTenantBankAccount")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	grpcRequest := tenantpb.AddBankAccountGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
		BankName:            utils.IfNotNilString(input.BankName),
		BankTransferEnabled: utils.IfNotNilBool(input.BankTransferEnabled),
		AllowInternational:  utils.IfNotNilBool(input.AllowInternational),
		Currency:            utils.IfNotNilString(input.Currency.String()),
		Iban:                utils.IfNotNilString(input.Iban),
		Bic:                 utils.IfNotNilString(input.Bic),
		SortCode:            utils.IfNotNilString(input.SortCode),
		AccountNumber:       utils.IfNotNilString(input.AccountNumber),
		RoutingNumber:       utils.IfNotNilString(input.RoutingNumber),
		OtherDetails:        utils.IfNotNilString(input.OtherDetails),
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*commonpb.IdResponse](func() (*commonpb.IdResponse, error) {
		return s.grpcClients.TenantClient.AddBankAccount(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, model2.NodeLabelBankAccount, span)

	return response.Id, nil
}

func (s *bankAccountService) UpdateTenantBankAccount(ctx context.Context, input *model.BankAccountUpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountService.UpdateTenantBankAccount")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	if input.ID == "" {
		err := fmt.Errorf("bank account id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	bankAccountExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), input.ID, model2.NodeLabelBankAccount)
	if !bankAccountExists {
		err := fmt.Errorf("bank account with id {%s} not found", input.ID)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}
	bankAccountEntity, err := s.GetTenantBankAccount(ctx, input.ID)

	var fieldsMask []tenantpb.BankAccountFieldMask
	updateRequest := tenantpb.UpdateBankAccountGrpcRequest{
		Tenant:              common.GetTenantFromContext(ctx),
		Id:                  input.ID,
		LoggedInUserId:      common.GetUserIdFromContext(ctx),
		AppSource:           constants.AppSourceCustomerOsApi,
		BankName:            utils.IfNotNilString(input.BankName),
		BankTransferEnabled: utils.IfNotNilBool(input.BankTransferEnabled),
		AllowInternational:  utils.IfNotNilBool(input.AllowInternational),
		Iban:                utils.IfNotNilString(input.Iban),
		Bic:                 utils.IfNotNilString(input.Bic),
		SortCode:            utils.IfNotNilString(input.SortCode),
		AccountNumber:       utils.IfNotNilString(input.AccountNumber),
		RoutingNumber:       utils.IfNotNilString(input.RoutingNumber),
		OtherDetails:        utils.IfNotNilString(input.OtherDetails),
	}

	if input.BankName != nil {
		fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_BANK_NAME)
	}
	if input.BankTransferEnabled != nil {
		fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_BANK_TRANSFER_ENABLED)
	}
	if input.AllowInternational != nil {
		fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ALLOW_INTERNATIONAL)
	}
	if input.Currency != nil {
		updateRequest.Currency = input.Currency.String()
		fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_CURRENCY)
	}
	if input.Currency != nil && bankAccountEntity.Currency.String() != input.Currency.String() {
		if *input.Currency == model.CurrencyUsd {
			updateRequest.SortCode = ""
			updateRequest.Iban = ""
			updateRequest.Bic = ""
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_SORT_CODE)
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_IBAN)
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_BIC)
		} else if *input.Currency == model.CurrencyGbp {
			updateRequest.RoutingNumber = ""
			updateRequest.Iban = ""
			updateRequest.Bic = ""
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ROUTING_NUMBER)
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_IBAN)
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_BIC)
		} else if *input.Currency == model.CurrencyEur {
			updateRequest.RoutingNumber = ""
			updateRequest.AccountNumber = ""
			updateRequest.SortCode = ""
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ROUTING_NUMBER)
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ACCOUNT_NUMBER)
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_SORT_CODE)
		} else {
			updateRequest.Iban = ""
			updateRequest.SortCode = ""
			updateRequest.RoutingNumber = ""
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_IBAN)
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_SORT_CODE)
			fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ROUTING_NUMBER)
		}
	}
	if input.Iban != nil && !utils.ContainsElement(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_IBAN) {
		fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_IBAN)
	}
	if input.Bic != nil && !utils.ContainsElement(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_BIC) {
		fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_BIC)
	}
	if input.SortCode != nil && !utils.ContainsElement(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_SORT_CODE) {
		fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_SORT_CODE)
	}
	if input.AccountNumber != nil && !utils.ContainsElement(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ACCOUNT_NUMBER) {
		fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ACCOUNT_NUMBER)
	}
	if input.RoutingNumber != nil && !utils.ContainsElement(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ROUTING_NUMBER) {
		fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_ROUTING_NUMBER)
	}
	if input.OtherDetails != nil && !utils.ContainsElement(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_OTHER_DETAILS) {
		fieldsMask = append(fieldsMask, tenantpb.BankAccountFieldMask_BANK_ACCOUNT_FIELD_OTHER_DETAILS)
	}
	if len(fieldsMask) == 0 {
		span.LogFields(log.String("result", "No fields to update"))
		return nil
	}
	updateRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*commonpb.IdResponse](func() (*commonpb.IdResponse, error) {
		return s.grpcClients.TenantClient.UpdateBankAccount(ctx, &updateRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *bankAccountService) DeleteTenantBankAccount(ctx context.Context, bankAccountId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountService.DeleteTenantBankAccount")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("bankAccountId", bankAccountId))

	bankAccountExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), bankAccountId, model2.NodeLabelBankAccount)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error on checking if bank account exists: %s", err.Error())
		return false, err
	}
	if !bankAccountExists {
		err := fmt.Errorf("bank account with id {%s} not found", bankAccountId)
		tracing.TraceErr(span, err)
		s.log.Errorf(err.Error())
		return false, err
	}

	deleteRequest := tenantpb.DeleteBankAccountGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             bankAccountId,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		AppSource:      constants.AppSourceCustomerOsApi,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*emptypb.Empty](func() (*emptypb.Empty, error) {
		return s.grpcClients.TenantClient.DeleteBankAccount(ctx, &deleteRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return false, err
	}

	// wait for service line item to be deleted from graph db
	neo4jrepository.WaitForNodeDeletedFromNeo4j(ctx, s.repositories.Neo4jRepositories, bankAccountId, model2.NodeLabelBankAccount, span)

	bankAccountExists, err = s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), bankAccountId, model2.NodeLabelBankAccount)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error on checking if bank account exists: %s", err.Error())
		return false, err
	}

	return !bankAccountExists, nil
}
