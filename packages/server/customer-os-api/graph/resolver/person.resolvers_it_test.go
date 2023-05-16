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

func TestQueryResolver_PersonByEmailProvider(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	otherTenant := "other"
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, otherTenant)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{model.RoleUser.String(), model.RoleOwner.String()},
	})
	userId2 := neo4jt.CreateUser(ctx, driver, otherTenant, entity.UserEntity{
		FirstName: "otherFirst",
		LastName:  "otherLast",
		Roles:     []string{model.RoleUser.String()},
	})

	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId1, "test@openline.com", true, "MAIN")
	neo4jt.AddEmailTo(ctx, driver, entity.USER, otherTenant, userId2, "test@openline.com", true, "MAIN")

	personId1 := neo4jt.CreatePersonWithId(ctx, driver, "", entity.PersonEntity{
		Email:      "test@openline.com",
		Provider:   "dummy_provider",
		IdentityId: utils.StringPtr("123456789"),
	})

	neo4jt.LinkPersonToUser(ctx, driver, personId1, userId1, true)
	neo4jt.LinkPersonToUser(ctx, driver, personId1, userId2, false)

	rawResponse, err := c.RawPost(getQuery("person/get_person_by_email_provider"),
		client.Var("email", "test@openline.com"),
		client.Var("provider", "dummy_provider"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var person struct {
		Person_ByEmailProvider model.Person
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &person)
	require.Nil(t, err)
	require.NotNil(t, person)
	require.Equal(t, personId1, person.Person_ByEmailProvider.ID)
	require.Equal(t, len(person.Person_ByEmailProvider.Users), 2)
	for _, user := range person.Person_ByEmailProvider.Users {
		if user.User.ID == userId1 {
			require.True(t, user.Default)
			require.Contains(t, user.User.Roles, model.RoleUser)
			require.Contains(t, user.User.Roles, model.RoleOwner)
			require.Equal(t, user.Tenant, tenantName)

		} else if user.User.ID == userId2 {
			require.False(t, user.Default)
			require.Contains(t, user.User.Roles, model.RoleUser)
			require.NotContains(t, user.User.Roles, model.RoleOwner)
			require.Equal(t, user.Tenant, otherTenant)
		} else {
			t.Errorf("Unexpected user %s", user.User.ID)
		}
	}
}

func TestQueryResolver_PersonGetUsers(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	otherTenant := "other"
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, otherTenant)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{model.RoleUser.String(), model.RoleOwner.String()},
	})
	userId2 := neo4jt.CreateUser(ctx, driver, otherTenant, entity.UserEntity{
		FirstName: "otherFirst",
		LastName:  "otherLast",
		Roles:     []string{model.RoleUser.String()},
	})

	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId1, "test@openline.com", true, "MAIN")
	neo4jt.AddEmailTo(ctx, driver, entity.USER, otherTenant, userId2, "test@openline.com", true, "MAIN")

	personId1 := neo4jt.CreateDefaultPerson(ctx, driver, "test@openline.com", "dummy_provider")

	neo4jt.LinkPersonToUser(ctx, driver, personId1, userId1, true)
	neo4jt.LinkPersonToUser(ctx, driver, personId1, userId2, false)

	rawResponse, err := c.RawPost(getQuery("person/get_users"))
	assertRawResponseSuccess(t, rawResponse, err)

	var users struct {
		Person_GetUsers []model.PersonUser
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &users)
	require.Nil(t, err)
	require.NotNil(t, users)
	for _, user := range users.Person_GetUsers {
		if user.User.ID == userId1 {
			require.True(t, user.Default)
			require.Contains(t, user.User.Roles, model.RoleUser)
			require.Contains(t, user.User.Roles, model.RoleOwner)
			require.Equal(t, user.Tenant, tenantName)

		} else if user.User.ID == userId2 {
			require.False(t, user.Default)
			require.Contains(t, user.User.Roles, model.RoleUser)
			require.NotContains(t, user.User.Roles, model.RoleOwner)
			require.Equal(t, user.Tenant, otherTenant)
		} else {
			t.Errorf("Unexpected user %s", user.User.ID)
		}
	}
}

func TestMutationResolver_PersonSetDefaultUser(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	otherTenant := "other"
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, otherTenant)
	userId1 := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{model.RoleUser.String(), model.RoleOwner.String()},
	})
	userId2 := neo4jt.CreateUser(ctx, driver, otherTenant, entity.UserEntity{
		FirstName: "otherFirst",
		LastName:  "otherLast",
		Roles:     []string{model.RoleUser.String()},
	})

	neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId1, "test@openline.com", true, "MAIN")
	neo4jt.AddEmailTo(ctx, driver, entity.USER, otherTenant, userId2, "test@openline.com", true, "MAIN")

	personId1 := neo4jt.CreatePersonWithId(ctx, driver, "", entity.PersonEntity{
		Email:      "test@openline.com",
		Provider:   "dummy_provider",
		IdentityId: utils.StringPtr("123456789"),
	})

	neo4jt.LinkPersonToUser(ctx, driver, personId1, userId1, true)
	neo4jt.LinkPersonToUser(ctx, driver, personId1, userId2, false)

	rawResponse, err := c.RawPost(getQuery("person/set_default_user"),
		client.Var("personId", personId1),
		client.Var("userId", userId2),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var person struct {
		Person_SetDefaultUser model.Person
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &person)
	bytes, err := json.Marshal(rawResponse.Data)
	log.Printf("Response: %s", string(bytes))
	require.Nil(t, err)
	require.NotNil(t, person)
	require.Equal(t, personId1, person.Person_SetDefaultUser.ID)
	require.Equal(t, len(person.Person_SetDefaultUser.Users), 2)
	for _, user := range person.Person_SetDefaultUser.Users {
		if user.User.ID == userId1 {
			require.False(t, user.Default) // should not be default anymore
			require.Contains(t, user.User.Roles, model.RoleUser)
			require.Contains(t, user.User.Roles, model.RoleOwner)
			require.Equal(t, user.Tenant, tenantName)

		} else if user.User.ID == userId2 {
			require.True(t, user.Default) // should be default now
			require.Contains(t, user.User.Roles, model.RoleUser)
			require.NotContains(t, user.User.Roles, model.RoleOwner)
			require.Equal(t, user.Tenant, otherTenant)
		} else {
			t.Errorf("Unexpected user %s", user.User.ID)
		}
	}
}

func TestMutationResolver_PersonMerge(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	rawResponse, err := cOwner.RawPost(getQuery("person/person_merge"),
		client.Var("email", "test@openline.ai"),
		client.Var("provider", "dummy_provider"),
		client.Var("appSource", "dummy_app"))
	assertRawResponseSuccess(t, rawResponse, err)
	bytes, err := json.Marshal(rawResponse.Data)
	log.Printf("Response: %s", string(bytes))

	var person struct {
		Person_Merge model.Person
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &person)
	require.Nil(t, err)
	require.NotNil(t, person)

	createdPerson := person.Person_Merge
	require.NotNil(t, createdPerson.ID)
	require.NotNil(t, createdPerson.CreatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPerson.CreatedAt)
	require.NotNil(t, createdPerson.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPerson.UpdatedAt)
	require.Equal(t, createdPerson.UpdatedAt, createdPerson.CreatedAt)
	require.Equal(t, "test@openline.ai", createdPerson.Email)
	require.Equal(t, "dummy_provider", createdPerson.Provider)
	require.Equal(t, "dummy_app", createdPerson.AppSource)
	require.Nil(t, createdPerson.IdentityID)

	require.Equal(t, model.DataSourceOpenline, createdPerson.Source)
	require.Equal(t, model.DataSourceOpenline, createdPerson.SourceOfTruth)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Person"))
}

func TestMutationResolver_PersonMergeAccessControlled(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	rawResponse, err := c.RawPost(getQuery("person/person_merge"),
		client.Var("email", "test@openline.ai"),
		client.Var("provider", "dummy_provider"),
		client.Var("appSource", "dummy_app"))
	require.Nil(t, err)
	require.NotNil(t, rawResponse.Errors)

}

func TestMutationResolver_PersonUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	personId1 := neo4jt.CreateDefaultPerson(ctx, driver, "test@openline.ai", "dummy_provider")

	rawResponse, err := cOwner.RawPost(getQuery("person/person_update"),
		client.Var("personId", personId1),
		client.Var("identityId", "123456789"),
		client.Var("appSource", "dummy_app2"))
	assertRawResponseSuccess(t, rawResponse, err)
	bytes, err := json.Marshal(rawResponse.Data)
	log.Printf("Response: %s", string(bytes))

	var person struct {
		Person_Update model.Person
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &person)
	require.Nil(t, err)
	require.NotNil(t, person)

	createdPerson := person.Person_Update
	require.NotNil(t, createdPerson.ID)
	require.NotNil(t, createdPerson.CreatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPerson.CreatedAt)
	require.NotNil(t, createdPerson.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPerson.UpdatedAt)
	require.Equal(t, createdPerson.UpdatedAt, createdPerson.CreatedAt)
	require.Equal(t, "test@openline.ai", createdPerson.Email)
	require.Equal(t, "dummy_provider", createdPerson.Provider)
	require.Equal(t, "dummy_app2", createdPerson.AppSource)
	require.NotNil(t, createdPerson.IdentityID)
	require.Equal(t, *createdPerson.IdentityID, "123456789")

	require.Equal(t, model.DataSourceNa, createdPerson.Source) // test provisioning sets NA
	require.Equal(t, model.DataSourceOpenline, createdPerson.SourceOfTruth)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Person"))
}

func TestMutationResolver_PersonUpdateAccessControlled(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	personId1 := neo4jt.CreateDefaultPerson(ctx, driver, "test@openline.ai", "dummy_provider")

	rawResponse, err := c.RawPost(getQuery("person/person_update"),
		client.Var("personId", personId1),
		client.Var("identityId", "123456789"),
		client.Var("appSource", "dummy_app2"))
	require.Nil(t, err)
	require.NotNil(t, rawResponse.Errors)

}
