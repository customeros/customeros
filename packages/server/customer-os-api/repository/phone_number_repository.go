package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type PhoneNumberRepository interface {
	GetAllForIds(ctx context.Context, tenant string, entityType entity.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	RemoveRelationship(ctx context.Context, entityType entity.EntityType, tenant, entityId, phoneNumber string) error
	RemoveRelationshipById(ctx context.Context, entityType entity.EntityType, tenant, entityId, phoneNumberId string) error
	Exists(ctx context.Context, tenant string, e164 string) (bool, error)
	GetByPhoneNumber(ctx context.Context, tenant, e164 string) (*dbtype.Node, error)
}

type phoneNumberRepository struct {
	driver *neo4j.DriverWithContext
}

func NewPhoneNumberRepository(driver *neo4j.DriverWithContext) PhoneNumberRepository {
	return &phoneNumberRepository{
		driver: driver,
	}
}

func (r *phoneNumberRepository) GetAllForIds(ctx context.Context, tenant string, entityType entity.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.GetAllForIds")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := ""
	switch entityType {
	case entity.CONTACT:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(entity:Contact)`
	case entity.USER:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(entity:User)`
	case entity.ORGANIZATION:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(entity:Organization)`
	}
	cypher = cypher + `, (entity)-[rel:HAS]->(p:PhoneNumber)-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t)
					WHERE entity.id IN $entityIds
					RETURN p, rel, entity.id ORDER BY p.createdAt`
	params := map[string]any{
		"tenant":    tenant,
		"entityIds": entityIds,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *phoneNumberRepository) RemoveRelationship(ctx context.Context, entityType entity.EntityType, tenant, entityId, phoneNumber string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.RemoveRelationship")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := ""
	switch entityType {
	case entity.CONTACT:
		cypher = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.USER:
		cypher = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
		cypher = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}
	cypher += `MATCH (entity)-[rel:HAS]->(p:PhoneNumber)
			WHERE p.e164 = $phoneNumber OR p.rawPhoneNumber = $phoneNumber
            DELETE rel`
	params := map[string]any{
		"entityId":    entityId,
		"phoneNumber": phoneNumber,
		"tenant":      tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
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

	cypher := ""
	switch entityType {
	case entity.CONTACT:
		cypher = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.USER:
		cypher = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
		cypher = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}
	cypher += `MATCH (entity)-[rel:HAS]->(p:PhoneNumber {id:$phoneNumberId})
            DELETE rel`
	params := map[string]any{
		"entityId":      entityId,
		"phoneNumberId": phoneNumberId,
		"tenant":        tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
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

	cypher := fmt.Sprintf("MATCH (p:PhoneNumber_%s) WHERE p.e164 = $e164 OR p.rawPhoneNumber = $e164 RETURN p LIMIT 1", tenant)
	params := map[string]any{
		"e164": e164,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		return false, err
	}
	return result.(bool), err
}

func (r *phoneNumberRepository) GetByPhoneNumber(ctx context.Context, tenant, e164 string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.GetByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf("MATCH (p:PhoneNumber_%s) WHERE p.e164 = $e164 OR p.rawPhoneNumber = $e164 RETURN p LIMIT 1", tenant)
	params := map[string]any{
		"e164": e164,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
