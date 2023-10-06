package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OrganizationRepository interface {
	GetById(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
	GetMatchedOrganizationId(ctx context.Context, tenant, externalSystem, externalId, customerOsId string, domains []string) (string, error)
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) GetById(parentCtx context.Context, tenant, organizationId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "OrganizationRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId}) RETURN u`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}

func (r *organizationRepository) GetMatchedOrganizationId(ctx context.Context, tenant, externalSystem, externalId, customerOsId string, domains []string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetMatchedOrganizationId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", externalSystem), log.String("externalId", externalId),
		log.String("customerOsId", customerOsId), log.Object("domains", domains))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o1:Organization)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
				OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o2:Organization {customerOsId:$customerOsId})
				OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o3:Organization)-[:HAS_DOMAIN]->(d:Domain)
					WHERE d.domain in $domains
				with coalesce(o1, o2, o3) as organization
				where organization is not null
				return organization.id limit 1`

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"externalId":     externalId,
				"domains":        domains,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	orgIDs := dbRecords.([]*db.Record)
	if len(orgIDs) == 1 {
		return orgIDs[0].Values[0].(string), nil
	}
	return "", nil
}
