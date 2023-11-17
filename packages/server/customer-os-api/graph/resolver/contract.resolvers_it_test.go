package resolver

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contract"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestMutationResolver_ContractCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	contractId := uuid.New().String()
	calledCreateContract := false

	contractServiceCallbacks := events_platform.MockContractServiceCallbacks{
		CreateContract: func(context context.Context, contract *contractpb.CreateContractGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
			require.Equal(t, tenantName, contract.Tenant)
			require.Equal(t, orgId, contract.OrganizationId)
			require.Equal(t, testUserId, contract.LoggedInUserId)
			require.Equal(t, string(entity.DataSourceOpenline), contract.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, contract.SourceFields.AppSource)
			require.Equal(t, "Contract 1", contract.Name)
			require.Equal(t, "https://contract.com", contract.ContractUrl)
			require.Equal(t, contractpb.RenewalCycle_MONTHLY_RENEWAL, contract.RenewalCycle)
			expectedServiceStartedAt, err := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
			if err != nil {
				t.Fatalf("Failed to parse expected timestamp: %v", err)
			}
			require.Equal(t, timestamppb.New(expectedServiceStartedAt), contract.ServiceStartedAt)
			expectedSignedAt, err := time.Parse(time.RFC3339, "2019-02-01T00:00:00Z")
			if err != nil {
				t.Fatalf("Failed to parse expected timestamp: %v", err)
			}
			require.Equal(t, timestamppb.New(expectedSignedAt), contract.SignedAt)
			calledCreateContract = true
			neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
				Id: contractId,
			})
			return &contractpb.ContractIdGrpcResponse{
				Id: contractId,
			}, nil
		},
	}
	events_platform.SetContractCallbacks(&contractServiceCallbacks)

	rawResponse := callGraphQL(t, "contract/create_contract", map[string]interface{}{
		"orgId": orgId,
	})

	var contractStruct struct {
		Contract_Create model.Contract
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &contractStruct)
	require.Nil(t, err)
	contract := contractStruct.Contract_Create
	require.Equal(t, contractId, contract.ID)

	require.True(t, calledCreateContract)
}

func TestMutationResolver_ContractUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})
	calledUpdateContract := false

	contractServiceCallbacks := events_platform.MockContractServiceCallbacks{
		UpdateContract: func(context context.Context, contract *contractpb.UpdateContractGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
			require.Equal(t, tenantName, contract.Tenant)
			require.Equal(t, contractId, contract.Id)
			require.Equal(t, testUserId, contract.LoggedInUserId)
			require.Equal(t, string(entity.DataSourceOpenline), contract.SourceFields.Source)
			require.Equal(t, "test app source", contract.SourceFields.AppSource)
			require.Equal(t, "Updated Contract", contract.Name)
			require.Equal(t, "https://contract.com/updated", contract.ContractUrl)
			require.Equal(t, contractpb.RenewalCycle_ANNUALLY_RENEWAL, contract.RenewalCycle)
			expectedServiceStartedAt, err := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
			if err != nil {
				t.Fatalf("Failed to parse expected timestamp: %v", err)
			}
			require.Equal(t, timestamppb.New(expectedServiceStartedAt), contract.ServiceStartedAt)
			expectedSignedAt, err := time.Parse(time.RFC3339, "2019-02-01T00:00:00Z")
			if err != nil {
				t.Fatalf("Failed to parse expected timestamp: %v", err)
			}
			require.Equal(t, timestamppb.New(expectedSignedAt), contract.SignedAt)
			expectedEndedAt, err := time.Parse(time.RFC3339, "2019-03-01T00:00:00Z")
			if err != nil {
				t.Fatalf("Failed to parse expected timestamp: %v", err)
			}
			require.Equal(t, timestamppb.New(expectedEndedAt), contract.EndedAt)
			calledUpdateContract = true
			return &contractpb.ContractIdGrpcResponse{
				Id: contractId,
			}, nil
		},
	}
	events_platform.SetContractCallbacks(&contractServiceCallbacks)

	rawResponse := callGraphQL(t, "contract/update_contract", map[string]interface{}{
		"contractId": contractId,
	})

	var contractStruct struct {
		Contract_Update model.Contract
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &contractStruct)
	require.Nil(t, err)
	contract := contractStruct.Contract_Update
	require.Equal(t, contractId, contract.ID)

	require.True(t, calledUpdateContract)
}
