package subscriptions

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"golang.org/x/net/context"
	"time"
)

func WaitCheckNodeExistsInNeo4j(ctx context.Context, neo4jRepository *neo4jrepository.Repositories, tenant, id, nodeLabel string) bool {
	for i := 1; i <= constants.MaxRetriesCheckDataInNeo4j; i++ {
		found, findErr := neo4jRepository.CommonReadRepository.ExistsById(ctx, tenant, id, nodeLabel)
		if found && findErr == nil {
			return true
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}
	return false
}
