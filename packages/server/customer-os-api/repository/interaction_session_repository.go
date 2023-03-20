package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type InteractionSessionRepository interface {
	GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.InteractionSessionEntity) (*dbtype.Node, error)
}

type interactionSessionRepository struct {
	driver *neo4j.DriverWithContext
}

func NewInteractionSessionRepository(driver *neo4j.DriverWithContext) InteractionSessionRepository {
	return &interactionSessionRepository{
		driver: driver,
	}
}

func (r *interactionSessionRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.InteractionSessionEntity) (*dbtype.Node, error) {
	query := "MERGE (is:InteractionSession_%s {id:randomUUID()}) " +
		" ON CREATE SET is:InteractionSession, " +
		"				is.identifier=$identifier, " +
		"				is.source=$source, " +
		"				is.channel=$channel, " +
		"				is.createdAt=$now, " +
		"				is.updatedAt=$now, " +
		"				is.name=$name, " +
		" 				is.status=$status, " +
		"				is.type=$type, " +
		"				is.sourceOfTruth=$sourceOfTruth, " +
		"				is.appSource=$appSource " +
		" RETURN is"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
		map[string]any{
			"identifier":    entity.SessionIdentifier,
			"source":        entity.Source,
			"channel":       entity.Channel,
			"now":           entity.CreatedAt,
			"name":          entity.Name,
			"status":        entity.Status,
			"type":          entity.Type,
			"sourceOfTruth": entity.SourceOfTruth,
			"appSource":     entity.AppSource,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *interactionSessionRepository) GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (e:InteractionEvent_%s)-[:PART_OF]->(s:InteractionSession) " +
		" WHERE e.id IN $ids " +
		" RETURN s, e.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}
