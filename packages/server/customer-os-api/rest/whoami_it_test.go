package rest

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var whoamIRouter *gin.Engine

func setup() {
	whoamIRouter = gin.Default()

	whoamIRouter.GET("/whoami",
		WhoamiHandler(serviceContainer))
}

func TestGet_Whoami(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	setup()
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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/whoami", nil)
	req.Header.Add("X-Openline-Identity-Id", testPlayerId)

	whoamIRouter.ServeHTTP(w, req)
	require.Equal(t, w.Code, 200)
	var resp WhoAmIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("WhoAmi, unable to decode json: %v\n", err.Error())
		return
	}
	require.Equal(t, len(resp.Users), 2)
	require.Equal(t, len(resp.Users[0].Emails), 1)
	require.Equal(t, len(resp.Users[1].Emails), 1)
}
