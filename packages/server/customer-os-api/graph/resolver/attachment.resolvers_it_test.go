package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"log"
	"testing"
	"time"
)

func TestMutationResolver_AttachmentCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	rawResponse, err := c.RawPost(getQuery("attachment/create_attachment"),
		client.Var("mimeType", "text/plain"),
		client.Var("extension", "txt"),
		client.Var("name", "readme.txt"),
		client.Var("size", 123),
	)

	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("attachment: %v", rawResponse.Data)
	var attachmentCreate struct {
		Attachment_Create struct {
			ID        string `json:"id"`
			MimeType  string `json:"mimeType"`
			Extension string `json:"extension"`
			Name      string `json:"name"`
			Size      int64  `json:"size"`
			AppSource string `json:"appSource"`
		} `json:"attachment_Create"`
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &attachmentCreate)
	require.Nil(t, err)
	require.Equal(t, "text/plain", attachmentCreate.Attachment_Create.MimeType)
	require.Equal(t, "txt", attachmentCreate.Attachment_Create.Extension)
	require.Equal(t, "readme.txt", attachmentCreate.Attachment_Create.Name)
	require.Equal(t, int64(123), attachmentCreate.Attachment_Create.Size)
	require.Equal(t, "Oasis", attachmentCreate.Attachment_Create.AppSource)

	rawResponse, err = c.RawPost(getQuery("attachment/get_attachment"),
		client.Var("attachmentId", attachmentCreate.Attachment_Create.ID),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var attachmentGet struct {
		Attachment struct {
			ID        string `json:"id"`
			MimeType  string `json:"mimeType"`
			Extension string `json:"extension"`
			Name      string `json:"name"`
			Size      int64  `json:"size"`
			AppSource string `json:"appSource"`
		} `json:"attachment"`
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &attachmentGet)
	require.Nil(t, err)
	require.Equal(t, "text/plain", attachmentGet.Attachment.MimeType)
	require.Equal(t, "txt", attachmentGet.Attachment.Extension)
	require.Equal(t, "readme.txt", attachmentGet.Attachment.Name)
	require.Equal(t, int64(123), attachmentGet.Attachment.Size)
	require.Equal(t, "Oasis", attachmentGet.Attachment.AppSource)
}

func TestQueryResolver_Attachment(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := time.Now().UTC()

	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "CALL", "ACTIVE", "VOICE", now, false)

	analysis1 := neo4jt.CreateAnalysis(ctx, driver, tenantName, "This is a summary of the conversation", "text/plain", "SUMMARY", now)
	neo4jt.ActionDescribes(ctx, driver, tenantName, analysis1, interactionSession1, repository.INTERACTION_SESSION)

	rawResponse, err := c.RawPost(getQuery("analysis/get_analysis"),
		client.Var("analysisId", analysis1))
	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("response: %v", rawResponse.Data)

	var analysis struct {
		Analysis struct {
			ID           string              `json:"id"`
			ContentType  string              `json:"contentType"`
			Content      string              `json:"content"`
			AnalysisType string              `json:"analysisType"`
			AppSource    string              `json:"appSource"`
			Describes    []map[string]string `json:"describes"`
		} `json:"analysis"`
	}
	_, ok := rawResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("response is not a map")
	}

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &analysis)
	require.Nil(t, err)
	require.Equal(t, "text/plain", analysis.Analysis.ContentType)
	require.Equal(t, "This is a summary of the conversation", analysis.Analysis.Content)
	require.Equal(t, "test", analysis.Analysis.AppSource)

	require.Len(t, analysis.Analysis.Describes, 1)
	log.Printf("Describe: %v", analysis.Analysis.Describes[0])
	require.Equal(t, "InteractionSession", analysis.Analysis.Describes[0]["__typename"])
	require.Equal(t, interactionSession1, analysis.Analysis.Describes[0]["id"])
	require.Equal(t, "mySessionIdentifier", analysis.Analysis.Describes[0]["sessionIdentifier"])
	require.Equal(t, "session1", analysis.Analysis.Describes[0]["name"])

}
