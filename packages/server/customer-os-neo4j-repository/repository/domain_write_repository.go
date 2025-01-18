package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type DomainWriteRepository interface {
	MergeDomain(ctx context.Context, tx *neo4j.ManagedTransaction, domain, source, appSource string) (bool, error)
	SetPrimaryDetails(ctx context.Context, domain, primaryDomain string, primary, accessible bool) error
}

type domainWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewDomainWriteRepository(driver *neo4j.DriverWithContext, database string) DomainWriteRepository {
	return &domainWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r domainWriteRepository) MergeDomain(ctx context.Context, tx *neo4j.ManagedTransaction, domain, source, appSource string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainWriteRepository.MergeDomain")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.SetTag(tracing.SpanTagEntityId, domain)

	cypher := fmt.Sprintf(`
	MERGE (d:Domain {domain:$domain})
	ON CREATE SET
		d.createdAt=datetime(),
		d.updatedAt=datetime(),
		d.source=$source,
		d.appSource=$appSource
	RETURN d.createdAt = datetime() AS justCreated`)

	params := map[string]interface{}{
		"domain":    domain,
		"source":    source,
		"appSource": appSource,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	span.LogFields(log.Bool("result.justCreated", result.(bool)))
	return result.(bool), nil
}

func (r domainWriteRepository) SetPrimaryDetails(ctx context.Context, domain, primaryDomain string, primary, accessible bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainWriteRepository.SetPrimaryDetails")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.SetTag(tracing.SpanTagEntityId, domain)
	span.LogFields(log.String("primaryDomain", primaryDomain), log.Bool("primary", primary))

	cypher := `MATCH (d:Domain {domain:$domain}) 
				SET d.primary=$primary, 
					d.accessible=$accessible, 
					d.primaryDomain=$primaryDomain, 
					d.techPrimaryDomainCheckRequestedAt=datetime()`

	params := map[string]interface{}{
		"domain":        domain,
		"primaryDomain": primaryDomain,
		"primary":       primary,
		"accessible":    accessible,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
