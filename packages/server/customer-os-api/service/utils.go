package service

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"time"
)

func WaitForObjectCreationAndLogSpan(ctx context.Context, s *repository.Repositories, id, nodeLabel string, span opentracing.Span) {
	for i := 1; i <= constants.MaxRetriesCheckDataInNeo4jAfterEventRequest; i++ {
		found, findErr := s.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), id, nodeLabel)
		if found && findErr == nil {
			span.LogFields(log.Bool(fmt.Sprintf("response - %s saved in db", nodeLabel), true))
			break
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}

	span.LogFields(log.String(fmt.Sprintf("response - created %s with id", nodeLabel), id))
}
