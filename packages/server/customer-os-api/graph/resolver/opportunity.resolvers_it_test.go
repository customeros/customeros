package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"testing"
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
