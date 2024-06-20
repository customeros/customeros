package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type EmailRepository interface {
	GetAllFor(ctx context.Context, tenant string, entityType entity.EntityType, entityId string) ([]*db.Record, error)
	GetAllForIds(ctx context.Context, tenant string, entityType entity.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	RemoveRelationship(ctx context.Context, entityType entity.EntityType, tenant, entityId, email string) error
	RemoveRelationshipById(ctx context.Context, entityType entity.EntityType, tenant, entityId, emailId string) error
	DeleteById(ctx context.Context, tenant, emailId string) error
	Exists(ctx context.Context, tenant, email string) (bool, error)
}

type emailRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewEmailRepository(driver *neo4j.DriverWithContext, database string) EmailRepository {
	return &emailRepository{
		driver:   driver,
		database: database,
	}
}

func (r *emailRepository) GetAllFor(ctx context.Context, tenant string, entityType entity.EntityType, entityId string) ([]*db.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.GetAllFor")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := ""
	switch entityType {
	case entity.CONTACT:
		cypher = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.USER:
		cypher = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
		cypher = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}
	cypher += `, (entity)-[rel:HAS]->(e:Email) RETURN e, rel`
	params := map[string]interface{}{
		"entityId": entityId,
		"tenant":   tenant,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))
	result, err := r.executeQuery(ctx, cypher, params, span)
	if err != nil {
		return nil, err
	}
	return result.Records, nil
}

func (r *emailRepository) GetAllForIds(ctx context.Context, tenant string, entityType entity.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.GetAllForIds")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
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
	query = query + `, (entity)-[rel:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
					WHERE entity.id IN $entityIds
					RETURN e, rel, entity.id ORDER BY e.email, e.rawEmail`

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

func (r *emailRepository) RemoveRelationship(ctx context.Context, entityType entity.EntityType, tenant, entityId, email string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.RemoveRelationship")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
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
		_, err := tx.Run(ctx, query+`MATCH (entity)-[rel:HAS]->(e:Email)
			WHERE e.email = $email OR e.rawEmail = $email
            DELETE rel`,
			map[string]any{
				"entityId": entityId,
				"email":    email,
				"tenant":   tenant,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *emailRepository) RemoveRelationshipById(ctx context.Context, entityType entity.EntityType, tenant, entityId, emailId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.RemoveRelationshipById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
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
		_, err := tx.Run(ctx, query+`MATCH (entity)-[rel:HAS]->(e:Email {id:$emailId})
            DELETE rel`,
			map[string]any{
				"entityId": entityId,
				"emailId":  emailId,
				"tenant":   tenant,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *emailRepository) DeleteById(ctx context.Context, tenant, emailId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.DeleteById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `MATCH (e:Email {id:$emailId})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
            DETACH DELETE e`,
			map[string]any{
				"tenant":  tenant,
				"emailId": emailId,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *emailRepository) Exists(ctx context.Context, tenant string, email string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailRepository.Exists")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	query := "MATCH (e:Email_%s) WHERE e.rawEmail = $email OR e.email = $email RETURN e LIMIT 1"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"email": email,
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

func (r *emailRepository) executeQuery(ctx context.Context, cypher string, params map[string]any, span opentracing.Span) (*neo4j.EagerResult, error) {
	return utils.ExecuteQuery(ctx, *r.driver, r.database, cypher, params, func(err error) {
		tracing.TraceErr(span, err)
	})
}
