package repository

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
)

func ContractWasInserted(ctx context.Context, contractId, organizationId string) {
	cid := test.CreateContract(ctx, driver, tenantName, organizationId, entity.ContractEntity{
		Id: contractId,
	})
	fmt.Sprintf(cid)
}
