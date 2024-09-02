package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_AddAttachmentToNote(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	noteId := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "Note content", "text/plain", utils.Now())
	attachmentId := neo4jt.CreateAttachment(ctx, driver, tenantName, entity.AttachmentEntity{
		Id:            "",
		MimeType:      "text/plain",
		FileName:      "readme.txt",
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

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	noteId := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "Note content", "text/plain", utils.Now())
	attachmentId := neo4jt.CreateAttachment(ctx, driver, tenantName, entity.AttachmentEntity{
		Id:            "",
		MimeType:      "text/plain",
		FileName:      "readme.txt",
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

	neo4jtest.CreateTenant(ctx, driver, tenantName)
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
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Note_"+tenantName))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "NOTED"))

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName,
		"Note", "Note_" + tenantName, "TimelineEvent", "TimelineEvent_" + tenantName})
}

func TestMutationResolver_NoteDelete(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	noteId := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "Note content", "text/plain", utils.Now())

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Note_"+tenantName))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "NOTED"))

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
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 0, neo4jtest.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 0, neo4jtest.GetCountOfNodes(ctx, driver, "Note_"+tenantName))
	require.Equal(t, 0, neo4jtest.GetCountOfRelationships(ctx, driver, "NOTED"))

	// Check the labels on the nodes in the Neo4j database
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName})
}
