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
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/mocked_grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/neo4j"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/model"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphIssueEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	externalSystemId := "sf"
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, externalSystemId)
	reporterOrgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	submitterOrgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	submitterUserId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 1, "Organization": 2, "ExternalSystem": 1, "Issue": 0, "TimelineEvent": 0})

	// prepare grpc mock
	lastTouchpointInvoked := false
	organizationServiceCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		RefreshLastTouchpoint: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, reporterOrgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, org.AppSource)
			lastTouchpointInvoked = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: reporterOrgId,
			}, nil
		},
	}
	mocked_grpc.SetOrganizationCallbacks(&organizationServiceCallbacks)

	// prepare event handler
	issueEventHandler := &IssueEventHandler{
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	now := utils.Now()
	issueId := uuid.New().String()
	issueAggregate := aggregate.NewIssueAggregateWithTenantAndID(tenantName, issueId)
	createEvent, err := event.NewIssueCreateEvent(issueAggregate, model.IssueDataFields{
		Subject:                   "test subject",
		Description:               "test description",
		Status:                    "open",
		Priority:                  "high",
		ReportedByOrganizationId:  utils.StringPtr(reporterOrgId),
		SubmittedByOrganizationId: utils.StringPtr(submitterOrgId),
		SubmittedByUserId:         utils.StringPtr(submitterUserId),
	}, cmnmod.Source{
		Source:        constants.SourceOpenline,
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		SourceOfTruth: constants.SourceOpenline,
	}, cmnmod.ExternalSystem{
		ExternalSystemId: "sf",
		ExternalId:       "123",
	}, now, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = issueEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 1, "User_" + tenantName: 1,
		"Organization": 2, "Organization_" + tenantName: 2,
		"ExternalSystem": 1, "ExternalSystem_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"REPORTED_BY":    1,
		"SUBMITTED_BY":   2,
		"IS_LINKED_WITH": 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, issueId, "REPORTED_BY", reporterOrgId)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, issueId, "SUBMITTED_BY", submitterOrgId)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, issueId, "SUBMITTED_BY", submitterUserId)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, issueId, "IS_LINKED_WITH", externalSystemId)

	issueDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Issue_"+tenantName, issueId)
	require.Nil(t, err)
	require.NotNil(t, issueDbNode)

	// verify issue
	issue := neo4jmapper.MapDbNodeToIssueEntity(issueDbNode)
	require.Equal(t, issueId, issue.Id)
	require.Equal(t, "test subject", issue.Subject)
	require.Equal(t, "test description", issue.Description)
	require.Equal(t, "open", issue.Status)
	require.Equal(t, "high", issue.Priority)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), issue.Source)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), issue.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, issue.AppSource)
	require.Equal(t, now, issue.CreatedAt)
	test.AssertRecentTime(t, issue.UpdatedAt)

	// Check refresh last touchpoint
	require.Truef(t, lastTouchpointInvoked, "RefreshLastTouchpoint was not invoked")
}

func TestGraphIssueEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	issueId := neo4jt.CreateIssue(ctx, testDatabase.Driver, tenantName, neo4jentity.IssueEntity{
		Subject:     "test subject",
		Description: "test description",
		Status:      "open",
		Priority:    "high",
	})
	neo4jt.LinkIssueReportedBy(ctx, testDatabase.Driver, issueId, orgId)
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Issue": 1, "TimelineEvent": 1})

	// prepare event handler
	issueEventHandler := &IssueEventHandler{
		repositories: testDatabase.Repositories,
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

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	issueDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Issue_"+tenantName, issueId)
	require.Nil(t, err)
	require.NotNil(t, issueDbNode)

	// verify issue
	issue := neo4jmapper.MapDbNodeToIssueEntity(issueDbNode)
	require.Equal(t, issueId, issue.Id)
	require.Equal(t, "test subject updated", issue.Subject)
	require.Equal(t, "test description updated", issue.Description)
	require.Equal(t, "closed", issue.Status)
	require.Equal(t, "low", issue.Priority)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), issue.SourceOfTruth)
	test.AssertRecentTime(t, issue.UpdatedAt)
}

func TestGraphIssueEventHandler_OnUpdate_CurrentSourceOpenline_UpdateSourceNonOpenline_UpdateOnlyEmptyFields(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	issueId := neo4jt.CreateIssue(ctx, testDatabase.Driver, tenantName, neo4jentity.IssueEntity{
		Subject:       "test subject",
		Description:   "test description",
		SourceOfTruth: constants.SourceOpenline,
	})
	neo4jt.LinkIssueReportedBy(ctx, testDatabase.Driver, issueId, orgId)
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Issue": 1, "TimelineEvent": 1})

	// prepare event handler
	issueEventHandler := &IssueEventHandler{
		repositories: testDatabase.Repositories,
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

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	issueDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Issue_"+tenantName, issueId)
	require.Nil(t, err)
	require.NotNil(t, issueDbNode)

	// verify issue
	issue := neo4jmapper.MapDbNodeToIssueEntity(issueDbNode)
	require.Equal(t, issueId, issue.Id)
	require.Equal(t, "test subject", issue.Subject)
	require.Equal(t, "test description", issue.Description)
	require.Equal(t, "closed", issue.Status)
	require.Equal(t, "low", issue.Priority)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), issue.SourceOfTruth)
	test.AssertRecentTime(t, issue.UpdatedAt)
}

func TestGraphIssueEventHandler_OnAddUserAssignee(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	issueId := neo4jt.CreateIssue(ctx, testDatabase.Driver, tenantName, neo4jentity.IssueEntity{})
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "Issue": 1, "TimelineEvent": 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{"ASSIGNED_TO": 0})

	// prepare event handler
	issueEventHandler := &IssueEventHandler{
		repositories: testDatabase.Repositories,
	}
	updatedAt := utils.Now().Add(time.Duration(-1) * time.Minute)
	issueAggregate := aggregate.NewIssueAggregateWithTenantAndID(tenantName, issueId)
	addUserAssigneeEvent, err := event.NewIssueAddUserAssigneeEvent(issueAggregate, userId, updatedAt)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = issueEventHandler.OnAddUserAssignee(context.Background(), addUserAssigneeEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 1, "User_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, issueId, "ASSIGNED_TO", userId)

	issueDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Issue_"+tenantName, issueId)
	require.Nil(t, err)
	require.NotNil(t, issueDbNode)
	issue := neo4jmapper.MapDbNodeToIssueEntity(issueDbNode)
	test.AssertRecentTime(t, issue.UpdatedAt)
}

func TestGraphIssueEventHandler_OnRemoveUserAssignee(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	issueId := neo4jt.CreateIssue(ctx, testDatabase.Driver, tenantName, neo4jentity.IssueEntity{})
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
	neo4jt.LinkIssueAssignedTo(ctx, testDatabase.Driver, issueId, userId)
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "Issue": 1, "TimelineEvent": 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{"ASSIGNED_TO": 1})

	// prepare event handler
	issueEventHandler := &IssueEventHandler{
		repositories: testDatabase.Repositories,
	}
	updatedAt := utils.Now().Add(time.Duration(-1) * time.Hour)
	issueAggregate := aggregate.NewIssueAggregateWithTenantAndID(tenantName, issueId)
	removeUserAssigneeEvent, err := event.NewIssueRemoveUserAssigneeEvent(issueAggregate, userId, updatedAt)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = issueEventHandler.OnRemoveUserAssignee(context.Background(), removeUserAssigneeEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 1, "User_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{"ASSIGNED_TO": 0})

	issueDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Issue_"+tenantName, issueId)
	require.Nil(t, err)
	require.NotNil(t, issueDbNode)
	issue := neo4jmapper.MapDbNodeToIssueEntity(issueDbNode)
	test.AssertRecentTime(t, issue.UpdatedAt)
}

func TestGraphIssueEventHandler_OnAddUserFollower(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	issueId := neo4jt.CreateIssue(ctx, testDatabase.Driver, tenantName, neo4jentity.IssueEntity{})
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "Issue": 1, "TimelineEvent": 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{"FOLLOWED_BY": 0})

	// prepare event handler
	issueEventHandler := &IssueEventHandler{
		repositories: testDatabase.Repositories,
	}
	updatedAt := utils.Now().Add(time.Duration(-10) * time.Minute)
	issueAggregate := aggregate.NewIssueAggregateWithTenantAndID(tenantName, issueId)
	addUserFollowerEvent, err := event.NewIssueAddUserFollowerEvent(issueAggregate, userId, updatedAt)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = issueEventHandler.OnAddUserFollower(context.Background(), addUserFollowerEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 1, "User_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, issueId, "FOLLOWED_BY", userId)

	issueDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Issue_"+tenantName, issueId)
	require.Nil(t, err)
	require.NotNil(t, issueDbNode)
	issue := neo4jmapper.MapDbNodeToIssueEntity(issueDbNode)
	test.AssertRecentTime(t, issue.UpdatedAt)
}

func TestGraphIssueEventHandler_OnRemoveUserFollower(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	issueId := neo4jt.CreateIssue(ctx, testDatabase.Driver, tenantName, neo4jentity.IssueEntity{})
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
	neo4jt.LinkIssueFollowedBy(ctx, testDatabase.Driver, issueId, userId)
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "Issue": 1, "TimelineEvent": 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{"FOLLOWED_BY": 1})

	// prepare event handler
	issueEventHandler := &IssueEventHandler{
		repositories: testDatabase.Repositories,
	}
	updatedAt := utils.Now().Add(time.Duration(-1) * time.Hour)
	issueAggregate := aggregate.NewIssueAggregateWithTenantAndID(tenantName, issueId)
	removeUserFollowerEvent, err := event.NewIssueRemoveUserFollowerEvent(issueAggregate, userId, updatedAt)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = issueEventHandler.OnRemoveUserFollower(context.Background(), removeUserFollowerEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 1, "User_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{"FOLLOWED_BY": 0})

	issueDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Issue_"+tenantName, issueId)
	require.Nil(t, err)
	require.NotNil(t, issueDbNode)
	issue := neo4jmapper.MapDbNodeToIssueEntity(issueDbNode)
	test.AssertRecentTime(t, issue.UpdatedAt)
}
