package neo4j_repository

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

func WaitForNodeCreatedInNeo4j(ctx context.Context, repositories *Repositories, id, nodeLabel string, span opentracing.Span) {
	operation := func() error {
		found, findErr := repositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), id, nodeLabel)
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

func WaitForNodeCreatedInNeo4jWithConfig(ctx context.Context, span opentracing.Span, repositories *Repositories, id, nodeLabel string, maxWaitTime time.Duration) {
	operation := func() error {
		found, findErr := repositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), id, nodeLabel)
		if findErr != nil {
			return findErr
		}
		if !found {
			return errors.New(fmt.Sprintf("Node %s with id %s not found in Neo4j", nodeLabel, id))
		}
		return nil
	}

	err := backoff.Retry(operation, utils.BackOffConfig(250*time.Millisecond, 1, 500*time.Millisecond, maxWaitTime, 50))
	if err != nil {
		span.LogFields(log.Bool("result.created", false))
	} else {
		span.LogFields(log.Bool("result.created", true))
	}
}

func WaitForNodeDeletedFromNeo4j(ctx context.Context, repositories *Repositories, id, nodeLabel string, span opentracing.Span) {
	operation := func() error {
		found, findErr := repositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), id, nodeLabel)
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

func LogAndExecuteWriteQuery(ctx context.Context, driver neo4j.DriverWithContext, cypher string, params map[string]any, span opentracing.Span) error {
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func LogAndExecuteWriteQueryInTx(ctx context.Context, tx *neo4j.ManagedTransaction, driver *neo4j.DriverWithContext, database, cypher string, params map[string]any, span opentracing.Span) error {
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, driver, database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
