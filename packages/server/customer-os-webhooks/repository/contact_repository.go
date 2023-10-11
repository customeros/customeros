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

type ContactRepository interface {
	GetById(ctx context.Context, tenant, contactId string) (*dbtype.Node, error)
	GetMatchedContactId(ctx context.Context, tenant, externalSystem, externalId string, emails []string) (string, error)
}

type contactRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactRepository(driver *neo4j.DriverWithContext) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) GetById(parentCtx context.Context, tenant, contactId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "ContactRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}) RETURN c`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}

func (r *contactRepository) GetMatchedContactId(ctx context.Context, tenant, externalSystem, externalId string, emails []string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetMatchedContactId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", externalSystem), log.String("externalId", externalId), log.Object("emails", emails))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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
	span.LogFields(log.String("query", query), log.Object("params", params))

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
