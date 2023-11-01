package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_Issue(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	issueId := neo4jt.CreateIssue(ctx, driver, tenantName, entity.IssueEntity{
		Subject:     "testSubject",
		Status:      "testStatus",
		Priority:    "testPriority",
		Description: "testDescription",
		CreatedAt:   utils.Now(),
	})

	tagId1 := neo4jt.CreateTag(ctx, driver, tenantName, "critical")
	tagId2 := neo4jt.CreateTag(ctx, driver, tenantName, "issue-tag")

	neo4jt.CreateHubspotExternalSystem(ctx, driver, tenantName)
	syncDate := utils.Now()
	neo4jt.LinkWithHubspotExternalSystem(ctx, driver, issueId, "1234567890", utils.StringPtr("www.external.com"), utils.StringPtr("ticket"), syncDate)

	neo4jt.TagIssue(ctx, driver, issueId, tagId1)
	neo4jt.TagIssue(ctx, driver, issueId, tagId2)

	channel := "EMAIL"
	interactionEventId := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", &channel, utils.Now())
	neo4jt.InteractionEventPartOfIssue(ctx, driver, interactionEventId, issueId)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Issue":            1,
		"Tag":              2,
		"InteractionEvent": 1,
		"ExternalSystem":   1,
	})
	assertRelationship(ctx, t, driver, issueId, "IS_LINKED_WITH", string(entity.Hubspot))

	rawResponse, err := c.RawPost(getQuery("issue/get_issue"),
		client.Var("issueId", issueId))
	assertRawResponseSuccess(t, rawResponse, err)

	var issueStruct struct {
		Issue model.Issue
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &issueStruct)
	issue := issueStruct.Issue

	require.Nil(t, err)
	require.NotNil(t, issue)
	require.NotNil(t, issue.CreatedAt)
	require.Equal(t, "testSubject", *issue.Subject)
	require.Equal(t, "testStatus", issue.Status)
	require.Equal(t, "testDescription", *issue.Description)
	require.Equal(t, "testPriority", *issue.Priority)
	require.Equal(t, 2, len(issue.Tags))
	require.ElementsMatch(t, []string{tagId1, tagId2}, []string{issue.Tags[0].ID, issue.Tags[1].ID})
	require.ElementsMatch(t, []string{"critical", "issue-tag"}, []string{issue.Tags[0].Name, issue.Tags[1].Name})
	require.Equal(t, 1, len(issue.InteractionEvents))
	require.Equal(t, interactionEventId, issue.InteractionEvents[0].ID)
	require.Equal(t, 1, len(issue.ExternalLinks))
	require.Equal(t, "1234567890", *issue.ExternalLinks[0].ExternalID)
	require.Equal(t, "www.external.com", *issue.ExternalLinks[0].ExternalURL)
	require.Equal(t, "ticket", *issue.ExternalLinks[0].ExternalSource)
	require.Equal(t, syncDate, *issue.ExternalLinks[0].SyncDate)
}

func TestQueryResolver_Issue_WithParticipants(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	issueId := neo4jt.CreateIssue(ctx, driver, tenantName, entity.IssueEntity{
		Subject:     "testSubject",
		Status:      "testStatus",
		Priority:    "testPriority",
		Description: "testDescription",
		CreatedAt:   utils.Now(),
	})

	userId := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{})
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	contactId := neo4jt.CreateContact(ctx, driver, tenantName, entity.ContactEntity{})

	neo4jt.IssueSubmittedBy(ctx, driver, issueId, userId)
	neo4jt.IssueReportedBy(ctx, driver, issueId, orgId)
	neo4jt.IssueAssignedTo(ctx, driver, issueId, userId)
	neo4jt.IssueFollowedBy(ctx, driver, issueId, userId)
	neo4jt.IssueFollowedBy(ctx, driver, issueId, contactId)
	neo4jt.IssueFollowedBy(ctx, driver, issueId, orgId)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Issue":        1,
		"User":         1,
		"Organization": 1,
		"Contact":      1,
	})

	rawResponse := callGraphQL(t, "issue/get_issue_with_participants", map[string]interface{}{"issueId": issueId})

	issue := rawResponse.Data.(map[string]interface{})["issue"]
	submittedBy := issue.(map[string]interface{})["submittedBy"].(map[string]interface{})
	require.Equal(t, userId, submittedBy["userParticipant"].(map[string]interface{})["id"])
	require.Equal(t, "UserParticipant", submittedBy["__typename"])

	reportedBy := issue.(map[string]interface{})["reportedBy"].(map[string]interface{})
	require.Equal(t, orgId, reportedBy["organizationParticipant"].(map[string]interface{})["id"])
	require.Equal(t, "OrganizationParticipant", reportedBy["__typename"])

	assignedTo := issue.(map[string]interface{})["assignedTo"].([]interface{})
	require.Equal(t, 1, len(assignedTo))
	require.Equal(t, userId, assignedTo[0].(map[string]interface{})["userParticipant"].(map[string]interface{})["id"])
	require.Equal(t, "UserParticipant", assignedTo[0].(map[string]interface{})["__typename"])

	followedBy := issue.(map[string]interface{})["followedBy"].([]interface{})
	require.Equal(t, 3, len(followedBy))
	require.ElementsMatch(t, []string{userId, contactId, orgId}, []string{
		extractParticipantId(followedBy[0].(map[string]interface{})),
		extractParticipantId(followedBy[1].(map[string]interface{})),
		extractParticipantId(followedBy[2].(map[string]interface{})),
	})
}

func extractParticipantId(participant map[string]interface{}) string {
	switch participant["__typename"] {
	case "UserParticipant":
		return participant["userParticipant"].(map[string]interface{})["id"].(string)
	case "OrganizationParticipant":
		return participant["organizationParticipant"].(map[string]interface{})["id"].(string)
	case "ContactParticipant":
		return participant["contactParticipant"].(map[string]interface{})["id"].(string)
	default:
		return ""
	}
}
