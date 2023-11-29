package resolver

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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

func TestQueryResolver_Contract_WithServiceLineItems(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	yesterday := now.Add(time.Duration(-24) * time.Hour)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})

	serviceLineItemId1 := neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:      "service line item 1",
		CreatedAt: yesterday,
		UpdatedAt: yesterday,
		Billed:    entity.BilledTypeAnnually,
		Price:     13,
		Quantity:  2,
		Source:    entity.DataSourceOpenline,
		AppSource: "test1",
	})
	serviceLineItemId2 := neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Name:      "service line item 2",
		CreatedAt: now,
		UpdatedAt: now,
		Billed:    entity.BilledTypeUsage,
		Price:     255,
		Quantity:  23,
		Source:    entity.DataSourceOpenline,
		AppSource: "test2",
	})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Organization":    1,
		"Contract":        1,
		"ServiceLineItem": 2,
	})
	assertRelationship(ctx, t, driver, contractId, "HAS_SERVICE", serviceLineItemId1)
	assertRelationship(ctx, t, driver, contractId, "HAS_SERVICE", serviceLineItemId2)

	rawResponse := callGraphQL(t, "contract/get_contract_with_service_line_items",
		map[string]interface{}{"contractId": contractId})

	var contractStruct struct {
		Contract model.Contract
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contractStruct)
	require.Nil(t, err)

	contract := contractStruct.Contract
	require.NotNil(t, contract)
	require.Equal(t, 2, len(contract.ServiceLineItems))

	firstServiceLineItem := contract.ServiceLineItems[0]
	require.Equal(t, serviceLineItemId1, firstServiceLineItem.ID)
	require.Equal(t, "service line item 1", firstServiceLineItem.Name)
	require.Equal(t, yesterday, firstServiceLineItem.CreatedAt)
	require.Equal(t, yesterday, firstServiceLineItem.UpdatedAt)
	require.Equal(t, model.BilledTypeAnnually, firstServiceLineItem.Billed)
	require.Equal(t, float64(13), firstServiceLineItem.Price)
	require.Equal(t, int64(2), firstServiceLineItem.Quantity)
	require.Equal(t, model.DataSourceOpenline, firstServiceLineItem.Source)
	require.Equal(t, "test1", firstServiceLineItem.AppSource)

	secondServiceLineItem := contract.ServiceLineItems[1]
	require.Equal(t, serviceLineItemId2, secondServiceLineItem.ID)
	require.Equal(t, "service line item 2", secondServiceLineItem.Name)
	require.Equal(t, now, secondServiceLineItem.CreatedAt)
	require.Equal(t, now, secondServiceLineItem.UpdatedAt)
	require.Equal(t, model.BilledTypeUsage, secondServiceLineItem.Billed)
	require.Equal(t, float64(255), secondServiceLineItem.Price)
	require.Equal(t, int64(23), secondServiceLineItem.Quantity)
	require.Equal(t, model.DataSourceOpenline, secondServiceLineItem.Source)
	require.Equal(t, "test2", secondServiceLineItem.AppSource)
}

func TestQueryResolver_Contract_WithOpportunities(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	yesterday := now.Add(time.Duration(-24) * time.Hour)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})

	opportunityId1 := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:          "oppo 1",
		CreatedAt:     now,
		UpdatedAt:     now,
		Amount:        49,
		InternalType:  entity.InternalTypeUpsell,
		InternalStage: entity.InternalStageOpen,
		Source:        entity.DataSourceOpenline,
		GeneralNotes:  "test notes 1",
		Comments:      "test comments 1",
		AppSource:     "test1",
	})
	opportunityId2 := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:          "oppo 2",
		CreatedAt:     yesterday,
		UpdatedAt:     yesterday,
		Amount:        1239,
		InternalType:  entity.InternalTypeNbo,
		InternalStage: entity.InternalStageEvaluating,
		Source:        entity.DataSourceOpenline,
		GeneralNotes:  "test notes 2",
		Comments:      "test comments 2",
		AppSource:     "test2",
	})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Organization": 1,
		"Contract":     1,
		"Opportunity":  2,
	})
	assertRelationship(ctx, t, driver, contractId, "HAS_OPPORTUNITY", opportunityId1)

	rawResponse := callGraphQL(t, "contract/get_contract_with_opportunities",
		map[string]interface{}{"contractId": contractId})

	var contractStruct struct {
		Contract model.Contract
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contractStruct)
	require.Nil(t, err)

	contract := contractStruct.Contract
	require.NotNil(t, contract)
	require.Equal(t, 2, len(contract.Opportunities))

	firstOpportunity := contract.Opportunities[0]
	require.Equal(t, opportunityId1, firstOpportunity.ID)
	require.Equal(t, "oppo 1", firstOpportunity.Name)
	require.Equal(t, now, firstOpportunity.CreatedAt)
	require.Equal(t, now, firstOpportunity.UpdatedAt)
	require.Equal(t, float64(49), firstOpportunity.Amount)
	require.Equal(t, model.InternalStageOpen, firstOpportunity.InternalStage)
	require.Equal(t, model.InternalTypeUpsell, firstOpportunity.InternalType)
	require.Equal(t, model.DataSourceOpenline, firstOpportunity.Source)
	require.Equal(t, "test notes 1", firstOpportunity.GeneralNotes)
	require.Equal(t, "test comments 1", firstOpportunity.Comments)
	require.Equal(t, "test1", firstOpportunity.AppSource)

	secondOpportunity := contract.Opportunities[1]
	require.Equal(t, opportunityId2, secondOpportunity.ID)
	require.Equal(t, "oppo 2", secondOpportunity.Name)
	require.Equal(t, yesterday, secondOpportunity.CreatedAt)
	require.Equal(t, yesterday, secondOpportunity.UpdatedAt)
	require.Equal(t, float64(1239), secondOpportunity.Amount)
	require.Equal(t, model.InternalStageEvaluating, secondOpportunity.InternalStage)
	require.Equal(t, model.InternalTypeNbo, secondOpportunity.InternalType)
	require.Equal(t, model.DataSourceOpenline, secondOpportunity.Source)
	require.Equal(t, "test notes 2", secondOpportunity.GeneralNotes)
	require.Equal(t, "test comments 2", secondOpportunity.Comments)
	require.Equal(t, "test2", secondOpportunity.AppSource)
}
