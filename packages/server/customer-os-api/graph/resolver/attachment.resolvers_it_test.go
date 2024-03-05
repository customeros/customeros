package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestMutationResolver_AttachmentCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	rawResponse, err := c.RawPost(getQuery("attachment/create_attachment"),
		client.Var("mimeType", "text/plain"),
		client.Var("fileName", "readme.txt"),
		client.Var("basePath", "/GLOBAL"),
		client.Var("cdnUrl", "test-url"),
		client.Var("size", 123),
	)

	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("attachment: %v", rawResponse.Data)
	var attachmentCreate struct {
		Attachment_Create model.Attachment
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &attachmentCreate)
	require.Nil(t, err)
	require.Equal(t, "test-url", attachmentCreate.Attachment_Create.CdnURL)
	require.Equal(t, "/GLOBAL", attachmentCreate.Attachment_Create.BasePath)
	require.Equal(t, "text/plain", attachmentCreate.Attachment_Create.MimeType)
	require.Equal(t, "readme.txt", attachmentCreate.Attachment_Create.FileName)
	require.Equal(t, int64(123), attachmentCreate.Attachment_Create.Size)
	require.Equal(t, "Oasis", attachmentCreate.Attachment_Create.AppSource)

	rawResponse, err = c.RawPost(getQuery("attachment/get_attachment"),
		client.Var("attachmentId", attachmentCreate.Attachment_Create.ID),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var attachmentGet struct {
		Attachment model.Attachment
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &attachmentGet)
	require.Nil(t, err)
	require.Equal(t, "/GLOBAL", attachmentGet.Attachment.BasePath)
	require.Equal(t, "text/plain", attachmentGet.Attachment.MimeType)
	require.Equal(t, "readme.txt", attachmentGet.Attachment.FileName)
	require.Equal(t, int64(123), attachmentGet.Attachment.Size)
	require.Equal(t, "Oasis", attachmentGet.Attachment.AppSource)
}
