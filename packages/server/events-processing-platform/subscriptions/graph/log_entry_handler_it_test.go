package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/mocked_grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGraphLogEntryEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{})
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, "sf")
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "User": 1, "ExternalSystem": 1, "LogEntry": 0, "TimelineEvent": 0})

	// prepare grpc mock
	calledEventsPlatform := false
	organizationServiceCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		RefreshLastTouchpoint: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, org.AppSource)
			calledEventsPlatform = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
	}
	mocked_grpc.SetOrganizationCallbacks(&organizationServiceCallbacks)

	// prepare event handler
	logEntryEventHandler := &LogEntryEventHandler{
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
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
		AppSource:     constants.AppSourceEventProcessingPlatform,
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
	logEntry := graph_db.MapDbNodeToLogEntryEntity(*logEntryDbNode)
	require.Equal(t, logEntryId, logEntry.Id)
	require.Equal(t, "test content", logEntry.Content)
	require.Equal(t, "test content type", logEntry.ContentType)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), logEntry.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, logEntry.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), logEntry.SourceOfTruth)
	require.Equal(t, now, logEntry.CreatedAt)
	require.Equal(t, now, logEntry.UpdatedAt)
	require.Equal(t, now, logEntry.StartedAt)

	// Check refresh last touch point
	require.Truef(t, calledEventsPlatform, "RefreshLastTouchpoint was not invoked")
}

func TestGraphLogEntryEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	logEntryId := neo4jt.CreateLogEntryForOrg(ctx, testDatabase.Driver, tenantName, orgId, entity.LogEntryEntity{})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "LogEntry": 1, "TimelineEvent": 1})

	// prepare event handler
	logEntryEventHandler := &LogEntryEventHandler{
		repositories: testDatabase.Repositories,
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
	logEntry := graph_db.MapDbNodeToLogEntryEntity(*logEntryDbNode)
	require.Equal(t, logEntryId, logEntry.Id)
	require.Equal(t, "test content", logEntry.Content)
	require.Equal(t, "test content type", logEntry.ContentType)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), logEntry.SourceOfTruth)
	require.Equal(t, now, logEntry.UpdatedAt)
	require.Equal(t, now, logEntry.StartedAt)
}

func TestGraphLogEntryEventHandler_OnAddTag(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	logEntryId := neo4jt.CreateLogEntryForOrg(ctx, testDatabase.Driver, tenantName, orgId, entity.LogEntryEntity{})
	tagId := neo4jt.CreateTag(ctx, testDatabase.Driver, tenantName, entity.TagEntity{})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "LogEntry": 1, "TimelineEvent": 1, "Tag": 1})

	// prepare event handler
	logEntryEventHandler := &LogEntryEventHandler{
		repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	logEntryAggregate := aggregate.NewLogEntryAggregateWithTenantAndID(tenantName, logEntryId)
	addTagEvent, err := event.NewLogEntryAddTagEvent(logEntryAggregate, tagId, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = logEntryEventHandler.OnAddTag(context.Background(), addTagEvent)
	require.Nil(t, err, "failed to execute event handler")

	// CHECK
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"LogEntry": 1, "LogEntry_" + tenantName: 1,
		"Tag": 1, "Tag_" + tenantName: 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"TAGGED": 1,
	})
	relationship, err := neo4jtest.GetRelationship(ctx, testDatabase.Driver, logEntryId, tagId)
	require.Nil(t, err)
	require.NotNil(t, relationship)
	require.Equal(t, "TAGGED", relationship.Type)
	require.Equal(t, now, relationship.Props["taggedAt"])
}

func TestGraphLogEntryEventHandler_OnRemoveTag(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	logEntryId := neo4jt.CreateLogEntryForOrg(ctx, testDatabase.Driver, tenantName, orgId, entity.LogEntryEntity{})
	tagId := neo4jt.CreateTag(ctx, testDatabase.Driver, tenantName, entity.TagEntity{})
	neo4jt.LinkTag(ctx, testDatabase.Driver, tagId, logEntryId)
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "LogEntry": 1, "TimelineEvent": 1, "Tag": 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"TAGGED": 1,
	})

	// prepare event handler
	logEntryEventHandler := &LogEntryEventHandler{
		repositories: testDatabase.Repositories,
	}
	logEntryAggregate := aggregate.NewLogEntryAggregateWithTenantAndID(tenantName, logEntryId)
	removeTagEvent, err := event.NewLogEntryRemoveTagEvent(logEntryAggregate, tagId)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = logEntryEventHandler.OnRemoveTag(context.Background(), removeTagEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"LogEntry": 1, "LogEntry_" + tenantName: 1,
		"Tag": 1, "Tag_" + tenantName: 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"TAGGED": 0,
	})
}
