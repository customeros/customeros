package resolver

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	masterplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestMutationResolver_MasterPlanCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
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
