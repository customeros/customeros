package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryResolver_LogEntry(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	secAgo60 := utils.Now().Add(-60 * time.Second)
	secAgo30 := utils.Now().Add(-30 * time.Second)
	orgId := neo4jt.CreateOrganization(ctx, driver, tenantName, "testOrganization")
	logEntryId := neo4jt.CreateLogEntryForOrganization(ctx, driver, tenantName, orgId, entity.LogEntryEntity{
		StartedAt:   secAgo60,
		Content:     "log entry content",
		ContentType: "text/plain",
	})

	tagId1 := neo4jt.CreateTag(ctx, driver, tenantName, "red")
	tagId2 := neo4jt.CreateTag(ctx, driver, tenantName, "blue")

	neo4jt.TagLogEntry(ctx, driver, logEntryId, tagId1, &secAgo30)
	neo4jt.TagLogEntry(ctx, driver, logEntryId, tagId2, nil)

	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	neo4jt.LogEntryCreatedByUser(ctx, driver, logEntryId, userId)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Organization": 1,
		"LogEntry":     1,
		"Tag":          2,
		"User":         1,
	})
	assertNeo4jRelationCount(ctx, t, driver, map[string]int{
		"TAGGED":     2,
		"LOGGED":     1,
		"CREATED_BY": 1,
	})

	rawResponse := callGraphQL(t, "log_entry/get_log_entry", map[string]interface{}{
		"id": logEntryId,
	})

	var logEntryStruct struct {
		LogEntry model.LogEntry
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &logEntryStruct)
	logEntry := logEntryStruct.LogEntry

	require.Nil(t, err)
	require.NotNil(t, logEntry)
	require.NotNil(t, logEntry.CreatedAt)
	require.NotNil(t, logEntry.UpdatedAt)
	require.Equal(t, "log entry content", *logEntry.Content)
	require.Equal(t, "text/plain", *logEntry.ContentType)
	require.Equal(t, secAgo60, logEntry.StartedAt)
	require.Equal(t, userId, logEntry.CreatedBy.ID)
	require.Equal(t, 2, len(logEntry.Tags))
	firstTag := logEntry.Tags[0]
	require.Equal(t, tagId1, firstTag.ID)
	require.Equal(t, "red", firstTag.Name)
	secondTag := logEntry.Tags[1]
	require.Equal(t, tagId2, secondTag.ID)
	require.Equal(t, "blue", secondTag.Name)
}
