package graph

import (
	"testing"

	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	event "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestOrganizationPlanEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{
		Name: "test org",
	})
	mpid := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatform,
		Name:          "master plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     utils.Now(),
		UpdatedAt:     utils.Now(),
		Retired:       false,
	})

	neo4jtest.CreateMasterPlanMilestone(ctx, testDatabase.Driver, tenantName, mpid, neo4jentity.MasterPlanMilestoneEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatform,
		Name:          "milestone name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     utils.Now(),
		UpdatedAt:     utils.Now(),
		Retired:       false,
		Order:         0,
		DurationHours: 24,
		Items:         []string{"item1", "item2"},
		Optional:      false,
	})

	// Prepare the event handler
	orgPlanEventHandler := &OrganizationPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanCreateEvent
	orgPlanId := uuid.New().String()
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	timeNow := utils.Now()
	createEvent, err := event.NewOrganizationPlanCreateEvent(
		orgAggregate,
		orgPlanId,
		mpid,
		"org plan name",
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
	)
	require.Nil(t, err)

	// EXECUTE
	err = orgPlanEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelOrganizationPlan:                    1,
		neo4jutil.NodeLabelOrganizationPlan + "_" + tenantName: 1})

	orgPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelOrganizationPlan, orgPlanId)
	require.Nil(t, err)
	require.NotNil(t, orgPlanDbNode)

	// verify org plan node
	orgPlan := neo4jmapper.MapDbNodeToOrganizationPlanEntity(orgPlanDbNode)
	require.Equal(t, orgPlanId, orgPlan.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), orgPlan.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, orgPlan.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), orgPlan.SourceOfTruth)
	require.Equal(t, timeNow, orgPlan.CreatedAt)
	test.AssertRecentTime(t, orgPlan.UpdatedAt)
	require.Equal(t, "org plan name", orgPlan.Name)
	require.Equal(t, model.NotStarted.String(), orgPlan.StatusDetails.Status)

	createdMilestones := neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "HAS_MILESTONE")
	// should be 2 => 1 master plan milestone + 1 org plan milestone
	require.Equal(t, 2, createdMilestones)
}

// func TestOrganizationPlanEventHandler_OnCreateMilestone(t *testing.T) {
// 	ctx := context.Background()
// 	defer tearDownTestCase(ctx, testDatabase)(t)

// 	// prepare neo4j data
// 	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
// 	masterPlanId := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{})

// 	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
// 		neo4jentity.NodeLabelMasterPlan: 1,
// 	})

// 	// Prepare the event handler
// 	masterPlanEventHandler := &OrganizationPlanEventHandler{
// 		log:          testLogger,
// 		repositories: testDatabase.Repositories,
// 	}

// 	// Create an MasterPlanMilestoneCreateEvent
// 	milestoneId := uuid.New().String()
// 	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenantName, masterPlanId)
// 	timeNow := utils.Now()
// 	createEvent, err := event.NewMasterPlanMilestoneCreateEvent(
// 		masterPlanAggregate,
// 		milestoneId,
// 		"milestone name",
// 		24,
// 		10,
// 		[]string{"item1", "item2"},
// 		true,
// 		commonmodel.Source{
// 			Source:    constants.SourceOpenline,
// 			AppSource: constants.AppSourceEventProcessingPlatform,
// 		},
// 		timeNow,
// 	)
// 	require.Nil(t, err)

// 	// EXECUTE
// 	err = masterPlanEventHandler.OnCreateMilestone(context.Background(), createEvent)
// 	require.Nil(t, err)

// 	// verify nodes and relationships
// 	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
// 		neo4jentity.NodeLabelMasterPlan:                             1,
// 		neo4jentity.NodeLabelMasterPlanMilestone:                    1,
// 		neo4jentity.NodeLabelMasterPlanMilestone + "_" + tenantName: 1})
// 	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, masterPlanId, "HAS_MILESTONE", milestoneId)

// 	// verify master plan milestone node
// 	masterPlanMilestoneDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jentity.NodeLabelMasterPlanMilestone, milestoneId)
// 	require.Nil(t, err)
// 	require.NotNil(t, masterPlanMilestoneDbNode)

// 	milestone := neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneDbNode)
// 	require.Equal(t, milestoneId, milestone.Id)
// 	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), milestone.Source)
// 	require.Equal(t, constants.AppSourceEventProcessingPlatform, milestone.AppSource)
// 	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), milestone.SourceOfTruth)
// 	require.Equal(t, timeNow, milestone.CreatedAt)
// 	test.AssertRecentTime(t, milestone.UpdatedAt)
// 	require.Equal(t, "milestone name", milestone.Name)
// 	require.Equal(t, int64(10), milestone.Order)
// 	require.Equal(t, int64(24), milestone.DurationHours)
// 	require.Equal(t, []string{"item1", "item2"}, milestone.Items)
// 	require.Equal(t, true, milestone.Optional)
// }

// func TestOrganizationPlanEventHandler_OnUpdate(t *testing.T) {
// 	ctx := context.Background()
// 	defer tearDownTestCase(ctx, testDatabase)(t)

// 	// prepare neo4j data
// 	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
// 	masterPlanId := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{})

// 	// Prepare the event handler
// 	masterPlanEventHandler := &OrganizationPlanEventHandler{
// 		log:          testLogger,
// 		repositories: testDatabase.Repositories,
// 	}

// 	// Create an MasterPlanUpdateEvent
// 	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenantName, masterPlanId)
// 	timeNow := utils.Now()
// 	updateEvent, err := event.NewMasterPlanUpdateEvent(
// 		masterPlanAggregate,
// 		"master plan updated name",
// 		true,
// 		timeNow,
// 		[]string{event.FieldMaskName, event.FieldMaskRetired},
// 	)
// 	require.Nil(t, err)

// 	// EXECUTE
// 	err = masterPlanEventHandler.OnUpdate(context.Background(), updateEvent)
// 	require.Nil(t, err)

// 	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
// 		neo4jentity.NodeLabelMasterPlan:                    1,
// 		neo4jentity.NodeLabelMasterPlan + "_" + tenantName: 1})

// 	masterPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jentity.NodeLabelMasterPlan, masterPlanId)
// 	require.Nil(t, err)
// 	require.NotNil(t, masterPlanDbNode)

// 	// verify master plan node
// 	masterPlan := neo4jmapper.MapDbNodeToMasterPlanEntity(masterPlanDbNode)
// 	require.Equal(t, masterPlanId, masterPlan.Id)
// 	require.Equal(t, timeNow, masterPlan.UpdatedAt)
// 	require.Equal(t, "master plan updated name", masterPlan.Name)
// 	require.Equal(t, true, masterPlan.Retired)
// }

// func TestOrganizationPlanEventHandler_OnUpdateMilestone(t *testing.T) {
// 	ctx := context.Background()
// 	defer tearDownTestCase(ctx, testDatabase)(t)

// 	// prepare neo4j data
// 	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
// 	masterPlanId := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{})
// 	milestoneId := neo4jtest.CreateMasterPlanMilestone(ctx, testDatabase.Driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{})

// 	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
// 		neo4jentity.NodeLabelMasterPlan:          1,
// 		neo4jentity.NodeLabelMasterPlanMilestone: 1,
// 	})

// 	// Prepare the event handler
// 	masterPlanEventHandler := &OrganizationPlanEventHandler{
// 		log:          testLogger,
// 		repositories: testDatabase.Repositories,
// 	}

// 	// Create an MasterPlanMilestoneCreateEvent
// 	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenantName, masterPlanId)
// 	timeNow := utils.Now()
// 	updateEvent, err := event.NewMasterPlanMilestoneUpdateEvent(
// 		masterPlanAggregate,
// 		milestoneId,
// 		"new name",
// 		24,
// 		10,
// 		[]string{"item1", "item2"},
// 		[]string{event.FieldMaskName, event.FieldMaskOptional, event.FieldMaskItems, event.FieldMaskDurationHours, event.FieldMaskOrder},
// 		true,
// 		true,
// 		timeNow,
// 	)
// 	require.Nil(t, err)

// 	// EXECUTE
// 	err = masterPlanEventHandler.OnUpdateMilestone(context.Background(), updateEvent)
// 	require.Nil(t, err)

// 	// verify nodes and relationships
// 	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
// 		neo4jentity.NodeLabelMasterPlan:          1,
// 		neo4jentity.NodeLabelMasterPlanMilestone: 1})

// 	// verify master plan milestone node
// 	masterPlanMilestoneDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jentity.NodeLabelMasterPlanMilestone, milestoneId)
// 	require.Nil(t, err)
// 	require.NotNil(t, masterPlanMilestoneDbNode)

// 	milestone := neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneDbNode)
// 	require.Equal(t, milestoneId, milestone.Id)
// 	require.Equal(t, timeNow, milestone.UpdatedAt)
// 	require.Equal(t, "new name", milestone.Name)
// 	require.Equal(t, int64(10), milestone.Order)
// 	require.Equal(t, int64(24), milestone.DurationHours)
// 	require.Equal(t, []string{"item1", "item2"}, milestone.Items)
// 	require.Equal(t, true, milestone.Optional)
// 	require.Equal(t, false, milestone.Retired)
// }

// func TestOrganizationPlanEventHandler_OnReorderMilestones(t *testing.T) {
// 	ctx := context.Background()
// 	defer tearDownTestCase(ctx, testDatabase)(t)

// 	// prepare neo4j data
// 	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
// 	masterPlanId := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{})
// 	milestoneId1 := neo4jtest.CreateMasterPlanMilestone(ctx, testDatabase.Driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{Order: 0})
// 	milestoneId2 := neo4jtest.CreateMasterPlanMilestone(ctx, testDatabase.Driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{Order: 1})

// 	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
// 		neo4jentity.NodeLabelMasterPlan:          1,
// 		neo4jentity.NodeLabelMasterPlanMilestone: 2,
// 	})

// 	// Prepare the event handler
// 	masterPlanEventHandler := &OrganizationPlanEventHandler{
// 		log:          testLogger,
// 		repositories: testDatabase.Repositories,
// 	}

// 	// Create an MasterPlanMilestoneCreateEvent
// 	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenantName, masterPlanId)
// 	timeNow := utils.Now()
// 	reorderEvent, err := event.NewMasterPlanMilestoneReorderEvent(
// 		masterPlanAggregate,
// 		[]string{milestoneId2, milestoneId1},
// 		timeNow,
// 	)
// 	require.Nil(t, err)

// 	// EXECUTE
// 	err = masterPlanEventHandler.OnReorderMilestones(context.Background(), reorderEvent)
// 	require.Nil(t, err)

// 	// verify nodes and relationships
// 	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
// 		neo4jentity.NodeLabelMasterPlan:          1,
// 		neo4jentity.NodeLabelMasterPlanMilestone: 2})

// 	// verify master plan milestone nodes
// 	masterPlanMilestoneDbNode1, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jentity.NodeLabelMasterPlanMilestone, milestoneId1)
// 	require.Nil(t, err)
// 	require.NotNil(t, masterPlanMilestoneDbNode1)
// 	milestone1 := neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneDbNode1)
// 	require.Equal(t, int64(1), milestone1.Order)

// 	masterPlanMilestoneDbNode2, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jentity.NodeLabelMasterPlanMilestone, milestoneId2)
// 	require.Nil(t, err)
// 	require.NotNil(t, masterPlanMilestoneDbNode2)
// 	milestone2 := neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneDbNode2)
// 	require.Equal(t, int64(0), milestone2.Order)
// }
