package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
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

//func TestMutationResolver_TenantAddBillingProfile(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//
//	neo4jtest.CreateTenant(ctx, driver, tenantName)
//	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
//	profileId := uuid.New().String()
//	calledAddTenanBillingProfile := false
//
//	tenantServiceCallbacks := events_platform.MockTenantServiceCallbacks{
//		AddBillingProfile: func(context context.Context, profile *tenantpb.AddBillingProfileRequest) (*commonpb.IdResponse, error) {
//			require.Equal(t, tenantName, profile.Tenant)
//			require.Equal(t, testUserId, profile.LoggedInUserId)
//			require.Equal(t, neo4jentity.DataSourceOpenline.String(), profile.SourceFields.Source)
//			require.Equal(t, constants.AppSourceCustomerOsApi, profile.SourceFields.AppSource)
//			require.Equal(t, "phone", profile.Phone)
//			require.Equal(t, "legalName", profile.LegalName)
//			require.Equal(t, "addressLine1", profile.AddressLine1)
//			require.Equal(t, "addressLine2", profile.AddressLine2)
//			require.Equal(t, "addressLine3", profile.AddressLine3)
//			require.Equal(t, "locality", profile.Locality)
//			require.Equal(t, "country", profile.Country)
//			require.Equal(t, "zip", profile.Zip)
//			require.Equal(t, "domesticPaymentsBankInfo", profile.DomesticPaymentsBankInfo)
//			require.Equal(t, "internationalPaymentsBankInfo", profile.InternationalPaymentsBankInfo)
//			require.Equal(t, "vatNumber", profile.VatNumber)
//			require.Equal(t, "sendInvoicesFrom", profile.SendInvoicesFrom)
//			require.Equal(t, "sendInvoicesBcc", profile.SendInvoicesBcc)
//			require.Equal(t, true, profile.CanPayWithCard)
//			require.Equal(t, true, profile.CanPayWithDirectDebitSEPA)
//			require.Equal(t, true, profile.CanPayWithDirectDebitACH)
//			require.Equal(t, true, profile.CanPayWithDirectDebitBacs)
//			require.Equal(t, true, profile.CanPayWithPigeon)
//
//			calledAddTenanBillingProfile = true
//			neo4jtest.CreateTenantBillingProfile(ctx, driver, tenantName, neo4jentity.TenantBillingProfileEntity{Id: profileId})
//			return &commonpb.IdResponse{
//				Id: profileId,
//			}, nil
//		},
//	}
//	events_platform.SetTenantCallbacks(&tenantServiceCallbacks)
//
//	rawResponse := callGraphQL(t, "tenant/add_tenant_billing_profile", map[string]interface{}{})
//	require.Nil(t, rawResponse.Errors)
//
//	var billingProfileStruct struct {
//		Tenant_AddBillingProfile model.TenantBillingProfile
//	}
//
//	err := decode.Decode(rawResponse.Data.(map[string]any), &billingProfileStruct)
//	require.Nil(t, err)
//
//	profile := billingProfileStruct.Tenant_AddBillingProfile
//	require.Equal(t, profileId, profile.ID)
//
//	require.True(t, calledAddTenanBillingProfile)
//}
//
//func TestMutationResolver_TenantUpdateBillingProfile(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//
//	neo4jtest.CreateTenant(ctx, driver, tenantName)
//	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
//	profileId := neo4jtest.CreateTenantBillingProfile(ctx, driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
//	calledUpdateTenantBillingProfile := false
//
//	tenantServiceCallbacks := events_platform.MockTenantServiceCallbacks{
//		UpdateBillingProfile: func(context context.Context, profile *tenantpb.UpdateBillingProfileRequest) (*commonpb.IdResponse, error) {
//			require.Equal(t, tenantName, profile.Tenant)
//			require.Equal(t, profileId, profile.Id)
//			require.Equal(t, testUserId, profile.LoggedInUserId)
//			require.Equal(t, constants.AppSourceCustomerOsApi, profile.AppSource)
//			require.Equal(t, "phone", profile.Phone)
//			require.Equal(t, "legalName", profile.LegalName)
//			require.Equal(t, "addressLine1", profile.AddressLine1)
//			require.Equal(t, "addressLine2", profile.AddressLine2)
//			require.Equal(t, "addressLine3", profile.AddressLine3)
//			require.Equal(t, "locality", profile.Locality)
//			require.Equal(t, "country", profile.Country)
//			require.Equal(t, "zip", profile.Zip)
//			require.Equal(t, "domesticPaymentsBankInfo", profile.DomesticPaymentsBankInfo)
//			require.Equal(t, "internationalPaymentsBankInfo", profile.InternationalPaymentsBankInfo)
//			require.Equal(t, "vatNumber", profile.VatNumber)
//			require.Equal(t, "sendInvoicesFrom", profile.SendInvoicesFrom)
//			require.Equal(t, "sendInvoicesBcc", profile.SendInvoicesBcc)
//			require.Equal(t, true, profile.CanPayWithCard)
//			require.Equal(t, true, profile.CanPayWithDirectDebitSEPA)
//			require.Equal(t, true, profile.CanPayWithDirectDebitACH)
//			require.Equal(t, true, profile.CanPayWithDirectDebitBacs)
//			require.Equal(t, true, profile.CanPayWithPigeon)
//			require.ElementsMatch(t, []tenantpb.TenantBillingProfileFieldMask{
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_PHONE,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_LEGAL_NAME,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_1,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_2,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_3,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_LOCALITY,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_COUNTRY,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ZIP,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_DOMESTIC_PAYMENTS_BANK_INFO,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_INTERNATIONAL_PAYMENTS_BANK_INFO,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_VAT_NUMBER,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_SEND_INVOICES_FROM,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_SEND_INVOICES_BCC,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_CARD,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_DIRECT_DEBIT_SEPA,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_DIRECT_DEBIT_ACH,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_DIRECT_DEBIT_BACS,
//				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_PIGEON,
//			},
//				profile.FieldsMask)
//			calledUpdateTenantBillingProfile = true
//			return &commonpb.IdResponse{
//				Id: profileId,
//			}, nil
//		},
//	}
//	events_platform.SetTenantCallbacks(&tenantServiceCallbacks)
//
//	rawResponse := callGraphQL(t, "tenant/update_tenant_billing_profile", map[string]interface{}{"id": profileId})
//	require.Nil(t, rawResponse.Errors)
//
//	var billingProfileStruct struct {
//		Tenant_UpdateBillingProfile model.TenantBillingProfile
//	}
//
//	err := decode.Decode(rawResponse.Data.(map[string]any), &billingProfileStruct)
//	require.Nil(t, err)
//
//	profile := billingProfileStruct.Tenant_UpdateBillingProfile
//	require.Equal(t, profileId, profile.ID)
//
//	require.True(t, calledUpdateTenantBillingProfile)
//}
