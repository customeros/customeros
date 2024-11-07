package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_OpportunityForContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	creatorUserId := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{})
	ownerUserId := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	now := utils.Now()
	opportunityId := neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		Name:      "test opportunity",
		Amount:    float64(100),
		MaxAmount: float64(200),
		CreatedAt: now,
		UpdatedAt: now,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalApproved:     true,
			RenewalAdjustedRate: 33,
		},
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
	require.Equal(t, opportunityId, opportunity.Metadata.ID)
	require.Equal(t, "test opportunity", opportunity.Name)
	require.Equal(t, float64(100), opportunity.Amount)
	require.Equal(t, float64(200), opportunity.MaxAmount)
	require.Equal(t, creatorUserId, opportunity.CreatedBy.ID)
	require.Equal(t, ownerUserId, opportunity.Owner.ID)
	require.Equal(t, now, opportunity.Metadata.Created)
	require.Equal(t, now, opportunity.Metadata.LastUpdated)
	require.True(t, opportunity.RenewalApproved)
	require.Equal(t, int64(33), opportunity.RenewalAdjustedRate)
	require.Nil(t, opportunity.Organization)
}

func TestQueryResolver_OpportunityForOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	opportunityId := neo4jtest.CreateOpportunityForOrganization(ctx, driver, tenantName, orgId, neo4jentity.OpportunityEntity{
		Name:          "test opportunity",
		CreatedAt:     utils.Now(),
		UpdatedAt:     utils.Now(),
		InternalType:  neo4jenum.OpportunityInternalTypeNBO,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		ExternalType:  "external type",
		ExternalStage: "external stage",
	})

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
	require.Equal(t, opportunityId, opportunity.Metadata.ID)
	require.Equal(t, "test opportunity", opportunity.Name)
	require.Equal(t, model.InternalTypeNbo, opportunity.InternalType)
	require.Equal(t, model.InternalStageOpen, opportunity.InternalStage)
	require.Equal(t, "external type", opportunity.ExternalType)
	require.Equal(t, "external stage", opportunity.ExternalStage)
	require.NotNil(t, opportunity.Organization)
	require.Equal(t, orgId, opportunity.Organization.Metadata.ID)
}

func TestMutationResolver_OpportunityRenewalUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	ownerUserId := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{})
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	opportunityId := neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{})
	calledUpdateRenewalOpportunity := false

	opportunityServiceCallbacks := events_platform.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunity: func(context context.Context, renewalOpportunity *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, renewalOpportunity.Tenant)
			require.Equal(t, opportunityId, renewalOpportunity.Id)
			require.Equal(t, testUserId, renewalOpportunity.LoggedInUserId)
			require.Equal(t, string(neo4jentity.DataSourceOpenline), renewalOpportunity.SourceFields.Source)
			require.Equal(t, "test app source", renewalOpportunity.SourceFields.AppSource)
			require.Equal(t, float64(100), renewalOpportunity.Amount)
			require.Equal(t, opportunitypb.RenewalLikelihood_HIGH_RENEWAL, renewalOpportunity.RenewalLikelihood)
			require.Equal(t, "test comments", renewalOpportunity.Comments)
			require.Equal(t, ownerUserId, renewalOpportunity.OwnerUserId)
			require.Equal(t, int64(50), renewalOpportunity.RenewalAdjustedRate)
			require.ElementsMatch(t, []opportunitypb.OpportunityMaskField{
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_COMMENTS,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_ADJUSTED_RATE,
			},
				renewalOpportunity.FieldsMask)
			calledUpdateRenewalOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	events_platform.SetOpportunityCallbacks(&opportunityServiceCallbacks)
	neo4jt.OpportunityOwnedBy(ctx, driver, opportunityId, ownerUserId)

	rawResponse := callGraphQL(t, "opportunity/update_renewal_opportunity", map[string]interface{}{
		"opportunityId": opportunityId,
		"ownerUserId":   ownerUserId,
	})

	var opportunityRenewalUpdateStruct struct {
		OpportunityRenewalUpdate model.Opportunity
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &opportunityRenewalUpdateStruct)
	require.Nil(t, err)
	renewalOpportunity := opportunityRenewalUpdateStruct.OpportunityRenewalUpdate
	require.Equal(t, opportunityId, renewalOpportunity.ID)
	require.Equal(t, ownerUserId, renewalOpportunity.Owner.ID)

	require.True(t, calledUpdateRenewalOpportunity)
}

func TestMutationResolver_OpportunityRenewalUpdateAllForOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId1 := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})
	contractId2 := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.ContractEntity{})
	opportunityId1 := neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId1, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	opportunityId2 := neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId2, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodZero,
		},
	})
	calledUpdateRenewalOpportunityCounter := 0
	opportunityServiceCallbacks := events_platform.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunity: func(context context.Context, renewalOpportunity *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, renewalOpportunity.Tenant)
			if renewalOpportunity.Id == opportunityId1 {
				require.Equal(t, opportunityId1, renewalOpportunity.Id)
			} else {
				require.Equal(t, opportunityId2, renewalOpportunity.Id)
			}
			require.Equal(t, string(neo4jentity.DataSourceOpenline), renewalOpportunity.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, renewalOpportunity.SourceFields.AppSource)
			require.Equal(t, opportunitypb.RenewalLikelihood_MEDIUM_RENEWAL, renewalOpportunity.RenewalLikelihood)
			require.Equal(t, "", renewalOpportunity.OwnerUserId)
			require.Equal(t, int64(50), renewalOpportunity.RenewalAdjustedRate)
			require.ElementsMatch(t, []opportunitypb.OpportunityMaskField{
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_ADJUSTED_RATE},
				renewalOpportunity.FieldsMask)
			calledUpdateRenewalOpportunityCounter++
			return &opportunitypb.OpportunityIdGrpcResponse{}, nil
		},
	}
	events_platform.SetOpportunityCallbacks(&opportunityServiceCallbacks)

	rawResponse := callGraphQL(t, "opportunity/update_renewal_opportunities_for_organization", map[string]interface{}{
		"organizationId": organizationId,
	})

	var opportunityRenewalUpdateStruct struct {
		OpportunityRenewal_UpdateAllForOrganization model.Organization
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &opportunityRenewalUpdateStruct)
	require.Nil(t, err)
	organization := opportunityRenewalUpdateStruct.OpportunityRenewal_UpdateAllForOrganization
	require.Equal(t, organizationId, organization.Metadata.ID)
	require.Equal(t, 2, calledUpdateRenewalOpportunityCounter)
}
