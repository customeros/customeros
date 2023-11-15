package resolver

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	contractgrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contract"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_ContractCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org")

	createdContractId := uuid.New().String()

	calledCreateContract := false

	contractServiceCallbacks := events_platform.MockContractServiceCallbacks{
		CreateContract: func(context context.Context, contact *contractgrpc.CreateContractGrpcRequest) (*contractgrpc.ContractIdGrpcResponse, error) {
			require.Equal(t, "name", contact.Name)
			require.Equal(t, "MONTHLY_RENEWAL", contact.RenewalCycle)
			require.Equal(t, organizationId, organizationId)
			require.Equal(t, "", contact.SignedAt)
			require.Equal(t, "", contact.SignedAt)
			require.Equal(t, tenantName, contact.Tenant)
			require.Equal(t, string(entity.DataSourceOpenline), contact.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, contact.SourceFields.AppSource)
			require.Equal(t, testUserId, contact.LoggedInUserId)
			calledCreateContract = true
			neo4jt.CreateContract(ctx, driver, tenantName, entity.ContractEntity{
				ID: createdContractId,
			})
			return &contractgrpc.ContractIdGrpcResponse{
				Id: createdContractId,
			}, nil
		},
	}
	events_platform.SetContractCallbacks(&contractServiceCallbacks)

	rawResponse := callGraphQL(t, "contract/create_contract", map[string]interface{}{})

	var contractStruct struct {
		ContractCreate model.Contract
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &contractStruct)
	require.Nil(t, err)
	contract := contractStruct.ContractCreate
	require.Equal(t, createdContractId, contract.ID)
	require.True(t, calledCreateContract)
	//TODO: assert contract fields

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Contract": 1,
	})
}
