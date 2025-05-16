package repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type ContactRepository interface {
	// Deprecated
	GetById(ctx context.Context, tenant, contactId string) (*dbtype.Node, error)
	// Deprecated
	GetMatchedContactId(ctx context.Context, tenant, externalSystem, externalId string, emails []string) (string, error)
	// Deprecated
	GetContactIdById(ctx context.Context, tenant, id string) (string, error)
	// Deprecated
	GetContactIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
	// Deprecated
	GetJobRoleId(ctx context.Context, tenant, contactId, organizationId string) (string, error)
}

type contactRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewContactRepository(driver *neo4j.DriverWithContext, database string) ContactRepository {
	return &contactRepository{
		driver:   driver,
		database: database,
	}
}

func (r *contactRepository) GetById(parentCtx context.Context, tenant, contactId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "ContactRepository.GetById")
	defer spans.Finish()
	spans.LogKV("contactId", contactId)

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}) RETURN c`
	spans.LogKV("query", query)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}

func (r *contactRepository) GetMatchedContactId(ctx context.Context, tenant, externalSystem, externalId string, emails []string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactRepository.GetMatchedContactId")
	defer spans.Finish()
	spans.LogKV("externalSystem", externalSystem)
	spans.LogKV("externalId", externalId)
	spans.LogObjectAsJson("emails", emails)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c1:Contact)-[:IS_LINKED_WITH {externalId:$contactExternalId}]->(e)
				OPTIONAL MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c2:Contact)-[:HAS]->(e2:Email)
					WHERE (e2.rawEmail in $emails OR e2.email in $emails) AND size($emails) > 0
				WITH coalesce(c1, c2) as contacts
				WHERE contacts IS NOT NULL
				RETURN contacts.id LIMIT 1`
	params := map[string]interface{}{
		"tenant":            tenant,
		"externalSystem":    externalSystem,
		"contactExternalId": externalId,
		"emails":            emails,
	}
	spans.LogKV("query", query)
	spans.LogObjectAsJson("params", params)

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	contactIDs := dbRecords.([]*db.Record)
	if len(contactIDs) == 1 {
		return contactIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *contactRepository) GetContactIdById(ctx context.Context, tenant, id string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactRepository.GetContactIdById")
	defer spans.Finish()
	spans.LogKV("contactId", id)

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})
				return c.id order by c.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactRepository.GetContactIdByExternalId")
	defer spans.Finish()
	spans.LogKV("externalId", externalId)
	spans.LogKV("externalSystemId", externalSystemId)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
				return c.id order by c.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactRepository.GetJobRoleId")
	defer spans.Finish()
	spans.LogKV("contactId", contactId)
	spans.LogKV("organizationId", organizationId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})
				MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				MATCH (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org) 
				return j.id order by j.createdAt`
	params := map[string]interface{}{
		"tenant":         tenant,
		"contactId":      contactId,
		"organizationId": organizationId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
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
