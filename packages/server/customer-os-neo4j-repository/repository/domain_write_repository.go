package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type DomainWriteRepository interface {
	CreateDomain(ctx context.Context, domain, source, appSource string, now time.Time) error
}

type domainWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func (d domainWriteRepository) CreateDomain(ctx context.Context, domain, source, appSource string, time time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventWriteRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, domain)
	span.SetTag(tracing.SpanTagEntityId, domain)
	tracing.LogObjectAsJson(span, "data", domain)

	cypher := fmt.Sprintf(`
	MERGE (d:Domain {domain:$domain})
	ON CREATE SET
		d.createdAt=$createdAt,
		d.updatedAt=$updatedAt,
		d.source=$source,
		d.appSource=$appSource
	RETURN d
`)

	params := map[string]interface{}{
		"domain":    domain,
		"createdAt": time,
		"updatedAt": time,
		"source":    source,
		"appSource": appSource,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *d.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func NewDomainWriteRepository(driver *neo4j.DriverWithContext, database string) DomainWriteRepository {
	return &domainWriteRepository{
		driver:   driver,
		database: database,
	}
}
