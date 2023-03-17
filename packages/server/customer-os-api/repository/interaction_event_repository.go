package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type InteractionEventRepository interface {
	GetAllForInteractionSessions(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetSentByParticipantsForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetSentToParticipantsForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetReplyToInteractionEventsForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
}

type interactionEventRepository struct {
	driver *neo4j.DriverWithContext
}

func NewInteractionEventRepository(driver *neo4j.DriverWithContext) InteractionEventRepository {
	return &interactionEventRepository{
		driver: driver,
	}
}

func (r *interactionEventRepository) GetAllForInteractionSessions(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (s:InteractionSession_%s)<-[:PART_OF]-(e:InteractionEvent) " +
		" WHERE s.id IN $ids " +
		" RETURN e, s.id ORDER BY e.createdAt ASC"

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

func (r *interactionEventRepository) GetSentByParticipantsForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (ie:InteractionEvent_%s)-[rel:SENT_BY]->(p) " +
		" WHERE ie.id IN $ids " +
		" RETURN p, rel, ie.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
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

func (r *interactionEventRepository) GetReplyToInteractionEventsForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (ie:InteractionEvent_%s)-[rel:REPLIES_TO]->(rie:InteractionEvent_%s) " +
		" WHERE ie.id IN $ids " +
		" RETURN rie, ie.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
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

func (r *interactionEventRepository) GetSentToParticipantsForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (ie:InteractionEvent_%s)-[rel:SENT_TO]->(p) " +
		" WHERE ie.id IN $ids " +
		" RETURN p, rel, ie.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
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
