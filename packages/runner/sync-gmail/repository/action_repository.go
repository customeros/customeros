package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type ActionRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityId string, entityType entity.EntityType, actionType entity.ActionType, source, appSource string) (*dbtype.Node, error)
}

type actionRepository struct {
	driver *neo4j.DriverWithContext
}

func NewActionRepository(driver *neo4j.DriverWithContext) ActionRepository {
	return &actionRepository{
		driver: driver,
	}
}

func (r *actionRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityId string, entityType entity.EntityType, actionType entity.ActionType, source, appSource string) (*dbtype.Node, error) {
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
