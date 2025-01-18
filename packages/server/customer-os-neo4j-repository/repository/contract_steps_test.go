package repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
)

func ContractWasInserted(ctx context.Context, contractId, organizationId string) {
	cid := test.CreateContractForOrganization(ctx, driver, tenantName, organizationId, entity.ContractEntity{
		Id: contractId,
	})
	fmt.Sprintf(cid)
}
