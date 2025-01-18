package resolver

import (
	"context"
	"testing"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
)

func TestQueryGlobalCache_GCliCache_IsOwnerTrue(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		Id:        testUserId,
		FirstName: "a",
		LastName:  "b",
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "org1"})

	neo4jtest.UserOwnsOrganization(ctx, driver, testUserId, organizationId)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))

	rawResponse, err := c.RawPost(getQuery("cache/global_Cache"))
	assertRawResponseSuccess(t, rawResponse, err)

	var gcliCacheResponse struct {
		Global_Cache struct {
			User    model.User `json:"user"`
			IsOwner bool       `json:"isOwner"`
		}
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &gcliCacheResponse)
	require.Nil(t, err)
	require.NotNil(t, gcliCacheResponse)

	require.Equal(t, true, gcliCacheResponse.Global_Cache.IsOwner)
}

func TestQueryGlobalCache_GCliCache_HasContracts_False(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		Id:        testUserId,
		FirstName: "a",
		LastName:  "b",
	})

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 0, neo4jtest.GetCountOfNodes(ctx, driver, "Contract"))

	rawResponse, err := c.RawPost(getQuery("cache/global_Cache"))
	assertRawResponseSuccess(t, rawResponse, err)

	var gcliCacheResponse struct {
		Global_Cache struct {
			ContractsExist bool `json:"contractsExist"`
		}
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &gcliCacheResponse)
	require.Nil(t, err)
	require.NotNil(t, gcliCacheResponse)

	require.Equal(t, false, gcliCacheResponse.Global_Cache.ContractsExist)
}

func TestQueryGlobalCache_GCliCache_HasContracts_True(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		Id:        testUserId,
		FirstName: "a",
		LastName:  "b",
	})

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contract"))

	rawResponse, err := c.RawPost(getQuery("cache/global_Cache"))
	assertRawResponseSuccess(t, rawResponse, err)

	var gcliCacheResponse struct {
		Global_Cache struct {
			ContractsExist bool `json:"contractsExist"`
		}
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &gcliCacheResponse)
	require.Nil(t, err)
	require.NotNil(t, gcliCacheResponse)

	require.Equal(t, true, gcliCacheResponse.Global_Cache.ContractsExist)
}

func TestQueryGlobalCache_No_Logo(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		Id:        testUserId,
		FirstName: "a",
		LastName:  "b",
	})

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "User"))

	rawResponse, err := c.RawPost(getQuery("cache/global_Cache"))
	assertRawResponseSuccess(t, rawResponse, err)

	var gcliCacheResponse struct {
		Global_Cache struct {
			CdnLogoUrl string `json:"cdnLogoUrl"`
		}
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &gcliCacheResponse)
	require.Nil(t, err)
	require.NotNil(t, gcliCacheResponse)

	require.Equal(t, "", gcliCacheResponse.Global_Cache.CdnLogoUrl)
}

func TestQueryGlobalCache_Has_Logo(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		Id:        testUserId,
		FirstName: "a",
		LastName:  "b",
	})

	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
		LogoRepositoryFileId: "1",
	})
	neo4jt.CreateAttachment(ctx, driver, tenantName, neo4jentity.AttachmentEntity{
		Id:     "1",
		CdnUrl: "https://cdn.openline.com/1",
	})

	rawResponse, err := c.RawPost(getQuery("cache/global_Cache"))
	assertRawResponseSuccess(t, rawResponse, err)

	var gcliCacheResponse struct {
		Global_Cache struct {
			CdnLogoUrl string `json:"cdnLogoUrl"`
		}
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &gcliCacheResponse)
	require.Nil(t, err)
	require.NotNil(t, gcliCacheResponse)

	require.Equal(t, "https://cdn.openline.com/1", gcliCacheResponse.Global_Cache.CdnLogoUrl)
}
