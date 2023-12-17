package graph

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/model"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestGraphCommentEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	externalSystemId := "hubspot"
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, externalSystemId)
	commentedIssueId := neo4jt.CreateIssue(ctx, testDatabase.Driver, tenantName, entity.IssueEntity{})
	authorUserId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{})

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "Issue": 1, "ExternalSystem": 1, "Comment": 0, "TimelineEvent": 1})

	// prepare event handler
	commentEventHandler := &CommentEventHandler{
		repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	commentId := uuid.New().String()
	commentAggregate := aggregate.NewCommentAggregateWithTenantAndID(tenantName, commentId)
	createEvent, err := event.NewCommentCreateEvent(commentAggregate, model.CommentDataFields{
		Content:          "test content",
		ContentType:      "text",
		AuthorUserId:     utils.StringPtr(authorUserId),
		CommentedIssueId: utils.StringPtr(commentedIssueId),
	}, commonmodel.Source{
		Source:        constants.SourceOpenline,
		AppSource:     constants.AppSourceEventProcessingPlatform,
		SourceOfTruth: constants.SourceOpenline,
	}, commonmodel.ExternalSystem{
		ExternalSystemId: externalSystemId,
		ExternalId:       "123",
	}, now, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = commentEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 1, "User_" + tenantName: 1,
		"ExternalSystem": 1, "ExternalSystem_" + tenantName: 1,
		"Issue": 1, "Issue_" + tenantName: 1,
		"Comment": 1, "Comment_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, commentId, "COMMENTED", commentedIssueId)
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, commentId, "CREATED_BY", authorUserId)

	commentDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Comment_"+tenantName, commentId)
	require.Nil(t, err)
	require.NotNil(t, commentDbNode)

	// verify comment
	comment := graph_db.MapDbNodeToCommentEntity(*commentDbNode)
	require.Equal(t, commentId, comment.Id)
	require.Equal(t, "test content", comment.Content)
	require.Equal(t, "text", comment.ContentType)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), comment.Source)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), comment.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, comment.AppSource)
	require.Equal(t, now, comment.CreatedAt)
	require.Equal(t, now, comment.UpdatedAt)
}

func TestGraphCommentEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	commentId := neo4jt.CreateComment(ctx, testDatabase.Driver, tenantName, entity.CommentEntity{
		Content:     "test content",
		ContentType: "text",
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Comment": 1})

	// prepare event handler
	commentEventHandler := &CommentEventHandler{
		repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	commentAggregate := aggregate.NewCommentAggregateWithTenantAndID(tenantName, commentId)
	updateEvent, err := event.NewCommentUpdateEvent(commentAggregate, "test content update", "html", constants.SourceOpenline, commonmodel.ExternalSystem{}, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = commentEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Comment": 1, "Comment_" + tenantName: 1})

	commentDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Comment_"+tenantName, commentId)
	require.Nil(t, err)
	require.NotNil(t, commentDbNode)

	// verify comment
	comment := graph_db.MapDbNodeToCommentEntity(*commentDbNode)
	require.Equal(t, commentId, comment.Id)
	require.Equal(t, "test content update", comment.Content)
	require.Equal(t, "html", comment.ContentType)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), comment.SourceOfTruth)
	require.Equal(t, now, comment.UpdatedAt)
}

func TestGraphCommentEventHandler_OnUpdate_CurrentSourceOpenline_UpdateSourceNonOpenline_UpdateOnlyEmptyFields(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	commentId := neo4jt.CreateComment(ctx, testDatabase.Driver, tenantName, entity.CommentEntity{
		Content:       "original content",
		Source:        constants.SourceOpenline,
		SourceOfTruth: constants.SourceOpenline,
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Comment": 1})

	// prepare event handler
	commentEventHandler := &CommentEventHandler{
		repositories: testDatabase.Repositories,
	}
	now := utils.Now()
	commentAggregate := aggregate.NewCommentAggregateWithTenantAndID(tenantName, commentId)
	updateEvent, err := event.NewCommentUpdateEvent(commentAggregate, "test content updated", "type updated", "hubspot", commonmodel.ExternalSystem{}, now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = commentEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Comment": 1, "Comment_" + tenantName: 1})

	commentDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Comment_"+tenantName, commentId)
	require.Nil(t, err)
	require.NotNil(t, commentDbNode)

	// verify comment
	comment := graph_db.MapDbNodeToCommentEntity(*commentDbNode)
	require.Equal(t, commentId, comment.Id)
	require.Equal(t, "original content", comment.Content)
	require.Equal(t, "type updated", comment.ContentType)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), comment.SourceOfTruth)
	require.Equal(t, now, comment.UpdatedAt)
}
