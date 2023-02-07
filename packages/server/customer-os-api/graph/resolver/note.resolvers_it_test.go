package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_NoteCreateForContact(t *testing.T) {
	defer tearDownTestCase()(t)

	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateDefaultUserWithId(driver, tenantName, testUserId)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_note_for_contact"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var note struct {
		Note_CreateForContact model.Note
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &note)
	require.Nil(t, err)

	createdNote := note.Note_CreateForContact

	require.NotNil(t, createdNote.ID)
	require.NotNil(t, createdNote.CreatedAt)
	require.NotNil(t, createdNote.UpdatedAt)
	require.Equal(t, "Note content", createdNote.HTML)
	require.Equal(t, model.DataSourceOpenline, createdNote.Source)
	require.Equal(t, model.DataSourceOpenline, createdNote.SourceOfTruth)
	require.Equal(t, common.AppSourceCustomerOsApi, createdNote.AppSource)
	require.Equal(t, testUserId, createdNote.CreatedBy.ID)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Note_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "NOTED"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "CREATED"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "User", "Note", "Note_" + tenantName})
}

func TestMutationResolver_NoteCreateForOrganization(t *testing.T) {
	defer tearDownTestCase()(t)

	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateDefaultUserWithId(driver, tenantName, testUserId)
	organizationId := neo4jt.CreateOrganization(driver, tenantName, "test org")

	rawResponse, err := c.RawPost(getQuery("create_note_for_organization"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var note struct {
		Note_CreateForOrganization model.Note
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &note)
	require.Nil(t, err)

	createdNote := note.Note_CreateForOrganization

	require.NotNil(t, createdNote.ID)
	require.NotNil(t, createdNote.CreatedAt)
	require.NotNil(t, createdNote.UpdatedAt)
	require.Equal(t, "Note content", createdNote.HTML)
	require.Equal(t, model.DataSourceOpenline, createdNote.Source)
	require.Equal(t, model.DataSourceOpenline, createdNote.SourceOfTruth)
	require.Equal(t, "test", createdNote.AppSource)
	require.Equal(t, testUserId, createdNote.CreatedBy.ID)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Organization_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Note_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "NOTED"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Organization", "Organization_" + tenantName, "User", "Note", "Note_" + tenantName})
}

func TestMutationResolver_NoteUpdate(t *testing.T) {
	defer tearDownTestCase()(t)

	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	noteId := neo4jt.CreateNoteForContact(driver, tenantName, contactId, "Note content")

	rawResponse, err := c.RawPost(getQuery("update_note"),
		client.Var("noteId", noteId))
	assertRawResponseSuccess(t, rawResponse, err)

	var note struct {
		Note_Update model.Note
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &note)
	require.Nil(t, err)

	updatedNote := note.Note_Update

	require.NotNil(t, updatedNote.ID)
	require.NotNil(t, updatedNote.UpdatedAt)
	require.Equal(t, "updated content", updatedNote.HTML)
	require.Equal(t, model.DataSourceOpenline, updatedNote.SourceOfTruth)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Note_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "NOTED"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "Note", "Note_" + tenantName})
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
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact_"+tenantName))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Note"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Note_"+tenantName))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "NOTED"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName})
}
