package repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type OrganizationRepository interface {
	// Deprecated
	GetById(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
	// Deprecated
	GetMatchedOrganizationId(ctx context.Context, tenant, externalSystem, externalId, customerOsId string, domains []string) (string, error)
	// Deprecated
	GetOrganizationIdById(ctx context.Context, tenant, id string) (string, error)
	// Deprecated
	GetOrganizationIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
	// Deprecated
	GetOrganizationIdByDomain(ctx context.Context, tenant, domain string) (string, error)
	// Deprecated
	IsDomainUsedByOrganization(ctx context.Context, tenant, domain, skipOrganizationIds string) (bool, error)
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
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "OrganizationRepository.GetById")
	defer spans.Finish()
	spans.LogKV("organizationId", organizationId)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId}) RETURN o`
	spans.LogKV("query", query)

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationRepository.GetMatchedOrganizationId")
	defer spans.Finish()
	spans.LogKV("externalSystem", externalSystem)
	spans.LogKV("externalId", externalId)
	spans.LogKV("customerOsId", customerOsId)
	spans.LogKV("domains", domains)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})
				OPTIONAL MATCH (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) 
					WHERE $externalSystem <> ''
				OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o1:Organization)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
				OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o2:Organization {customerOsId:$customerOsId})
					WHERE $customerOsId <> ''
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
				"customerOsId":   customerOsId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	orgIDs := dbRecords.([]*db.Record)
	if len(orgIDs) == 1 {
		spans.LogKV("result.organizationId", orgIDs[0].Values[0])
		return orgIDs[0].Values[0].(string), nil
	}
	spans.LogKV("result.organizationId", "")
	return "", nil
}

func (r *organizationRepository) GetOrganizationIdById(ctx context.Context, tenant, id string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationRepository.GetOrganizationIdById")
	defer spans.Finish()
	spans.LogKV("organizationId", id)

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
		spans.TraceError(err)
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *organizationRepository) GetOrganizationIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationRepository.GetOrganizationIdByExternalId")
	defer spans.Finish()
	spans.LogKV("externalId", externalId)
	spans.LogKV("externalSystemId", externalSystemId)

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
		spans.TraceError(err)
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *organizationRepository) GetOrganizationIdByDomain(ctx context.Context, tenant, domain string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationRepository.GetOrganizationIdByDomain")
	defer spans.Finish()
	spans.LogKV("tenant", tenant)
	spans.LogKV("domain", domain)

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

func (r *organizationRepository) IsDomainUsedByOrganization(ctx context.Context, tenant, domain, skipOrganizationId string) (bool, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationRepository.IsDomainUsedByOrganization")
	defer spans.Finish()
	spans.LogKV("tenant", tenant)
	spans.LogKV("domain", domain)
	spans.LogKV("skipOrganizationId", skipOrganizationId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_DOMAIN]->(d:Domain {domain:$domain})
				WHERE org.id <> $skipOrganizationId
				return org.id limit 1`
	params := map[string]any{
		"tenant":             tenant,
		"domain":             domain,
		"skipOrganizationId": skipOrganizationId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return false, err
		} else {
			return queryResult.Next(ctx), nil
		}
	})
	if err != nil {
		spans.TraceError(err)
		return false, err
	}
	spans.LogKV("result", result.(bool))
	return result.(bool), err
}
