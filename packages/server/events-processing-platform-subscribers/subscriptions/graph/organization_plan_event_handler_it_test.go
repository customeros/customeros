package graph

import (
	"fmt"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	orgmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	event "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/model"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestOrganizationPlanEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	timeNow := utils.Now()
	so := int64(0)
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{
		Name: "test org",
		OnboardingDetails: neo4jentity.OnboardingDetails{
			Status:       orgmodel.NotApplicable.String(),
			UpdatedAt:    &timeNow,
			Comments:     "",
			SortingOrder: &so,
		},
	})
	mpid := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{
		Source:        constants.SourceOpenline,
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "master plan name",
		SourceOfTruth: constants.SourceOpenline,
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
	})

	neo4jtest.CreateMasterPlanMilestone(ctx, testDatabase.Driver, tenantName, mpid, neo4jentity.MasterPlanMilestoneEntity{
		Source:        constants.SourceOpenline,
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "milestone name",
		SourceOfTruth: constants.SourceOpenline,
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
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
	createEvent, err := event.NewOrganizationPlanCreateEvent(
		orgAggregate,
		orgPlanId,
		mpid,
		orgId,
		"org plan name",
		events.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
	)
	require.Nil(t, err)

	// EXECUTE
	err = orgPlanEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:                    1,
		model2.NodeLabelOrganizationPlan + "_" + tenantName: 1})

	orgPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlan, orgPlanId)
	require.Nil(t, err)
	require.NotNil(t, orgPlanDbNode)

	// verify org plan node
	orgPlan := neo4jmapper.MapDbNodeToOrganizationPlanEntity(orgPlanDbNode)
	require.Equal(t, orgPlanId, orgPlan.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), orgPlan.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, orgPlan.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), orgPlan.SourceOfTruth)
	require.Equal(t, timeNow, orgPlan.CreatedAt)
	test.AssertRecentTime(t, orgPlan.UpdatedAt)
	require.Equal(t, "org plan name", orgPlan.Name)
	require.Equal(t, model.NotStarted.String(), orgPlan.StatusDetails.Status)

	createdMilestones := neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "HAS_MILESTONE")
	// should be 1 => 1 master plan milestone + 0 org plan milestone
	require.Equal(t, 1, createdMilestones)

	// double check there is only one organization plan created
	organizationPlansCount := neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "ORGANIZATION_PLAN_BELONGS_TO_ORGANIZATION")
	require.Equal(t, 1, organizationPlansCount)
	opForOrgNodes, err := neo4jtest.GetAllNodesByLabel(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlan+"_"+tenantName)
	require.Nil(t, err)
	require.Len(t, opForOrgNodes, 1)

	// Check onboarding status updated
	orgDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganization, orgId)
	require.Nil(t, err)
	org := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	require.Equal(t, orgmodel.NotStarted.String(), org.OnboardingDetails.Status)
}

func TestOrganizationPlanEventHandler_OnCreateMilestone(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	timeNow := utils.Now()
	mpid := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "master plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{Name: "test org"})
	opid := neo4jtest.CreateOrganizationPlan(ctx, testDatabase.Driver, tenantName, mpid, orgId, neo4jentity.OrganizationPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "org plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		StatusDetails: neo4jentity.OrganizationPlanStatusDetails{
			Status:    model.NotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan: 1,
	})

	// Prepare the event handler
	orgPlanEventHandler := &OrganizationPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an OrgPlanMilestoneCreateEvent
	milestoneId := uuid.New().String()
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	createEvent, err := event.NewOrganizationPlanMilestoneCreateEvent(
		orgAggregate,
		opid,
		milestoneId,
		"milestone name",
		10,
		[]string{"item1", "item2"},
		true,  // optional
		false, // adhoc
		events.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
		timeNow.Add(time.Hour*24), // due date
	)
	require.Nil(t, err)

	// EXECUTE
	err = orgPlanEventHandler.OnCreateMilestone(context.Background(), createEvent)
	require.Nil(t, err)

	// verify nodes and relationships
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:                             1,
		model2.NodeLabelOrganizationPlanMilestone:                    1,
		model2.NodeLabelOrganizationPlanMilestone + "_" + tenantName: 1})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, opid, "HAS_MILESTONE", milestoneId)

	// verify org plan milestone node
	orgPlanMilestoneDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlanMilestone, milestoneId)
	require.Nil(t, err)
	require.NotNil(t, orgPlanMilestoneDbNode)

	milestone := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(orgPlanMilestoneDbNode)
	require.Equal(t, milestoneId, milestone.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), milestone.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, milestone.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), milestone.SourceOfTruth)
	require.Equal(t, timeNow, milestone.CreatedAt)
	test.AssertRecentTime(t, milestone.UpdatedAt)
	require.Equal(t, "milestone name", milestone.Name)
	require.Equal(t, int64(10), milestone.Order)
	require.Equal(t, timeNow.Add(time.Hour*24), milestone.DueDate)
	require.Equal(t, true, milestone.Optional)
	require.Equal(t, model.NotStarted.String(), milestone.StatusDetails.Status)
	for i, item := range milestone.Items {
		require.Equal(t, model.TaskNotDone.String(), item.Status)
		txt := fmt.Sprintf("item%d", i+1)
		require.Equal(t, txt, item.Text)
		require.NotEqual(t, "", item.Uuid) // have *some* uuid
	}
}

func TestOrganizationPlanEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	timeNow := utils.Now()
	mpid := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "master plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{Name: "test org"})
	opid := neo4jtest.CreateOrganizationPlan(ctx, testDatabase.Driver, tenantName, mpid, orgId, neo4jentity.OrganizationPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "org plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		StatusDetails: neo4jentity.OrganizationPlanStatusDetails{
			Status:    model.NotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan: 1,
	})
	// Prepare the event handler
	orgPlanEventHandler := &OrganizationPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an OrgPlanUpdateEvent
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	updateTime := utils.Now()
	updateEvent, err := event.NewOrganizationPlanUpdateEvent(
		orgAggregate,
		opid,
		"org plan updated name",
		true,
		updateTime,
		[]string{event.FieldMaskName, event.FieldMaskRetired, event.FieldMaskStatusDetails},
		model.OrganizationPlanDetails{
			Status:    model.Late.String(),
			Comments:  "comments",
			UpdatedAt: updateTime,
		},
	)
	require.Nil(t, err)

	// EXECUTE
	err = orgPlanEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:                    1,
		model2.NodeLabelOrganizationPlan + "_" + tenantName: 1})

	orgPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlan, opid)
	require.Nil(t, err)
	require.NotNil(t, orgPlanDbNode)

	// verify org plan node
	orgPlan := neo4jmapper.MapDbNodeToOrganizationPlanEntity(orgPlanDbNode)
	require.Equal(t, opid, orgPlan.Id)
	test.AssertRecentTime(t, orgPlan.UpdatedAt)
	require.Equal(t, "org plan updated name", orgPlan.Name)
	require.Equal(t, true, orgPlan.Retired)
	require.Equal(t, model.Late.String(), orgPlan.StatusDetails.Status)
	require.Equal(t, "comments", orgPlan.StatusDetails.Comments)
	require.Equal(t, updateTime, orgPlan.StatusDetails.UpdatedAt)
}

func TestOrganizationPlanEventHandler_OnUpdateMilestone(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	timeNow := utils.Now()
	mpid := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "master plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{Name: "test org"})
	opid := neo4jtest.CreateOrganizationPlan(ctx, testDatabase.Driver, tenantName, mpid, orgId, neo4jentity.OrganizationPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "org plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		StatusDetails: neo4jentity.OrganizationPlanStatusDetails{
			Status:    model.NotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	milestoneId := neo4jtest.CreateOrganizationPlanMilestone(ctx, testDatabase.Driver, tenantName, opid, neo4jentity.OrganizationPlanMilestoneEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "milestone name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		Order:         0,
		DueDate:       timeNow.Add(time.Hour * 24),
		Items:         []neo4jentity.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskNotDone.String(), UpdatedAt: timeNow, Uuid: "item1"}, {Text: "item2", Status: model.TaskNotDone.String(), UpdatedAt: timeNow, Uuid: "item2"}},
		Optional:      false,
		StatusDetails: neo4jentity.OrganizationPlanMilestoneStatusDetails{
			Status:    model.MilestoneNotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:                             1,
		model2.NodeLabelOrganizationPlanMilestone:                    1,
		model2.NodeLabelOrganizationPlanMilestone + "_" + tenantName: 1,
	})

	// Prepare the event handler
	orgPlanEventHandler := &OrganizationPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanMilestoneCreateEvent
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	updateTime := utils.Now()
	updateEvent, err := event.NewOrganizationPlanMilestoneUpdateEvent(
		orgAggregate,
		opid,
		milestoneId,
		"new name",
		10,
		[]model.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskDone.String(), UpdatedAt: updateTime, Uuid: "item1"}, {Text: "item2Change", Status: model.TaskNotDone.String(), UpdatedAt: updateTime, Uuid: "item2"}},
		[]string{event.FieldMaskName, event.FieldMaskOptional, event.FieldMaskItems, event.FieldMaskDueDate, event.FieldMaskOrder, event.FieldMaskStatusDetails},
		true,  // optional
		false, // adhoc
		true,  // retired
		updateTime,
		timeNow.Add(time.Hour*48), // due date
		model.OrganizationPlanDetails{
			Status:    model.MilestoneNotStarted.String(),
			Comments:  "comments",
			UpdatedAt: updateTime,
		},
	)
	require.Nil(t, err)

	// EXECUTE
	err = orgPlanEventHandler.OnUpdateMilestone(context.Background(), updateEvent)
	require.Nil(t, err)

	// verify nodes and relationships
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:          1,
		model2.NodeLabelOrganizationPlanMilestone: 1})

	// verify master plan milestone node
	orgPlanMilestoneDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlanMilestone, milestoneId)
	require.Nil(t, err)
	require.NotNil(t, orgPlanMilestoneDbNode)

	milestone := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(orgPlanMilestoneDbNode)
	require.Equal(t, milestoneId, milestone.Id)
	test.AssertRecentTime(t, milestone.UpdatedAt)
	require.Equal(t, "new name", milestone.Name)
	require.Equal(t, int64(10), milestone.Order)
	require.Equal(t, timeNow.Add(time.Hour*48), milestone.DueDate)
	require.Equal(t, true, milestone.Optional)
	require.Equal(t, false, milestone.Retired)                                        // mask not passed so we ignore this field update
	require.Equal(t, model.MilestoneStarted.String(), milestone.StatusDetails.Status) // automatic update
	require.Equal(t, "comments", milestone.StatusDetails.Comments)
	require.Equal(t, updateTime, milestone.StatusDetails.UpdatedAt)
	for i, item := range milestone.Items {
		if i == 0 {
			require.Equal(t, model.TaskDone.String(), item.Status)
			require.Equal(t, "item1", item.Text)
			require.Equal(t, "item1", item.Uuid)
		} else {
			require.Equal(t, model.TaskNotDone.String(), item.Status)
			require.Equal(t, "item2Change", item.Text)
			require.Equal(t, "item2", item.Uuid)
		}
	}
	organizationPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlan, opid)
	require.Nil(t, err)
	op := neo4jmapper.MapDbNodeToOrganizationPlanEntity(organizationPlanDbNode)
	require.Equal(t, model.OnTrack.String(), op.StatusDetails.Status) // automatic update

	// Check onboarding status updated
	orgDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganization, orgId)
	require.Nil(t, err)
	org := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	require.Equal(t, orgmodel.OnTrack.String(), org.OnboardingDetails.Status)
}

func TestOrganizationPlanEventHandler_OnReorderMilestones(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	timeNow := utils.Now()
	mpid := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "master plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{Name: "test org"})
	opid := neo4jtest.CreateOrganizationPlan(ctx, testDatabase.Driver, tenantName, mpid, orgId, neo4jentity.OrganizationPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "org plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		StatusDetails: neo4jentity.OrganizationPlanStatusDetails{
			Status:    model.NotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	milestoneId1 := neo4jtest.CreateOrganizationPlanMilestone(ctx, testDatabase.Driver, tenantName, opid, neo4jentity.OrganizationPlanMilestoneEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "milestone name 1",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		Order:         0,
		DueDate:       timeNow.Add(time.Hour * 24),
		Items:         []neo4jentity.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskNotDone.String(), UpdatedAt: timeNow, Uuid: "item1"}, {Text: "item2", Status: model.TaskNotDone.String(), UpdatedAt: timeNow, Uuid: "item2"}},
		Optional:      false,
		StatusDetails: neo4jentity.OrganizationPlanMilestoneStatusDetails{
			Status:    model.MilestoneNotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})
	milestoneId2 := neo4jtest.CreateOrganizationPlanMilestone(ctx, testDatabase.Driver, tenantName, opid, neo4jentity.OrganizationPlanMilestoneEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "milestone name 2",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		Order:         1,
		DueDate:       timeNow.Add(time.Hour * 24),
		Items:         []neo4jentity.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskNotDone.String(), UpdatedAt: timeNow, Uuid: "item1"}, {Text: "item2", Status: model.TaskNotDone.String(), UpdatedAt: timeNow, Uuid: "item2"}},
		Optional:      false,
		StatusDetails: neo4jentity.OrganizationPlanMilestoneStatusDetails{
			Status:    model.MilestoneNotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:          1,
		model2.NodeLabelOrganizationPlanMilestone: 2,
	})

	// Prepare the event handler
	orgPlanEventHandler := &OrganizationPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanMilestoneCreateEvent
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	reorderEvent, err := event.NewOrganizationPlanMilestoneReorderEvent(
		orgAggregate,
		opid,
		[]string{milestoneId2, milestoneId1},
		timeNow,
	)
	require.Nil(t, err)

	// EXECUTE
	err = orgPlanEventHandler.OnReorderMilestones(context.Background(), reorderEvent)
	require.Nil(t, err)

	// verify nodes and relationships
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:          1,
		model2.NodeLabelOrganizationPlanMilestone: 2})

	// verify master plan milestone nodes
	orgPlanMilestoneDbNode1, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlanMilestone, milestoneId1)
	require.Nil(t, err)
	require.NotNil(t, orgPlanMilestoneDbNode1)
	milestone1 := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(orgPlanMilestoneDbNode1)
	require.Equal(t, int64(1), milestone1.Order)
	require.Equal(t, model.MilestoneNotStarted.String(), milestone1.StatusDetails.Status) // don't update status automatically

	orgPlanMilestoneDbNode2, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlanMilestone, milestoneId2)
	require.Nil(t, err)
	require.NotNil(t, orgPlanMilestoneDbNode2)
	milestone2 := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(orgPlanMilestoneDbNode2)
	require.Equal(t, int64(0), milestone2.Order)
	require.Equal(t, model.MilestoneNotStarted.String(), milestone2.StatusDetails.Status) // don't update status automatically

	organizationPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlan, opid)
	require.Nil(t, err)
	op := neo4jmapper.MapDbNodeToOrganizationPlanEntity(organizationPlanDbNode)
	require.Equal(t, model.NotStarted.String(), op.StatusDetails.Status) // no automatic update
}

func TestOrganizationPlanEventHandler_OnUpdateMilestoneLate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	timeNow := utils.Now()
	mpid := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "master plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{Name: "test org"})
	opid := neo4jtest.CreateOrganizationPlan(ctx, testDatabase.Driver, tenantName, mpid, orgId, neo4jentity.OrganizationPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "org plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		StatusDetails: neo4jentity.OrganizationPlanStatusDetails{
			Status:    model.NotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	milestoneId := neo4jtest.CreateOrganizationPlanMilestone(ctx, testDatabase.Driver, tenantName, opid, neo4jentity.OrganizationPlanMilestoneEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "milestone name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		Order:         0,
		DueDate:       timeNow.Add(time.Hour * 24),
		Items:         []neo4jentity.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskNotDone.String(), UpdatedAt: timeNow, Uuid: "item1"}, {Text: "item2", Status: model.TaskNotDone.String(), UpdatedAt: timeNow, Uuid: "item2"}},
		Optional:      false,
		StatusDetails: neo4jentity.OrganizationPlanMilestoneStatusDetails{
			Status:    model.MilestoneNotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:                             1,
		model2.NodeLabelOrganizationPlanMilestone:                    1,
		model2.NodeLabelOrganizationPlanMilestone + "_" + tenantName: 1,
	})

	// Prepare the event handler
	orgPlanEventHandler := &OrganizationPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanMilestoneCreateEvent
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	lateUpdateTime := timeNow.Add(time.Hour * 48)
	updateEvent, err := event.NewOrganizationPlanMilestoneUpdateEvent(
		orgAggregate,
		opid,
		milestoneId,
		"new name",
		10,
		[]model.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskDone.String(), UpdatedAt: lateUpdateTime, Uuid: "item1"}, {Text: "item2Change", Status: model.TaskNotDone.String(), UpdatedAt: lateUpdateTime, Uuid: "item2"}},
		[]string{event.FieldMaskName, event.FieldMaskOptional, event.FieldMaskItems, event.FieldMaskOrder, event.FieldMaskStatusDetails},
		true,  // optional
		false, // adhoc
		true,  // retired
		lateUpdateTime,
		timeNow.Add(time.Hour*24), // due date
		model.OrganizationPlanDetails{
			Status:    model.MilestoneNotStarted.String(),
			Comments:  "comments",
			UpdatedAt: lateUpdateTime,
		},
	)
	require.Nil(t, err)

	// EXECUTE
	err = orgPlanEventHandler.OnUpdateMilestone(context.Background(), updateEvent)
	require.Nil(t, err)

	// verify nodes and relationships
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:          1,
		model2.NodeLabelOrganizationPlanMilestone: 1})

	// verify master plan milestone node
	orgPlanMilestoneDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlanMilestone, milestoneId)
	require.Nil(t, err)
	require.NotNil(t, orgPlanMilestoneDbNode)

	milestone := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(orgPlanMilestoneDbNode)
	require.Equal(t, milestoneId, milestone.Id)
	test.AssertRecentTime(t, milestone.UpdatedAt)
	require.Equal(t, "new name", milestone.Name)
	require.Equal(t, int64(10), milestone.Order)
	require.Equal(t, timeNow.Add(time.Hour*24), milestone.DueDate) // no due date update
	require.Equal(t, true, milestone.Optional)
	require.Equal(t, false, milestone.Retired)                                            // mask not passed so we ignore this field update
	require.Equal(t, model.MilestoneStartedLate.String(), milestone.StatusDetails.Status) // automatic update
	require.Equal(t, "comments", milestone.StatusDetails.Comments)
	require.Equal(t, lateUpdateTime, milestone.StatusDetails.UpdatedAt)
	for i, item := range milestone.Items {
		if i == 0 {
			require.Equal(t, model.TaskDone.String(), item.Status)
			require.Equal(t, "item1", item.Text)
			require.Equal(t, "item1", item.Uuid)
		} else {
			require.Equal(t, model.TaskNotDone.String(), item.Status)
			require.Equal(t, "item2Change", item.Text)
			require.Equal(t, "item2", item.Uuid)
		}
	}
	organizationPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlan, opid)
	require.Nil(t, err)
	op := neo4jmapper.MapDbNodeToOrganizationPlanEntity(organizationPlanDbNode)
	require.Equal(t, model.Late.String(), op.StatusDetails.Status) // automatic update

	// Check onboarding status updated
	orgDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganization, orgId)
	require.Nil(t, err)
	org := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	require.Equal(t, orgmodel.Late.String(), org.OnboardingDetails.Status)
}

func TestOrganizationPlanEventHandler_OnUpdateMilestoneAllDoneLate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	timeNow := utils.Now()
	mpid := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "master plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{Name: "test org"})
	opid := neo4jtest.CreateOrganizationPlan(ctx, testDatabase.Driver, tenantName, mpid, orgId, neo4jentity.OrganizationPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "org plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		StatusDetails: neo4jentity.OrganizationPlanStatusDetails{
			Status:    model.NotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	milestoneId := neo4jtest.CreateOrganizationPlanMilestone(ctx, testDatabase.Driver, tenantName, opid, neo4jentity.OrganizationPlanMilestoneEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "milestone name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		Order:         0,
		DueDate:       timeNow.Add(time.Hour * 24),
		Items:         []neo4jentity.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskNotDone.String(), UpdatedAt: timeNow, Uuid: "item1"}, {Text: "item2", Status: model.TaskNotDone.String(), UpdatedAt: timeNow, Uuid: "item2"}},
		Optional:      false,
		StatusDetails: neo4jentity.OrganizationPlanMilestoneStatusDetails{
			Status:    model.MilestoneNotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:                             1,
		model2.NodeLabelOrganizationPlanMilestone:                    1,
		model2.NodeLabelOrganizationPlanMilestone + "_" + tenantName: 1,
	})

	// Prepare the event handler
	orgPlanEventHandler := &OrganizationPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanMilestoneCreateEvent
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	lateUpdateTime := timeNow.Add(time.Hour * 48)
	updateEvent, err := event.NewOrganizationPlanMilestoneUpdateEvent(
		orgAggregate,
		opid,
		milestoneId,
		"new name",
		10,
		[]model.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskDone.String(), UpdatedAt: lateUpdateTime, Uuid: "item1"}, {Text: "item2Change", Status: model.TaskDoneLate.String(), UpdatedAt: lateUpdateTime, Uuid: "item2"}},
		[]string{event.FieldMaskName, event.FieldMaskOptional, event.FieldMaskItems, event.FieldMaskOrder, event.FieldMaskStatusDetails},
		true,  // optional
		false, // adhoc
		true,  // retired
		lateUpdateTime,
		timeNow.Add(time.Hour*24), // due date
		model.OrganizationPlanDetails{
			Status:    model.MilestoneNotStarted.String(),
			Comments:  "comments",
			UpdatedAt: lateUpdateTime,
		},
	)
	require.Nil(t, err)

	// EXECUTE
	err = orgPlanEventHandler.OnUpdateMilestone(context.Background(), updateEvent)
	require.Nil(t, err)

	// verify nodes and relationships
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:          1,
		model2.NodeLabelOrganizationPlanMilestone: 1})

	// verify master plan milestone node
	orgPlanMilestoneDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlanMilestone, milestoneId)
	require.Nil(t, err)
	require.NotNil(t, orgPlanMilestoneDbNode)

	milestone := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(orgPlanMilestoneDbNode)
	require.Equal(t, milestoneId, milestone.Id)
	test.AssertRecentTime(t, milestone.UpdatedAt)

	require.Equal(t, model.MilestoneDoneLate.String(), milestone.StatusDetails.Status) // automatic update
	require.Equal(t, "comments", milestone.StatusDetails.Comments)
	require.Equal(t, lateUpdateTime, milestone.StatusDetails.UpdatedAt)
	for i, item := range milestone.Items {
		if i == 0 {
			require.Equal(t, model.TaskDone.String(), item.Status)
			require.Equal(t, "item1", item.Text)
			require.Equal(t, "item1", item.Uuid)
		} else {
			require.Equal(t, model.TaskDoneLate.String(), item.Status)
			require.Equal(t, "item2Change", item.Text)
			require.Equal(t, "item2", item.Uuid)
		}
	}
	organizationPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlan, opid)
	require.Nil(t, err)
	op := neo4jmapper.MapDbNodeToOrganizationPlanEntity(organizationPlanDbNode)
	require.Equal(t, model.DoneLate.String(), op.StatusDetails.Status) // automatic update

	// Check onboarding status updated
	orgDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganization, orgId)
	require.Nil(t, err)
	org := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	require.Equal(t, orgmodel.Done.String(), org.OnboardingDetails.Status)
}

func TestOrganizationPlanEventHandler_OnUpdateMilestoneDueDateLate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	timeNow := utils.Now()
	mpid := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "master plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{Name: "test org"})
	opid := neo4jtest.CreateOrganizationPlan(ctx, testDatabase.Driver, tenantName, mpid, orgId, neo4jentity.OrganizationPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "org plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		StatusDetails: neo4jentity.OrganizationPlanStatusDetails{
			Status:    model.NotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	milestoneId := neo4jtest.CreateOrganizationPlanMilestone(ctx, testDatabase.Driver, tenantName, opid, neo4jentity.OrganizationPlanMilestoneEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "milestone name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		Order:         0,
		DueDate:       timeNow.Add(time.Hour * 24),
		Items:         []neo4jentity.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskDone.String(), UpdatedAt: timeNow, Uuid: "item1"}, {Text: "item2", Status: model.TaskNotDone.String(), UpdatedAt: timeNow, Uuid: "item2"}},
		Optional:      false,
		StatusDetails: neo4jentity.OrganizationPlanMilestoneStatusDetails{
			Status:    model.MilestoneStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:                             1,
		model2.NodeLabelOrganizationPlanMilestone:                    1,
		model2.NodeLabelOrganizationPlanMilestone + "_" + tenantName: 1,
	})

	// Prepare the event handler
	orgPlanEventHandler := &OrganizationPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanMilestoneCreateEvent
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	updateTime := utils.Now()
	updateEvent, err := event.NewOrganizationPlanMilestoneUpdateEvent(
		orgAggregate,
		opid,
		milestoneId,
		"new name",
		10,
		[]model.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskDone.String(), UpdatedAt: updateTime, Uuid: "item1"}, {Text: "item2Change", Status: model.TaskNotDone.String(), UpdatedAt: updateTime, Uuid: "item2"}},
		[]string{event.FieldMaskName, event.FieldMaskOptional, event.FieldMaskItems, event.FieldMaskDueDate, event.FieldMaskOrder, event.FieldMaskStatusDetails},
		true,  // optional
		false, // adhoc
		true,  // retired
		updateTime,
		timeNow.AddDate(0, 0, -2), // due date
		model.OrganizationPlanDetails{
			Status:    model.MilestoneStarted.String(),
			Comments:  "comments",
			UpdatedAt: updateTime,
		},
	)
	require.Nil(t, err)

	// EXECUTE
	err = orgPlanEventHandler.OnUpdateMilestone(context.Background(), updateEvent)
	require.Nil(t, err)

	// verify nodes and relationships
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:          1,
		model2.NodeLabelOrganizationPlanMilestone: 1})

	// verify master plan milestone node
	orgPlanMilestoneDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlanMilestone, milestoneId)
	require.Nil(t, err)
	require.NotNil(t, orgPlanMilestoneDbNode)

	milestone := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(orgPlanMilestoneDbNode)
	require.Equal(t, milestoneId, milestone.Id)
	test.AssertRecentTime(t, milestone.UpdatedAt)
	require.Equal(t, "new name", milestone.Name)
	require.Equal(t, int64(10), milestone.Order)
	require.Equal(t, timeNow.AddDate(0, 0, -2), milestone.DueDate)
	require.Equal(t, true, milestone.Optional)
	require.Equal(t, false, milestone.Retired)                                            // mask not passed so we ignore this field update
	require.Equal(t, model.MilestoneStartedLate.String(), milestone.StatusDetails.Status) // automatic update
	require.Equal(t, "comments", milestone.StatusDetails.Comments)
	require.Equal(t, updateTime, milestone.StatusDetails.UpdatedAt)
	for i, item := range milestone.Items {
		if i == 0 {
			require.Equal(t, model.TaskDoneLate.String(), item.Status)
			require.Equal(t, "item1", item.Text)
			require.Equal(t, "item1", item.Uuid)
		} else {
			require.Equal(t, model.TaskNotDoneLate.String(), item.Status)
			require.Equal(t, "item2Change", item.Text)
			require.Equal(t, "item2", item.Uuid)
		}
	}
	organizationPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlan, opid)
	require.Nil(t, err)
	op := neo4jmapper.MapDbNodeToOrganizationPlanEntity(organizationPlanDbNode)
	require.Equal(t, model.Late.String(), op.StatusDetails.Status) // automatic update

	// Check onboarding status updated
	orgDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganization, orgId)
	require.Nil(t, err)
	org := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	require.Equal(t, orgmodel.Late.String(), org.OnboardingDetails.Status)
}

func TestOrganizationPlanEventHandler_OnUpdateMilestoneDueDateOnTrack(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	timeNow := utils.Now()
	mpid := neo4jtest.CreateMasterPlan(ctx, testDatabase.Driver, tenantName, neo4jentity.MasterPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "master plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{Name: "test org"})
	opid := neo4jtest.CreateOrganizationPlan(ctx, testDatabase.Driver, tenantName, mpid, orgId, neo4jentity.OrganizationPlanEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "org plan name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		StatusDetails: neo4jentity.OrganizationPlanStatusDetails{
			Status:    model.NotStarted.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	milestoneId := neo4jtest.CreateOrganizationPlanMilestone(ctx, testDatabase.Driver, tenantName, opid, neo4jentity.OrganizationPlanMilestoneEntity{
		Source:        neo4jentity.DataSource(constants.SourceOpenline),
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		Name:          "milestone name",
		SourceOfTruth: neo4jentity.DataSource(constants.SourceOpenline),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		Retired:       false,
		Order:         0,
		DueDate:       timeNow.AddDate(0, 0, -2),
		Items:         []neo4jentity.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskDoneLate.String(), UpdatedAt: timeNow, Uuid: "item1"}, {Text: "item2", Status: model.TaskNotDoneLate.String(), UpdatedAt: timeNow, Uuid: "item2"}},
		Optional:      false,
		StatusDetails: neo4jentity.OrganizationPlanMilestoneStatusDetails{
			Status:    model.MilestoneStartedLate.String(),
			Comments:  "",
			UpdatedAt: timeNow,
		},
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:                             1,
		model2.NodeLabelOrganizationPlanMilestone:                    1,
		model2.NodeLabelOrganizationPlanMilestone + "_" + tenantName: 1,
	})

	// Prepare the event handler
	orgPlanEventHandler := &OrganizationPlanEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an MasterPlanMilestoneCreateEvent
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	updateTime := utils.Now()
	updateEvent, err := event.NewOrganizationPlanMilestoneUpdateEvent(
		orgAggregate,
		opid,
		milestoneId,
		"new name",
		10,
		[]model.OrganizationPlanMilestoneItem{{Text: "item1", Status: model.TaskDoneLate.String(), UpdatedAt: updateTime, Uuid: "item1"}, {Text: "item2Change", Status: model.TaskNotDoneLate.String(), UpdatedAt: updateTime, Uuid: "item2"}},
		[]string{event.FieldMaskName, event.FieldMaskOptional, event.FieldMaskItems, event.FieldMaskDueDate, event.FieldMaskOrder, event.FieldMaskStatusDetails},
		true,  // optional
		false, // adhoc
		true,  // retired
		updateTime,
		timeNow.Add(time.Hour*48), // due date in the future
		model.OrganizationPlanDetails{
			Status:    model.MilestoneStartedLate.String(),
			Comments:  "comments",
			UpdatedAt: updateTime,
		},
	)
	require.Nil(t, err)

	// EXECUTE
	err = orgPlanEventHandler.OnUpdateMilestone(context.Background(), updateEvent)
	require.Nil(t, err)

	// verify nodes and relationships
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganizationPlan:          1,
		model2.NodeLabelOrganizationPlanMilestone: 1})

	// verify master plan milestone node
	orgPlanMilestoneDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlanMilestone, milestoneId)
	require.Nil(t, err)
	require.NotNil(t, orgPlanMilestoneDbNode)

	milestone := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(orgPlanMilestoneDbNode)
	require.Equal(t, milestoneId, milestone.Id)
	test.AssertRecentTime(t, milestone.UpdatedAt)
	require.Equal(t, "new name", milestone.Name)
	require.Equal(t, int64(10), milestone.Order)
	require.Equal(t, timeNow.Add(time.Hour*48), milestone.DueDate)
	require.Equal(t, true, milestone.Optional)
	require.Equal(t, false, milestone.Retired)                                        // mask not passed so we ignore this field update
	require.Equal(t, model.MilestoneStarted.String(), milestone.StatusDetails.Status) // automatic update
	require.Equal(t, "comments", milestone.StatusDetails.Comments)
	require.Equal(t, updateTime, milestone.StatusDetails.UpdatedAt)
	for i, item := range milestone.Items {
		if i == 0 {
			require.Equal(t, model.TaskDone.String(), item.Status)
			require.Equal(t, "item1", item.Text)
			require.Equal(t, "item1", item.Uuid)
		} else {
			require.Equal(t, model.TaskNotDone.String(), item.Status)
			require.Equal(t, "item2Change", item.Text)
			require.Equal(t, "item2", item.Uuid)
		}
	}
	organizationPlanDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganizationPlan, opid)
	require.Nil(t, err)
	op := neo4jmapper.MapDbNodeToOrganizationPlanEntity(organizationPlanDbNode)
	require.Equal(t, model.OnTrack.String(), op.StatusDetails.Status) // automatic update

	// Check onboarding status updated
	orgDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelOrganization, orgId)
	require.Nil(t, err)
	org := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	require.Equal(t, orgmodel.OnTrack.String(), org.OnboardingDetails.Status)
}
