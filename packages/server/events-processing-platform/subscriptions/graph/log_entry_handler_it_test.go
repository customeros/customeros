package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGraphLogEntryEventHandler_OnCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{})
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "User": 1, "LogEntry": 0, "TimelineEvent": 0})

	// prepare event handler
	logEntryEventHandler := &GraphLogEntryEventHandler{
		Repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	logEntryId := uuid.New().String()
	logEntryAggregate := aggregate.NewLogEntryAggregateWithTenantAndID(tenantName, logEntryId)
	event, err := events.NewLogEntryCreateEvent(logEntryAggregate, models.LogEntryDataFields{
		Content:              "test content",
		ContentType:          "test content type",
		AuthorUserId:         utils.StringPtr(userId),
		LoggedOrganizationId: utils.StringPtr(orgId),
	}, commonModels.Source{
		Source:        constants.SourceOpenline,
		AppSource:     constants.AppSourceEventProcessingPlatform,
		SourceOfTruth: constants.SourceOpenline,
	}, now, now, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = logEntryEventHandler.OnCreate(context.Background(), event)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"User": 1, "User_" + tenantName: 1,
		"LogEntry": 1, "LogEntry_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jt.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"LOGGED":     1,
		"CREATED_BY": 1,
	})

	logEntryDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "LogEntry_"+tenantName, logEntryId)
	require.Nil(t, err)
	require.NotNil(t, logEntryDbNode)

	// verify log entry
	logEntry := graph_db.MapDbNodeToLogEntryEntity(*logEntryDbNode)
	require.Equal(t, logEntryId, logEntry.Id)
	require.Equal(t, "test content", logEntry.Content)
	require.Equal(t, "test content type", logEntry.ContentType)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), logEntry.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, logEntry.AppSource)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), logEntry.SourceOfTruth)
	require.Equal(t, now, logEntry.CreatedAt)
	require.Equal(t, now, logEntry.UpdatedAt)
	require.Equal(t, now, logEntry.StartedAt)
}

func TestGraphLogEntryEventHandler_OnUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	logEntryId := neo4jt.CreateLogEntryForOrg(ctx, testDatabase.Driver, tenantName, orgId, entity.LogEntryEntity{})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "LogEntry": 1, "TimelineEvent": 1})

	// prepare event handler
	logEntryEventHandler := &GraphLogEntryEventHandler{
		Repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	logEntryAggregate := aggregate.NewLogEntryAggregateWithTenantAndID(tenantName, logEntryId)
	event, err := events.NewLogEntryUpdateEvent(logEntryAggregate, "test content", "test content type", "openline", now, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = logEntryEventHandler.OnUpdate(context.Background(), event)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"LogEntry": 1, "LogEntry_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	logEntryDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "LogEntry_"+tenantName, logEntryId)
	require.Nil(t, err)
	require.NotNil(t, logEntryDbNode)

	// verify log entry
	logEntry := graph_db.MapDbNodeToLogEntryEntity(*logEntryDbNode)
	require.Equal(t, logEntryId, logEntry.Id)
	require.Equal(t, "test content", logEntry.Content)
	require.Equal(t, "test content type", logEntry.ContentType)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), logEntry.SourceOfTruth)
	require.Equal(t, now, logEntry.UpdatedAt)
	require.Equal(t, now, logEntry.StartedAt)
}

func TestGraphLogEntryEventHandler_OnAddTag(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	logEntryId := neo4jt.CreateLogEntryForOrg(ctx, testDatabase.Driver, tenantName, orgId, entity.LogEntryEntity{})
	tagId := neo4jt.CreateTag(ctx, testDatabase.Driver, tenantName, entity.TagEntity{})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "LogEntry": 1, "TimelineEvent": 1, "Tag": 1})

	// prepare event handler
	logEntryEventHandler := &GraphLogEntryEventHandler{
		Repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	logEntryAggregate := aggregate.NewLogEntryAggregateWithTenantAndID(tenantName, logEntryId)
	event, err := events.NewLogEntryAddTagEvent(logEntryAggregate, tagId, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = logEntryEventHandler.OnAddTag(context.Background(), event)
	require.Nil(t, err, "failed to execute event handler")

	// CHECK
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"LogEntry": 1, "LogEntry_" + tenantName: 1,
		"Tag": 1, "Tag_" + tenantName: 1})
	neo4jt.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"TAGGED": 1,
	})
	relationship, err := neo4jt.GetRelationship(ctx, testDatabase.Driver, logEntryId, tagId)
	require.Nil(t, err)
	require.NotNil(t, relationship)
	require.Equal(t, "TAGGED", relationship.Type)
	require.Equal(t, now, relationship.Props["taggedAt"])
}

func TestGraphLogEntryEventHandler_OnRemoveTag(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	logEntryId := neo4jt.CreateLogEntryForOrg(ctx, testDatabase.Driver, tenantName, orgId, entity.LogEntryEntity{})
	tagId := neo4jt.CreateTag(ctx, testDatabase.Driver, tenantName, entity.TagEntity{})
	neo4jt.LinkTag(ctx, testDatabase.Driver, tagId, logEntryId)
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "LogEntry": 1, "TimelineEvent": 1, "Tag": 1})
	neo4jt.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"TAGGED": 1,
	})

	// prepare event handler
	logEntryEventHandler := &GraphLogEntryEventHandler{
		Repositories: testDatabase.Repositories,
	}
	logEntryAggregate := aggregate.NewLogEntryAggregateWithTenantAndID(tenantName, logEntryId)
	event, err := events.NewLogEntryRemoveTagEvent(logEntryAggregate, tagId)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = logEntryEventHandler.OnRemoveTag(context.Background(), event)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"LogEntry": 1, "LogEntry_" + tenantName: 1,
		"Tag": 1, "Tag_" + tenantName: 1})
	neo4jt.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"TAGGED": 0,
	})
}
