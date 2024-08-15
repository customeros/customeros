package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type BankAccountReadRepository interface {
	GetBankAccounts(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetBankAccountById(ctx context.Context, tenant, id string) (*dbtype.Node, error)
}

type bankAccountReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewBankAccountReadRepository(driver *neo4j.DriverWithContext, database string) BankAccountReadRepository {
	return &bankAccountReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *bankAccountReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *bankAccountReadRepository) GetBankAccounts(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountReadRepository.GetBankAccounts")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})-[:HAS_BANK_ACCOUNT]->(ba:BankAccount)
			RETURN ba ORDER BY ba.createdAt ASC`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})

	if err != nil {
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	if result == nil {
		return nil, nil
	}
	return result.([]*dbtype.Node), nil
}

func (r *bankAccountReadRepository) GetBankAccountById(ctx context.Context, tenant, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountReadRepository.GetBankAccountById")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})-[:HAS_BANK_ACCOUNT]->(ba:BankAccount {id:$id})
			RETURN ba`
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})

	if err != nil {
		tracing.TraceErr(span, err)
		span.LogFields(log.Bool("result.found", false))
		return nil, err
	}

	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}
