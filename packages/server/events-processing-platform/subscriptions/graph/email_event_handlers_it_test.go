package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphEmailEventHandler_OnEmailCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	emailEventHandler := &GraphEmailEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myMailId, _ := uuid.NewUUID()
	emailAggregate := aggregate.NewEmailAggregateWithTenantAndID(tenantName, myMailId.String())
	email := "test@test.com"
	curTime := time.Now().UTC()
	event, err := events.NewEmailCreateEvent(emailAggregate, tenantName, email, commonModels.Source{
		Source:        "N/A",
		SourceOfTruth: "N/A",
		AppSource:     "unit-test",
	}, curTime, curTime)
	require.Nil(t, err)
	err = emailEventHandler.OnEmailCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Email_"+tenantName), "Incorrect number of Email_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "EMAIL_ADDRESS_BELONGS_TO_TENANT"), "Incorrect number of EMAIL_ADDRESS_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, myMailId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, myMailId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, email, utils.GetStringPropOrEmpty(props, "rawEmail"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(props, "appSource"))

}
