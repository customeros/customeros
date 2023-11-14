package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	jobRoleProto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/job_role"
	userProto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_UserByEmail(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	otherTenant := "other"
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, otherTenant)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
	})
	userId2 := neo4jt.CreateUser(ctx, driver, otherTenant, entity.UserEntity{
		FirstName: "otherFirst",
		LastName:  "otherLast",
	})

	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId1, "test@openline.com", true, "MAIN")
	neo4jt.AddEmailTo(ctx, driver, entity.USER, otherTenant, userId2, "test@openline.com", true, "MAIN")

	rawResponse, err := c.RawPost(getQuery("user/get_user_by_email"),
		client.Var("email", "test@openline.com"))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		User_ByEmail model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, userId1, user.User_ByEmail.ID)
}

func TestQueryResolver_Users(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	otherTenant := "other"
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, otherTenant)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{"OWNER", "USER"},
	})
	userId2 := neo4jt.CreateUser(ctx, driver, otherTenant, entity.UserEntity{
		FirstName: "otherFirst",
		LastName:  "otherLast",
	})

	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId1, "test@openline.com", true, "MAIN")
	neo4jt.AddEmailTo(ctx, driver, entity.USER, otherTenant, userId2, "test@openline.com", true, "MAIN")

	rawResponse, err := c.RawPost(getQuery("user/get_users"))
	assertRawResponseSuccess(t, rawResponse, err)

	var users struct {
		Users model.UserPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &users)
	require.Nil(t, err)
	require.NotNil(t, users)
	require.Equal(t, 1, users.Users.TotalPages)
	require.Equal(t, int64(1), users.Users.TotalElements)
	require.Equal(t, "first", users.Users.Content[0].FirstName)
	require.Equal(t, "last", users.Users.Content[0].LastName)
	require.Equal(t, "test@openline.com", *users.Users.Content[0].Emails[0].Email)
	require.NotNil(t, users.Users.Content[0].CreatedAt)
	require.Contains(t, users.Users.Content[0].Roles, model.RoleOwner)
	require.Contains(t, users.Users.Content[0].Roles, model.RoleUser)
}

func TestQueryResolver_Users_FilteredAndSorted(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first_f_internal",
		LastName:  "first_l_internal",
		Internal:  true,
	})
	neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first_f",
		LastName:  "first_l",
	})
	neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "second_f",
		LastName:  "second_l",
	})
	neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "third_f",
		LastName:  "third_l",
	})
	neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "fourth_f",
		LastName:  "fourth_l",
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"User": 5})

	rawResponse, err := c.RawPost(getQuery("user/get_users_filtered_and_sorted"))
	assertRawResponseSuccess(t, rawResponse, err)

	var users struct {
		Users model.UserPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &users)
	require.Nil(t, err)
	require.NotNil(t, users)
	require.Equal(t, 1, users.Users.TotalPages)
	require.Equal(t, int64(2), users.Users.TotalElements)
	require.Equal(t, 2, len(users.Users.Content))
	require.Equal(t, "second_f", users.Users.Content[0].FirstName)
	require.Equal(t, "first_f", users.Users.Content[1].FirstName)
}

func TestQueryResolver_User(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	userId := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName:       "first",
		LastName:        "user",
		ProfilePhotoUrl: "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png",
	})
	neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "second",
		LastName:  "user",
	})

	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId, "test@openline.com", true, "MAIN")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))

	rawResponse := callGraphQL(t, "user/get_user_by_id", map[string]interface{}{"userId": userId})

	var user struct {
		User model.User
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, userId, user.User.ID)
	require.Equal(t, "first", user.User.FirstName)
	require.Equal(t, "user", user.User.LastName)
	require.Equal(t, "test@openline.com", *user.User.Emails[0].Email)
	require.Equal(t, "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png", *user.User.ProfilePhotoURL)
}

func TestQueryResolver_User_WithPhoneNumbers(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	phoneNumberId1 := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, userId, "+1111", true, "MAIN")
	phoneNumberId2 := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, userId, "+2222", false, "WORK")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))

	rawResponse, err := c.RawPost(getQuery("user/get_user_with_phone_numbers"),
		client.Var("userId", userId))
	assertRawResponseSuccess(t, rawResponse, err)

	var userStruct struct {
		User model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &userStruct)
	require.Nil(t, err)
	user := userStruct.User

	require.Equal(t, userId, user.ID)
	phoneNumbers := user.PhoneNumbers
	require.Equal(t, 2, len(phoneNumbers))
	var phoneNumber1, phoneNumber2 *model.PhoneNumber
	if phoneNumberId1 == phoneNumbers[0].ID {
		phoneNumber1 = phoneNumbers[0]
		phoneNumber2 = phoneNumbers[1]
	} else {
		phoneNumber1 = phoneNumbers[1]
		phoneNumber2 = phoneNumbers[0]
	}
	require.Equal(t, phoneNumberId1, phoneNumber1.ID)
	require.NotNil(t, phoneNumber1.CreatedAt)
	require.Equal(t, true, phoneNumber1.Primary)
	require.Equal(t, "+1111", *phoneNumber1.RawPhoneNumber)
	require.Equal(t, "+1111", *phoneNumber1.E164)
	require.Equal(t, model.PhoneNumberLabelMain, *phoneNumber1.Label)

	require.Equal(t, phoneNumberId2, phoneNumber2.ID)
	require.NotNil(t, phoneNumber2.CreatedAt)
	require.Equal(t, false, phoneNumber2.Primary)
	require.Equal(t, "+2222", *phoneNumber2.RawPhoneNumber)
	require.Equal(t, "+2222", *phoneNumber2.E164)
	require.Equal(t, model.PhoneNumberLabelWork, *phoneNumber2.Label)
}

func TestMutationResolver_AddJobRoleInTenant(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{"USER"},
	})
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "OPENLINE")
	calledCreateJobRole := false
	calledLinkJobRole := false

	jobRoleId, _ := uuid.NewUUID()
	jobRoleServiceCallbacks := events_platform.MockJobRoleServiceCallbacks{
		CreateJobRole: func(context context.Context, jobRole *jobRoleProto.CreateJobRoleGrpcRequest) (*jobRoleProto.JobRoleIdGrpcResponse, error) {
			require.Equal(t, "openline", jobRole.Tenant)
			require.Equal(t, "jobTitle", jobRole.JobTitle)
			require.Equal(t, "some description", *jobRole.Description)
			require.Equal(t, true, *jobRole.Primary)
			calledCreateJobRole = true
			return &jobRoleProto.JobRoleIdGrpcResponse{
				Id: jobRoleId.String(),
			}, nil
		},
	}
	userServiceCallbacks := events_platform.MockUserServiceCallbacks{
		LinkJobRoleToUser: func(context context.Context, request *userProto.LinkJobRoleToUserGrpcRequest) (*userProto.UserIdGrpcResponse, error) {
			require.Equal(t, "openline", request.Tenant)
			require.Equal(t, userId1, request.UserId)
			require.Equal(t, jobRoleId.String(), request.JobRoleId)
			calledLinkJobRole = true
			return &userProto.UserIdGrpcResponse{
				Id: userId1,
			}, nil
		},
	}
	events_platform.SetJobRoleCallbacks(&jobRoleServiceCallbacks)
	events_platform.SetUserCallbacks(&userServiceCallbacks)

	title := "jobTitle"
	isPrimary := true
	appSrc := "testApp"
	desr := "some description"
	rawResponse, err := cAdminWithTenant.RawPost(getQuery("user/customer_user_add_job_role"),
		client.Var("userId", userId1),
		client.Var("jobRoleInput", model.JobRoleInput{
			OrganizationID: &organizationId,
			JobTitle:       &title,
			Primary:        &isPrimary,
			AppSource:      &appSrc,
			Description:    &desr,
		}),
		client.Var("tenant", "otherTenant"))
	assertRawResponseSuccess(t, rawResponse, err)

	var jobRole struct {
		Customer_user_AddJobRole model.CustomerUser
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &jobRole)
	require.Nil(t, err)
	require.Equal(t, userId1, jobRole.Customer_user_AddJobRole.ID)
	require.True(t, calledCreateJobRole)
	require.True(t, calledLinkJobRole)
}

func TestMutationResolver_GetUserJobRoleInTenant(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{"USER"},
	})
	roleId := neo4jt.UserWorksAs(ctx, driver, userId1, "jobTitle", "some description", true)
	getRawResponse, err := c.RawPost(getQuery("user/get_users_with_job_roles"))
	assertRawResponseSuccess(t, getRawResponse, err)

	var users struct {
		Users model.UserPage
	}

	err = decode.Decode(getRawResponse.Data.(map[string]any), &users)
	require.Nil(t, err)
	require.Equal(t, roleId, users.Users.Content[0].JobRoles[0].ID)
	require.Equal(t, "some description", *users.Users.Content[0].JobRoles[0].Description)
	require.Equal(t, "jobTitle", *users.Users.Content[0].JobRoles[0].JobTitle)
}

func TestMutationResolver_GetUserCalendarInTenant(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{"USER"},
	})
	link := "https://cal.com/first-last"
	roleId := neo4jt.UserHasCalendar(ctx, driver, userId1, link, "CALCOM", true)
	getRawResponse, err := c.RawPost(getQuery("user/get_users_with_calendars"))
	assertRawResponseSuccess(t, getRawResponse, err)

	var users struct {
		Users model.UserPage
	}

	err = decode.Decode(getRawResponse.Data.(map[string]any), &users)
	require.Nil(t, err)
	require.Equal(t, roleId, users.Users.Content[0].Calendars[0].ID)
	require.Equal(t, link, *users.Users.Content[0].Calendars[0].Link)
	require.Equal(t, true, users.Users.Content[0].Calendars[0].Primary)
	require.Equal(t, "CALCOM", users.Users.Content[0].Calendars[0].CalType.String())
}
