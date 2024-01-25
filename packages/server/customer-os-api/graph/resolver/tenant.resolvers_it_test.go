package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestMutationResolver_TenantMerge(t *testing.T) {
	ctx := context.TODO()
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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	rawResponse, err := c.RawPost(getQuery("tenant/merge_tenant"),
		client.Var("name", "testtenant"),
	)

	require.Nil(t, err)
	require.NotNil(t, rawResponse.Errors)

}

func TestMutationResolver_TenantMerge_CheckDefaultData(t *testing.T) {
	ctx := context.TODO()
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
	ctx := context.TODO()
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
	ctx := context.TODO()
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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	profileId := neo4jtest.CreateTenantBillingProfile(ctx, driver, tenantName, neo4jentity.TenantBillingProfileEntity{
		Email:                         "test@gmail.com",
		Phone:                         "123456789",
		LegalName:                     "test",
		AddressLine1:                  "address1",
		AddressLine2:                  "address2",
		AddressLine3:                  "address3",
		Locality:                      "locality",
		Country:                       "country",
		Zip:                           "zip",
		DomesticPaymentsBankInfo:      "domesticPaymentsBankInfo",
		InternationalPaymentsBankInfo: "internationalPaymentsBankInfo",
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
	require.Equal(t, "test@gmail.com", tenantBillingProfile.Email)
	require.Equal(t, "123456789", tenantBillingProfile.Phone)
	require.Equal(t, "test", tenantBillingProfile.LegalName)
	require.Equal(t, "address1", tenantBillingProfile.AddressLine1)
	require.Equal(t, "address2", tenantBillingProfile.AddressLine2)
	require.Equal(t, "address3", tenantBillingProfile.AddressLine3)
	require.Equal(t, "locality", tenantBillingProfile.Locality)
	require.Equal(t, "country", tenantBillingProfile.Country)
	require.Equal(t, "zip", tenantBillingProfile.Zip)
	require.Equal(t, "domesticPaymentsBankInfo", tenantBillingProfile.DomesticPaymentsBankInfo)
	require.Equal(t, "internationalPaymentsBankInfo", tenantBillingProfile.InternationalPaymentsBankInfo)
}
