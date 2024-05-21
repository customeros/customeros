package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_BillableInfo(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	whiteOrg1 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Hide: false,
	})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Hide: false,
	})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Hide: true,
	})
	whiteContact1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.LinkContactWithOrganization(ctx, driver, whiteContact1, whiteOrg1)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 3, "Contact": 1})

	rawResponse := callGraphQL(t, "tenant/get_billable_info", map[string]interface{}{})

	var billableInfoStruct struct {
		BillableInfo model.TenantBillableInfo
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &billableInfoStruct)
	require.Nil(t, err)

	require.NotNil(t, billableInfoStruct.BillableInfo)
	require.Equal(t, int64(2), billableInfoStruct.BillableInfo.WhitelistedOrganizations)
	require.Equal(t, int64(1), billableInfoStruct.BillableInfo.GreylistedOrganizations)
	require.Equal(t, int64(1), billableInfoStruct.BillableInfo.WhitelistedContacts)
	require.Equal(t, int64(0), billableInfoStruct.BillableInfo.GreylistedContacts)
}
