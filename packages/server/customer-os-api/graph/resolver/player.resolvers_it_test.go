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

func TestQueryResolver_PlayerByAuthIdProvider(t *testing.T) {
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

	playerId1 := neo4jt.CreatePlayerWithId(ctx, driver, "", entity.PlayerEntity{
		AuthId:     "test@openline.com",
		Provider:   "dummy_provider",
		IdentityId: utils.StringPtr("123456789"),
	})

	neo4jt.LinkPlayerToUser(ctx, driver, playerId1, userId1, true)
	neo4jt.LinkPlayerToUser(ctx, driver, playerId1, userId2, false)

	rawResponse, err := c.RawPost(getQuery("player/get_player_by_authid_provider"),
		client.Var("authId", "test@openline.com"),
		client.Var("provider", "dummy_provider"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var player struct {
		Player_ByAuthIdProvider model.Player
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &player)
	require.Nil(t, err)
	require.NotNil(t, player)
	require.Equal(t, playerId1, player.Player_ByAuthIdProvider.ID)
	require.Equal(t, len(player.Player_ByAuthIdProvider.Users), 2)
	for _, user := range player.Player_ByAuthIdProvider.Users {
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

func TestQueryResolver_PlayerGetUsers(t *testing.T) {
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

	playerId1 := neo4jt.CreateDefaultPlayer(ctx, driver, "test@openline.com", "dummy_provider")

	neo4jt.LinkPlayerToUser(ctx, driver, playerId1, userId1, true)
	neo4jt.LinkPlayerToUser(ctx, driver, playerId1, userId2, false)

	rawResponse, err := c.RawPost(getQuery("player/get_users"))
	assertRawResponseSuccess(t, rawResponse, err)

	var users struct {
		Player_GetUsers []model.PlayerUser
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &users)
	require.Nil(t, err)
	require.NotNil(t, users)
	for _, user := range users.Player_GetUsers {
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

func TestMutationResolver_PlayerSetDefaultUser(t *testing.T) {
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

	playerId1 := neo4jt.CreatePlayerWithId(ctx, driver, "", entity.PlayerEntity{
		AuthId:     "test@openline.com",
		Provider:   "dummy_provider",
		IdentityId: utils.StringPtr("123456789"),
	})

	neo4jt.LinkPlayerToUser(ctx, driver, playerId1, userId1, true)
	neo4jt.LinkPlayerToUser(ctx, driver, playerId1, userId2, false)

	rawResponse, err := c.RawPost(getQuery("player/set_default_user"),
		client.Var("playerId", playerId1),
		client.Var("userId", userId2),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var player struct {
		Player_SetDefaultUser model.Player
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &player)
	bytes, err := json.Marshal(rawResponse.Data)
	log.Printf("Response: %s", string(bytes))
	require.Nil(t, err)
	require.NotNil(t, player)
	require.Equal(t, playerId1, player.Player_SetDefaultUser.ID)
	require.Equal(t, len(player.Player_SetDefaultUser.Users), 2)
	for _, user := range player.Player_SetDefaultUser.Users {
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

func TestMutationResolver_PlayerMerge(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	rawResponse, err := cOwner.RawPost(getQuery("player/player_merge"),
		client.Var("authId", "test@openline.ai"),
		client.Var("provider", "dummy_provider"),
		client.Var("appSource", "dummy_app"))
	assertRawResponseSuccess(t, rawResponse, err)
	bytes, err := json.Marshal(rawResponse.Data)
	log.Printf("Response: %s", string(bytes))

	var player struct {
		Player_Merge model.Player
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &player)
	require.Nil(t, err)
	require.NotNil(t, player)

	createdPlayer := player.Player_Merge
	require.NotNil(t, createdPlayer.ID)
	require.NotNil(t, createdPlayer.CreatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPlayer.CreatedAt)
	require.NotNil(t, createdPlayer.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPlayer.UpdatedAt)
	require.Equal(t, createdPlayer.UpdatedAt, createdPlayer.CreatedAt)
	require.Equal(t, "test@openline.ai", createdPlayer.AuthID)
	require.Equal(t, "dummy_provider", createdPlayer.Provider)
	require.Equal(t, "dummy_app", createdPlayer.AppSource)
	require.Nil(t, createdPlayer.IdentityID)

	require.Equal(t, model.DataSourceOpenline, createdPlayer.Source)
	require.Equal(t, model.DataSourceOpenline, createdPlayer.SourceOfTruth)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Player"))
}

func TestMutationResolver_PlayerMergeAccessControlled(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	rawResponse, err := c.RawPost(getQuery("player/player_merge"),
		client.Var("authId", "test@openline.ai"),
		client.Var("provider", "dummy_provider"),
		client.Var("appSource", "dummy_app"))
	require.Nil(t, err)
	require.NotNil(t, rawResponse.Errors)

}

func TestMutationResolver_PlayerUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	playerId1 := neo4jt.CreateDefaultPlayer(ctx, driver, "test@openline.ai", "dummy_provider")

	rawResponse, err := cOwner.RawPost(getQuery("player/player_update"),
		client.Var("playerId", playerId1),
		client.Var("identityId", "123456789"),
		client.Var("appSource", "dummy_app2"))
	assertRawResponseSuccess(t, rawResponse, err)
	bytes, err := json.Marshal(rawResponse.Data)
	log.Printf("Response: %s", string(bytes))

	var player struct {
		Player_Update model.Player
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &player)
	require.Nil(t, err)
	require.NotNil(t, player)

	createdPlayer := player.Player_Update
	require.NotNil(t, createdPlayer.ID)
	require.NotNil(t, createdPlayer.CreatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPlayer.CreatedAt)
	require.NotNil(t, createdPlayer.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdPlayer.UpdatedAt)
	require.Equal(t, createdPlayer.UpdatedAt, createdPlayer.CreatedAt)
	require.Equal(t, "test@openline.ai", createdPlayer.AuthID)
	require.Equal(t, "dummy_provider", createdPlayer.Provider)
	require.Equal(t, "dummy_app2", createdPlayer.AppSource)
	require.NotNil(t, createdPlayer.IdentityID)
	require.Equal(t, *createdPlayer.IdentityID, "123456789")

	require.Equal(t, model.DataSourceNa, createdPlayer.Source) // test provisioning sets NA
	require.Equal(t, model.DataSourceOpenline, createdPlayer.SourceOfTruth)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Player"))
}

func TestMutationResolver_PlayerUpdateAccessControlled(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	playerId1 := neo4jt.CreateDefaultPlayer(ctx, driver, "test@openline.ai", "dummy_provider")

	rawResponse, err := c.RawPost(getQuery("player/player_update"),
		client.Var("playerId", playerId1),
		client.Var("identityId", "123456789"),
		client.Var("appSource", "dummy_app2"))
	require.Nil(t, err)
	require.NotNil(t, rawResponse.Errors)

}
