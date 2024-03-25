package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
	"testing"
)

func TestMutationResolver_TenantMerge(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenant(ctx, driver, "other")

	rawResponse, err := cAdmin.RawPost(getQuery("tenant/merge_tenant"),
		client.Var("name", "testtenant"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantResponse struct {
		Tenant_Merge *string `json:"tenant_Merge"`
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse)

	require.NotNil(t, tenantResponse.Tenant_Merge)
	require.Equal(t, "testtenant", *tenantResponse.Tenant_Merge)

	rawResponse2, err := cAdmin.RawPost(getQuery("tenant/merge_tenant"),
		client.Var("name", "testtenant"),
	)
	assertRawResponseSuccess(t, rawResponse2, err)

	var tenantResponse2 struct {
		Tenant_Merge *string `json:"tenant_Merge"`
	}

	err = decode.Decode(rawResponse2.Data.(map[string]any), &tenantResponse2)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse2)

	require.NotNil(t, tenantResponse2.Tenant_Merge)
	require.NotEqual(t, "testtenant", *tenantResponse2.Tenant_Merge)
	require.True(t, strings.HasPrefix(*tenantResponse2.Tenant_Merge, "testtenant"))

}

func TestMutationResolver_TenantMerge_AccessControlled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	rawResponse, err := c.RawPost(getQuery("tenant/merge_tenant"),
		client.Var("name", "testtenant"),
	)

	require.Nil(t, err)
	require.NotNil(t, rawResponse.Errors)

}

func TestMutationResolver_TenantMerge_CheckDefaultData(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateOrganizationRelationship(ctx, driver, "Customer")
	neo4jt.CreateOrganizationRelationship(ctx, driver, "Supplier")

	newTenantName := "test_tenant"
	rawResponse, err := cAdmin.RawPost(getQuery("tenant/merge_tenant"),
		client.Var("name", newTenantName),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Tenant"))
}

func TestMutationResolver_GetByWorkspace(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenant(ctx, driver, "other")
	neo4jtest.CreateWorkspace(ctx, driver, "testworkspace", "testprovider", tenantName)

	rawResponse, err := cAdmin.RawPost(getQuery("tenant/get_by_workspace"),
		client.Var("name", "testworkspace"),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantResponse struct {
		Tenant_ByWorkspace *string
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse)

	require.NotNil(t, tenantResponse.Tenant_ByWorkspace)
	require.Equal(t, tenantName, *tenantResponse.Tenant_ByWorkspace)

	rawResponse2, err := cAdmin.RawPost(getQuery("tenant/get_by_workspace"),
		client.Var("name", "testworkspace2"),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantResponse2 struct {
		Tenant_ByWorkspace *string
	}

	err = decode.Decode(rawResponse2.Data.(map[string]any), &tenantResponse2)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse2)
	require.Nil(t, tenantResponse2.Tenant_ByWorkspace)

}

func TestMutationResolver_GetByEmail(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{})
	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId, "test@openline.ai", false, "test")

	rawResponse, err := cAdmin.RawPost(getQuery("tenant/get_by_email"),
		client.Var("email", "test@openline.ai"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantResponse struct {
		Tenant_ByEmail *string
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse)

	require.NotNil(t, tenantResponse.Tenant_ByEmail)
	require.Equal(t, tenantName, *tenantResponse.Tenant_ByEmail)

	rawResponse2, err := cAdmin.RawPost(getQuery("tenant/get_by_email"),
		client.Var("email", "other@openline.ai"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantResponse2 struct {
		Tenant_ByEmail *string
	}
	err = decode.Decode(rawResponse2.Data.(map[string]any), &tenantResponse2)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse2)
	require.Nil(t, tenantResponse2.Tenant_ByEmail)

}

func TestQueryResolver_GetTenantBillingProfiles(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	profileId := neo4jtest.CreateTenantBillingProfile(ctx, driver, tenantName, neo4jentity.TenantBillingProfileEntity{
		Phone:                  "123456789",
		LegalName:              "test",
		AddressLine1:           "address1",
		AddressLine2:           "address2",
		AddressLine3:           "address3",
		Locality:               "locality",
		Country:                "country",
		Region:                 "region",
		Zip:                    "zip",
		VatNumber:              "vatNumber",
		SendInvoicesFrom:       "sendInvoicesFrom",
		SendInvoicesBcc:        "sendInvoicesBcc",
		CanPayWithPigeon:       true,
		CanPayWithBankTransfer: true,
		Check:                  true,
	})

	rawResponse, err := c.RawPost(getQuery("tenant/get_tenant_billing_profiles"))
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantGraphqlResponse struct {
		TenantBillingProfiles []model.TenantBillingProfile
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantGraphqlResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantGraphqlResponse)

	require.Equal(t, 1, len(tenantGraphqlResponse.TenantBillingProfiles))
	tenantBillingProfile := tenantGraphqlResponse.TenantBillingProfiles[0]
	require.Equal(t, profileId, tenantBillingProfile.ID)
	require.Equal(t, "123456789", tenantBillingProfile.Phone)
	require.Equal(t, "test", tenantBillingProfile.LegalName)
	require.Equal(t, "address1", tenantBillingProfile.AddressLine1)
	require.Equal(t, "address2", tenantBillingProfile.AddressLine2)
	require.Equal(t, "address3", tenantBillingProfile.AddressLine3)
	require.Equal(t, "locality", tenantBillingProfile.Locality)
	require.Equal(t, "country", tenantBillingProfile.Country)
	require.Equal(t, "region", tenantBillingProfile.Region)
	require.Equal(t, "zip", tenantBillingProfile.Zip)
	require.Equal(t, "vatNumber", tenantBillingProfile.VatNumber)
	require.Equal(t, "sendInvoicesFrom", tenantBillingProfile.SendInvoicesFrom)
	require.Equal(t, "sendInvoicesBcc", tenantBillingProfile.SendInvoicesBcc)
	require.True(t, tenantBillingProfile.CanPayWithPigeon)
	require.True(t, tenantBillingProfile.CanPayWithBankTransfer)
	require.True(t, tenantBillingProfile.Check)
}

func TestQueryResolver_GetTenantBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	profileId := neo4jtest.CreateTenantBillingProfile(ctx, driver, tenantName, neo4jentity.TenantBillingProfileEntity{
		Phone:                  "123456789",
		LegalName:              "test",
		AddressLine1:           "address1",
		AddressLine2:           "address2",
		AddressLine3:           "address3",
		Locality:               "locality",
		Country:                "country",
		Region:                 "region",
		Zip:                    "zip",
		VatNumber:              "vatNumber",
		SendInvoicesFrom:       "sendInvoicesFrom",
		SendInvoicesBcc:        "sendInvoicesBcc",
		CanPayWithPigeon:       true,
		CanPayWithBankTransfer: true,
		Check:                  true,
	})

	rawResponse, err := c.RawPost(getQuery("tenant/get_tenant_billing_profile"), client.Var("id", profileId))
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantGraphqlResponse struct {
		TenantBillingProfile model.TenantBillingProfile
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantGraphqlResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantGraphqlResponse)

	tenantBillingProfile := tenantGraphqlResponse.TenantBillingProfile
	require.Equal(t, profileId, tenantBillingProfile.ID)
	require.Equal(t, "123456789", tenantBillingProfile.Phone)
	require.Equal(t, "test", tenantBillingProfile.LegalName)
	require.Equal(t, "address1", tenantBillingProfile.AddressLine1)
	require.Equal(t, "address2", tenantBillingProfile.AddressLine2)
	require.Equal(t, "address3", tenantBillingProfile.AddressLine3)
	require.Equal(t, "locality", tenantBillingProfile.Locality)
	require.Equal(t, "country", tenantBillingProfile.Country)
	require.Equal(t, "region", tenantBillingProfile.Region)
	require.Equal(t, "zip", tenantBillingProfile.Zip)
	require.Equal(t, "vatNumber", tenantBillingProfile.VatNumber)
	require.Equal(t, "sendInvoicesFrom", tenantBillingProfile.SendInvoicesFrom)
	require.Equal(t, "sendInvoicesBcc", tenantBillingProfile.SendInvoicesBcc)
	require.True(t, tenantBillingProfile.CanPayWithPigeon)
	require.True(t, tenantBillingProfile.CanPayWithBankTransfer)
	require.True(t, tenantBillingProfile.Check)
}

func TestMutationResolver_TenantAddBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	profileId := uuid.New().String()
	calledAddTenanBillingProfile := false

	tenantServiceCallbacks := events_platform.MockTenantServiceCallbacks{
		AddBillingProfile: func(context context.Context, profile *tenantpb.AddBillingProfileRequest) (*commonpb.IdResponse, error) {
			require.Equal(t, tenantName, profile.Tenant)
			require.Equal(t, testUserId, profile.LoggedInUserId)
			require.Equal(t, neo4jentity.DataSourceOpenline.String(), profile.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, profile.SourceFields.AppSource)
			require.Equal(t, "phone", profile.Phone)
			require.Equal(t, "legalName", profile.LegalName)
			require.Equal(t, "addressLine1", profile.AddressLine1)
			require.Equal(t, "addressLine2", profile.AddressLine2)
			require.Equal(t, "addressLine3", profile.AddressLine3)
			require.Equal(t, "locality", profile.Locality)
			require.Equal(t, "country", profile.Country)
			require.Equal(t, "region", profile.Region)
			require.Equal(t, "zip", profile.Zip)
			require.Equal(t, "vatNumber", profile.VatNumber)
			require.Equal(t, "sendInvoicesFrom", profile.SendInvoicesFrom)
			require.Equal(t, "sendInvoicesBcc", profile.SendInvoicesBcc)
			require.True(t, profile.CanPayWithPigeon)
			require.True(t, profile.CanPayWithBankTransfer)
			require.True(t, profile.Check)

			calledAddTenanBillingProfile = true
			neo4jtest.CreateTenantBillingProfile(ctx, driver, tenantName, neo4jentity.TenantBillingProfileEntity{Id: profileId})
			return &commonpb.IdResponse{
				Id: profileId,
			}, nil
		},
	}
	events_platform.SetTenantCallbacks(&tenantServiceCallbacks)

	rawResponse := callGraphQL(t, "tenant/add_tenant_billing_profile", map[string]interface{}{})
	require.Nil(t, rawResponse.Errors)

	var billingProfileStruct struct {
		Tenant_AddBillingProfile model.TenantBillingProfile
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &billingProfileStruct)
	require.Nil(t, err)

	profile := billingProfileStruct.Tenant_AddBillingProfile
	require.Equal(t, profileId, profile.ID)

	require.True(t, calledAddTenanBillingProfile)
}

func TestMutationResolver_TenantUpdateBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	profileId := neo4jtest.CreateTenantBillingProfile(ctx, driver, tenantName, neo4jentity.TenantBillingProfileEntity{})
	calledUpdateTenantBillingProfile := false

	tenantServiceCallbacks := events_platform.MockTenantServiceCallbacks{
		UpdateBillingProfile: func(context context.Context, profile *tenantpb.UpdateBillingProfileRequest) (*commonpb.IdResponse, error) {
			require.Equal(t, tenantName, profile.Tenant)
			require.Equal(t, profileId, profile.Id)
			require.Equal(t, testUserId, profile.LoggedInUserId)
			require.Equal(t, constants.AppSourceCustomerOsApi, profile.AppSource)
			require.Equal(t, "phone", profile.Phone)
			require.Equal(t, "legalName", profile.LegalName)
			require.Equal(t, "addressLine1", profile.AddressLine1)
			require.Equal(t, "addressLine2", profile.AddressLine2)
			require.Equal(t, "addressLine3", profile.AddressLine3)
			require.Equal(t, "locality", profile.Locality)
			require.Equal(t, "country", profile.Country)
			require.Equal(t, "region", profile.Region)
			require.Equal(t, "zip", profile.Zip)
			require.Equal(t, "vatNumber", profile.VatNumber)
			require.Equal(t, "sendInvoicesFrom", profile.SendInvoicesFrom)
			require.Equal(t, "sendInvoicesBcc", profile.SendInvoicesBcc)
			require.True(t, profile.CanPayWithPigeon)
			require.True(t, profile.CanPayWithBankTransfer)
			require.True(t, profile.Check)
			require.ElementsMatch(t, []tenantpb.TenantBillingProfileFieldMask{
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_PHONE,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_LEGAL_NAME,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_1,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_2,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ADDRESS_LINE_3,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_LOCALITY,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_COUNTRY,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_ZIP,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_VAT_NUMBER,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_SEND_INVOICES_FROM,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_SEND_INVOICES_BCC,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_PIGEON,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CAN_PAY_WITH_BANK_TRANSFER,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_REGION,
				tenantpb.TenantBillingProfileFieldMask_TENANT_BILLING_PROFILE_FIELD_CHECK,
			},
				profile.FieldsMask)
			calledUpdateTenantBillingProfile = true
			return &commonpb.IdResponse{
				Id: profileId,
			}, nil
		},
	}
	events_platform.SetTenantCallbacks(&tenantServiceCallbacks)

	rawResponse := callGraphQL(t, "tenant/update_tenant_billing_profile", map[string]interface{}{"id": profileId})
	require.Nil(t, rawResponse.Errors)

	var billingProfileStruct struct {
		Tenant_UpdateBillingProfile model.TenantBillingProfile
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &billingProfileStruct)
	require.Nil(t, err)

	profile := billingProfileStruct.Tenant_UpdateBillingProfile
	require.Equal(t, profileId, profile.ID)

	require.True(t, calledUpdateTenantBillingProfile)
}

func TestQueryResolver_GetTenantSettings(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
		LogoRepositoryFileId: "logoRepositoryFileId",
		BaseCurrency:         neo4jenum.CurrencyUSD,
		InvoicingEnabled:     true,
	})

	rawResponse, err := c.RawPost(getQuery("tenant/get_tenant_settings"))
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantGraphqlResponse struct {
		TenantSettings model.TenantSettings
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantGraphqlResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantGraphqlResponse)

	tenantSettings := tenantGraphqlResponse.TenantSettings
	require.Equal(t, "logoRepositoryFileId", *tenantSettings.LogoRepositoryFileID)
	require.Equal(t, model.CurrencyUsd, *tenantSettings.BaseCurrency)
	require.Equal(t, true, tenantSettings.BillingEnabled)
}

func TestMutationResolver_TenantUpdateSettings(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{})
	calledUpdateTenantSettings := false

	tenantServiceCallbacks := events_platform.MockTenantServiceCallbacks{
		UpdateTenantSettings: func(context context.Context, profile *tenantpb.UpdateTenantSettingsRequest) (*emptypb.Empty, error) {
			require.Equal(t, tenantName, profile.Tenant)
			require.Equal(t, testUserId, profile.LoggedInUserId)
			require.Equal(t, constants.AppSourceCustomerOsApi, profile.AppSource)
			require.Equal(t, "123-456-789", profile.LogoRepositoryFileId)
			require.Equal(t, "EUR", profile.BaseCurrency)
			require.Equal(t, true, profile.InvoicingEnabled)
			require.ElementsMatch(t, []tenantpb.TenantSettingsFieldMask{
				tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_INVOICING_ENABLED,
				tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_BASE_CURRENCY,
				tenantpb.TenantSettingsFieldMask_TENANT_SETTINGS_FIELD_LOGO_REPOSITORY_FILE_ID,
			},
				profile.FieldsMask)
			calledUpdateTenantSettings = true
			return &emptypb.Empty{}, nil
		},
	}
	events_platform.SetTenantCallbacks(&tenantServiceCallbacks)

	rawResponse := callGraphQL(t, "tenant/update_tenant_settings", map[string]interface{}{})
	require.Nil(t, rawResponse.Errors)

	var responseStruct struct {
		Tenant_UpdateSettings model.TenantSettings
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &responseStruct)
	require.Nil(t, err)

	require.True(t, calledUpdateTenantSettings)
}
