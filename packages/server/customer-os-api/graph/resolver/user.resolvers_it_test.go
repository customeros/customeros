package resolver

import (
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
