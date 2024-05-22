package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type DomainReadRepository interface {
	GetDomain(ctx context.Context, domain string) (*dbtype.Node, error)
}

type domainReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewDomainReadRepository(driver *neo4j.DriverWithContext, database string) DomainReadRepository {
	return &domainReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *domainReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *domainReadRepository) GetDomain(ctx context.Context, domain string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainReadRepository.GetDomain")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "neo4jRepository")

	cypher := fmt.Sprintf(`MATCH (d:Domain {domain:$domain}) RETURN d`)
	params := map[string]any{
		"domain": domain,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
