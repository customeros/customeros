package graph

import (
	"github.com/google/uuid"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestMasterPlanEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)

	// Prepare the event handler
	masterPlanEventHandler := &MasterPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanCreateEvent
	masterPlanId := uuid.New().String()
	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenantName, masterPlanId)
	timeNow := utils.Now()
	createEvent, err := event.NewMasterPlanCreateEvent(
		masterPlanAggregate,
		"master plan name",
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
	)
	require.Nil(t, err)

	// EXECUTE
	err = masterPlanEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jentity.NodeLabel_MasterPlan:                    1,
		neo4jentity.NodeLabel_MasterPlan + "_" + tenantName: 1})

	masterPlanDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, neo4jentity.NodeLabel_MasterPlan, masterPlanId)
	require.Nil(t, err)
	require.NotNil(t, masterPlanDbNode)

	// verify master plan node
	masterPlan := graph_db.MapDbNodeToMasterPlanEntity(masterPlanDbNode)
	require.Equal(t, masterPlanId, masterPlan.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), masterPlan.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, masterPlan.AppSource)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), masterPlan.SourceOfTruth)
	require.Equal(t, timeNow, masterPlan.CreatedAt)
	test.AssertRecentTime(t, masterPlan.UpdatedAt)
	require.Equal(t, "master plan name", masterPlan.Name)
}
