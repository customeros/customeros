package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
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

	orgId := neo4jt.CreateOrganization(ctx, driver, tenantName, "testOrganization")

	noteId := neo4jt.CreateNoteForOrganization(ctx, driver, tenantName, orgId, "note", utils.Now())

	neo4jt.TagIssue(ctx, driver, issueId, tagId1)
	neo4jt.TagIssue(ctx, driver, issueId, tagId2)

	neo4jt.NoteMentionsTag(ctx, driver, noteId, tagId2)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Issue"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Tag"))

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
	require.Equal(t, 1, len(issue.MentionedByNotes))
	require.Equal(t, noteId, issue.MentionedByNotes[0].ID)
	require.Equal(t, "note", issue.MentionedByNotes[0].HTML)
}
