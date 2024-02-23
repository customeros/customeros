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
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestMutationResolver_ContractCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := uuid.New().String()
	calledCreateContract := false

	contractServiceCallbacks := events_platform.MockContractServiceCallbacks{
		CreateContract: func(context context.Context, contract *contractpb.CreateContractGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
			require.Equal(t, tenantName, contract.Tenant)
			require.Equal(t, orgId, contract.OrganizationId)
			require.Equal(t, testUserId, contract.LoggedInUserId)
			require.Equal(t, string(neo4jentity.DataSourceOpenline), contract.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, contract.SourceFields.AppSource)
			require.Equal(t, "Contract 1", contract.Name)
			require.Equal(t, "https://contract.com", contract.ContractUrl)
			require.Equal(t, contractpb.RenewalCycle_MONTHLY_RENEWAL, contract.RenewalCycle)
			require.Equal(t, "USD", contract.Currency)
			require.Equal(t, int64(7), *contract.RenewalPeriods)
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
			neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
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

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	calledUpdateContract := false

	contractServiceCallbacks := events_platform.MockContractServiceCallbacks{
		UpdateContract: func(context context.Context, contract *contractpb.UpdateContractGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
			require.Equal(t, tenantName, contract.Tenant)
			require.Equal(t, contractId, contract.Id)
			require.Equal(t, testUserId, contract.LoggedInUserId)
			require.Equal(t, string(neo4jentity.DataSourceOpenline), contract.SourceFields.Source)
			require.Equal(t, "test app source", contract.SourceFields.AppSource)
			require.Equal(t, "Updated Contract", contract.Name)
			require.Equal(t, "https://contract.com/updated", contract.ContractUrl)
			require.Equal(t, contractpb.RenewalCycle_ANNUALLY_RENEWAL, contract.RenewalCycle)
			require.Equal(t, int64(3), *contract.RenewalPeriods)
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
			require.Equal(t, commonpb.BillingCycle_ANNUALLY_BILLING, contract.BillingCycle)
			require.Equal(t, "USD", contract.Currency)
			require.Equal(t, "test address line 1", contract.AddressLine1)
			require.Equal(t, "test address line 2", contract.AddressLine2)
			require.Equal(t, "test locality", contract.Locality)
			require.Equal(t, "test country", contract.Country)
			require.Equal(t, "test zip", contract.Zip)
			require.Equal(t, "test organization legal name", contract.OrganizationLegalName)
			require.Equal(t, "test invoice email", contract.InvoiceEmail)
			require.Equal(t, "test invoice note", contract.InvoiceNote)
			require.Equal(t, true, contract.CanPayWithCard)
			require.Equal(t, true, contract.CanPayWithDirectDebit)
			require.Equal(t, true, contract.CanPayWithBankTransfer)
			require.Equal(t, 20, len(contract.FieldsMask))
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

func TestMutationResolver_ContractUpdate_NullDateFields(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	calledUpdateContract := false

	contractServiceCallbacks := events_platform.MockContractServiceCallbacks{
		UpdateContract: func(context context.Context, contract *contractpb.UpdateContractGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
			require.Equal(t, tenantName, contract.Tenant)
			require.Equal(t, contractId, contract.Id)
			require.Equal(t, testUserId, contract.LoggedInUserId)
			require.Equal(t, string(neo4jentity.DataSourceOpenline), contract.SourceFields.Source)
			require.Equal(t, "customer-os-api", contract.SourceFields.AppSource)

			require.Nil(t, contract.SignedAt)
			require.Nil(t, contract.ServiceStartedAt)
			require.Nil(t, contract.EndedAt)
			require.Nil(t, contract.InvoicingStartDate)

			require.Equal(t, 4, len(contract.FieldsMask))
			calledUpdateContract = true
			return &contractpb.ContractIdGrpcResponse{
				Id: contractId,
			}, nil
		},
	}
	events_platform.SetContractCallbacks(&contractServiceCallbacks)

	rawResponse := callGraphQL(t, "contract/update_contract_null_dates", map[string]interface{}{
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
	require.Nil(t, contract.SignedAt)
	require.Nil(t, contract.ServiceStartedAt)
	require.Nil(t, contract.EndedAt)
	require.Nil(t, contract.InvoicingStartDate)

	require.True(t, calledUpdateContract)
}

func TestQueryResolver_Contract_WithServiceLineItems(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	yesterday := now.Add(time.Duration(-24) * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		AddressLine1:           "address line 1",
		AddressLine2:           "address line 2",
		Zip:                    "zip",
		Locality:               "locality",
		Country:                "country",
		OrganizationLegalName:  "organization legal name",
		InvoiceEmail:           "invoice email",
		InvoiceNote:            "invoice note",
		BillingCycle:           neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:     &now,
		InvoicingEnabled:       true,
		CanPayWithCard:         true,
		CanPayWithDirectDebit:  true,
		CanPayWithBankTransfer: true,
	})

	serviceLineItemId1 := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:      "service line item 1",
		CreatedAt: yesterday,
		UpdatedAt: yesterday,
		Billed:    neo4jenum.BilledTypeAnnually,
		Price:     13,
		Quantity:  2,
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: "test1",
		VatRate:   0.1,
	})
	serviceLineItemId2 := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:      "service line item 2",
		CreatedAt: now,
		UpdatedAt: now,
		Billed:    neo4jenum.BilledTypeUsage,
		Price:     255,
		Quantity:  23,
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: "test2",
		VatRate:   0.2,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Organization":    1,
		"Contract":        1,
		"ServiceLineItem": 2,
	})
	neo4jtest.AssertRelationship(ctx, t, driver, contractId, "HAS_SERVICE", serviceLineItemId1)
	neo4jtest.AssertRelationship(ctx, t, driver, contractId, "HAS_SERVICE", serviceLineItemId2)

	rawResponse := callGraphQL(t, "contract/get_contract_with_service_line_items",
		map[string]interface{}{"contractId": contractId})

	var contractStruct struct {
		Contract model.Contract
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contractStruct)
	require.Nil(t, err)

	contract := contractStruct.Contract
	require.NotNil(t, contract)
	require.Equal(t, contractId, contract.Metadata.ID)
	require.True(t, contract.BillingEnabled)

	billingDetails := contract.BillingDetails
	require.Equal(t, "address line 1", *billingDetails.AddressLine1)
	require.Equal(t, "address line 2", *billingDetails.AddressLine2)
	require.Equal(t, "zip", *billingDetails.PostalCode)
	require.Equal(t, "locality", *billingDetails.Locality)
	require.Equal(t, "country", *billingDetails.Country)
	require.Equal(t, "organization legal name", *billingDetails.OrganizationLegalName)
	require.Equal(t, "invoice email", *billingDetails.BillingEmail)
	require.Equal(t, "invoice note", *billingDetails.InvoiceNote)
	require.Equal(t, model.ContractBillingCycleMonthlyBilling, *billingDetails.BillingCycle)
	require.True(t, *billingDetails.CanPayWithCard)
	require.True(t, *billingDetails.CanPayWithDirectDebit)
	require.True(t, *billingDetails.CanPayWithBankTransfer)

	require.Equal(t, 2, len(contract.ContractLineItems))

	firstContractLineItem := contract.ContractLineItems[0]
	require.Equal(t, serviceLineItemId1, firstContractLineItem.ID)
	require.Equal(t, "service line item 1", firstContractLineItem.Name)
	require.Equal(t, yesterday, firstContractLineItem.CreatedAt)
	require.Equal(t, yesterday, firstContractLineItem.UpdatedAt)
	require.Equal(t, model.BilledTypeAnnually, firstContractLineItem.Billed)
	require.Equal(t, float64(13), firstContractLineItem.Price)
	require.Equal(t, int64(2), firstContractLineItem.Quantity)
	require.Equal(t, model.DataSourceOpenline, firstContractLineItem.Source)
	require.Equal(t, "test1", firstContractLineItem.AppSource)
	require.Equal(t, 0.1, firstContractLineItem.VatRate)

	secondContractLineItem := contract.ContractLineItems[1]
	require.Equal(t, serviceLineItemId2, secondContractLineItem.ID)
	require.Equal(t, "service line item 2", secondContractLineItem.Name)
	require.Equal(t, now, secondContractLineItem.CreatedAt)
	require.Equal(t, now, secondContractLineItem.UpdatedAt)
	require.Equal(t, model.BilledTypeUsage, secondContractLineItem.Billed)
	require.Equal(t, float64(255), secondContractLineItem.Price)
	require.Equal(t, int64(23), secondContractLineItem.Quantity)
	require.Equal(t, model.DataSourceOpenline, secondContractLineItem.Source)
	require.Equal(t, "test2", secondContractLineItem.AppSource)
	require.Equal(t, 0.2, secondContractLineItem.VatRate)
}

func TestQueryResolver_Contract_WithOpportunities(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	yesterday := now.Add(time.Duration(-24) * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})

	opportunityId1 := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:          "oppo 1",
		CreatedAt:     now,
		UpdatedAt:     now,
		Amount:        49,
		InternalType:  entity.InternalTypeUpsell,
		InternalStage: entity.InternalStageOpen,
		Source:        neo4jentity.DataSourceOpenline,
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
		Source:        neo4jentity.DataSourceOpenline,
		GeneralNotes:  "test notes 2",
		Comments:      "test comments 2",
		AppSource:     "test2",
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Organization": 1,
		"Contract":     1,
		"Opportunity":  2,
	})
	neo4jtest.AssertRelationship(ctx, t, driver, contractId, "HAS_OPPORTUNITY", opportunityId1)

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
