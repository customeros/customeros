package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_UserCreate(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	createTenant(driver, "other")

	rawResponse, err := c.RawPost(getQuery("create_user"))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		UserCreate model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, "first", user.UserCreate.FirstName)
	require.Equal(t, "last", user.UserCreate.LastName)
	require.Equal(t, "user@openline.ai", user.UserCreate.Email)
	require.NotNil(t, user.UserCreate.CreatedAt)
	require.NotNil(t, user.UserCreate.ID)
}

func TestQueryResolver_Users(t *testing.T) {
	defer setupTestCase()(t)
	otherTenant := "other"
	createTenant(driver, tenantName)
	createTenant(driver, otherTenant)
	createUser(driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Email:     "test@openline.ai",
	})
	createUser(driver, otherTenant, entity.UserEntity{
		FirstName: "otherFirst",
		LastName:  "otherLast",
		Email:     "otherEmail",
	})

	rawResponse, err := c.RawPost(getQuery("get_users"))
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
	require.Equal(t, "test@openline.ai", users.Users.Content[0].Email)
	require.NotNil(t, users.Users.Content[0].CreatedAt)
}

func TestQueryResolver_Users_FilteredAndSorted(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)

	createUser(driver, tenantName, entity.UserEntity{
		FirstName: "first_f",
		LastName:  "first_l",
		Email:     "user1@openline.ai",
	})
	createUser(driver, tenantName, entity.UserEntity{
		FirstName: "second_f",
		LastName:  "second_l",
		Email:     "user2@openline.ai",
	})
	createUser(driver, tenantName, entity.UserEntity{
		FirstName: "third_f",
		LastName:  "third_l",
		Email:     "user3@openline.ai",
	})
	createUser(driver, tenantName, entity.UserEntity{
		FirstName: "fourth_f",
		LastName:  "fourth_l",
		Email:     "user4@openline.ai",
	})

	require.Equal(t, 4, getCountOfNodes(driver, "User"))

	rawResponse, err := c.RawPost(getQuery("get_users_filtered_and_sorted"))
	assertRawResponseSuccess(t, rawResponse, err)

	var users struct {
		Users model.UserPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &users)
	require.Nil(t, err)
	require.NotNil(t, users)
	require.Equal(t, 1, users.Users.TotalPages)
	require.Equal(t, int64(3), users.Users.TotalElements)
	require.Equal(t, 3, len(users.Users.Content))
	require.Equal(t, "third_f", users.Users.Content[0].FirstName)
	require.Equal(t, "second_f", users.Users.Content[1].FirstName)
	require.Equal(t, "first_f", users.Users.Content[2].FirstName)
}

func TestQueryResolver_User(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)

	userId1 := createUser(driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "user",
		Email:     "user1@openline.ai",
	})
	createUser(driver, tenantName, entity.UserEntity{
		FirstName: "second",
		LastName:  "user",
		Email:     "user2@openline.ai",
	})

	require.Equal(t, 2, getCountOfNodes(driver, "User"))

	rawResponse, err := c.RawPost(getQuery("get_user_by_id"),
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
	require.Equal(t, "user1@openline.ai", user.User.Email)
}

func TestQueryResolver_User_WithConversations(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)

	user1 := createDefaultUser(driver, tenantName)
	user2 := createDefaultUser(driver, tenantName)
	contact1 := createDefaultContact(driver, tenantName)
	contact2 := createDefaultContact(driver, tenantName)
	contact3 := createDefaultContact(driver, tenantName)

	conv1_1 := createConversation(driver, user1, contact1)
	conv1_2 := createConversation(driver, user1, contact2)
	conv2_1 := createConversation(driver, user2, contact1)
	conv2_3 := createConversation(driver, user2, contact3)

	require.Equal(t, 2, getCountOfNodes(driver, "User"))
	require.Equal(t, 3, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 4, getCountOfNodes(driver, "Conversation"))

	rawResponse, err := c.RawPost(getQuery("get_user_with_conversations"),
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
	require.ElementsMatch(t, []string{contact1, contact2}, []string{conversations[0].Contact.ID, conversations[1].Contact.ID})
	require.ElementsMatch(t, []string{contact1, contact2}, []string{conversations[0].ContactID, conversations[1].ContactID})
	require.Equal(t, user1, conversations[0].User.ID)
	require.Equal(t, user1, conversations[1].User.ID)
	require.Equal(t, user1, conversations[0].UserID)
	require.Equal(t, user1, conversations[1].UserID)

	require.NotNil(t, conv2_1)
	require.NotNil(t, conv2_3)
}
