package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ActionRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityId string, entityType entity.EntityType, actionType entity.ActionType, source entity.DataSource, appSource string) (*dbtype.Node, error)
}

type actionRepository struct {
	driver *neo4j.DriverWithContext
}

func NewActionRepository(driver *neo4j.DriverWithContext) ActionRepository {
	return &actionRepository{
		driver: driver,
	}
}

func (r *actionRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityId string, entityType entity.EntityType, actionType entity.ActionType, source entity.DataSource, appSource string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := ""
	switch entityType {
	case entity.ORGANIZATION:
		query = fmt.Sprintf(`MATCH (p:Organization_%s {id:$entityId}) `, tenant)
	}

	query += fmt.Sprintf("MERGE (p)<-[:ACTION_ON]-(a:Action {id:randomUUID()}) "+
		"		ON CREATE SET 	a.type=$type, "+
		"						a.createdAt=$createdAt, "+
		"						a.updatedAt=$createdAt, "+
		"						a.source=$source, "+
		"						a.appSource=$appSource, "+
		"						a:Action_%s, "+
		"						a:TimelineEvent, "+
		"						a:TimelineEvent_%s", tenant, tenant)

	query += ` return a `

	span.LogFields(log.String("query", query))

	if queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
		map[string]interface{}{
			"tenant":    tenant,
			"entityId":  entityId,
			"createdAt": utils.Now(),
			"type":      actionType,
			"source":    source,
			"appSource": appSource,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}
