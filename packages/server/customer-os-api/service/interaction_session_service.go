package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type InteractionSessionService interface {
	mapDbNodeToInteractionSessionEntity(node *dbtype.Node) *entity.InteractionSessionEntity
}

type interactionSessionService struct {
	repositories *repository.Repositories
}

func NewInteractionSessionService(repositories *repository.Repositories) InteractionSessionService {
	return &interactionSessionService{
		repositories: repositories,
	}
}

func (s *interactionSessionService) mapDbNodeToInteractionSessionEntity(node *dbtype.Node) *entity.InteractionSessionEntity {
	props := utils.GetPropsFromNode(*node)
	interactionSessionEntity := entity.InteractionSessionEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		StartedAt:     utils.GetTimePropOrNow(props, "startedAt"),
		EndedAt:       utils.GetTimePropOrNil(props, "endedAt"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		Status:        utils.GetStringPropOrEmpty(props, "status"),
		Type:          utils.GetStringPropOrEmpty(props, "type"),
		Channel:       utils.GetStringPropOrEmpty(props, "channel"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &interactionSessionEntity
}
