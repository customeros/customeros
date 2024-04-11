package service

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

func WaitForNodeCreatedInNeo4j(ctx context.Context, repositories *repository.Repositories, id, nodeLabel string, span opentracing.Span) {
	operation := func() error {
		found, findErr := repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), id, nodeLabel)
		if findErr != nil {
			return findErr
		}
		if !found {
			return errors.New(fmt.Sprintf("Node %s with id %s not found in Neo4j", nodeLabel, id))
		}
		return nil
	}

	err := backoff.Retry(operation, utils.BackOffConfig(100*time.Millisecond, 1.5, 1*time.Second, 5*time.Second, 10))
	if err != nil {
		span.LogFields(log.Bool("result.created", false))
	} else {
		span.LogFields(log.Bool("result.created", true))
	}
}

func WaitForNodeDeletedFromNeo4j(ctx context.Context, repositories *repository.Repositories, id, nodeLabel string, span opentracing.Span) {
	operation := func() error {
		found, findErr := repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), id, nodeLabel)
		if findErr != nil {
			return findErr
		}
		if found {
			return errors.New(fmt.Sprintf("Node %s with id %s still exists in Neo4j", nodeLabel, id))
		}
		return nil
	}

	err := backoff.Retry(operation, utils.BackOffConfig(100*time.Millisecond, 1.5, 1*time.Second, 5*time.Second, 10))
	if err != nil {
		span.LogFields(log.Bool("result.deleted", false))
	} else {
		span.LogFields(log.Bool("result.deleted", true))
	}
}
