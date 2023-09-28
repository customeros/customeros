package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	//"github.com/mrdulin/gqlgen-cnode/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	//"github.com/mrdulin/gqlgen-cnode/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/stretchr/testify/mock"

	//"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/resolver"
	srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"

	"testing"
)

var (
	ID        = "123"
	MimeType  = "text/plain"
	Extension = "txt"
	Name      = "readme.txt"
	Size      = 123
	AppSource = "test app"
)

func TestMutationResolver_AttachmentCreate_UT(t *testing.T) {

	t.Run("should create attachment correctly", func(t *testing.T) {
		testAttachmentService := new(MockedAttachmentService)
		mockedServices := srv.Services{
			AttachmentService: testAttachmentService,
		}
		resolvers := Resolver{Services: &mockedServices}
		ue := model.AttachmentInput{
			MimeType:  MimeType,
			Name:      Name,
			Size:      int64(Size),
			Extension: Extension,
			AppSource: AppSource}
		testAttachmentService.On("attachment_Create", mock.AnythingOfType("string")).Return(&ue)
		q := `
		 mutation {
		   attachment_Create(input: {mimeType: "text/plain", name: "readme.txt", size: 123, extension: "txt", appSource: "test app"}) {
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
		`
		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))
		rawResponse, err := c.RawPost(q)
		require.Nil(t, err)

		var attachmentStruct struct {
			Attachment_Create model.Attachment
		}

		err = decode.Decode(rawResponse.Data.(map[string]any), &attachmentStruct)
		require.Nil(t, err)
		require.NotNil(t, attachmentStruct)

		//testAttachmentService.AssertExpectations(t)
		attachment := attachmentStruct.Attachment_Create
		require.Equal(t, "text/plain", attachment.MimeType)
	})

}

//func TestQueryResolver_Attachment_UT(t *testing.T) {
//	t.Run("should query user correctly", func(t *testing.T) {
//		testUserService := new(mocks.MockedUserService)
//		resolvers := resolver.Resolver{UserService: testUserService}
//		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))
//		u := model.UserDetail{User: model.User{Loginname: &loginname, AvatarURL: &avatarURL}, Score: &score, CreateAt: &createAt}
//		testUserService.On("GetUserByLoginname", mock.AnythingOfType("string")).Return(&u)
//		var resp struct {
//			User struct {
//				Loginname, AvatarURL, CreateAt string
//				Score                          int
//			}
//		}
//		q := `
//      query GetUser($loginname: String!) {
//        user(loginname: $loginname) {
//          loginname
//          avatarUrl
//          createAt
//          score
//        }
//      }
//    `
//		c.MustPost(q, &resp, client.Var("loginname", "mrdulin"))
//		testUserService.AssertCalled(t, "GetUserByLoginname", "mrdulin")
//		require.Equal(t, "mrdulin", resp.User.Loginname)
//		require.Equal(t, "avatar.jpg", resp.User.AvatarURL)
//		require.Equal(t, 50, resp.User.Score)
//		require.Equal(t, "1900-01-01", resp.User.CreateAt)
//	})
//}
