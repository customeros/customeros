package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type ActionItemRepository interface {
	ActionsItemsExistsForInteractionEvent(ctx context.Context, tenant, interactionEventId string) (bool, error)
	CreateActionItemForEmail(ctx context.Context, tenant, interactionEventId, content, source, appSource string, createdAt time.Time) (*dbtype.Node, error)
}

type actionItemRepository struct {
	driver *neo4j.DriverWithContext
}

func NewActionItemRepository(driver *neo4j.DriverWithContext) ActionItemRepository {
	return &actionItemRepository{
		driver: driver,
	}
}
func (r *actionItemRepository) ActionsItemsExistsForInteractionEvent(ctx context.Context, tenant, interactionEventId string) (bool, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (i:InteractionEvent_%s{id:$interactionEventId})-[:INCLUDES]->(a:ActionItem_%s) RETURN count(a)`, tenant, tenant)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":             tenant,
				"interactionEventId": interactionEventId,
			})
		if err != nil {
			return nil, err
		}
		count, err := queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return count.Values[0].(int64), err
	})

	if err != nil {
		return false, err
	}
	return result.(int64) > 0, nil
}

func (r *actionItemRepository) CreateActionItemForEmail(ctx context.Context, tenant, interactionEventId, content, source, appSource string, createdAt time.Time) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (i:InteractionEvent_%s{id:$interactionEventId}) `, tenant)
	query += fmt.Sprintf(`MERGE (i)-[r:INCLUDES]->(a:ActionItem_%s{id:randomUUID()}) `, tenant)
	query += fmt.Sprintf("ON CREATE SET " +
		" a:ActionItem, " +
		" a.createdAt=$createdAt, " +
		" a.content=$content, " +
		" a.source=$source, " +
		" a.sourceOfTruth=$sourceOfTruth, " +
		" a.appSource=$appSource " +
		" RETURN a")

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":             tenant,
				"interactionEventId": interactionEventId,
				"createdAt":          createdAt,
				"content":            content,
				"source":             source,
				"sourceOfTruth":      source,
				"appSource":          appSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
