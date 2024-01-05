package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestMutationResolver_AttachmentCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
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
