package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type ContactRepository interface {
	GetContactIdsForEmail(ctx context.Context, tenant, emailId string) ([]string, error)
	GetAllCrossTenantsNotSynced(ctx context.Context, size int) ([]*utils.DbNodeAndId, error)
	GetContactIdById(ctx context.Context, tenant, id string) (string, error)
	GetContactIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
	GetJobRoleId(ctx context.Context, tenant, contactId, organizationId string) (string, error)
}

type contactRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactRepository(driver *neo4j.DriverWithContext) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) GetContactIdsForEmail(ctx context.Context, tenant, emailId string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactIdsForEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(:Email {id:$emailId})
		RETURN c.id`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":  tenant,
				"emailId": emailId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return []string{}, err
	}
	contactIDs := make([]string, 0)
	for _, v := range dbRecords.([]*db.Record) {
		contactIDs = append(contactIDs, v.Values[0].(string))
	}
	return contactIDs, nil
}

func (r *contactRepository) GetAllCrossTenantsNotSynced(ctx context.Context, size int) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetAllCrossTenantsNotSynced")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant)
 			WHERE (c.syncedWithEventStore is null or c.syncedWithEventStore=false)
			RETURN c, t.name limit $size`,
			map[string]any{
				"size": size,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *contactRepository) GetContactIdById(ctx context.Context, tenant, id string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactIdById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})
				return c.id order by c.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":    tenant,
			"contactId": id,
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

func (r *contactRepository) GetContactIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactIdByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
				return c.id order by c.createdAt`

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

func (r *contactRepository) GetJobRoleId(ctx context.Context, tenant, contactId, organizationId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetJobRoleId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})
				MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				MATCH (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org) 
				return j.id order by j.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":         tenant,
			"contactId":      contactId,
			"organizationId": organizationId,
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
