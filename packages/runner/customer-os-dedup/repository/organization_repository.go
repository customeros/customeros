package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type OrganizationRepository interface {
	SuggestOrganizationsMerge(ctx context.Context, tenant, primaryOrgId, secondaryOrgId, suggestedBy, suggestedByDtls string, confidence float64) error
	OrgsAlreadyComparedForDuplicates(ctx context.Context, tenant, org1Id, org2Id string) (bool, error)
	GetOrganizationsForDedupComparison(ctx context.Context, tenant string, limit int) ([]*dbtype.Node, error)
	ExistsNewOrganizationsCreatedAfter(ctx context.Context, tenant string, after time.Time) (bool, error)
}

type organizationRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext, database string) OrganizationRepository {
	return &organizationRepository{
		driver:   driver,
		database: database,
	}
}

func (r *organizationRepository) SuggestOrganizationsMerge(ctx context.Context, tenant, primaryOrgId, secondaryOrgId, suggestedBy, suggestedByDtls string, confidence float64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetTenantsWithOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(span)

	query := `MATCH (t:Tenant {name:$tenant}),
				(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(primary:Organization {id:$primaryOrgId}),
				(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(secondary:Organization {id:$secondaryOrgId})
				WHERE NOT (secondary)-[:SUGGESTED_MERGE]-(primary)
				MERGE (secondary)-[rel:SUGGESTED_MERGE]->(primary)
				SET rel.suggestedBy = $suggestedBy,
					rel.suggestedByInfo = $suggestedByDtls,
					rel.confidence = $confidence, 
					rel.suggestedAt = $now`
	span.LogFields(log.String("query", query))

	_, err := utils.ExecuteQuery(ctx, *r.driver, r.database, query, map[string]any{
		"tenant":          tenant,
		"now":             utils.Now(),
		"primaryOrgId":    primaryOrgId,
		"secondaryOrgId":  secondaryOrgId,
		"suggestedBy":     suggestedBy,
		"suggestedByDtls": suggestedByDtls,
		"confidence":      confidence,
	})
	return err
}

func (r *organizationRepository) OrgsAlreadyComparedForDuplicates(ctx context.Context, tenant, org1Id, org2Id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.OrgsComparedForDuplicates")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(span)

	query := `MATCH (t:Tenant {name:$tenant}),
				(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(first:Organization {id:$org1Id}),
				(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(second:Organization {id:$org2Id}),
				(first)-[rel:SUGGESTED_MERGE]-(second)
				RETURN count(rel) > 0 as result`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant": tenant,
			"org1Id": org1Id,
			"org2Id": org2Id,
		})
		return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, queryResult, err)
	})
	span.LogFields(log.Bool("result", result.(bool)))
	return result.(bool), err
}

func (r *organizationRepository) GetOrganizationsCountForDedupComparison(ctx context.Context, tenant string, limit int) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(span)

	query := `MATCH (t:Tenant {name:'openline'})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization) 
				WHERE NOT (org)-[:SUGGESTED_MERGE]->(:Organization)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(t)
				WITH org, rand() AS r ORDER BY r 
				RETURN org LIMIT $limit`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant": tenant,
			"limit":  limit,
		})
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.([]*dbtype.Node), nil
}

func (r *organizationRepository) GetOrganizationsForDedupComparison(ctx context.Context, tenant string, limit int) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization) 
				WHERE NOT (org)-[:SUGGESTED_MERGE]->(:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
				AND org.name <> '' 
				AND NOT org.name IS NULL
				WITH org, rand() AS r ORDER BY r 
				RETURN org LIMIT $limit`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant": tenant,
			"limit":  limit,
		})
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.([]*dbtype.Node), nil
}

func (r *organizationRepository) ExistsNewOrganizationsCreatedAfter(ctx context.Context, tenant string, after time.Time) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.ExistsNewOrganizationsCreatedAfter")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)
				WHERE org.createdAt > $after
				RETURN count(org) > 0 as result`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant": tenant,
			"after":  after,
		})
		return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, queryResult, err)
	})
	span.LogFields(log.Bool("result", result.(bool)))
	return result.(bool), err
}
