package resolver

import (
	//"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/resolver"
	srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"

	"testing"
)

var (
	loginname = "mrdulin"
	avatarURL = "avatar.jpg"
	score     = 50
	createAt  = "1900-01-01"
)

func TestMutationResolver_AttachmentCreate_UT(t *testing.T) {

	t.Run("should create attachment correctly", func(t *testing.T) {
		testAttachmentService := new(MockedAttachmentService)
		mockedServices := srv.Services{
			AttachmentService: testAttachmentService,
		}
		resolvers := Resolver{Services: &mockedServices}
		//	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers})))
		//	ue := model.UserEntity{ID: "123", User: model.User{Loginname: &loginname, AvatarURL: &avatarURL}}
		//	testAttachmentService.On("ValidateAccessToken", mock.AnythingOfType("string")).Return(&ue)
		//	var resp struct {
		//		ValidateAccessToken struct{ ID, Loginname, AvatarUrl string }
		//	}
		//	q := `
		//  mutation {
		//    validateAccessToken(accesstoken: "abc") {
		//      id,
		//      loginname,
		//      avatarUrl
		//    }
		//  }
		//`
		//	c.MustPost(q, &resp)
		//	testAttachmentService.AssertExpectations(t)
		//require.Equal(t, "123", resp.ValidateAccessToken.ID)
		//require.Equal(t, "mrdulin", resp.ValidateAccessToken.Loginname)
		//require.Equal(t, "avatar.jpg", resp.ValidateAccessToken.AvatarUrl)
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
