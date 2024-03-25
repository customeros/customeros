package resolver

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

// TODO
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
			require.True(t, bankAccount.AllowInternational)
			require.Equal(t, "IBAN-123456789", bankAccount.Iban)
			require.Equal(t, "BIC-123456789", bankAccount.Bic)
			require.Equal(t, "ACC-123456789", bankAccount.AccountNumber)
			require.Equal(t, "routing-123456789", bankAccount.RoutingNumber)
			require.Equal(t, "sort-123456789", bankAccount.SortCode)
			require.Equal(t, "otherDetails-123", bankAccount.OtherDetails)
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
			require.True(t, bankAccount.AllowInternational)
			require.Equal(t, "", bankAccount.Iban)
			require.Equal(t, "", bankAccount.Bic)
			require.Equal(t, "ACC-123456789", bankAccount.AccountNumber)
			require.Equal(t, "routing-123456789", bankAccount.RoutingNumber)
			require.Equal(t, "otherDetails-123", bankAccount.OtherDetails)
			require.Equal(t, "", bankAccount.SortCode)
			require.Equal(t, 10, len(bankAccount.FieldsMask))
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
