package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type Neo4jRepository interface {
	ExecuteIntegrityCheckerQuery(ctx context.Context, name, query string) (int64, error)
}

type neo4jRepository struct {
	driver *neo4j.DriverWithContext
}

func NewNeo4jRepository(driver *neo4j.DriverWithContext) Neo4jRepository {
	return &neo4jRepository{
		driver: driver,
	}
}

func (r *neo4jRepository) ExecuteIntegrityCheckerQuery(ctx context.Context, name, query string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Neo4jRepository.ExecuteIntegrityCheckerQuery")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(span)
	span.SetTag("checker-name", name)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	countFoundRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{})
		return utils.ExtractSingleRecordFirstValueAsType[int64](ctx, queryResult, err)
	})
	span.LogFields(log.Int64("found records", countFoundRecords.(int64)))
	return countFoundRecords.(int64), err
}
