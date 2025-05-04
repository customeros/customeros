package neo4j_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "DomainWriteRepository.MergeDomain")
	defer spans.Finish()

	spans.TagEntity(domain)

	cypher := `
	MERGE (d:Domain {domain:$domain})
	ON CREATE SET
		d.createdAt=datetime(),
		d.updatedAt=datetime(),
		d.source=$source,
		d.appSource=$appSource
	RETURN d.createdAt = datetime() AS justCreated`

	params := map[string]interface{}{
		"domain":    domain,
		"source":    source,
		"appSource": appSource,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return false, err
	}
	spans.LogKV("result.justCreated", result.(bool))
	return result.(bool), nil
}

func (r domainWriteRepository) SetPrimaryDetails(ctx context.Context, domain, primaryDomain string, primary, accessible bool) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "DomainWriteRepository.SetPrimaryDetails")
	defer spans.Finish()

	spans.TagEntity(domain)
	spans.LogKV("primaryDomain", primaryDomain)
	spans.LogKV("primary", primary)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}
