package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_NoteCreateForContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("note/create_note_for_contact"),
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
	require.Equal(t, "Note content", *createdNote.Content)
	require.Equal(t, "text/markdown", *createdNote.ContentType)
	require.Equal(t, model.DataSourceOpenline, createdNote.Source)
	require.Equal(t, model.DataSourceOpenline, createdNote.SourceOfTruth)
	require.Equal(t, constants.AppSourceCustomerOsApi, createdNote.AppSource)
	require.Equal(t, testUserId, createdNote.CreatedBy.ID)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Note_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "NOTED"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "CREATED"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "User", "User_" + tenantName,
		"Note", "Note_" + tenantName, "TimelineEvent", "TimelineEvent_" + tenantName})
}

func TestMutationResolver_NoteCreateForOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org")

	rawResponse, err := c.RawPost(getQuery("note/create_note_for_organization"),
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
	require.Equal(t, "Note content", *createdNote.Content)
	require.Equal(t, "text/html", *createdNote.ContentType)
	require.Equal(t, model.DataSourceOpenline, createdNote.Source)
	require.Equal(t, model.DataSourceOpenline, createdNote.SourceOfTruth)
	require.Equal(t, "test", createdNote.AppSource)
	require.Equal(t, testUserId, createdNote.CreatedBy.ID)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Note_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "NOTED"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Organization", "Organization_" + tenantName, "User", "User_" + tenantName,
		"Note", "Note_" + tenantName, "TimelineEvent", "TimelineEvent_" + tenantName})
}

func TestMutationResolver_AddAttachmentToNote(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	noteId := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "Note content", "text/plain", utils.Now())
	attachmentId := neo4jt.CreateAttachment(ctx, driver, tenantName, entity.AttachmentEntity{
		Id:            "",
		MimeType:      "text/plain",
		Name:          "readme.txt",
		Extension:     "txt",
		Size:          123,
		Source:        "",
		SourceOfTruth: "",
		AppSource:     "",
	})

	rawResponse, err := c.RawPost(getQuery("note/add_attachment_to_note"),
		client.Var("noteId", noteId),
		client.Var("attachmentId", attachmentId))
	assertRawResponseSuccess(t, rawResponse, err)

	var note struct {
		Note_LinkAttachment model.Note
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &note)
	require.Nil(t, err)

	require.NotNil(t, note.Note_LinkAttachment.ID)
	require.Len(t, note.Note_LinkAttachment.Includes, 1)
	require.Equal(t, note.Note_LinkAttachment.Includes[0].ID, attachmentId)

}

func TestMutationResolver_RemoveAttachmentFromNote(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	noteId := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "Note content", "text/plain", utils.Now())
	attachmentId := neo4jt.CreateAttachment(ctx, driver, tenantName, entity.AttachmentEntity{
		Id:            "",
		MimeType:      "text/plain",
		Name:          "readme.txt",
		Extension:     "txt",
		Size:          123,
		Source:        "",
		SourceOfTruth: "",
		AppSource:     "",
	})

	rawResponse, err := c.RawPost(getQuery("note/add_attachment_to_note"),
		client.Var("noteId", noteId),
		client.Var("attachmentId", attachmentId))
	assertRawResponseSuccess(t, rawResponse, err)

	var note struct {
		Note_LinkAttachment model.Note
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &note)
	require.Nil(t, err)

	require.NotNil(t, note.Note_LinkAttachment.ID)
	require.Len(t, note.Note_LinkAttachment.Includes, 1)
	require.Equal(t, note.Note_LinkAttachment.Includes[0].ID, attachmentId)

	rawRemoveResponse, err := c.RawPost(getQuery("note/remove_attachment_from_note"),
		client.Var("noteId", noteId),
		client.Var("attachmentId", attachmentId))
	assertRawResponseSuccess(t, rawResponse, err)

	var note_unlink struct {
		Note_UnlinkAttachment model.Note
	}

	err = decode.Decode(rawRemoveResponse.Data.(map[string]any), &note_unlink)
	require.Nil(t, err)

	require.NotNil(t, note_unlink.Note_UnlinkAttachment.ID)
	require.Len(t, note_unlink.Note_UnlinkAttachment.Includes, 0)

}

func TestMutationResolver_NoteUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	noteId := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "Note content", "text/plain", utils.Now())

	rawResponse, err := c.RawPost(getQuery("note/update_note"),
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
	require.Equal(t, "updated content", *updatedNote.Content)
	require.Equal(t, "text/markdown", *updatedNote.ContentType)
	require.Equal(t, model.DataSourceOpenline, updatedNote.SourceOfTruth)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Note_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "NOTED"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName,
		"Note", "Note_" + tenantName, "TimelineEvent", "TimelineEvent_" + tenantName})
}

func TestMutationResolver_NoteDelete(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	noteId := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "Note content", "text/plain", utils.Now())

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Note_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "NOTED"))

	rawResponse, err := c.RawPost(getQuery("note/delete_note"),
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
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "Note_"+tenantName))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(ctx, driver, "NOTED"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName})
}

func TestQueryResolver_GetNote_WithNotedEntities(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org")
	noteId := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "Note content", "text/plain", utils.Now())
	neo4jt.LinkNoteWithOrganization(ctx, driver, noteId, organizationId)

	rawResponse, err := c.RawPost(getQuery("note/get_note_with_noted_entities_via_organization_query"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	organization := rawResponse.Data.(map[string]interface{})["organization"]
	require.Equal(t, organizationId, organization.(map[string]interface{})["id"])

	note := organization.(map[string]interface{})["notes"].(map[string]interface{})["content"].([]interface{})[0]

	var notedContact, notedOrganization interface{}

	require.Equal(t, noteId, note.(map[string]interface{})["id"])
	require.NotNil(t, note.(map[string]interface{})["createdAt"])
	require.NotNil(t, note.(map[string]interface{})["updatedAt"])

	if note.(map[string]interface{})["noted"].([]interface{})[0].(map[string]interface{})["__typename"] == "Contact" {
		notedContact = note.(map[string]interface{})["noted"].([]interface{})[0]
	} else if note.(map[string]interface{})["noted"].([]interface{})[1].(map[string]interface{})["__typename"] == "Contact" {
		notedContact = note.(map[string]interface{})["noted"].([]interface{})[1]
	}
	if note.(map[string]interface{})["noted"].([]interface{})[0].(map[string]interface{})["__typename"] == "Organization" {
		notedOrganization = note.(map[string]interface{})["noted"].([]interface{})[0]
	} else if note.(map[string]interface{})["noted"].([]interface{})[1].(map[string]interface{})["__typename"] == "Organization" {
		notedOrganization = note.(map[string]interface{})["noted"].([]interface{})[1]
	}

	require.Equal(t, "Contact", notedContact.(map[string]interface{})["__typename"])
	require.Equal(t, "Organization", notedOrganization.(map[string]interface{})["__typename"])

	require.Equal(t, contactId, notedContact.(map[string]interface{})["id"])
	require.Equal(t, "first", notedContact.(map[string]interface{})["firstName"])
	require.Equal(t, "last", notedContact.(map[string]interface{})["lastName"])
	require.Equal(t, organizationId, notedOrganization.(map[string]interface{})["id"])
	require.Equal(t, "test org", notedOrganization.(map[string]interface{})["name"])

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Note_"+tenantName))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "NOTED"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName,
		"Organization", "Organization_" + tenantName, "Note", "Note_" + tenantName,
		"TimelineEvent", "TimelineEvent_" + tenantName})
}
