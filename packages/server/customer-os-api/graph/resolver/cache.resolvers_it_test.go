package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_GCliCache_IsOwnerFalse(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateUserWithId(ctx, driver, tenantName, testUserId, entity.UserEntity{
		FirstName: "a",
		LastName:  "b",
	})
	neo4jt.CreateOrganization(ctx, driver, tenantName, "org1")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(ctx, driver, "OWNS"))

	rawResponse, err := c.RawPost(getQuery("cache/global_Cache"))
	assertRawResponseSuccess(t, rawResponse, err)

	var gcliCacheResponse struct {
		Global_Cache struct {
			User      model.User       `json:"user"`
			IsOwner   bool             `json:"isOwner"`
			GCliCache []model.GCliItem `json:"gcliCache"`
		}
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &gcliCacheResponse)
	require.Nil(t, err)
	require.NotNil(t, gcliCacheResponse)

	require.Equal(t, false, gcliCacheResponse.Global_Cache.IsOwner)
}

func TestQueryResolver_GCliCache_IsOwnerTrue(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	userId := neo4jt.CreateUserWithId(ctx, driver, tenantName, testUserId, entity.UserEntity{
		FirstName: "a",
		LastName:  "b",
	})
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org1")

	neo4jt.UserOwnsOrganization(ctx, driver, userId, organizationId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "OWNS"))

	rawResponse, err := c.RawPost(getQuery("cache/global_Cache"))
	assertRawResponseSuccess(t, rawResponse, err)

	var gcliCacheResponse struct {
		Global_Cache struct {
			User      model.User       `json:"user"`
			IsOwner   bool             `json:"isOwner"`
			GCliCache []model.GCliItem `json:"gcliCache"`
		}
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &gcliCacheResponse)
	require.Nil(t, err)
	require.NotNil(t, gcliCacheResponse)

	require.Equal(t, true, gcliCacheResponse.Global_Cache.IsOwner)
}

func TestQueryResolver_GCliCache(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateUserWithId(ctx, driver, tenantName, testUserId, entity.UserEntity{
		FirstName: "a",
		LastName:  "b",
	})

	neo4jt.CreateCountryWith(ctx, driver, "1", "USA", "United States")
	neo4jt.CreateState(ctx, driver, "USA", "Alabama", "AL")
	neo4jt.CreateState(ctx, driver, "USA", "Louisiana", "LA")

	services.Cache.InitCache()

	//neo4jt.CreateContactWith(ctx, driver, tenantName, "a", "1")
	//neo4jt.CreateContactWith(ctx, driver, tenantName, "ab", "2")
	//neo4jt.CreateContactWith(ctx, driver, tenantName, "abc", "3")
	//neo4jt.CreateContactWith(ctx, driver, tenantName, "abcd", "4")
	//neo4jt.CreateContactWith(ctx, driver, tenantName, "b", "1")

	//neo4jt.CreateOrganization(ctx, driver, tenantName, "a")
	//neo4jt.CreateOrganization(ctx, driver, tenantName, "ab")
	//neo4jt.CreateOrganization(ctx, driver, tenantName, "abc")
	//neo4jt.CreateOrganization(ctx, driver, tenantName, "abcd")
	//neo4jt.CreateOrganization(ctx, driver, tenantName, "b")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Country"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "State"))
	//require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	//require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))

	rawResponse, err := c.RawPost(getQuery("cache/global_Cache"))
	assertRawResponseSuccess(t, rawResponse, err)

	var gcliCacheResponse struct {
		Global_Cache struct {
			User      model.User       `json:"user"`
			IsOwner   bool             `json:"isOwner"`
			GCliCache []model.GCliItem `json:"gcliCache"`
		}
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &gcliCacheResponse)
	require.Nil(t, err)
	require.NotNil(t, gcliCacheResponse)

	require.Equal(t, 2, len(gcliCacheResponse.Global_Cache.GCliCache))

	require.Equal(t, "STATE", gcliCacheResponse.Global_Cache.GCliCache[0].Type.String())
	require.Equal(t, "STATE", gcliCacheResponse.Global_Cache.GCliCache[1].Type.String())
	//require.Equal(t, "CONTACT", gcliCacheResponse.Global_Cache.GCliCache[2].Type.String())
	//require.Equal(t, "CONTACT", gcliCacheResponse.Global_Cache.GCliCache[3].Type.String())
	//require.Equal(t, "CONTACT", gcliCacheResponse.Global_Cache.GCliCache[4].Type.String())
	//require.Equal(t, "CONTACT", gcliCacheResponse.Global_Cache.GCliCache[5].Type.String())
	//require.Equal(t, "ORGANIZATION", gcliCacheResponse.Global_Cache.GCliCache[6].Type.String())
	//require.Equal(t, "ORGANIZATION", gcliCacheResponse.Global_Cache.GCliCache[7].Type.String())
	//require.Equal(t, "ORGANIZATION", gcliCacheResponse.Global_Cache.GCliCache[8].Type.String())
	//require.Equal(t, "ORGANIZATION", gcliCacheResponse.Global_Cache.GCliCache[9].Type.String())
}

func TestQueryResolver_GCliCache_HasContracts_False(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateUserWithId(ctx, driver, tenantName, testUserId, entity.UserEntity{
		FirstName: "a",
		LastName:  "b",
	})

	neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "Contract"))

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

func TestQueryResolver_GCliCache_HasContracts_True(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateUserWithId(ctx, driver, tenantName, testUserId, entity.UserEntity{
		FirstName: "a",
		LastName:  "b",
	})

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contract"))

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
