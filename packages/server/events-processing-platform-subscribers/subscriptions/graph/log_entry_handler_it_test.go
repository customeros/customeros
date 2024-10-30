package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/model"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGraphLogEntryEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, "sf")
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{
		Name: "test org",
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "User": 1, "ExternalSystem": 1, "LogEntry": 0, "TimelineEvent": 0})

	// prepare event handler
	logEntryEventHandler := &LogEntryEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	now := utils.Now()
	logEntryId := uuid.New().String()
	logEntryAggregate := aggregate.NewLogEntryAggregateWithTenantAndID(tenantName, logEntryId)
	createEvent, err := event.NewLogEntryCreateEvent(logEntryAggregate, model.LogEntryDataFields{
		Content:              "test content",
		ContentType:          "test content type",
		AuthorUserId:         utils.StringPtr(userId),
		LoggedOrganizationId: utils.StringPtr(orgId),
	}, cmnmod.Source{
		Source:        constants.SourceOpenline,
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		SourceOfTruth: constants.SourceOpenline,
	}, cmnmod.ExternalSystem{
		ExternalSystemId: "sf",
		ExternalId:       "123",
	}, now, now, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = logEntryEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"User": 1, "User_" + tenantName: 1,
		"ExternalSystem": 1, "ExternalSystem_" + tenantName: 1,
		"LogEntry": 1, "LogEntry_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"LOGGED":         1,
		"CREATED_BY":     1,
		"IS_LINKED_WITH": 1,
	})

	logEntryDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "LogEntry_"+tenantName, logEntryId)
	require.Nil(t, err)
	require.NotNil(t, logEntryDbNode)

	// verify log entry
	logEntry := neo4jmapper.MapDbNodeToLogEntryEntity(logEntryDbNode)
	require.Equal(t, logEntryId, logEntry.Id)
	require.Equal(t, "test content", logEntry.Content)
	require.Equal(t, "test content type", logEntry.ContentType)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), logEntry.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, logEntry.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), logEntry.SourceOfTruth)
	require.Equal(t, now, logEntry.CreatedAt)
	test.AssertRecentTime(t, logEntry.UpdatedAt)
	require.Equal(t, now, logEntry.StartedAt)
}

func TestGraphLogEntryEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	logEntryId := neo4jtest.CreateLogEntryForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.LogEntryEntity{})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "LogEntry": 1, "TimelineEvent": 1})

	// prepare event handler
	logEntryEventHandler := &LogEntryEventHandler{
		services: testDatabase.Services,
	}
	now := utils.Now()
	logEntryAggregate := aggregate.NewLogEntryAggregateWithTenantAndID(tenantName, logEntryId)
	updateEvent, err := event.NewLogEntryUpdateEvent(logEntryAggregate, "test content", "test content type", "openline", now, now, nil)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = logEntryEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"LogEntry": 1, "LogEntry_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	logEntryDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "LogEntry_"+tenantName, logEntryId)
	require.Nil(t, err)
	require.NotNil(t, logEntryDbNode)

	// verify log entry
	logEntry := neo4jmapper.MapDbNodeToLogEntryEntity(logEntryDbNode)
	require.Equal(t, logEntryId, logEntry.Id)
	require.Equal(t, "test content", logEntry.Content)
	require.Equal(t, "test content type", logEntry.ContentType)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), logEntry.SourceOfTruth)
	test.AssertRecentTime(t, logEntry.UpdatedAt)
}
