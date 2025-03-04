package resolver

import (
	"context"
	"testing"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

	"github.com/stretchr/testify/assert"

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

	var tasks struct {
		Tasks_Search model.TaskSearchResult
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tasks)
	require.Nil(t, err)
	require.NotNil(t, tasks)

	searchResult := tasks.Tasks_Search
	require.Equal(t, totalAvailable, searchResult.TotalAvailable)
	require.Equal(t, totalElements, searchResult.TotalElements)
	require.Equal(t, int(totalElements), len(searchResult.Tasks))
}

func verifyTaskSortOrder(t *testing.T, sortBy postgresEntity.ColumnViewType, direction commonmodel.SortingDirection, expectedOrder []string) {
	sortedResult := assertTaskSort(t, sortBy, direction)
	assert.Equal(t, len(expectedOrder), len(sortedResult), "Mismatch in result length")
	for i, expected := range expectedOrder {
		assert.Equal(t, expected, sortedResult[i], "Mismatch at index %d", i)
	}
}

func assertTaskSort(t *testing.T, sortBy postgresEntity.ColumnViewType, sortDirection commonmodel.SortingDirection) []string {
	rawResponse, err := c.RawPost(getQuery("task/tasks_sort"),
		client.Var("limit", 10),
		client.Var("sortByField", string(sortBy)),
		client.Var("sortByDirection", sortDirection),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var tasks struct {
		Tasks_Search model.TaskSearchResult
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tasks)
	require.Nil(t, err)
	require.NotNil(t, tasks)

	return tasks.Tasks_Search.Tasks
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

	expectedAsc := []string{"2", "1", "empty"}
	expectedDesc := []string{"1", "2", "empty"}

	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeTasksSubject, commonmodel.SortingDirectionAsc, expectedAsc)
	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeTasksSubject, commonmodel.SortingDirectionDesc, expectedDesc)
}

func TestTaskResolver_SearchTasks_FilterByDescription(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Description: "aaAAaa description"})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Description: "bbBBbb description"})

	// Test case 1: Search with "aa" (case insensitive)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksDescription, "aa", commonmodel.ComparisonOperatorContains, 2, 1)

	// Test case 2: Search with "bb" (case insensitive)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksDescription, "bb", commonmodel.ComparisonOperatorContains, 2, 1)

	// Test case 3: Search with "AA" (case insensitive)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksDescription, "AA", commonmodel.ComparisonOperatorContains, 2, 1)

	// Test case 4: Search with "BB" (case insensitive)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksDescription, "BB", commonmodel.ComparisonOperatorContains, 2, 1)

	// Test case 5: Search with empty string (should return all)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksDescription, "", commonmodel.ComparisonOperatorContains, 2, 2)

	// Test case 6: Search with non-existent text
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksDescription, "xyz", commonmodel.ComparisonOperatorContains, 2, 0)
}

func TestTaskResolver_SearchTasks_SortByDescription(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	// Create tasks with different descriptions
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "1", Description: "Zebra description"})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "2", Description: "Apple description"})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "empty", Description: ""})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, commonmodel.NodeLabelTask))

	expectedAsc := []string{"2", "1", "empty"}
	expectedDesc := []string{"1", "2", "empty"}

	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeTasksDescription, commonmodel.SortingDirectionAsc, expectedAsc)
	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeTasksDescription, commonmodel.SortingDirectionDesc, expectedDesc)
}

func TestTaskResolver_SearchTasks_FilterByStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	// Create tasks with different statuses
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Status: enum.TaskStatusTodo})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Status: enum.TaskStatusDone})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Status: enum.TaskStatusInProgress})

	// Test case 1: Search with single status
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksStatus, []string{enum.TaskStatusTodo.String()}, commonmodel.ComparisonOperatorIn, 3, 1)

	// Test case 2: Search with multiple statuses
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksStatus, []string{enum.TaskStatusTodo.String(), enum.TaskStatusInProgress.String()}, commonmodel.ComparisonOperatorIn, 3, 2)

	// Test case 3: Search with all statuses
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksStatus, []string{enum.TaskStatusTodo.String(), enum.TaskStatusDone.String(), enum.TaskStatusInProgress.String()}, commonmodel.ComparisonOperatorIn, 3, 3)
}

func TestTaskResolver_SearchTasks_SortByStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	// Create tasks with different statuses
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "todo", Status: enum.TaskStatusTodo})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "in_progress", Status: enum.TaskStatusInProgress})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "done", Status: enum.TaskStatusDone})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, commonmodel.NodeLabelTask))

	// In ascending order: Todo -> In Progress -> Done
	expectedAsc := []string{"todo", "in_progress", "done"}
	// In descending order: Done -> In Progress -> Todo
	expectedDesc := []string{"done", "in_progress", "todo"}

	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeTasksStatus, commonmodel.SortingDirectionAsc, expectedAsc)
	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeTasksStatus, commonmodel.SortingDirectionDesc, expectedDesc)
}

func TestTaskResolver_SearchTasks_FilterByCreator(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	user1Id := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{})
	user2Id := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{})
	taskA := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{})
	taskB := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{})
	taskC := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{})
	neo4jtest.TaskCreatedBy(ctx, driver, taskA, user1Id)
	neo4jtest.TaskCreatedBy(ctx, driver, taskB, user2Id)
	neo4jtest.TaskCreatedBy(ctx, driver, taskC, user2Id)

	// Test case 1: Search with user1 ID (should return 1 task)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksAuthor, []string{user1Id}, commonmodel.ComparisonOperatorIn, 4, 1)

	// Test case 2: Search with user2 ID (should return 2 tasks)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksAuthor, []string{user2Id}, commonmodel.ComparisonOperatorIn, 4, 2)

	// Test case 3: Search with non-existent creator
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksAuthor, []string{"non-existent-id"}, commonmodel.ComparisonOperatorIn, 4, 0)

	// Test case 4: Search with not empty creator ID
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksAuthor, "", commonmodel.ComparisonOperatorIsNotEmpty, 4, 3)

	// Test case 5: Search with empty creator ID
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksAuthor, "", commonmodel.ComparisonOperatorIsEmpty, 4, 1)
}

func TestTaskResolver_SearchTasks_FilterByAssignee(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	user1Id := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{})
	user2Id := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{})
	taskA := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{})
	taskB := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{})
	taskC := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{})
	neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{})
	neo4jtest.TaskAssignedTo(ctx, driver, taskA, user1Id)
	neo4jtest.TaskAssignedTo(ctx, driver, taskB, user2Id)
	neo4jtest.TaskAssignedTo(ctx, driver, taskC, user2Id)

	// Test case 1: Search with user1 ID (should return 1 task)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksAssignees, []string{user1Id}, commonmodel.ComparisonOperatorIn, 4, 1)

	// Test case 2: Search with user2 ID (should return 2 tasks)
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksAssignees, []string{user2Id}, commonmodel.ComparisonOperatorIn, 4, 2)

	// Test case 3: Search with non-existent assignee
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksAssignees, []string{"non-existent-id"}, commonmodel.ComparisonOperatorIn, 4, 0)

	// Test case 4: Search with not empty assignee ID
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksAssignees, "", commonmodel.ComparisonOperatorIsNotEmpty, 4, 3)

	// Test case 5: Search with empty assignee ID
	assertTaskSearch(t, postgresEntity.ColumnViewTypeTasksAssignees, "", commonmodel.ComparisonOperatorIsEmpty, 4, 1)
}

func TestTaskResolver_SearchTasks_SortByCreator(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	// Create users with different first names
	user1Id := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{FirstName: "Zebra"})
	user2Id := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{FirstName: "Apple"})
	user3Id := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{FirstName: ""})

	// Create tasks and assign creators
	taskA := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "1"})
	taskB := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "2"})
	taskC := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "empty"})
	neo4jtest.TaskCreatedBy(ctx, driver, taskA, user1Id)
	neo4jtest.TaskCreatedBy(ctx, driver, taskB, user2Id)
	neo4jtest.TaskCreatedBy(ctx, driver, taskC, user3Id)

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, commonmodel.NodeLabelTask))

	// In ascending order: Apple -> Zebra -> empty
	expectedAsc := []string{"2", "1", "empty"}
	// In descending order: Zebra -> Apple -> empty
	expectedDesc := []string{"1", "2", "empty"}

	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeTasksAuthor, commonmodel.SortingDirectionAsc, expectedAsc)
	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeTasksAuthor, commonmodel.SortingDirectionDesc, expectedDesc)
}

func TestTaskResolver_SearchTasks_SortByAssignee(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	// Create users with different first names
	user1Id := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{FirstName: "Zebra"})
	user2Id := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{FirstName: "Apple"})
	user3Id := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{FirstName: ""})

	// Create tasks and assign users
	taskA := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "1"})
	taskB := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "2"})
	taskC := neo4jtest.CreateTask(ctx, driver, tenantName, neo4jentity.TaskEntity{Id: "empty"})
	neo4jtest.TaskAssignedTo(ctx, driver, taskA, user1Id)
	neo4jtest.TaskAssignedTo(ctx, driver, taskB, user2Id)
	neo4jtest.TaskAssignedTo(ctx, driver, taskC, user3Id)

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, commonmodel.NodeLabelTask))

	// In ascending order: Apple -> Zebra -> empty
	expectedAsc := []string{"2", "1", "empty"}
	// In descending order: Zebra -> Apple -> empty
	expectedDesc := []string{"1", "2", "empty"}

	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeTasksAssignees, commonmodel.SortingDirectionAsc, expectedAsc)
	verifyTaskSortOrder(t, postgresEntity.ColumnViewTypeTasksAssignees, commonmodel.SortingDirectionDesc, expectedDesc)
}
