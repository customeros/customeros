package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type OrganizationRepository interface {
	GetOrganizationIdById(ctx context.Context, tenant, id string) (string, error)
	GetOrganizationIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
	GetOrganizationIdByDomain(ctx context.Context, tenant, domain string) (string, error)
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
	log    logger.Logger
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext, log logger.Logger) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
		log:    log,
	}
}

func (r *organizationRepository) GetOrganizationIdById(ctx context.Context, tenant, id string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationIdById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				return org.id order by org.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":         tenant,
			"organizationId": id,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *organizationRepository) GetOrganizationIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationIdByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
				return org.id order by org.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":           tenant,
			"externalId":       externalId,
			"externalSystemId": externalSystemId,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *organizationRepository) GetOrganizationIdByDomain(ctx context.Context, tenant, domain string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationIdByDomain")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_DOMAIN]->(d:Domain {domain:$domain})
				return org.id order by org.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant": tenant,
			"domain": domain,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}
