package graph

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
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
		events.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
	)
	require.Nil(t, err)

	// EXECUTE
	err = masterPlanEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model.NodeLabelMasterPlan:                    1,
		model.NodeLabelMasterPlan + "_" + tenantName: 1})

	masterPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model.NodeLabelMasterPlan, masterPlanId)
	require.Nil(t, err)
	require.NotNil(t, masterPlanDbNode)

	// verify master plan node
	masterPlan := neo4jmapper.MapDbNodeToMasterPlanEntity(masterPlanDbNode)
	require.Equal(t, masterPlanId, masterPlan.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), masterPlan.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, masterPlan.AppSource)
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
		model.NodeLabelMasterPlan: 1,
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
		events.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
	)
	require.Nil(t, err)

	// EXECUTE
	err = masterPlanEventHandler.OnCreateMilestone(context.Background(), createEvent)
	require.Nil(t, err)

	// verify nodes and relationships
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model.NodeLabelMasterPlan:                             1,
		model.NodeLabelMasterPlanMilestone:                    1,
		model.NodeLabelMasterPlanMilestone + "_" + tenantName: 1})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, masterPlanId, "HAS_MILESTONE", milestoneId)

	// verify master plan milestone node
	masterPlanMilestoneDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model.NodeLabelMasterPlanMilestone, milestoneId)
	require.Nil(t, err)
	require.NotNil(t, masterPlanMilestoneDbNode)

	milestone := neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneDbNode)
	require.Equal(t, milestoneId, milestone.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), milestone.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, milestone.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), milestone.SourceOfTruth)
	require.Equal(t, timeNow, milestone.CreatedAt)
	test.AssertRecentTime(t, milestone.UpdatedAt)
	require.Equal(t, "milestone name", milestone.Name)
	require.Equal(t, int64(10), milestone.Order)
	require.Equal(t, int64(24), milestone.DurationHours)
	require.Equal(t, []string{"item1", "item2"}, milestone.Items)
	require.Equal(t, true, milestone.Optional)
}

func TestMasterPlanEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	masterPlanId := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{})

	// Prepare the event handler
	masterPlanEventHandler := &MasterPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanUpdateEvent
	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenantName, masterPlanId)
	timeNow := utils.Now()
	updateEvent, err := event.NewMasterPlanUpdateEvent(
		masterPlanAggregate,
		"master plan updated name",
		true,
		timeNow,
		[]string{event.FieldMaskName, event.FieldMaskRetired},
	)
	require.Nil(t, err)

	// EXECUTE
	err = masterPlanEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model.NodeLabelMasterPlan:                    1,
		model.NodeLabelMasterPlan + "_" + tenantName: 1})

	masterPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model.NodeLabelMasterPlan, masterPlanId)
	require.Nil(t, err)
	require.NotNil(t, masterPlanDbNode)

	// verify master plan node
	masterPlan := neo4jmapper.MapDbNodeToMasterPlanEntity(masterPlanDbNode)
	require.Equal(t, masterPlanId, masterPlan.Id)
	test.AssertRecentTime(t, masterPlan.UpdatedAt)
	require.Equal(t, "master plan updated name", masterPlan.Name)
	require.Equal(t, true, masterPlan.Retired)
}

func TestMasterPlanEventHandler_OnUpdateMilestone(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	masterPlanId := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{})
	milestoneId := neo4jtest.CreateMasterPlanMilestone(ctx, testDatabase.Driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model.NodeLabelMasterPlan:          1,
		model.NodeLabelMasterPlanMilestone: 1,
	})

	// Prepare the event handler
	masterPlanEventHandler := &MasterPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanMilestoneCreateEvent
	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenantName, masterPlanId)
	timeNow := utils.Now()
	updateEvent, err := event.NewMasterPlanMilestoneUpdateEvent(
		masterPlanAggregate,
		milestoneId,
		"new name",
		24,
		10,
		[]string{"item1", "item2"},
		[]string{event.FieldMaskName, event.FieldMaskOptional, event.FieldMaskItems, event.FieldMaskDurationHours, event.FieldMaskOrder},
		true,
		true,
		timeNow,
	)
	require.Nil(t, err)

	// EXECUTE
	err = masterPlanEventHandler.OnUpdateMilestone(context.Background(), updateEvent)
	require.Nil(t, err)

	// verify nodes and relationships
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model.NodeLabelMasterPlan:          1,
		model.NodeLabelMasterPlanMilestone: 1})

	// verify master plan milestone node
	masterPlanMilestoneDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model.NodeLabelMasterPlanMilestone, milestoneId)
	require.Nil(t, err)
	require.NotNil(t, masterPlanMilestoneDbNode)

	milestone := neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneDbNode)
	require.Equal(t, milestoneId, milestone.Id)
	test.AssertRecentTime(t, milestone.UpdatedAt)
	require.Equal(t, "new name", milestone.Name)
	require.Equal(t, int64(10), milestone.Order)
	require.Equal(t, int64(24), milestone.DurationHours)
	require.Equal(t, []string{"item1", "item2"}, milestone.Items)
	require.Equal(t, true, milestone.Optional)
	require.Equal(t, false, milestone.Retired)
}

func TestMasterPlanEventHandler_OnReorderMilestones(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	masterPlanId := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{})
	milestoneId1 := neo4jtest.CreateMasterPlanMilestone(ctx, testDatabase.Driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{Order: 0})
	milestoneId2 := neo4jtest.CreateMasterPlanMilestone(ctx, testDatabase.Driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{Order: 1})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model.NodeLabelMasterPlan:          1,
		model.NodeLabelMasterPlanMilestone: 2,
	})

	// Prepare the event handler
	masterPlanEventHandler := &MasterPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanMilestoneCreateEvent
	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenantName, masterPlanId)
	timeNow := utils.Now()
	reorderEvent, err := event.NewMasterPlanMilestoneReorderEvent(
		masterPlanAggregate,
		[]string{milestoneId2, milestoneId1},
		timeNow,
	)
	require.Nil(t, err)

	// EXECUTE
	err = masterPlanEventHandler.OnReorderMilestones(context.Background(), reorderEvent)
	require.Nil(t, err)

	// verify nodes and relationships
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model.NodeLabelMasterPlan:          1,
		model.NodeLabelMasterPlanMilestone: 2})

	// verify master plan milestone nodes
	masterPlanMilestoneDbNode1, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model.NodeLabelMasterPlanMilestone, milestoneId1)
	require.Nil(t, err)
	require.NotNil(t, masterPlanMilestoneDbNode1)
	milestone1 := neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneDbNode1)
	require.Equal(t, int64(1), milestone1.Order)

	masterPlanMilestoneDbNode2, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model.NodeLabelMasterPlanMilestone, milestoneId2)
	require.Nil(t, err)
	require.NotNil(t, masterPlanMilestoneDbNode2)
	milestone2 := neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneDbNode2)
	require.Equal(t, int64(0), milestone2.Order)
}
