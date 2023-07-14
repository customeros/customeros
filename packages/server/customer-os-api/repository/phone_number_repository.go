package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type PhoneNumberRepository interface {
	GetAllForIds(ctx context.Context, tenant string, entityType entity.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetByIdAndRelatedEntity(ctx context.Context, entityType entity.EntityType, tenant, phoneNumberId, entityId string) (*dbtype.Node, error)
	MergePhoneNumberToInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType entity.EntityType, entityId string, phoneNumberEntity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error)
	UpdatePhoneNumberForInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType entity.EntityType, entityId string, phoneNumberEntity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error)
	SetOtherPhoneNumbersNonPrimaryInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType entity.EntityType, entityId, phoneNumberId string) error
	RemoveRelationship(ctx context.Context, entityType entity.EntityType, tenant, entityId, phoneNumber string) error
	RemoveRelationshipById(ctx context.Context, entityType entity.EntityType, tenant, entityId, phoneNumberId string) error
	Exists(ctx context.Context, tenant string, e164 string) (bool, error)
	GetByPhoneNumber(ctx context.Context, tenant, e164 string) (*dbtype.Node, error)
	GetById(ctx context.Context, phoneNumberId string) (*dbtype.Node, error)
	LinkWithCountryInTx(ctx context.Context, tx neo4j.ManagedTransaction, phoneNumberId, countryCodeA2 string) error
}

type phoneNumberRepository struct {
	driver *neo4j.DriverWithContext
}

func NewPhoneNumberRepository(driver *neo4j.DriverWithContext) PhoneNumberRepository {
	return &phoneNumberRepository{
		driver: driver,
	}
}

func (r *phoneNumberRepository) MergePhoneNumberToInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType entity.EntityType, entityId string, phoneNumberEntity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.MergePhoneNumberToInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	query := ""

	switch entityType {
	case entity.CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case entity.USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	}

	var operation, onCreate = "MERGE", ""
	if phoneNumberEntity.RawPhoneNumber == "" {
		operation = "CREATE"
		onCreate = ""
	}

	query = query +
		" %s (t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {rawPhoneNumber: $rawPhoneNumber}) " +
		" %s SET p.id=randomUUID(), " +
		"				p.source=$source, " +
		"				p.sourceOfTruth=$sourceOfTruth, " +
		" 				p.appSource=$appSource, " +
		"				p.createdAt=$now, " +
		"				p.updatedAt=$now, " +
		"				p:%s " +
		" WITH p, entity " +
		" MERGE (entity)-[rel:HAS]->(p) " +
		" SET 	rel.label=$label, " +
		"		rel.primary=$primary, " +
		"		p.sourceOfTruth=$sourceOfTruth," +
		"		p.updatedAt=$now " +
		" RETURN p, rel"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, operation, onCreate, "PhoneNumber_"+tenant),
		map[string]interface{}{
			"tenant":         tenant,
			"entityId":       entityId,
			"rawPhoneNumber": phoneNumberEntity.RawPhoneNumber,
			"label":          phoneNumberEntity.Label,
			"primary":        phoneNumberEntity.Primary,
			"source":         phoneNumberEntity.Source,
			"sourceOfTruth":  phoneNumberEntity.SourceOfTruth,
			"appSource":      phoneNumberEntity.AppSource,
			"now":            utils.Now(),
		})
	return utils.ExtractSingleRecordNodeAndRelationship(ctx, queryResult, err)
}

func (r *phoneNumberRepository) UpdatePhoneNumberForInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType entity.EntityType, entityId string, phoneNumberEntity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.UpdatePhoneNumberForInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := ""

	switch entityType {
	case entity.CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case entity.USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	}

	updateRawPhoneNumberQuery := ""
	if phoneNumberEntity.RawPhoneNumber != "" {
		updateRawPhoneNumberQuery = ", p.rawPhoneNumber=$rawPhoneNumber"
	}

	query = query + `, (entity)-[rel:HAS]->(p:PhoneNumber {id: $phoneNumberId})
            SET rel.label=$label,
				rel.primary=$primary,
				p.sourceOfTruth=$sourceOfTruth,
				p.updatedAt=$now 
				%s
			RETURN p, rel`

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, updateRawPhoneNumberQuery),
		map[string]interface{}{
			"tenant":         tenant,
			"entityId":       entityId,
			"phoneNumberId":  phoneNumberEntity.Id,
			"label":          phoneNumberEntity.Label,
			"primary":        phoneNumberEntity.Primary,
			"sourceOfTruth":  phoneNumberEntity.SourceOfTruth,
			"now":            utils.Now(),
			"rawPhoneNumber": phoneNumberEntity.RawPhoneNumber,
		})
	return utils.ExtractSingleRecordNodeAndRelationship(ctx, queryResult, err)
}

func (r *phoneNumberRepository) GetAllForIds(ctx context.Context, tenant string, entityType entity.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.GetAllForIds")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
					RETURN p, rel, entity.id ORDER BY p.createdAt`

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

func (r *phoneNumberRepository) SetOtherPhoneNumbersNonPrimaryInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType entity.EntityType, entityId, phoneNumberId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.SetOtherPhoneNumbersNonPrimaryInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := ""

	switch entityType {
	case entity.CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}

	_, err := tx.Run(ctx, query+`, (entity)-[rel:HAS]->(p:PhoneNumber)
			WHERE p.id <> $phoneNumberId
            SET rel.primary=false, 
				p.updatedAt=$now`,
		map[string]interface{}{
			"tenant":        tenant,
			"entityId":      entityId,
			"phoneNumberId": phoneNumberId,
			"now":           utils.Now(),
		})
	return err
}

func (r *phoneNumberRepository) GetByIdAndRelatedEntity(ctx context.Context, entityType entity.EntityType, tenant, phoneNumberId, entityId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.GetByIdAndRelatedEntity")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("phoneNumberId", phoneNumberId), log.String("entityId", entityId))

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
	span.LogFields(log.String("query", query))

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

func (r *phoneNumberRepository) RemoveRelationship(ctx context.Context, entityType entity.EntityType, tenant, entityId, phoneNumber string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.RemoveRelationship")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := ""
	switch entityType {
	case entity.CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query+`MATCH (entity)-[rel:HAS]->(p:PhoneNumber)
			WHERE p.e164 = $phoneNumber OR p.rawPhoneNumber = $phoneNumber
            DELETE rel`,
			map[string]any{
				"entityId":    entityId,
				"phoneNumber": phoneNumber,
				"tenant":      tenant,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *phoneNumberRepository) RemoveRelationshipById(ctx context.Context, entityType entity.EntityType, tenant, entityId, phoneNumberId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.RemoveRelationshipById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := ""
	switch entityType {
	case entity.CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query+`MATCH (entity)-[rel:HAS]->(p:PhoneNumber {id:$phoneNumberId})
            DELETE rel`,
			map[string]any{
				"entityId":      entityId,
				"phoneNumberId": phoneNumberId,
				"tenant":        tenant,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *phoneNumberRepository) Exists(ctx context.Context, tenant string, e164 string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.Exists")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (p:PhoneNumber_%s) WHERE p.e164 = $e164 OR p.rawPhoneNumber = $e164 RETURN p LIMIT 1"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"e164": e164,
			}); err != nil {
			return false, err
		} else {
			return queryResult.Next(ctx), nil

		}
	})
	if err != nil {
		return false, err
	}
	return result.(bool), err
}

func (r *phoneNumberRepository) GetByPhoneNumber(ctx context.Context, tenant, e164 string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.GetByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf("MATCH (p:PhoneNumber_%s) WHERE p.e164 = $e164 OR p.rawPhoneNumber = $e164 RETURN p LIMIT 1", tenant)
	span.LogFields(log.String("query", query))

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"e164": e164,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *phoneNumberRepository) GetById(ctx context.Context, phoneNumberId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	query := "MATCH (:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId}) RETURN p"
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"phoneNumberId": phoneNumberId,
				"tenant":        common.GetTenantFromContext(ctx),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *phoneNumberRepository) LinkWithCountryInTx(ctx context.Context, tx neo4j.ManagedTransaction, phoneNumberId, countryCodeA2 string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.LinkWithCountryInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("phoneNumberId", phoneNumberId), log.String("countryCodeA2", countryCodeA2))

	query := `MATCH (p:PhoneNumber {id:$phoneNumberId})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
				OPTIONAL MATCH (:Country)<-[rel:LINKED_TO]-(p)
				DELETE rel
				WITH p
				MATCH (c:Country {codeA2:$countryCodeA2})
				MERGE (c)<-[:LINKED_TO]-(p)
				RETURN c`
	span.LogFields(log.String("query", query))

	result, err := tx.Run(ctx, query,
		map[string]any{
			"phoneNumberId": phoneNumberId,
			"countryCodeA2": countryCodeA2,
			"tenant":        common.GetTenantFromContext(ctx),
		})
	if err != nil {
		return err
	}
	_, err = result.Single(ctx)
	if err != nil {
		return err
	}
	return nil
}
