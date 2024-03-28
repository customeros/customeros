package resolver

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	offeringpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/offering"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_OfferingCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	offeringId := uuid.New().String()

	calledCreateOffering := false

	offeringCallbacks := events_platform.MockOfferingServiceCallbacks{
		CreateOffering: func(context context.Context, offering *offeringpb.CreateOfferingGrpcRequest) (*commonpb.IdResponse, error) {

			require.Equal(t, "1", offering.Name)
			require.Equal(t, true, offering.Active)
			require.Equal(t, "PRODUCT", offering.Type)
			require.Equal(t, "SUBSCRIPTION", offering.PricingModel)
			require.Equal(t, int64(2), offering.PricingPeriodInMonths)
			require.Equal(t, "AUD", offering.Currency)
			require.Equal(t, float64(3), offering.Price)
			require.Equal(t, true, offering.PriceCalculated)
			require.Equal(t, true, offering.Conditional)
			require.Equal(t, true, offering.Taxable)
			require.Equal(t, "REVENUE_SHARE", offering.PriceCalculationType)
			require.Equal(t, float64(4), offering.PriceCalculationRevenueSharePercentage)
			require.Equal(t, "MONTHLY", offering.ConditionalsMinimumChargePeriod)
			require.Equal(t, float64(5), offering.ConditionalsMinimumChargeAmount)

			calledCreateOffering = true

			return &commonpb.IdResponse{
				Id: offeringId,
			}, nil
		},
	}
	events_platform.SetOfferingCallbacks(&offeringCallbacks)

	rawResponse := callGraphQL(t, "offering/create_offering", map[string]interface{}{})
	require.Nil(t, rawResponse.Errors)

	var graphqlResponse struct {
		Offering_Create model.Offering
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &graphqlResponse)
	require.Nil(t, err)

	offering := graphqlResponse.Offering_Create
	require.Equal(t, offeringId, offering.Metadata.ID)

	require.True(t, calledCreateOffering)
}

//func TestMutationResolver_BankAccountUpdate_FromEurToUsd(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//
//	neo4jtest.CreateTenant(ctx, driver, tenantName)
//	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
//	bankAccountId := neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{
//		BankName: "Bank",
//		Currency: neo4jenum.CurrencyEUR,
//	})
//	calledUpdateBankAccount := false
//
//	tenantServiceCallbacks := events_platform.MockTenantServiceCallbacks{
//		UpdateBankAccount: func(context context.Context, bankAccount *tenantpb.UpdateBankAccountGrpcRequest) (*commonpb.IdResponse, error) {
//			require.Equal(t, tenantName, bankAccount.Tenant)
//			require.Equal(t, testUserId, bankAccount.LoggedInUserId)
//			require.Equal(t, constants.AppSourceCustomerOsApi, bankAccount.AppSource)
//			require.Equal(t, bankAccountId, bankAccount.Id)
//			require.Equal(t, "Bank of America", bankAccount.BankName)
//			require.Equal(t, "USD", bankAccount.Currency)
//			require.True(t, bankAccount.BankTransferEnabled)
//			require.True(t, bankAccount.AllowInternational)
//			require.Equal(t, "", bankAccount.Iban)
//			require.Equal(t, "", bankAccount.Bic)
//			require.Equal(t, "ACC-123456789", bankAccount.AccountNumber)
//			require.Equal(t, "routing-123456789", bankAccount.RoutingNumber)
//			require.Equal(t, "otherDetails-123", bankAccount.OtherDetails)
//			require.Equal(t, "", bankAccount.SortCode)
//			require.Equal(t, 10, len(bankAccount.FieldsMask))

//
//require.ElementsMatch(t, []offeringpb.OfferingFieldMask{
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_NAME,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_ACTIVE,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_TYPE,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICING_MODEL,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICING_PERIOD_IN_MONTHS,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_CURRENCY,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE_CALCULATED,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_CONDITIONAL,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_TAXABLE,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE_CALCULATION_TYPE,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE_CALCULATION_REVENUE_SHARE_PERCENTAGE,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_CONDITIONALS_MINIMUM_CHARGE_PERIOD,
//	offeringpb.OfferingFieldMask_OFFERING_FIELD_CONDITIONALS_MINIMUM_CHARGE_AMOUNT},
//	offering.FieldsMask)

//			calledUpdateBankAccount = true
//			neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{Id: bankAccountId})
//			return &commonpb.IdResponse{
//				Id: bankAccountId,
//			}, nil
//		},
//	}
//	events_platform.SetTenantCallbacks(&tenantServiceCallbacks)
//
//	rawResponse := callGraphQL(t, "bank_account/update_bank_account", map[string]interface{}{
//		"accountId": bankAccountId,
//	})
//	require.Nil(t, rawResponse.Errors)
//
//	var graphqlResponse struct {
//		BankAccount_Update model.BankAccount
//	}
//
//	err := decode.Decode(rawResponse.Data.(map[string]any), &graphqlResponse)
//	require.Nil(t, err)
//
//	bankAccount := graphqlResponse.BankAccount_Update
//	require.Equal(t, bankAccountId, bankAccount.Metadata.ID)
//
//	require.True(t, calledUpdateBankAccount)
//}
//
//func TestMutationResolver_BankAccountDelete(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//
//	neo4jtest.CreateTenant(ctx, driver, tenantName)
//	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
//	bankAccountId := neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{})
//
//	calledDeleteBankAccount := false
//
//	tenantServiceCallbacks := events_platform.MockTenantServiceCallbacks{
//		DeleteBankAccount: func(context context.Context, bankAccount *tenantpb.DeleteBankAccountGrpcRequest) (*emptypb.Empty, error) {
//			require.Equal(t, tenantName, bankAccount.Tenant)
//			require.Equal(t, testUserId, bankAccount.LoggedInUserId)
//			require.Equal(t, constants.AppSourceCustomerOsApi, bankAccount.AppSource)
//			require.Equal(t, bankAccountId, bankAccount.Id)
//			calledDeleteBankAccount = true
//			neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{Id: bankAccountId})
//			return &emptypb.Empty{}, nil
//		},
//	}
//	events_platform.SetTenantCallbacks(&tenantServiceCallbacks)
//
//	rawResponse := callGraphQL(t, "bank_account/delete_bank_account", map[string]interface{}{
//		"accountId": bankAccountId,
//	})
//	require.Nil(t, rawResponse.Errors)
//
//	var graphqlResponse struct {
//		BankAccount_Delete model.DeleteResponse
//	}
//
//	err := decode.Decode(rawResponse.Data.(map[string]any), &graphqlResponse)
//	require.Nil(t, err)
//
//	deleteResponse := graphqlResponse.BankAccount_Delete
//	require.True(t, deleteResponse.Accepted)
//	require.False(t, deleteResponse.Completed)
//
//	require.True(t, calledDeleteBankAccount)
//}
