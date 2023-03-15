package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type InteractionEventService interface {
	GetInteractionEventsForInteractionSessions(ctx context.Context, ids []string) (*entity.InteractionEventEntities, error)

	mapDbNodeToInteractionEventEntity(node dbtype.Node) *entity.InteractionEventEntity
}

type interactionEventService struct {
	repositories *repository.Repositories
}

func NewInteractionEventService(repositories *repository.Repositories) InteractionEventService {
	return &interactionEventService{
		repositories: repositories,
	}
}

func (s *interactionEventService) GetInteractionEventsForInteractionSessions(ctx context.Context, ids []string) (*entity.InteractionEventEntities, error) {
	interactionEvents, err := s.repositories.InteractionEventRepository.GetAllForInteractionSessions(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}
	interactionEventEntities := entity.InteractionEventEntities{}
	for _, v := range interactionEvents {
		interactionEventEntity := s.mapDbNodeToInteractionEventEntity(*v.Node)
		interactionEventEntity.DataloaderKey = v.LinkedNodeId
		interactionEventEntities = append(interactionEventEntities, *interactionEventEntity)
	}
	return &interactionEventEntities, nil
}

func (s *interactionEventService) mapDbNodeToInteractionEventEntity(node dbtype.Node) *entity.InteractionEventEntity {
	props := utils.GetPropsFromNode(node)
	interactionEventEntity := entity.InteractionEventEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		EventIdentifier: utils.GetStringPropOrEmpty(props, "identifier"),
		Channel:         utils.GetStringPropOrEmpty(props, "channel"),
		Content:         utils.GetStringPropOrEmpty(props, "content"),
		ContentType:     utils.GetStringPropOrEmpty(props, "contentType"),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
		Source:          entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:   entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &interactionEventEntity
}
