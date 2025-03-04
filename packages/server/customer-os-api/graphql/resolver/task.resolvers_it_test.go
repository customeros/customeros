package resolver

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	"github.com/customeros/customeros/packages/server/customer-os-api/utils/decode"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/stretchr/testify/require"

	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
)

func assertTaskSearch(t *testing.T, filterName postgresEntity.ColumnViewType, searchValue any, operator commonmodel.ComparisonOperator, totalAvailable int64, totalElements int64) {
	rawResponse, err := c.RawPost(getQuery("task/tasks_search"),
		client.Var("limit", 10),
		client.Var("filterName", filterName),
		client.Var("filterValue", searchValue),
		client.Var("filterOperation", operator),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var contacts struct {
		Tasks_Search model.TaskSearchResult
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contacts)
	require.Nil(t, err)
	require.NotNil(t, contacts)

	searchResult := contacts.Tasks_Search
	require.Equal(t, totalAvailable, searchResult.TotalAvailable)
	require.Equal(t, totalElements, searchResult.TotalElements)
	require.Equal(t, int(totalElements), len(searchResult.Tasks))
}

func verifyTaskSortOrder(t *testing.T, sortBy postgresEntity.ColumnViewType, direction commonmodel.SortingDirection, expectedOrder []string) {
	sortedResult := assertContactSort(t, sortBy, direction)
	assert.Equal(t, len(expectedOrder), len(sortedResult), "Mismatch in result length")
	for i, expected := range expectedOrder {
		assert.Equal(t, expected, sortedResult[i], "Mismatch at index %d", i)
	}
}

func assertTaskSort(t *testing.T, sortBy postgresEntity.ColumnViewType, sortDirection commonmodel.SortingDirection) []string {
	rawResponse, err := c.RawPost(getQuery("task/tasks_search"),
		client.Var("limit", 10),
		client.Var("sortByField", string(sortBy)),
		client.Var("sortByDirection", sortDirection),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var contacts struct {
		Tasks_Search model.ContactSearchResult
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contacts)
	require.Nil(t, err)
	require.NotNil(t, contacts)

	return contacts.Tasks_Search.Ids
}

func TestTaskResolver_SearchTasks_FilterBySubject(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Subject: "aaAAaa"})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Subject: "bbBBbb"})

	// Test case 1: Search with "aa" (case insensitive)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksSubject, "aa", commonmodel.ComparisonOperatorContains, 2, 1)

	// Test case 2: Search with "bb" (case insensitive)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksSubject, "bb", commonmodel.ComparisonOperatorContains, 2, 1)

	// Test case 3: Search with "AA" (case insensitive)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksSubject, "AA", commonmodel.ComparisonOperatorContains, 2, 1)

	// Test case 4: Search with "BB" (case insensitive)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksSubject, "BB", commonmodel.ComparisonOperatorContains, 2, 1)

	// Test case 5: Search with empty string (should return all)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksSubject, "", commonmodel.ComparisonOperatorContains, 2, 2)

	// Test case 6: Search with non-existent text
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksSubject, "xyz", commonmodel.ComparisonOperatorContains, 2, 0)
}

func TestTaskResolver_SearchTasks_SortBySubject(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	// Create tasks with different subjects
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "1", Subject: "Zebra"})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "2", Subject: "Apple"})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "empty", Subject: ""})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, commonmodel.NodeLabelTask))

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeContactsName, commonmodel.SortingDirectionAsc, expectedAsc)
	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeContactsName, commonmodel.SortingDirectionDesc, expectedDesc)
}
