package resolver

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestQueryResolver_BankAccounts(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	today := utils.Now()
	yesterday := today.AddDate(0, 0, -1)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	account1 := neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{
		CreatedAt:           today,
		UpdatedAt:           today,
		BankName:            "bankName1",
		BankTransferEnabled: true,
		Currency:            neo4jenum.CurrencyEUR,
		Iban:                "iban1",
		Bic:                 "bic1",
		SortCode:            "sortCode1",
		AccountNumber:       "accountNumber1",
		RoutingNumber:       "routingNumber1",
	})
	account2 := neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{
		CreatedAt:           yesterday,
		UpdatedAt:           yesterday,
		BankName:            "bankName2",
		BankTransferEnabled: false,
		Currency:            neo4jenum.CurrencyUSD,
		Iban:                "iban2",
		Bic:                 "bic2",
		SortCode:            "sortCode2",
		AccountNumber:       "accountNumber2",
		RoutingNumber:       "routingNumber2",
	})

	rawResponse, err := c.RawPost(getQuery("bank_account/get_bank_accounts"))
	assertRawResponseSuccess(t, rawResponse, err)

	var graphqlResponse struct {
		BankAccounts []model.BankAccount
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &graphqlResponse)
	require.Nil(t, err)
	require.NotNil(t, graphqlResponse)

	require.Equal(t, 2, len(graphqlResponse.BankAccounts))
	ba1 := graphqlResponse.BankAccounts[0]
	require.Equal(t, account2, ba1.Metadata.ID)
	require.Equal(t, yesterday, ba1.Metadata.Created)
	require.Equal(t, "bankName2", *ba1.BankName)
	require.Equal(t, false, ba1.BankTransferEnabled)
	require.Equal(t, model.CurrencyUsd, *ba1.Currency)
	require.Equal(t, "iban2", *ba1.Iban)
	require.Equal(t, "bic2", *ba1.Bic)
	require.Equal(t, "sortCode2", *ba1.SortCode)
	require.Equal(t, "accountNumber2", *ba1.AccountNumber)
	require.Equal(t, "routingNumber2", *ba1.RoutingNumber)

	ba2 := graphqlResponse.BankAccounts[1]
	require.Equal(t, account1, ba2.Metadata.ID)
	require.Equal(t, today, ba2.Metadata.Created)
	require.Equal(t, "bankName1", *ba2.BankName)
	require.Equal(t, true, ba2.BankTransferEnabled)
	require.Equal(t, model.CurrencyEur, *ba2.Currency)
	require.Equal(t, "iban1", *ba2.Iban)
	require.Equal(t, "bic1", *ba2.Bic)
	require.Equal(t, "sortCode1", *ba2.SortCode)
	require.Equal(t, "accountNumber1", *ba2.AccountNumber)
	require.Equal(t, "routingNumber1", *ba2.RoutingNumber)
}

func TestMutationResolver_BankAccountCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	bankAccountId := uuid.New().String()
	calledAddBankAccount := false

	tenantServiceCallbacks := events_platform.MockTenantServiceCallbacks{
		AddBankAccount: func(context context.Context, bankAccount *tenantpb.AddBankAccountGrpcRequest) (*commonpb.IdResponse, error) {
			require.Equal(t, tenantName, bankAccount.Tenant)
			require.Equal(t, testUserId, bankAccount.LoggedInUserId)
			require.Equal(t, neo4jentity.DataSourceOpenline.String(), bankAccount.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, bankAccount.SourceFields.AppSource)
			require.Equal(t, "Bank of America", bankAccount.BankName)
			require.Equal(t, "USD", bankAccount.Currency)
			require.True(t, bankAccount.BankTransferEnabled)
			require.Equal(t, "IBAN-123456789", bankAccount.Iban)
			require.Equal(t, "BIC-123456789", bankAccount.Bic)
			require.Equal(t, "ACC-123456789", bankAccount.AccountNumber)
			require.Equal(t, "routing-123456789", bankAccount.RoutingNumber)
			require.Equal(t, "sort-123456789", bankAccount.SortCode)
			calledAddBankAccount = true
			neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{Id: bankAccountId})
			return &commonpb.IdResponse{
				Id: bankAccountId,
			}, nil
		},
	}
	events_platform.SetTenantCallbacks(&tenantServiceCallbacks)

	rawResponse := callGraphQL(t, "bank_account/create_bank_account", map[string]interface{}{})
	require.Nil(t, rawResponse.Errors)

	var graphqlResponse struct {
		BankAccount_Create model.BankAccount
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &graphqlResponse)
	require.Nil(t, err)

	bankAccount := graphqlResponse.BankAccount_Create
	require.Equal(t, bankAccountId, bankAccount.Metadata.ID)

	require.True(t, calledAddBankAccount)
}

func TestMutationResolver_BankAccountUpdate_FromEurToUsd(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	bankAccountId := neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{
		BankName: "Bank",
		Currency: neo4jenum.CurrencyEUR,
	})
	calledUpdateBankAccount := false

	tenantServiceCallbacks := events_platform.MockTenantServiceCallbacks{
		UpdateBankAccount: func(context context.Context, bankAccount *tenantpb.UpdateBankAccountGrpcRequest) (*commonpb.IdResponse, error) {
			require.Equal(t, tenantName, bankAccount.Tenant)
			require.Equal(t, testUserId, bankAccount.LoggedInUserId)
			require.Equal(t, constants.AppSourceCustomerOsApi, bankAccount.AppSource)
			require.Equal(t, bankAccountId, bankAccount.Id)
			require.Equal(t, "Bank of America", bankAccount.BankName)
			require.Equal(t, "USD", bankAccount.Currency)
			require.True(t, bankAccount.BankTransferEnabled)
			require.Equal(t, "", bankAccount.Iban)
			require.Equal(t, "", bankAccount.Bic)
			require.Equal(t, "ACC-123456789", bankAccount.AccountNumber)
			require.Equal(t, "routing-123456789", bankAccount.RoutingNumber)
			require.Equal(t, "", bankAccount.SortCode)
			require.Equal(t, 8, len(bankAccount.FieldsMask))
			calledUpdateBankAccount = true
			neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{Id: bankAccountId})
			return &commonpb.IdResponse{
				Id: bankAccountId,
			}, nil
		},
	}
	events_platform.SetTenantCallbacks(&tenantServiceCallbacks)

	rawResponse := callGraphQL(t, "bank_account/update_bank_account", map[string]interface{}{
		"accountId": bankAccountId,
	})
	require.Nil(t, rawResponse.Errors)

	var graphqlResponse struct {
		BankAccount_Update model.BankAccount
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &graphqlResponse)
	require.Nil(t, err)

	bankAccount := graphqlResponse.BankAccount_Update
	require.Equal(t, bankAccountId, bankAccount.Metadata.ID)

	require.True(t, calledUpdateBankAccount)
}

func TestMutationResolver_BankAccountDelete(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	bankAccountId := neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{})

	calledDeleteBankAccount := false

	tenantServiceCallbacks := events_platform.MockTenantServiceCallbacks{
		DeleteBankAccount: func(context context.Context, bankAccount *tenantpb.DeleteBankAccountGrpcRequest) (*emptypb.Empty, error) {
			require.Equal(t, tenantName, bankAccount.Tenant)
			require.Equal(t, testUserId, bankAccount.LoggedInUserId)
			require.Equal(t, constants.AppSourceCustomerOsApi, bankAccount.AppSource)
			require.Equal(t, bankAccountId, bankAccount.Id)
			calledDeleteBankAccount = true
			neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{Id: bankAccountId})
			return &emptypb.Empty{}, nil
		},
	}
	events_platform.SetTenantCallbacks(&tenantServiceCallbacks)

	rawResponse := callGraphQL(t, "bank_account/delete_bank_account", map[string]interface{}{
		"accountId": bankAccountId,
	})
	require.Nil(t, rawResponse.Errors)

	var graphqlResponse struct {
		BankAccount_Delete model.DeleteResponse
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &graphqlResponse)
	require.Nil(t, err)

	deleteResponse := graphqlResponse.BankAccount_Delete
	require.True(t, deleteResponse.Accepted)
	require.False(t, deleteResponse.Completed)

	require.True(t, calledDeleteBankAccount)
}
