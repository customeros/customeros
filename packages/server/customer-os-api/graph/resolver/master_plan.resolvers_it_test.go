package resolver

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	masterplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestMutationResolver_MasterPlanCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	masterPlanId := uuid.New().String()
	calledCreateMasterPlan := false

	masterPlanServiceCallbacks := events_platform.MockMasterPlanServiceCallbacks{
		CreateMasterPlan: func(context context.Context, masterPlan *masterplanpb.CreateMasterPlanGrpcRequest) (*masterplanpb.MasterPlanIdGrpcResponse, error) {
			require.Equal(t, tenantName, masterPlan.Tenant)
			require.Equal(t, testUserId, masterPlan.LoggedInUserId)
			require.Equal(t, neo4jentity.DataSourceOpenline.String(), masterPlan.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, masterPlan.SourceFields.AppSource)
			require.Equal(t, "Draft plan", masterPlan.Name)
			calledCreateMasterPlan = true
			neo4jtest.CreateMasterPlan(ctx, driver, tenantName, neo4jentity.MasterPlanEntity{Id: masterPlanId})
			return &masterplanpb.MasterPlanIdGrpcResponse{
				Id: masterPlanId,
			}, nil
		},
	}
	events_platform.SetMasterPlanCallbacks(&masterPlanServiceCallbacks)

	rawResponse := callGraphQL(t, "master_plan/create_master_plan", map[string]interface{}{})
	require.Nil(t, rawResponse.Errors)

	var masterPlanStruct struct {
		MasterPlan_Create model.MasterPlan
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &masterPlanStruct)
	require.Nil(t, err)

	masterPlan := masterPlanStruct.MasterPlan_Create
	require.Equal(t, masterPlanId, masterPlan.ID)

	require.True(t, calledCreateMasterPlan)
}

func TestMutationResolver_MasterPlanMilestoneCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	masterPlanId := neo4jtest.CreateMasterPlan(ctx, driver, tenantName, neo4jentity.MasterPlanEntity{})
	masterPlanMilestoneId := uuid.New().String()

	calledCreateMasterPlanMilestone := false

	masterPlanServiceCallbacks := events_platform.MockMasterPlanServiceCallbacks{
		CreateMasterPlanMilestone: func(context context.Context, masterPlanMilestone *masterplanpb.CreateMasterPlanMilestoneGrpcRequest) (*masterplanpb.MasterPlanMilestoneIdGrpcResponse, error) {
			require.Equal(t, tenantName, masterPlanMilestone.Tenant)
			require.Equal(t, testUserId, masterPlanMilestone.LoggedInUserId)
			require.Equal(t, neo4jentity.DataSourceOpenline.String(), masterPlanMilestone.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, masterPlanMilestone.SourceFields.AppSource)
			require.Equal(t, "Draft milestone", masterPlanMilestone.Name)
			require.Equal(t, int64(10), masterPlanMilestone.Order)
			require.Equal(t, int64(48), masterPlanMilestone.DurationHours)
			require.Equal(t, []string{"do A", "do B", "do C"}, masterPlanMilestone.Items)
			require.Equal(t, true, masterPlanMilestone.Optional)
			calledCreateMasterPlanMilestone = true
			neo4jtest.CreateMasterPlanMilestone(ctx, driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{Id: masterPlanMilestoneId})
			return &masterplanpb.MasterPlanMilestoneIdGrpcResponse{
				Id: masterPlanMilestoneId,
			}, nil
		},
	}
	events_platform.SetMasterPlanCallbacks(&masterPlanServiceCallbacks)

	rawResponse := callGraphQL(t, "master_plan/create_master_plan_milestone", map[string]interface{}{"masterPlanId": masterPlanId})
	require.Nil(t, rawResponse.Errors)

	var masterPlanMilestoneStruct struct {
		MasterPlanMilestone_Create model.MasterPlanMilestone
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &masterPlanMilestoneStruct)
	require.Nil(t, err)

	masterPlanMilestone := masterPlanMilestoneStruct.MasterPlanMilestone_Create
	require.Equal(t, masterPlanMilestoneId, masterPlanMilestone.ID)

	require.True(t, calledCreateMasterPlanMilestone)
}

func TestQueryResolver_MasterPlan(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	timeNow := utils.Now()
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	masterPlanId := neo4jtest.CreateMasterPlan(ctx, driver, tenantName, neo4jentity.MasterPlanEntity{
		Name:      "Master plan 1",
		CreatedAt: timeNow,
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: "test",
	})

	rawResponse := callGraphQL(t, "master_plan/get_master_plan", map[string]interface{}{"id": masterPlanId})
	require.Nil(t, rawResponse.Errors)

	var masterPlanStruct struct {
		MasterPlan model.MasterPlan
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &masterPlanStruct)
	require.Nil(t, err)

	masterPlan := masterPlanStruct.MasterPlan
	require.Equal(t, masterPlanId, masterPlan.ID)
	require.Equal(t, "Master plan 1", masterPlan.Name)
	require.Equal(t, timeNow, masterPlan.CreatedAt)
	require.Equal(t, model.DataSourceOpenline, masterPlan.Source)
	require.Equal(t, "test", masterPlan.AppSource)
	require.False(t, masterPlan.Retired)
}

func TestQueryResolver_MasterPlan_WithMilestones(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	timeNow := utils.Now()
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	masterPlanId := neo4jtest.CreateMasterPlan(ctx, driver, tenantName, neo4jentity.MasterPlanEntity{
		CreatedAt: timeNow,
	})
	neo4jtest.CreateMasterPlanMilestone(ctx, driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{
		Name:          "Optional second",
		CreatedAt:     timeNow,
		Order:         2,
		Optional:      true,
		DurationHours: 24,
		Items:         []string{"A", "B"},
	})
	neo4jtest.CreateMasterPlanMilestone(ctx, driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{
		Name:          "Optional first",
		CreatedAt:     timeNow,
		Order:         1,
		Optional:      true,
		DurationHours: 48,
		Items:         []string{"C"},
	})
	neo4jtest.CreateMasterPlanMilestone(ctx, driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{
		Name:          "Retired optional",
		Retired:       true,
		CreatedAt:     timeNow,
		Order:         1,
		Optional:      true,
		DurationHours: 12,
		Items:         []string{"F", "G"},
	})
	neo4jtest.CreateMasterPlanMilestone(ctx, driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{
		Name:          "Default",
		CreatedAt:     timeNow,
		Order:         10,
		Optional:      false,
		DurationHours: 10,
		Items:         []string{"D"},
	})
	neo4jtest.CreateMasterPlanMilestone(ctx, driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{
		Name:          "Retired default",
		Retired:       true,
		CreatedAt:     timeNow,
		Order:         1,
		Optional:      false,
		DurationHours: 50,
		Items:         []string{"E"},
	})

	rawResponse := callGraphQL(t, "master_plan/get_master_plan", map[string]interface{}{"id": masterPlanId})
	require.Nil(t, rawResponse.Errors)

	var masterPlanStruct struct {
		MasterPlan model.MasterPlan
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &masterPlanStruct)
	require.Nil(t, err)

	masterPlan := masterPlanStruct.MasterPlan
	require.Equal(t, masterPlanId, masterPlan.ID)
	require.Equal(t, 3, len(masterPlan.Milestones))
	require.Equal(t, 2, len(masterPlan.RetiredMilestones))

	require.Equal(t, "Default", masterPlan.Milestones[0].Name)
	require.Equal(t, []string{"D"}, masterPlan.Milestones[0].Items)
	require.Equal(t, false, masterPlan.Milestones[0].Optional)
	require.Equal(t, int64(10), masterPlan.Milestones[0].DurationHours)

	require.Equal(t, "Optional first", masterPlan.Milestones[1].Name)
	require.Equal(t, []string{"C"}, masterPlan.Milestones[1].Items)
	require.Equal(t, true, masterPlan.Milestones[1].Optional)
	require.Equal(t, int64(48), masterPlan.Milestones[1].DurationHours)

	require.Equal(t, "Optional second", masterPlan.Milestones[2].Name)
	require.Equal(t, []string{"A", "B"}, masterPlan.Milestones[2].Items)
	require.Equal(t, true, masterPlan.Milestones[2].Optional)
	require.Equal(t, int64(24), masterPlan.Milestones[2].DurationHours)

	require.Equal(t, "Retired default", masterPlan.RetiredMilestones[0].Name)
	require.Equal(t, []string{"E"}, masterPlan.RetiredMilestones[0].Items)
	require.Equal(t, false, masterPlan.RetiredMilestones[0].Optional)
	require.Equal(t, int64(50), masterPlan.RetiredMilestones[0].DurationHours)

	require.Equal(t, "Retired optional", masterPlan.RetiredMilestones[1].Name)
	require.Equal(t, []string{"F", "G"}, masterPlan.RetiredMilestones[1].Items)
	require.Equal(t, true, masterPlan.RetiredMilestones[1].Optional)
	require.Equal(t, int64(12), masterPlan.RetiredMilestones[1].DurationHours)
}

func TestQueryResolver_MasterPlans(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	today := utils.Now()
	yesterday := today.Add(-24 * time.Hour)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	masterPlanId_today := neo4jtest.CreateMasterPlan(ctx, driver, tenantName, neo4jentity.MasterPlanEntity{
		Name:      "Today plan",
		CreatedAt: today,
	})
	masterPlanId_yday := neo4jtest.CreateMasterPlan(ctx, driver, tenantName, neo4jentity.MasterPlanEntity{
		Name:      "Yesterday plan",
		CreatedAt: yesterday,
	})

	rawResponse := callGraphQL(t, "master_plan/list_master_plans", map[string]interface{}{})
	require.Nil(t, rawResponse.Errors)

	var masterPlansStruct struct {
		MasterPlans []model.MasterPlan
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &masterPlansStruct)
	require.Nil(t, err)

	masterPlans := masterPlansStruct.MasterPlans
	require.Equal(t, masterPlanId_yday, masterPlans[0].ID)
	require.Equal(t, "Yesterday plan", masterPlans[0].Name)
	require.Equal(t, yesterday, masterPlans[0].CreatedAt)
	require.False(t, masterPlans[0].Retired)
	require.Equal(t, masterPlanId_today, masterPlans[1].ID)
	require.Equal(t, "Today plan", masterPlans[1].Name)
	require.Equal(t, today, masterPlans[1].CreatedAt)
	require.False(t, masterPlans[1].Retired)
}

func TestQueryResolver_MasterPlans_OnlyNonRetired(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	today := utils.Now()
	yesterday := today.Add(-24 * time.Hour)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	masterPlanId_today := neo4jtest.CreateMasterPlan(ctx, driver, tenantName, neo4jentity.MasterPlanEntity{
		Name:      "Today plan",
		CreatedAt: today,
	})
	neo4jtest.CreateMasterPlan(ctx, driver, tenantName, neo4jentity.MasterPlanEntity{
		Name:      "Yesterday plan",
		CreatedAt: yesterday,
		Retired:   true,
	})

	rawResponse := callGraphQL(t, "master_plan/list_master_plans_active", map[string]interface{}{})
	require.Nil(t, rawResponse.Errors)

	var masterPlansStruct struct {
		MasterPlans []model.MasterPlan
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &masterPlansStruct)
	require.Nil(t, err)

	masterPlans := masterPlansStruct.MasterPlans
	require.Equal(t, 1, len(masterPlans))
	require.Equal(t, masterPlanId_today, masterPlans[0].ID)
}

func TestMutationResolver_MasterPlanUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	masterPlanId := neo4jtest.CreateMasterPlan(ctx, driver, tenantName, neo4jentity.MasterPlanEntity{})
	calledUpdateMasterPlan := false

	masterPlanServiceCallbacks := events_platform.MockMasterPlanServiceCallbacks{
		UpdateMasterPlan: func(context context.Context, masterPlan *masterplanpb.UpdateMasterPlanGrpcRequest) (*masterplanpb.MasterPlanIdGrpcResponse, error) {
			require.Equal(t, tenantName, masterPlan.Tenant)
			require.Equal(t, masterPlanId, masterPlan.MasterPlanId)
			require.Equal(t, testUserId, masterPlan.LoggedInUserId)
			require.Equal(t, constants.AppSourceCustomerOsApi, masterPlan.AppSource)
			require.Equal(t, "Updated plan", masterPlan.Name)
			require.Equal(t, false, masterPlan.Retired)
			require.Equal(t, []masterplanpb.MasterPlanFieldMask{masterplanpb.MasterPlanFieldMask_MASTER_PLAN_PROPERTY_NAME}, masterPlan.FieldsMask)
			calledUpdateMasterPlan = true
			return &masterplanpb.MasterPlanIdGrpcResponse{
				Id: masterPlanId,
			}, nil
		},
	}
	events_platform.SetMasterPlanCallbacks(&masterPlanServiceCallbacks)

	rawResponse := callGraphQL(t, "master_plan/update_master_plan", map[string]interface{}{"id": masterPlanId})
	require.Nil(t, rawResponse.Errors)

	var masterPlanStruct struct {
		MasterPlan_Update model.MasterPlan
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &masterPlanStruct)
	require.Nil(t, err)

	masterPlan := masterPlanStruct.MasterPlan_Update
	require.Equal(t, masterPlanId, masterPlan.ID)

	require.True(t, calledUpdateMasterPlan)
}

func TestMutationResolver_MasterPlanMilestoneUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	masterPlanId := neo4jtest.CreateMasterPlan(ctx, driver, tenantName, neo4jentity.MasterPlanEntity{})
	milestoneId := neo4jtest.CreateMasterPlanMilestone(ctx, driver, tenantName, masterPlanId, neo4jentity.MasterPlanMilestoneEntity{})
	calledUpdateMasterPlanMilestone := false

	masterPlanServiceCallbacks := events_platform.MockMasterPlanServiceCallbacks{
		UpdateMasterPlanMilestone: func(context context.Context, milestone *masterplanpb.UpdateMasterPlanMilestoneGrpcRequest) (*masterplanpb.MasterPlanMilestoneIdGrpcResponse, error) {
			require.Equal(t, tenantName, milestone.Tenant)
			require.Equal(t, masterPlanId, milestone.MasterPlanId)
			require.Equal(t, milestoneId, milestone.MasterPlanMilestoneId)
			require.Equal(t, testUserId, milestone.LoggedInUserId)
			require.Equal(t, constants.AppSourceCustomerOsApi, milestone.AppSource)
			require.Equal(t, "Updated milestone", milestone.Name)
			require.Equal(t, true, milestone.Optional)
			require.Equal(t, false, milestone.Retired)
			require.Nil(t, milestone.Items)
			require.Equal(t, int64(10), milestone.Order)
			require.Equal(t, int64(0), milestone.DurationHours)
			require.ElementsMatch(t, []masterplanpb.MasterPlanMilestoneFieldMask{
				masterplanpb.MasterPlanMilestoneFieldMask_MASTER_PLAN_MILESTONE_PROPERTY_NAME,
				masterplanpb.MasterPlanMilestoneFieldMask_MASTER_PLAN_MILESTONE_PROPERTY_OPTIONAL,
				masterplanpb.MasterPlanMilestoneFieldMask_MASTER_PLAN_MILESTONE_PROPERTY_ORDER,
			}, milestone.FieldsMask)
			calledUpdateMasterPlanMilestone = true
			return &masterplanpb.MasterPlanMilestoneIdGrpcResponse{
				Id: masterPlanId,
			}, nil
		},
	}
	events_platform.SetMasterPlanCallbacks(&masterPlanServiceCallbacks)

	rawResponse := callGraphQL(t, "master_plan/update_master_plan_milestone", map[string]interface{}{"masterPlanId": masterPlanId, "id": milestoneId})
	require.Nil(t, rawResponse.Errors)

	var masterPlanMilestoneStruct struct {
		MasterPlanMilestone_Update model.MasterPlanMilestone
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &masterPlanMilestoneStruct)
	require.Nil(t, err)

	masterPlanMilestone := masterPlanMilestoneStruct.MasterPlanMilestone_Update
	require.Equal(t, milestoneId, masterPlanMilestone.ID)

	require.True(t, calledUpdateMasterPlanMilestone)
}

func TestMutationResolver_MasterPlanMilestoneReorder(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	masterPlanId := neo4jtest.CreateMasterPlan(ctx, driver, tenantName, neo4jentity.MasterPlanEntity{})
	calledReorderMasterPlanMilestones := false
	milestoneIds := []string{"1", "2", "3"}

	masterPlanServiceCallbacks := events_platform.MockMasterPlanServiceCallbacks{
		ReorderMasterPlanMilestones: func(context context.Context, request *masterplanpb.ReorderMasterPlanMilestonesGrpcRequest) (*masterplanpb.MasterPlanIdGrpcResponse, error) {
			require.Equal(t, tenantName, request.Tenant)
			require.Equal(t, masterPlanId, request.MasterPlanId)
			require.Equal(t, testUserId, request.LoggedInUserId)
			require.Equal(t, milestoneIds, request.MasterPlanMilestoneIds)
			require.Equal(t, constants.AppSourceCustomerOsApi, request.AppSource)
			calledReorderMasterPlanMilestones = true
			return &masterplanpb.MasterPlanIdGrpcResponse{
				Id: masterPlanId,
			}, nil
		},
	}
	events_platform.SetMasterPlanCallbacks(&masterPlanServiceCallbacks)

	rawResponse := callGraphQL(t, "master_plan/reorder_master_plan_milestones", map[string]interface{}{"masterPlanId": masterPlanId, "milestoneIds": milestoneIds})
	require.Nil(t, rawResponse.Errors)

	// GraphQL response
	var masterPlanStruct struct {
		MasterPlanMilestone_Reorder string
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &masterPlanStruct)
	require.Nil(t, err)

	// Verify
	require.Equal(t, masterPlanId, masterPlanStruct.MasterPlanMilestone_Reorder)

	require.True(t, calledReorderMasterPlanMilestones)
}
