package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type PlayerReadRepository interface {
	GetPlayerByAuthIdProvider(ctx context.Context, authId string, provider string) (*dbtype.Node, error)
	GetUsersForPlayer(ctx context.Context, ids []string) ([]*utils.DbNodeWithRelationIdAndTenant, error)
	GetPlayerByIdentityId(ctx context.Context, identityId string) (*dbtype.Node, error)
	GetPlayerForUser(ctx context.Context, userId string, relation entity.PlayerRelation) (*dbtype.Node, error)
}

type playerReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewPlayerReadRepository(driver *neo4j.DriverWithContext, database string) PlayerReadRepository {
	return &playerReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *playerReadRepository) GetPlayerByAuthIdProvider(ctx context.Context, authId string, provider string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerReadRepository.GetPlayerByAuthIdProvider")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := fmt.Sprintf("MATCH (p:Player {authId:$authId, provider:$provider}) RETURN p")
	params := map[string]any{
		"authId":   authId,
		"provider": provider,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})

	if err != nil && err.Error() == "Result contains no more records" {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if dbRecord == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))
	return dbRecord.(*dbtype.Node), err
}

func (r *playerReadRepository) GetUsersForPlayer(ctx context.Context, ids []string) ([]*utils.DbNodeWithRelationIdAndTenant, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerReadRepository.GetUsersForPlayer")
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

	data := result.([]*utils.DbNodeWithRelationIdAndTenant)
	if data == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Int("result.found", len(data)))

	return data, nil
}

func (r *playerReadRepository) GetPlayerForUser(ctx context.Context, userId string, relation entity.PlayerRelation) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerReadRepository.GetPlayerForUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

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

func (r *playerReadRepository) GetPlayerByIdentityId(ctx context.Context, identityId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PlayerReadRepository.GetPlayerByIdentityId")
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
