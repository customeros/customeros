package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
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
