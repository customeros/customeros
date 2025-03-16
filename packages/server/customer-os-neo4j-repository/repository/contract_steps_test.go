package neo4j_repository

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
)

func ContractWasInserted(ctx context.Context, contractId, organizationId string) string {
	cid := test.CreateContractForOrganization(ctx, driver, tenantName, organizationId, neo4j_entity.ContractEntity{
		Id: contractId,
	})
	return cid
}
