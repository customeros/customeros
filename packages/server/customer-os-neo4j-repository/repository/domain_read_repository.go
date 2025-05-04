package neo4j_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type DomainReadRepository interface {
	GetDomain(ctx context.Context, tx *neo4j.ManagedTransaction, domain string) (*dbtype.Node, error)
	GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	GetDomainsForPrimaryCheck(ctx context.Context, delayFromPreviousCheckInDays, limit int) ([]string, error)
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

func (r *domainReadRepository) GetDomain(ctx context.Context, tx *neo4j.ManagedTransaction, domain string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "DomainReadRepository.GetDomain")
	defer spans.Finish()

	cypher := `MATCH (d:Domain {domain:$domain}) RETURN d`
	params := map[string]any{
		"domain": domain,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if len(result.([]*dbtype.Node)) == 0 {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", true)
	return result.([]*dbtype.Node)[0], nil
}

func (r *domainReadRepository) GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "DomainRepository.GetForOrganizations")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_DOMAIN]->(d:Domain)
			WHERE o.id IN $organizationIds
			RETURN d, o.id ORDER BY d.domain ASC`
	params := map[string]any{
		"tenant":          tenant,
		"organizationIds": organizationIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteQuery(ctx, *r.driver, r.database, cypher, params, func(err error) {
		spans.TraceError(err)
	})
	if err != nil {
		return nil, err
	}
	return utils.ExtractAllRecordsAsDbNodeAndIdFromEagerResult(result), nil
}

func (r *domainReadRepository) GetDomainsForPrimaryCheck(ctx context.Context, delayFromPreviousCheckInDays, limit int) ([]string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "DomainReadRepository.GetDomainsForPrimaryCheck")
	defer spans.Finish()

	spans.LogKV("delayFromPreviousCheckInDays", delayFromPreviousCheckInDays)
	spans.LogKV("limit", limit)

	cypher := `MATCH (d:Domain)
		WHERE d.techPrimaryDomainCheckRequestedAt IS NULL OR d.techPrimaryDomainCheckRequestedAt < datetime() - duration({days:$delayFromPreviousCheckInDays})
		RETURN d.domain
		ORDER BY CASE WHEN d.techPrimaryDomainCheckRequestedAt IS NULL THEN 0 ELSE 1 END, d.techPrimaryDomainCheckRequestedAt ASC
				LIMIT $limit`
	params := map[string]any{
		"delayFromPreviousCheckInDays": delayFromPreviousCheckInDays,
		"limit":                        limit,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	domains, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	spans.LogKV("result.count", len(domains.([]string)))
	return domains.([]string), err
}
