package service

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
	"time"
)

type InteractionSessionService interface {
	GetInteractionEventsForInteractionSessions(ctx context.Context, ids []string) (*entity.InteractionSessionEntities, error)

	mapDbNodeToInteractionSessionEntity(node dbtype.Node) *entity.InteractionSessionEntity
	GetInteractionSessionById(ctx context.Context, id string) (*entity.InteractionSessionEntity, error)
	Create(ctx context.Context, entity *entity.InteractionSessionEntity) (*entity.InteractionSessionEntity, error)
	GetInteractionSessionBySessionIdentifier(ctx context.Context, sessionIdentifier string) (*entity.InteractionSessionEntity, error)
}

type interactionSessionService struct {
	repositories *repository.Repositories
}

func NewInteractionSessionService(repositories *repository.Repositories) InteractionSessionService {
	return &interactionSessionService{
		repositories: repositories,
	}
}

func (s *interactionSessionService) Create(ctx context.Context, entity *entity.InteractionSessionEntity) (*entity.InteractionSessionEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return s.repositories.InteractionSessionRepository.Create(ctx, tx, common.GetTenantFromContext(ctx), entity)
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToInteractionSessionEntity(*queryResult.(*dbtype.Node)), nil
}

func (s *interactionSessionService) GetInteractionSessionById(ctx context.Context, id string) (*entity.InteractionSessionEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (e:InteractionSession_%s {id:$id}) RETURN e`,
			common.GetTenantFromContext(ctx)),
			map[string]interface{}{
				"id": id,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToInteractionSessionEntity(queryResult.(dbtype.Node)), nil
}

func (s *interactionSessionService) GetInteractionSessionBySessionIdentifier(ctx context.Context, sessionIdentifier string) (*entity.InteractionSessionEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (e:InteractionSession_%s {identifier:$identifier}) RETURN e`,
			common.GetTenantFromContext(ctx)),
			map[string]interface{}{
				"identifier": sessionIdentifier,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToInteractionSessionEntity(queryResult.(dbtype.Node)), nil
}

func (s *interactionSessionService) GetInteractionEventsForInteractionSessions(ctx context.Context, ids []string) (*entity.InteractionSessionEntities, error) {
	interactionSessions, err := s.repositories.InteractionSessionRepository.GetAllForInteractionEvents(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}
	interactionSessionEntities := entity.InteractionSessionEntities{}
	for _, v := range interactionSessions {
		interactionSessionEntity := s.mapDbNodeToInteractionSessionEntity(*v.Node)
		interactionSessionEntity.DataloaderKey = v.LinkedNodeId
		interactionSessionEntities = append(interactionSessionEntities, *interactionSessionEntity)
	}
	return &interactionSessionEntities, nil
}

// createdAt takes priority over startedAt
func (s *interactionSessionService) migrateStartedAt(props map[string]any) time.Time {
	if props["createdAt"] != nil {
		return utils.GetTimePropOrNow(props, "createdAt")
	}
	if props["startedAt"] != nil {
		return utils.GetTimePropOrNow(props, "startedAt")
	}
	return time.Now()
}

func (s *interactionSessionService) mapDbNodeToInteractionSessionEntity(node dbtype.Node) *entity.InteractionSessionEntity {
	props := utils.GetPropsFromNode(node)
	interactionSessionEntity := entity.InteractionSessionEntity{
		Id:                utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:         s.migrateStartedAt(props),
		UpdatedAt:         utils.GetTimePropOrNow(props, "updatedAt"),
		EndedAt:           utils.GetTimePropOrNil(props, "endedAt"),
		SessionIdentifier: utils.GetStringPropOrNil(props, "identifier"),
		Name:              utils.GetStringPropOrEmpty(props, "name"),
		Status:            utils.GetStringPropOrEmpty(props, "status"),
		Type:              utils.GetStringPropOrNil(props, "type"),
		Channel:           utils.GetStringPropOrNil(props, "channel"),
		ChannelData:       utils.GetStringPropOrNil(props, "channelData"),
		AppSource:         utils.GetStringPropOrEmpty(props, "appSource"),
		Source:            entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:     entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &interactionSessionEntity
}

func (s *interactionSessionService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}
