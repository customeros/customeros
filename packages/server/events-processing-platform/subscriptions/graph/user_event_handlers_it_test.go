package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphUserEventHandler_OnUserCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userEventHandler := &GraphUserEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myUserId, _ := uuid.NewUUID()
	userAggregate := aggregate.NewUserAggregateWithTenantAndID(tenantName, myUserId.String())
	curTime := time.Now().UTC()

	event, err := events.NewUserCreateEvent(userAggregate, &models.UserDto{
		ID:        myUserId.String(),
		Tenant:    tenantName,
		FirstName: "Bob",
		LastName:  "Dole",
		Name:      "Bob Dole",
		Source: commonModels.Source{
			Source:        "N/A",
			SourceOfTruth: "N/A",
			AppSource:     "unit-test",
		},
		CreatedAt: nil,
		UpdatedAt: nil,
	}, curTime, curTime)
	require.Nil(t, err)
	err = userEventHandler.OnUserCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "User_"+tenantName), "Incorrect number of User_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "USER_BELONGS_TO_TENANT"), "Incorrect number of USER_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, myUserId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, myUserId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(props, "firstName"))
	require.Equal(t, "Dole", utils.GetStringPropOrEmpty(props, "lastName"))
	require.Equal(t, "Bob Dole", utils.GetStringPropOrEmpty(props, "name"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(props, "source"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(props, "sourceOfTruth"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(props, "appSource"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "syncedWithEventStore"))

}
