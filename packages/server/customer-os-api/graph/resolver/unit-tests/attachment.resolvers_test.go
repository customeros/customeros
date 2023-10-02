package unit_tests

import (
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/resolver"
	srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

type AttachmentInput struct {
	mimeType  string
	name      string
	size      int64
	extension string
	appSource string
}

func TestMutationResolver_AttachmentCreate_FullyPopulated(t *testing.T) {
	attachmentEntity := AttachmentInput{
		mimeType:  "text/plain",
		name:      "readme.txt",
		size:      123,
		extension: "txt",
		appSource: "test app",
	}
	t.Run("should create fully populated attachment correctly", func(t *testing.T) {
		testAttachmentService := new(MockedAttachmentService)
		mockedServices := srv.Services{
			AttachmentService: testAttachmentService,
		}
		resolvers := resolver.Resolver{Services: &mockedServices}
		q := fmt.Sprintf(`
		 mutation {
		   attachment_Create(input: {mimeType: "%s", name: "%s", size: %d, extension: "%s", appSource: "%s"}) {
			 	id,
				createdAt,
				mimeType,
				name,
				size,
				extension,
				source,
				sourceOfTruth,
				appSource
		   }
		 }
		`, attachmentEntity.mimeType, attachmentEntity.name, attachmentEntity.size, attachmentEntity.extension, attachmentEntity.appSource)
		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))
		rawResponse, err := c.RawPost(q)
		require.Nil(t, err)

		var attachmentStruct struct {
			Attachment_Create model.Attachment
		}

		err = decode.Decode(rawResponse.Data.(map[string]any), &attachmentStruct)
		require.Nil(t, err)
		require.NotNil(t, attachmentStruct)

		attachment := attachmentStruct.Attachment_Create

		require.Equal(t, "", attachment.ID)
		require.Equal(t, attachmentEntity.mimeType, attachment.MimeType)
		require.Equal(t, attachmentEntity.name, attachment.Name)
		require.Equal(t, attachmentEntity.size, attachment.Size)
		require.Equal(t, attachmentEntity.extension, attachment.Extension)
		require.Equal(t, model.DataSource("NA"), attachment.Source)
		require.Equal(t, model.DataSource("NA"), attachment.SourceOfTruth)
		require.Equal(t, attachmentEntity.appSource, attachment.AppSource)
	})
}

func TestMutationResolver_AttachmentCreate_EmptyPopulated(t *testing.T) {
	attachmentEntity := AttachmentInput{}
	t.Run("should create empty populated attachment correctly", func(t *testing.T) {
		testAttachmentService := new(MockedAttachmentService)
		mockedServices := srv.Services{
			AttachmentService: testAttachmentService,
		}
		resolvers := resolver.Resolver{Services: &mockedServices}
		q := fmt.Sprintf(`
		 mutation {
		   attachment_Create(input: {mimeType: "%s", name: "%s", size: %d, extension: "%s", appSource: "%s"}) {
			 	id,
				createdAt,
				mimeType,
				name,
				size,
				extension,
				source,
				sourceOfTruth,
				appSource
		   }
		 }
		`, attachmentEntity.mimeType, attachmentEntity.name, attachmentEntity.size, attachmentEntity.extension, attachmentEntity.appSource)
		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))
		rawResponse, err := c.RawPost(q)
		require.Nil(t, err)

		var attachmentStruct struct {
			Attachment_Create model.Attachment
		}

		err = decode.Decode(rawResponse.Data.(map[string]any), &attachmentStruct)
		require.Nil(t, err)
		require.NotNil(t, attachmentStruct)

		attachment := attachmentStruct.Attachment_Create

		require.Equal(t, "", attachment.ID)
		require.Equal(t, attachmentEntity.mimeType, attachment.MimeType)
		require.Equal(t, attachmentEntity.name, attachment.Name)
		require.Equal(t, attachmentEntity.size, attachment.Size)
		require.Equal(t, attachmentEntity.extension, attachment.Extension)
		require.Equal(t, model.DataSource("NA"), attachment.Source)
		require.Equal(t, model.DataSource("NA"), attachment.SourceOfTruth)
		require.Equal(t, attachmentEntity.appSource, attachment.AppSource)
	})
}
