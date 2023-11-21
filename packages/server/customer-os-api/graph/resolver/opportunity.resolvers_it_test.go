package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestQueryResolver_Opportunity(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	creatorUserId := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{})
	ownerUserId := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{})
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})
	now := utils.Now()
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:      "test opportunity",
		Amount:    float64(100),
		MaxAmount: float64(200),
		CreatedAt: now,
		UpdatedAt: now,
	})
	neo4jt.OpportunityCreatedBy(ctx, driver, opportunityId, creatorUserId)
	neo4jt.OpportunityOwnedBy(ctx, driver, opportunityId, ownerUserId)

	rawResponse := callGraphQL(t, "opportunity/get_opportunity", map[string]interface{}{
		"opportunityId": opportunityId,
	})

	var opportunityStruct struct {
		Opportunity model.Opportunity
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &opportunityStruct)
	require.Nil(t, err)
	opportunity := opportunityStruct.Opportunity

	require.NotNil(t, opportunity)
	require.Equal(t, opportunityId, opportunity.ID)
	require.Equal(t, "test opportunity", opportunity.Name)
	require.Equal(t, float64(100), opportunity.Amount)
	require.Equal(t, float64(200), opportunity.MaxAmount)
	require.Equal(t, creatorUserId, opportunity.CreatedBy.ID)
	require.Equal(t, ownerUserId, opportunity.Owner.ID)
	require.Equal(t, now, opportunity.CreatedAt)
	require.Equal(t, now, opportunity.UpdatedAt)
}

func TestMutationResolver_OpportunityUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{})
	calledUpdateOpportunity := false

	opportunityServiceCallbacks := events_platform.MockOpportunityServiceCallbacks{
		UpdateOpportunity: func(context context.Context, opportunity *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, opportunity.Tenant)
			require.Equal(t, opportunityId, opportunity.Id)
			require.Equal(t, testUserId, opportunity.LoggedInUserId)
			require.Equal(t, string(entity.DataSourceOpenline), opportunity.SourceFields.Source)
			require.Equal(t, "test app source", opportunity.SourceFields.AppSource)
			require.Equal(t, "Updated Opportunity", opportunity.Name)
			require.Equal(t, float64(100), opportunity.Amount)
			require.Equal(t, "external type", opportunity.ExternalType)
			require.Equal(t, "external stage", opportunity.ExternalStage)
			require.Equal(t, "general notes", opportunity.GeneralNotes)
			require.Equal(t, "next steps", opportunity.NextSteps)
			estimatedCloseAt, err := time.Parse(time.RFC3339, "2019-03-01T00:00:00Z")
			if err != nil {
				t.Fatalf("Failed to parse expected timestamp: %v", err)
			}
			require.Equal(t, timestamppb.New(estimatedCloseAt), opportunity.EstimatedCloseDate)
			calledUpdateOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	events_platform.SetOpportunityCallbacks(&opportunityServiceCallbacks)

	rawResponse := callGraphQL(t, "opportunity/update_opportunity", map[string]interface{}{
		"opportunityId": opportunityId,
	})

	var opportunityStruct struct {
		OpportunityUpdate model.Opportunity
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &opportunityStruct)
	require.Nil(t, err)
	opportunity := opportunityStruct.OpportunityUpdate
	require.Equal(t, opportunityId, opportunity.ID)

	require.True(t, calledUpdateOpportunity)
}

func TestMutationResolver_OpportunityRenewalUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{})
	calledUpdateRenewalOpportunity := false

	opportunityServiceCallbacks := events_platform.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunity: func(context context.Context, renewalOpportunity *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, renewalOpportunity.Tenant)
			require.Equal(t, opportunityId, renewalOpportunity.Id)
			require.Equal(t, testUserId, renewalOpportunity.LoggedInUserId)
			require.Equal(t, string(entity.DataSourceOpenline), renewalOpportunity.SourceFields.Source)
			require.Equal(t, "test app source", renewalOpportunity.SourceFields.AppSource)
			require.Equal(t, float64(100), renewalOpportunity.Amount)
			require.Equal(t, opportunitypb.RenewalLikelihood_HIGH_RENEWAL, renewalOpportunity.RenewalLikelihood)
			require.Equal(t, "test comments", renewalOpportunity.Comments)
			calledUpdateRenewalOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	events_platform.SetOpportunityCallbacks(&opportunityServiceCallbacks)

	rawResponse := callGraphQL(t, "update_renewal_opportunity", map[string]interface{}{
		"opportunityId": opportunityId,
	})

	var opportunityRenewalUpdateStruct struct {
		OpportunityRenewalUpdate model.Opportunity
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &opportunityRenewalUpdateStruct)
	require.Nil(t, err)
	renewalOpportunity := opportunityRenewalUpdateStruct.OpportunityRenewalUpdate
	require.Equal(t, opportunityId, renewalOpportunity.ID)

	require.True(t, calledUpdateRenewalOpportunity)
}
