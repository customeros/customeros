package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphEmailEventHandler_OnEmailCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	emailEventHandler := &GraphEmailEventHandler{
		Repositories: repositories,
	}
	myMailId, _ := uuid.NewUUID()
	emailAggregate := aggregate.NewEmailAggregateWithTenantAndID(tenantName, myMailId.String())
	email := "test@test.com"
	curTime := time.Now().UTC()

	event, err := events.NewEmailCreatedEvent(emailAggregate, tenantName, email, "N/A", "N/A", "unit-test", curTime, curTime)
	require.Nil(t, err)
	err = emailEventHandler.OnEmailCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email_"+tenantName), "Incorrect number of Email_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "EMAIL_ADDRESS_BELONGS_TO_TENANT"), "Incorrect number of EMAIL_ADDRESS_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jt.GetNodeById(ctx, driver, "Email_"+tenantName, myMailId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, myMailId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, email, utils.GetStringPropOrEmpty(props, "rawEmail"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(props, "appSource"))

}
