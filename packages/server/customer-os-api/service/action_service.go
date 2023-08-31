package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type ActionService interface {
	mapDbNodeToActionEntity(dbNode dbtype.Node) *entity.ActionEntity
}

type actionService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewActionService(log logger.Logger, repository *repository.Repositories) ActionService {
	return &actionService{
		log:          log,
		repositories: repository,
	}
}

func (s *actionService) mapDbNodeToActionEntity(dbNode dbtype.Node) *entity.ActionEntity {
	props := utils.GetPropsFromNode(dbNode)
	action := entity.ActionEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Type:      entity.GetActionType(utils.GetStringPropOrEmpty(props, "type")),
		Content:   utils.GetStringPropOrEmpty(props, "content"),
		Metadata:  utils.GetStringPropOrEmpty(props, "metadata"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		AppSource: utils.GetStringPropOrEmpty(props, "appSource"),
		Source:    entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
	}
	return &action
}
