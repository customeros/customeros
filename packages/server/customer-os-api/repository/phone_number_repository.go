package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type PhoneNumberRepository interface {
	GetAllForContact(ctx context.Context, tenant, contactId string) (any, error)
	GetAllForIds(ctx context.Context, tenant string, entityType entity.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error)

	MergePhoneNumberToContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error)
	UpdatePhoneNumberByContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error)
	SetOtherContactPhoneNumbersNonPrimaryInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, phoneNumberId string) error

	GetByIdAndRelatedEntity(ctx context.Context, entityType entity.EntityType, tenant, phoneNumberId, entityId string) (*dbtype.Node, error)
}

type phoneNumberRepository struct {
	driver *neo4j.DriverWithContext
}

func NewPhoneNumberRepository(driver *neo4j.DriverWithContext) PhoneNumberRepository {
	return &phoneNumberRepository{
		driver: driver,
	}
}

func (r *phoneNumberRepository) MergePhoneNumberToContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {rawPhoneNumber: $rawPhoneNumber}) " +
		" ON CREATE SET p.id=randomUUID(), " +
		"				p.source=$source, " +
		"				p.sourceOfTruth=$sourceOfTruth, " +
		" 				p.appSource=$appSource, " +
		"				p.createdAt=$now, " +
		"				p.updatedAt=$now, " +
		"				p:%s " +
		" WITH p, c " +
		" MERGE (c)-[rel:HAS]->(p) " +
		" SET 	rel.label=$label, " +
		"		rel.primary=$primary, " +
		"		p.sourceOfTruth=$sourceOfTruth," +
		"		p.updatedAt=$now " +
		" RETURN p, rel"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "PhoneNumber_"+tenant),
		map[string]interface{}{
			"tenant":         tenant,
			"contactId":      contactId,
			"rawPhoneNumber": entity.RawPhoneNumber,
			"label":          entity.Label,
			"primary":        entity.Primary,
			"source":         entity.Source,
			"sourceOfTruth":  entity.SourceOfTruth,
			"appSource":      entity.AppSource,
			"now":            utils.Now(),
		})
	return utils.ExtractSingleRecordNodeAndRelationship(ctx, queryResult, err)
}

func (r *phoneNumberRepository) UpdatePhoneNumberByContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				(c)-[rel:HAS]->(p:PhoneNumber {id: $phoneNumberId})
            SET rel.label=$label,
				rel.primary=$primary,
				p.sourceOfTruth=$sourceOfTruth,
				p.updatedAt=$now
			RETURN p, rel`,
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"phoneNumberId": entity.Id,
			"label":         entity.Label,
			"primary":       entity.Primary,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordNodeAndRelationship(ctx, queryResult, err)
}

func (r *phoneNumberRepository) GetAllForContact(ctx context.Context, tenant, contactId string) (any, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	return session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c)-[rel:HAS]->(p:PhoneNumber)
				RETURN p, rel`,
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    tenant,
			})
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		return records, nil
	})
}

func (r *phoneNumberRepository) GetAllForIds(ctx context.Context, tenant string, entityType entity.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := ""
	switch entityType {
	case entity.CONTACT:
		query = `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(entity:Contact)`
	case entity.USER:
		query = `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(entity:User)`
	case entity.ORGANIZATION:
		query = `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(entity:Organization)`
	}
	query = query + `, (entity)-[rel:HAS]->(p:PhoneNumber)-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t)
					WHERE entity.id IN $entityIds
					RETURN p, rel, entity.id ORDER BY p.e164`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":    tenant,
				"entityIds": entityIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *phoneNumberRepository) SetOtherContactPhoneNumbersNonPrimaryInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, phoneNumberId string) error {
	_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				 (c)-[rel:HAS]->(p:PhoneNumber)
			WHERE p.id <> $phoneNumberId
            SET rel.primary=false, p.updatedAt=$now`,
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"phoneNumberId": phoneNumberId,
			"now":           utils.Now(),
		})
	return err
}

func (r *phoneNumberRepository) GetByIdAndRelatedEntity(ctx context.Context, entityType entity.EntityType, tenant, phoneNumberId, entityId string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := ""
	switch entityType {
	case entity.CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case entity.USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	}
	query += ` MATCH (entity)-[rel:HAS]->(p:PhoneNumber {id:$phoneNumberId})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t)
			RETURN p`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":        tenant,
				"phoneNumberId": phoneNumberId,
				"entityId":      entityId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
