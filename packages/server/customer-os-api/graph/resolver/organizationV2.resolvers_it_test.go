package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_UIOrganizationsSearch_FilterByName(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: ""})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "A closed organization"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "OPENLINE"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "the openline"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "some other open organization"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "OpEnLiNe"})

	require.Equal(t, 6, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	searchBy := model.ColumnViewTypeOrganizationsName
	searchTerm := "open"

	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsEmpty, 6, 1)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsNotEmpty, 6, 5)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorContains, 6, 4)
	assertSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorNotContains, 6, 2)
}

func assertSearch(t *testing.T, filterName model.ColumnViewType, searchValue any, operator commonModel.ComparisonOperator, totalAvailable int64, totalElements int64) {
	rawResponse, err := c.RawPost(getQuery("organization/ui_organizations_search"),
		client.Var("limit", 10),
		client.Var("filterName", filterName),
		client.Var("filterValue", searchValue),
		client.Var("filterOperation", operator),
		client.Var("sortByField", model.ColumnViewTypeOrganizationsName),
		client.Var("sortByDirection", commonModel.SortingDirectionAsc),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organizations struct {
		Ui_Organizations_Search model.OrganizationSearchResult
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizations)
	require.Nil(t, err)
	require.NotNil(t, organizations)

	searchResult := organizations.Ui_Organizations_Search
	require.Equal(t, totalAvailable, searchResult.TotalAvailable)
	require.Equal(t, totalElements, searchResult.TotalElements)
}
