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

type TenantRepository interface {
	GetTenantsWithOrganizations(ctx context.Context, atLeastOrganizationsForTenant int) ([]string, error)
	GetTenantMetadata(ctx context.Context, tenantName string) (*dbtype.Node, error)
	UpdateTenantMetadataOrgDedupAt(ctx context.Context, tenantName string, time time.Time) error
}

type tenantRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewTenantRepository(driver *neo4j.DriverWithContext, database string) TenantRepository {
	return &tenantRepository{
		driver:   driver,
		database: database,
	}
}

func (r *tenantRepository) GetTenantsWithOrganizations(ctx context.Context, atLeastOrganizationsForTenant int) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.GetTenantsWithOrganizations")
	defer span.Finish()

	tracing.SetDefaultNeo4jRepositorySpanTags(span)

	query := `MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization) 
				WITH t, count(o) as orgsCount 
				WHERE orgsCount >= $limit
				RETURN t.name order by orgsCount asc`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"limit": atLeastOrganizationsForTenant,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	span.LogFields(log.Int("records", len(records.([]string))))
	return records.([]string), err
}

func (r *tenantRepository) GetTenantMetadata(ctx context.Context, tenantName string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.GetTenantMetadata")
	defer span.Finish()
	span.LogFields(log.String("tenantName", tenantName))

	tracing.SetDefaultNeo4jRepositorySpanTags(span)

	query := `MATCH (t:Tenant {name:$tenantName})
				MERGE (t)-[:HAS_METADATA]->(tm:TenantMetadata {tenantName:$tenantName}) RETURN tm`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenantName": tenantName,
		})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return records.(*dbtype.Node), err
}

func (r *tenantRepository) UpdateTenantMetadataOrgDedupAt(ctx context.Context, tenantName string, time time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.UpdateTenantMetadataOrgDedupAt")
	defer span.Finish()
	span.LogFields(log.String("tenantName", tenantName))

	tracing.SetDefaultNeo4jRepositorySpanTags(span)

	query := `MATCH (t:Tenant {name:$tenantName})
				MERGE (t)-[:HAS_METADATA]->(tm:TenantMetadata {tenantName:$tenantName})
				SET tm.orgDedupAt = $orgDedupAt`
	span.LogFields(log.String("query", query))

	_, err := utils.ExecuteQuery(ctx, *r.driver, r.database, query, map[string]any{
		"tenantName": tenantName,
		"orgDedupAt": time,
	})
	return err
}
