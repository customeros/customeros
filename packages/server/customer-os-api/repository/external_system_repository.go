package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ExternalSystemRepository interface {
	LinkNodeWithExternalSystemInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, nodeId, nodeType string, relationship entity.ExternalSystemEntity) error
	GetForEntities(ctx context.Context, tenant string, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error)
}

type externalSystemRepository struct {
	driver *neo4j.DriverWithContext
}

func NewExternalSystemRepository(driver *neo4j.DriverWithContext) ExternalSystemRepository {
	return &externalSystemRepository{
		driver: driver,
	}
}

func (r *externalSystemRepository) LinkNodeWithExternalSystemInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, nodeId, nodeType string, externalSystem entity.ExternalSystemEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemRepository.LinkContactWithExternalSystemInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})," +
		" (n:%s {id:$nodeId}) " +
		" MERGE (n)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e) " +
		" ON CREATE SET e:%s, " +
		"				r.syncDate=$syncDate, " +
		"				r.externalUrl=$externalUrl, " +
		"				r.externalSource=$externalSource, " +
		"				e.createdAt=datetime({timezone: 'UTC'}) " +
		" ON MATCH SET r.syncDate=$syncDate " +
		" RETURN r"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, nodeType, "ExternalSystem_"+tenant),
		map[string]any{
			"nodeId":           nodeId,
			"tenant":           tenant,
			"syncDate":         *externalSystem.Relationship.SyncDate,
			"externalId":       externalSystem.Relationship.ExternalId,
			"externalSystemId": externalSystem.ExternalSystemId,
			"externalUrl":      externalSystem.Relationship.ExternalUrl,
			"externalSource":   externalSystem.Relationship.ExternalSource,
		})

	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *externalSystemRepository) GetForEntities(ctx context.Context, tenant string, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemRepository.GetForEntities")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem)<-[rel:IS_LINKED_WITH]-(n:Issue|User|Contact|Organization|Meeting|InteractionEvent|InteractionSession|LogEntry|Comment)
			WHERE n.id IN $entityIds
			RETURN e, rel, n.id order by rel.syncDate`

	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
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
