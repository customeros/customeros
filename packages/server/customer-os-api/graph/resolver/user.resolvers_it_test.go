package resolver

import (
	"context"
	"encoding/json"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"log"
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

func TestMutationResolver_UserCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, "other")

	rawResponse, err := cAdminWithTenant.RawPost(getQuery("user/create_user"))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		User_Create model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)

	createdUser := user.User_Create
	require.NotNil(t, createdUser.ID)
	require.NotNil(t, createdUser.CreatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdUser.CreatedAt)
	require.NotNil(t, createdUser.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdUser.UpdatedAt)
	require.Equal(t, createdUser.UpdatedAt, createdUser.CreatedAt)
	require.Equal(t, "first", createdUser.FirstName)
	require.Equal(t, "last", createdUser.LastName)
	require.Equal(t, "user@openline.ai", *createdUser.Emails[0].Email)
	require.Equal(t, "user@openline.ai", *createdUser.Emails[0].RawEmail)
	require.Equal(t, false, *createdUser.Emails[0].Validated)
	require.NotNil(t, createdUser.Person)
	require.Equal(t, "user@openline.ai", createdUser.Person.Email)
	require.Equal(t, "dummy_provider", createdUser.Person.Provider)
	require.Equal(t, "dummy", createdUser.Person.AppSource)

	require.Equal(t, model.DataSourceOpenline, createdUser.Source)
	require.Equal(t, model.DataSourceOpenline, createdUser.SourceOfTruth)
	require.Equal(t, "dummy", createdUser.AppSource)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Person"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "User", "User_" + tenantName, "Email", "Email_" + tenantName, "Person"})
}

func TestMutationResolver_UserCreateAccessControlled(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, "other")

	rawResponse, err := c.RawPost(getQuery("user/create_user"))
	require.Nil(t, err)
	require.NotNil(t, rawResponse.Errors)

}

func TestMutationResolver_UserUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	userId := neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
	})

	rawResponse, err := c.RawPost(getQuery("user/update_user"),
		client.Var("userId", userId))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		User_Update model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)

	updatedUser := user.User_Update
	require.NotNil(t, updatedUser.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), updatedUser.UpdatedAt)
	require.Equal(t, userId, updatedUser.ID)
	require.Equal(t, "firstUpdated", updatedUser.FirstName)
	require.Equal(t, "lastUpdated", updatedUser.LastName)
	require.Equal(t, model.DataSourceOpenline, updatedUser.Source)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "USER_BELONGS_TO_TENANT"))

	// Users can't update other users
	rawResponse2, err := c.RawPost(getQuery("user/update_user"),
		client.Var("userId", userId1))

	bytes, _ := json.Marshal(rawResponse2)
	log.Print("JSON RESPONSE:" + string(bytes))
	require.Nil(t, rawResponse2.Data)
}

func TestMutationResolver_UserUpdateByOwner(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName:     "first",
		LastName:      "last",
		Source:        "openline",
		SourceOfTruth: "openline",
	})

	rawResponse, err := cOwner.RawPost(getQuery("user/update_user"),
		client.Var("userId", userId1))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		User_Update model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)

	updatedUser := user.User_Update
	require.NotNil(t, updatedUser.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), updatedUser.UpdatedAt)
	require.Equal(t, userId1, updatedUser.ID)
	require.Equal(t, "firstUpdated", updatedUser.FirstName)
	require.Equal(t, "lastUpdated", updatedUser.LastName)
	require.Equal(t, model.DataSourceOpenline, updatedUser.Source)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "USER_BELONGS_TO_TENANT"))

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

func TestMutationResolver_AddRole(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{"USER"},
	})

	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId1, "test@openline.com", true, "MAIN")

	rawResponse, err := cOwner.RawPost(getQuery("user/add_role"),
		client.Var("userId", userId1),
		client.Var("role", model.RoleOwner.String()))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		User_AddRole model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, "first", user.User_AddRole.FirstName)
	require.Equal(t, "last", user.User_AddRole.LastName)
	require.NotNil(t, user.User_AddRole.CreatedAt)
	require.Contains(t, user.User_AddRole.Roles, model.RoleOwner)
	require.Contains(t, user.User_AddRole.Roles, model.RoleUser)

	// Owners cannot give PlatformOwner role
	rawResponse2, err := cOwner.RawPost(getQuery("user/add_role"),
		client.Var("userId", userId1),
		client.Var("role", model.RoleCustomerOsPlatformOwner.String()))
	require.Nil(t, rawResponse2.Data)
}

func TestMutationResolver_AddRoleInTenant(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, "otherTenant")
	userId1 := neo4jt.CreateUser(ctx, driver, "otherTenant", entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{"USER"},
	})

	neo4jt.AddEmailTo(ctx, driver, entity.USER, "otherTenant", userId1, "test@openline.com", true, "MAIN")

	rawResponse, err := cCustomerOsPlatformOwner.RawPost(getQuery("user/add_role_in_tenant"),
		client.Var("userId", userId1),
		client.Var("role", model.RoleOwner.String()),
		client.Var("tenant", "otherTenant"))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		User_AddRoleInTenant model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, "first", user.User_AddRoleInTenant.FirstName)
	require.Equal(t, "last", user.User_AddRoleInTenant.LastName)
	require.NotNil(t, user.User_AddRoleInTenant.CreatedAt)
	require.Contains(t, user.User_AddRoleInTenant.Roles, model.RoleOwner)
	require.Contains(t, user.User_AddRoleInTenant.Roles, model.RoleUser)

	// Owners cannot use this method
	rawResponse2, err := cOwner.RawPost(getQuery("user/add_role_in_tenant"),
		client.Var("userId", userId1),
		client.Var("role", model.RoleUser.String()),
		client.Var("tenant", "otherTenant"))
	require.Nil(t, rawResponse2.Data)
}

func TestMutationResolver_RemoveRole(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{model.RoleOwner.String(), model.RoleUser.String(), model.RoleCustomerOsPlatformOwner.String()},
	})

	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId1, "test@openline.com", true, "MAIN")

	rawResponse, err := cOwner.RawPost(getQuery("user/remove_role"),
		client.Var("userId", userId1),
		client.Var("role", model.RoleOwner.String()))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		User_RemoveRole model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, "first", user.User_RemoveRole.FirstName)
	require.Equal(t, "last", user.User_RemoveRole.LastName)
	require.NotNil(t, user.User_RemoveRole.CreatedAt)
	require.NotContains(t, user.User_RemoveRole.Roles, model.RoleOwner)
	require.Contains(t, user.User_RemoveRole.Roles, model.RoleUser)

	// Owners cannot remove PlatformOwner role
	rawResponse2, err := cOwner.RawPost(getQuery("user/remove_role"),
		client.Var("userId", userId1),
		client.Var("role", model.RoleCustomerOsPlatformOwner.String()))
	require.Nil(t, rawResponse2.Data)
}

func TestMutationResolver_RemoveRoleInTenat(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, "otherTenant")
	userId1 := neo4jt.CreateUser(ctx, driver, "otherTenant", entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{model.RoleOwner.String(), model.RoleUser.String(), model.RoleCustomerOsPlatformOwner.String()},
	})

	neo4jt.AddEmailTo(ctx, driver, entity.USER, "otherTenant", userId1, "test@openline.com", true, "MAIN")

	rawResponse, err := cCustomerOsPlatformOwner.RawPost(getQuery("user/remove_role_in_tenant"),
		client.Var("userId", userId1),
		client.Var("role", model.RoleOwner.String()),
		client.Var("tenant", "otherTenant"))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		User_RemoveRoleInTenant model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, "first", user.User_RemoveRoleInTenant.FirstName)
	require.Equal(t, "last", user.User_RemoveRoleInTenant.LastName)
	require.NotNil(t, user.User_RemoveRoleInTenant.CreatedAt)
	require.NotContains(t, user.User_RemoveRoleInTenant.Roles, model.RoleOwner)
	require.Contains(t, user.User_RemoveRoleInTenant.Roles, model.RoleUser)

	// Owners cannot call cross-tenant methods
	rawResponse2, err := cOwner.RawPost(getQuery("user/remove_role_in_tenant"),
		client.Var("userId", userId1),
		client.Var("role", model.RoleUser.String()),
		client.Var("tenant", "otherTenant"))
	require.Nil(t, rawResponse2.Data)
}

func TestQueryResolver_Users_FilteredAndSorted(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

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

	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "User"))

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

	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "user",
	})
	neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "second",
		LastName:  "user",
	})

	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId1, "test@openline.com", true, "MAIN")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))

	rawResponse, err := c.RawPost(getQuery("user/get_user_by_id"),
		client.Var("userId", userId1))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		User model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, userId1, user.User.ID)
	require.Equal(t, "first", user.User.FirstName)
	require.Equal(t, "user", user.User.LastName)
	require.Equal(t, "test@openline.com", *user.User.Emails[0].Email)
}

func TestQueryResolver_User_WithConversations(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	user1 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	user2 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	contact1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contact2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contact3 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	conv1_1 := neo4jt.CreateConversation(ctx, driver, tenantName, user1, contact1, "subject 1", utils.Now())
	conv1_2 := neo4jt.CreateConversation(ctx, driver, tenantName, user1, contact2, "subject 2", utils.Now())
	conv2_1 := neo4jt.CreateConversation(ctx, driver, tenantName, user2, contact1, "subject 3", utils.Now())
	conv2_3 := neo4jt.CreateConversation(ctx, driver, tenantName, user2, contact3, "subject 4", utils.Now())

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "Conversation"))

	rawResponse, err := c.RawPost(getQuery("user/get_user_with_conversations"),
		client.Var("userId", user1))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		User model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, user1, user.User.ID)
	require.Equal(t, 1, user.User.Conversations.TotalPages)
	require.Equal(t, int64(2), user.User.Conversations.TotalElements)
	require.Equal(t, 2, len(user.User.Conversations.Content))
	conversations := user.User.Conversations.Content
	require.ElementsMatch(t, []string{conv1_1, conv1_2}, []string{conversations[0].ID, conversations[1].ID})
	require.ElementsMatch(t, []string{contact1, contact2}, []string{conversations[0].Contacts[0].ID, conversations[1].Contacts[0].ID})
	require.Equal(t, user1, conversations[0].Users[0].ID)
	require.Equal(t, user1, conversations[1].Users[0].ID)

	require.NotNil(t, conv2_1)
	require.NotNil(t, conv2_3)
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
