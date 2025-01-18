package neo4j_repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

type IndustryWriteRepository interface {
	ReplaceForOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, industryCode string) error
}

type industryWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewIndustryWriteRepository(driver *neo4j.DriverWithContext, database string) IndustryWriteRepository {
	return &industryWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *industryWriteRepository) ReplaceForOrganization(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId, industryCode string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IndustryWriteRepository.ReplaceForOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (o:Organization {id: $organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name: $tenant})
				MATCH (i:Industry {code: $industryCode})
				OPTIONAL MATCH (o)-[oldRel:HAS_INDUSTRY]->(other:Industry)
					WHERE other <> i
					DELETE oldRel
				MERGE (o)-[newRel:HAS_INDUSTRY]->(i)
					SET o.updatedAt = datetime()`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"industryCode":   industryCode,
	}

	return LogAndExecuteWriteQueryInTx(ctx, tx, r.driver, r.database, cypher, params, span)
}
