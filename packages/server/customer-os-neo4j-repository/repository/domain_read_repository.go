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
	GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
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
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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

func (r *domainReadRepository) GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainRepository.GetForOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_DOMAIN]->(d:Domain)
			WHERE o.id IN $organizationIds
			RETURN d, o.id ORDER BY d.domain ASC`
	params := map[string]any{
		"tenant":          tenant,
		"organizationIds": organizationIds,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteQuery(ctx, *r.driver, r.database, cypher, params, func(err error) {
		tracing.TraceErr(span, err)
	})
	if err != nil {
		return nil, err
	}
	return utils.ExtractAllRecordsAsDbNodeAndIdFromEagerResult(result), nil
}
