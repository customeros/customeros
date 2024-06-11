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
)

type PlayerRepository interface {
	Merge(ctx context.Context, tx neo4j.ManagedTransaction, entity *entity.PlayerEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tx neo4j.ManagedTransaction, entity *entity.PlayerEntity) (*dbtype.Node, error)
	SetDefaultUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, playerId, userId string, relation entity.PlayerRelation) (*dbtype.Node, error)
	LinkWithUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, playerId, userId, userTenant string, relation entity.PlayerRelation) error
	UnlinkUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, playerId, userId, userTenant string, relation entity.PlayerRelation) error
	GetUsersForPlayer(ctx context.Context, ids []string) ([]*utils.DbNodeWithRelationIdAndTenant, error)
	GetPlayerByAuthIdProvider(ctx context.Context, authId string, provider string) (*dbtype.Node, error)
	GetPlayerByIdentityId(ctx context.Context, identityId string) (*dbtype.Node, error)
	GetPlayerForUser(ctx context.Context, tenant string, userId string, relation entity.PlayerRelation) (*dbtype.Node, error)
}

type playerRepository struct {
	driver *neo4j.DriverWithContext
}

func NewPlayerRepository(driver *neo4j.DriverWithContext) PlayerRepository {
	return &playerRepository{
		driver: driver,
	}
}

func (r *playerRepository) Merge(ctx context.Context, tx neo4j.ManagedTransaction, entity *entity.PlayerEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MERGE (p:Player {authId:$authId, provider:$provider}) " +
		" ON CREATE SET p.id=RandomUUID(), " +
		"				p.identityId=$identityId, " +
		"				p.createdAt=$createdAt, " +
		"				p.updatedAt=datetime(), " +
		"				p.appSource=$appSource, " +
		"				p.source=$source, " +
		"				p.sourceOfTruth=$sourceOfTruth " +
		" RETURN p"

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"authId":        entity.AuthId,
			"provider":      entity.Provider,
			"identityId":    entity.IdentityId,
			"createdAt":     entity.CreatedAt,
			"appSource":     entity.AppSource,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *playerRepository) Update(ctx context.Context, tx neo4j.ManagedTransaction, entity *entity.PlayerEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerRepository.Update")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (p:Player {id: $id}) " +
		"  SET p.identityId=$identityId, " +
		"				p.updatedAt=datetime(), " +
		"				p.sourceOfTruth=$sourceOfTruth, " +
		"				p.appSource=$appSource " +
		" RETURN p"

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"id":            entity.Id,
			"identityId":    entity.IdentityId,
			"sourceOfTruth": entity.SourceOfTruth,
			"appSource":     entity.AppSource,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *playerRepository) SetDefaultUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, playerId, userId string, relation entity.PlayerRelation) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerRepository.SetDefaultUserInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	queryResult, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (p:Player {id:$playerId})-[r:%s]->(u:User)
			SET r.default=
				CASE u.id
					WHEN $userId THEN true
					ELSE false
				END
			RETURN DISTINCT(p)`, relation),
		map[string]interface{}{
			"playerId": playerId,
			"userId":   userId,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *playerRepository) LinkWithUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, playerId, userId, userTenant string, relation entity.PlayerRelation) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerRepository.LinkWithUserInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	queryResult, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (p:Player {id:$playerId}), (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$userTenant})
			MERGE (p)-[r:%s]->(u)
			SET r.default= CASE
				WHEN NOT EXISTS((p)-[:%s {default: true}]->(:User)) THEN true
				ELSE false
			END
			RETURN p`, relation, relation),
		map[string]interface{}{
			"playerId":   playerId,
			"userId":     userId,
			"userTenant": userTenant,
		})
	if err != nil {
		return fmt.Errorf("error linking player with user: %w", err)
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *playerRepository) UnlinkUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, playerId, userId, userTenant string, relation entity.PlayerRelation) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerRepository.UnlinkUserInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`
							MATCH (p:Player {id:$playerId}), (u:User_%s {id:$userId})
							MATCH (p)-[r:%s]->(u)
							DELETE r return p`, userTenant, relation)

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"playerId": playerId,
			"userId":   userId,
		})
	if err != nil {
		return fmt.Errorf("Error unlinking player with user: %w", err)
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *playerRepository) GetUsersForPlayer(ctx context.Context, ids []string) ([]*utils.DbNodeWithRelationIdAndTenant, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerRepository.GetUsersForPlayer")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (p:Player)-[rel:%s]->(u:User)-[:USER_BELONGS_TO_TENANT]->(t:Tenant) WHERE p.id IN $ids RETURN u, rel, p.id, t.name`, entity.IDENTIFIES)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
			map[string]any{
				"ids": ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationIdAndTenant(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error getting users for player: %w", err)
	}

	return result.([]*utils.DbNodeWithRelationIdAndTenant), nil
}

func (r *playerRepository) GetPlayerByAuthIdProvider(ctx context.Context, authId string, provider string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerRepository.GetPlayerByAuthIdProvider")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (p:Player {authId:$authId, provider:$provider}) RETURN p`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
			map[string]any{
				"authId":   authId,
				"provider": provider,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error getting player by authId and provider: %w", err)
	}

	return result.(*dbtype.Node), nil

}

func (r *playerRepository) GetPlayerForUser(ctx context.Context, tenant string, userId string, relation entity.PlayerRelation) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerRepository.GetPlayerForUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (p:Player)-[:%s]->(u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) RETURN p`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, relation),
			map[string]any{
				"userId": userId,
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error getting player for user: %w", err)
	}

	return result.(*dbtype.Node), nil

}

func (r *playerRepository) GetPlayerByIdentityId(ctx context.Context, identityId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerRepository.GetPlayerByIdentityId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (p:Player {identityId:$identityId}) RETURN p`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
			map[string]any{
				"identityId": identityId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error getting player by identityId: %w", err)
	}

	return result.(*dbtype.Node), nil

}
