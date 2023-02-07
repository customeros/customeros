package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_UserByEmail(t *testing.T) {
	defer tearDownTestCase()(t)
	otherTenant := "other"
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, otherTenant)
	userId1 := neo4jt.CreateUser(driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
	})
	userId2 := neo4jt.CreateUser(driver, otherTenant, entity.UserEntity{
		FirstName: "otherFirst",
		LastName:  "otherLast",
	})

	neo4jt.AddEmailTo(driver, repository.USER, tenantName, userId1, "test@openline.com", true, "MAIN")
	neo4jt.AddEmailTo(driver, repository.USER, otherTenant, userId2, "test@openline.com", true, "MAIN")

	rawResponse, err := c.RawPost(getQuery("get_user_by_email"), client.Var("email", "test@openline.com"))
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
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, "other")

	rawResponse, err := c.RawPost(getQuery("create_user"))
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
	require.Equal(t, "first", createdUser.FirstName)
	require.Equal(t, "last", createdUser.LastName)
	require.Equal(t, "user@openline.ai", createdUser.Emails[0].Email)
	require.Equal(t, model.DataSourceOpenline, createdUser.Source)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "User_"+tenantName))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "User", "User_" + tenantName, "Email", "Email_" + tenantName})
}

func TestMutationResolver_UserUpdate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	userId := neo4jt.CreateDefaultUser(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("update_user"),
		client.Var("userId", userId))
	assertRawResponseSuccess(t, rawResponse, err)

	var user struct {
		User_Update model.User
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &user)
	require.Nil(t, err)
	require.NotNil(t, user)

	updatedUser := user.User_Update
	require.Equal(t, userId, updatedUser.ID)
	require.Equal(t, "firstUpdated", updatedUser.FirstName)
	require.Equal(t, "lastUpdated", updatedUser.LastName)
	require.Equal(t, model.DataSourceOpenline, updatedUser.Source)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "USER_BELONGS_TO_TENANT"))
}

func TestQueryResolver_Users(t *testing.T) {
	defer tearDownTestCase()(t)
	otherTenant := "other"
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, otherTenant)
	userId1 := neo4jt.CreateUser(driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
	})
	userId2 := neo4jt.CreateUser(driver, otherTenant, entity.UserEntity{
		FirstName: "otherFirst",
		LastName:  "otherLast",
	})

	neo4jt.AddEmailTo(driver, repository.USER, tenantName, userId1, "test@openline.com", true, "MAIN")
	neo4jt.AddEmailTo(driver, repository.USER, otherTenant, userId2, "test@openline.com", true, "MAIN")

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
	require.Equal(t, "test@openline.com", users.Users.Content[0].Emails[0].Email)
	require.NotNil(t, users.Users.Content[0].CreatedAt)
}

func TestQueryResolver_Users_FilteredAndSorted(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	neo4jt.CreateUser(driver, tenantName, entity.UserEntity{
		FirstName: "first_f",
		LastName:  "first_l",
	})
	neo4jt.CreateUser(driver, tenantName, entity.UserEntity{
		FirstName: "second_f",
		LastName:  "second_l",
	})
	neo4jt.CreateUser(driver, tenantName, entity.UserEntity{
		FirstName: "third_f",
		LastName:  "third_l",
	})
	neo4jt.CreateUser(driver, tenantName, entity.UserEntity{
		FirstName: "fourth_f",
		LastName:  "fourth_l",
	})

	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "User"))

	rawResponse, err := c.RawPost(getQuery("get_users_filtered_and_sorted"))
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
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	userId1 := neo4jt.CreateUser(driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "user",
	})
	neo4jt.CreateUser(driver, tenantName, entity.UserEntity{
		FirstName: "second",
		LastName:  "user",
	})

	neo4jt.AddEmailTo(driver, repository.USER, tenantName, userId1, "test@openline.com", true, "MAIN")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "User"))

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
	require.Equal(t, "test@openline.com", user.User.Emails[0].Email)
}

func TestQueryResolver_User_WithConversations(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	user1 := neo4jt.CreateDefaultUser(driver, tenantName)
	user2 := neo4jt.CreateDefaultUser(driver, tenantName)
	contact1 := neo4jt.CreateDefaultContact(driver, tenantName)
	contact2 := neo4jt.CreateDefaultContact(driver, tenantName)
	contact3 := neo4jt.CreateDefaultContact(driver, tenantName)

	conv1_1 := neo4jt.CreateConversation(driver, user1, contact1)
	conv1_2 := neo4jt.CreateConversation(driver, user1, contact2)
	conv2_1 := neo4jt.CreateConversation(driver, user2, contact1)
	conv2_3 := neo4jt.CreateConversation(driver, user2, contact3)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "Conversation"))

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
	require.ElementsMatch(t, []string{contact1, contact2}, []string{conversations[0].Contacts[0].ID, conversations[1].Contacts[0].ID})
	require.Equal(t, user1, conversations[0].Users[0].ID)
	require.Equal(t, user1, conversations[1].Users[0].ID)

	require.NotNil(t, conv2_1)
	require.NotNil(t, conv2_3)
}
