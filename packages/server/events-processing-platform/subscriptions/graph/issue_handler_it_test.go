package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/model"
	orgaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	orgcmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGraphIssueEventHandler_OnCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	// prepare neo4j data
	externalSystemId := "sf"
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, externalSystemId)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "ExternalSystem": 1, "Issue": 0, "TimelineEvent": 0})

	// prepare event handler
	issueEventHandler := &GraphIssueEventHandler{
		Repositories:         testDatabase.Repositories,
		organizationCommands: orgcmdhnd.NewOrganizationCommands(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	issueId := uuid.New().String()
	issueAggregate := aggregate.NewIssueAggregateWithTenantAndID(tenantName, issueId)
	createEvent, err := event.NewIssueCreateEvent(issueAggregate, model.IssueDataFields{
		Subject:                  "test subject",
		Description:              "test description",
		Status:                   "open",
		Priority:                 "high",
		ReportedByOrganizationId: utils.StringPtr(orgId),
	}, cmnmod.Source{
		Source:        constants.SourceOpenline,
		AppSource:     constants.AppSourceEventProcessingPlatform,
		SourceOfTruth: constants.SourceOpenline,
	}, cmnmod.ExternalSystem{
		ExternalSystemId: "sf",
		ExternalId:       "123",
	}, now, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = issueEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"ExternalSystem": 1, "ExternalSystem_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jt.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"REPORTED_BY":    1,
		"IS_LINKED_WITH": 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, issueId, "REPORTED_BY", orgId)
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, issueId, "IS_LINKED_WITH", externalSystemId)

	issueDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Issue_"+tenantName, issueId)
	require.Nil(t, err)
	require.NotNil(t, issueDbNode)

	// verify issue
	issue := graph_db.MapDbNodeToIssueEntity(*issueDbNode)
	require.Equal(t, issueId, issue.Id)
	require.Equal(t, "test subject", issue.Subject)
	require.Equal(t, "test description", issue.Description)
	require.Equal(t, "open", issue.Status)
	require.Equal(t, "high", issue.Priority)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), issue.Source)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), issue.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, issue.AppSource)
	require.Equal(t, now, issue.CreatedAt)
	require.Equal(t, now, issue.UpdatedAt)

	// Check refresh last touchpoint event was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[orgAggregate.GetID()]
	require.Equal(t, 1, len(eventList))
	generatedEvent := eventList[0]
	require.Equal(t, orgevents.OrganizationRefreshLastTouchpointV1, generatedEvent.EventType)
	var eventData orgevents.OrganizationRefreshLastTouchpointEvent
	err = generatedEvent.GetJsonData(&eventData)
	require.Nil(t, err)
	require.Equal(t, tenantName, eventData.Tenant)
}

func TestGraphIssueEventHandler_OnUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	issueId := neo4jt.CreateIssue(ctx, testDatabase.Driver, tenantName, entity.IssueEntity{
		Subject:     "test subject",
		Description: "test description",
		Status:      "open",
		Priority:    "high",
	})
	neo4jt.LinkIssueReportedBy(ctx, testDatabase.Driver, issueId, orgId)
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Issue": 1, "TimelineEvent": 1})

	// prepare event handler
	issueEventHandler := &GraphIssueEventHandler{
		Repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	issueAggregate := aggregate.NewIssueAggregateWithTenantAndID(tenantName, issueId)
	updateEvent, err := event.NewIssueUpdateEvent(issueAggregate, model.IssueDataFields{
		Subject:     "test subject updated",
		Description: "test description updated",
		Status:      "closed",
		Priority:    "low",
	}, constants.SourceOpenline, cmnmod.ExternalSystem{}, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = issueEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	issueDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Issue_"+tenantName, issueId)
	require.Nil(t, err)
	require.NotNil(t, issueDbNode)

	// verify issue
	issue := graph_db.MapDbNodeToIssueEntity(*issueDbNode)
	require.Equal(t, issueId, issue.Id)
	require.Equal(t, "test subject updated", issue.Subject)
	require.Equal(t, "test description updated", issue.Description)
	require.Equal(t, "closed", issue.Status)
	require.Equal(t, "low", issue.Priority)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), issue.SourceOfTruth)
	require.Equal(t, now, issue.UpdatedAt)
}

func TestGraphIssueEventHandler_OnUpdate_CurrentSourceOpenline_UpdateSourceNonOpenline_UpdateOnlyEmptyFields(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	issueId := neo4jt.CreateIssue(ctx, testDatabase.Driver, tenantName, entity.IssueEntity{
		Subject:       "test subject",
		Description:   "test description",
		SourceOfTruth: constants.SourceOpenline,
	})
	neo4jt.LinkIssueReportedBy(ctx, testDatabase.Driver, issueId, orgId)
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Issue": 1, "TimelineEvent": 1})

	// prepare event handler
	issueEventHandler := &GraphIssueEventHandler{
		Repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	issueAggregate := aggregate.NewIssueAggregateWithTenantAndID(tenantName, issueId)
	updateEvent, err := event.NewIssueUpdateEvent(issueAggregate, model.IssueDataFields{
		Subject:     "test subject updated",
		Description: "test description updated",
		Status:      "closed",
		Priority:    "low",
	}, "hubspot", cmnmod.ExternalSystem{}, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = issueEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	issueDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Issue_"+tenantName, issueId)
	require.Nil(t, err)
	require.NotNil(t, issueDbNode)

	// verify issue
	issue := graph_db.MapDbNodeToIssueEntity(*issueDbNode)
	require.Equal(t, issueId, issue.Id)
	require.Equal(t, "test subject", issue.Subject)
	require.Equal(t, "test description", issue.Description)
	require.Equal(t, "closed", issue.Status)
	require.Equal(t, "low", issue.Priority)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), issue.SourceOfTruth)
	require.Equal(t, now, issue.UpdatedAt)
}
