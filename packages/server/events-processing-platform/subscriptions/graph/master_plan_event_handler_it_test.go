package graph

import (
	"github.com/google/uuid"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestMasterPlanEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

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

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jentity.NodeLabel_MasterPlan:                    1,
		neo4jentity.NodeLabel_MasterPlan + "_" + tenantName: 1})

	masterPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jentity.NodeLabel_MasterPlan, masterPlanId)
	require.Nil(t, err)
	require.NotNil(t, masterPlanDbNode)

	// verify master plan node
	masterPlan := graph_db.MapDbNodeToMasterPlanEntity(masterPlanDbNode)
	require.Equal(t, masterPlanId, masterPlan.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), masterPlan.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, masterPlan.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), masterPlan.SourceOfTruth)
	require.Equal(t, timeNow, masterPlan.CreatedAt)
	test.AssertRecentTime(t, masterPlan.UpdatedAt)
	require.Equal(t, "master plan name", masterPlan.Name)
}

func TestMasterPlanEventHandler_OnCreateMilestone(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	masterPlanId := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jentity.NodeLabel_MasterPlan: 1,
	})

	// Prepare the event handler
	masterPlanEventHandler := &MasterPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanMilestoneCreateEvent
	milestoneId := uuid.New().String()
	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenantName, masterPlanId)
	timeNow := utils.Now()
	createEvent, err := event.NewMasterPlanMilestoneCreateEvent(
		masterPlanAggregate,
		milestoneId,
		"milestone name",
		24,
		10,
		[]string{"item1", "item2"},
		true,
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
	)
	require.Nil(t, err)

	// EXECUTE
	err = masterPlanEventHandler.OnCreateMilestone(context.Background(), createEvent)
	require.Nil(t, err)

	// verify nodes and relationships
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jentity.NodeLabel_MasterPlan:                             1,
		neo4jentity.NodeLabel_MasterPlanMilestone:                    1,
		neo4jentity.NodeLabel_MasterPlanMilestone + "_" + tenantName: 1})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, masterPlanId, "HAS_MILESTONE", milestoneId)

	// verify master plan node
	masterPlanMilestoneDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jentity.NodeLabel_MasterPlanMilestone, milestoneId)
	require.Nil(t, err)
	require.NotNil(t, masterPlanMilestoneDbNode)

	milestone := graph_db.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneDbNode)
	require.Equal(t, milestoneId, milestone.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), milestone.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, milestone.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), milestone.SourceOfTruth)
	require.Equal(t, timeNow, milestone.CreatedAt)
	test.AssertRecentTime(t, milestone.UpdatedAt)
	require.Equal(t, "milestone name", milestone.Name)
	require.Equal(t, int64(10), milestone.Order)
	require.Equal(t, int64(24), milestone.DurationHours)
	require.Equal(t, []string{"item1", "item2"}, milestone.Items)
	require.Equal(t, true, milestone.Optional)
}
