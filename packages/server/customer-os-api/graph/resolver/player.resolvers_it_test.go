package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_PlayerByAuthIdProvider(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	otherTenant := "other"
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenant(ctx, driver, otherTenant)
	userId1 := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		FirstName: "first",
		LastName:  "last",
		Roles:     []string{model.RoleUser.String(), model.RoleOwner.String()},
	})
	userId2 := neo4jtest.CreateUser(ctx, driver, otherTenant, neo4jentity.UserEntity{
		FirstName: "otherFirst",
		LastName:  "otherLast",
		Roles:     []string{model.RoleUser.String()},
	})

	neo4jt.AddEmailTo(ctx, driver, commonModel.USER, tenantName, userId1, "test@openline.com", true, "MAIN")
	neo4jt.AddEmailTo(ctx, driver, commonModel.USER, otherTenant, userId2, "test@openline.com", true, "MAIN")

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

//TODO
//func TestMutationResolver_PlayerMerge(t *testing.T) {
//	ctx := context.TODO()
//	defer tearDownTestCase(ctx)(t)
//
//	rawResponse, err := cOwner.RawPost(getQuery("player/player_merge"),
//		client.Var("userId", "1"),
//		client.Var("authId", "test@openline.ai"),
//		client.Var("provider", "dummy_provider"),
//		client.Var("appSource", "dummy_app"))
//	assertRawResponseSuccess(t, rawResponse, err)
//	bytes, err := json.Marshal(rawResponse.Data)
//	log.Printf("Response: %s", string(bytes))
//
//	var player struct {
//		Player_Merge model.Result
//	}
//
//	err = decode.Decode(rawResponse.Data.(map[string]any), &player)
//	require.Nil(t, err)
//	require.NotNil(t, player)
//
//	createdPlayer := player.Player_Merge
//	require.Equal(t, true, createdPlayer.Result)
//
//	// Check the number of nodes and relationships in the Neo4j database
//	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Player"))
//}

func TestMutationResolver_PlayerMergeAccessControlled(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	rawResponse, err := c.RawPost(getQuery("player/player_merge"),
		client.Var("userId", "1"),
		client.Var("authId", "test@openline.ai"),
		client.Var("provider", "dummy_provider"),
		client.Var("appSource", "dummy_app"))
	require.Nil(t, err)
	require.NotNil(t, rawResponse.Errors)

}
