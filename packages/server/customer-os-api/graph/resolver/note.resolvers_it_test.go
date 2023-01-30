package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_NoteMergeToContact(t *testing.T) {
	defer tearDownTestCase()(t)

	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_note"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var note struct {
		Note_MergeToContact model.Note
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &note)
	require.Nil(t, err)

	createdNote := note.Note_MergeToContact

	require.NotNil(t, createdNote.ID)
	require.NotNil(t, createdNote.CreatedAt)
	require.Equal(t, "Note content", createdNote.HTML)
	require.Equal(t, model.DataSourceOpenline, createdNote.Source)
	require.Nil(t, createdNote.CreatedBy)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Note_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "NOTED"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Note", "Note_" + tenantName})
}

func TestMutationResolver_NoteDelete(t *testing.T) {
	defer tearDownTestCase()(t)

	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	noteId := neo4jt.CreateNoteForContact(driver, tenantName, contactId, "Note content")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Note_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "NOTED"))

	rawResponse, err := c.RawPost(getQuery("delete_note"),
		client.Var("noteId", noteId))
	assertRawResponseSuccess(t, rawResponse, err)

	var result struct {
		Note_Delete model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &result)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, true, result.Note_Delete.Result)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Note"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Note_"+tenantName))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "NOTED"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact"})
}
