package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

// TODO deprecate and remove
type CustomFieldTemplateService interface {
	FindLinkedWithCustomField(ctx context.Context, customFieldId string) (*neo4jentity.CustomFieldTemplateEntity, error)
}

type customFieldTemplateService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewCustomFieldTemplateService(log logger.Logger, repositories *repository.Repositories) CustomFieldTemplateService {
	return &customFieldTemplateService{
		log:          log,
		repositories: repositories,
	}
}

func (s *customFieldTemplateService) FindLinkedWithCustomField(ctx context.Context, customFieldId string) (*neo4jentity.CustomFieldTemplateEntity, error) {
	queryResult, err := s.repositories.CustomFieldTemplateRepository.FindByCustomFieldId(ctx, customFieldId)
	if err != nil {
		return nil, err
	}
	if len(queryResult.([]*db.Record)) == 0 {
		return nil, nil
	}
	return s.mapDbNodeToCustomFieldTemplate((queryResult.([]*db.Record))[0].Values[0].(dbtype.Node)), nil
}

// TODO delete it
func (s *customFieldTemplateService) mapDbNodeToCustomFieldTemplate(dbNode dbtype.Node) *neo4jentity.CustomFieldTemplateEntity {
	props := utils.GetPropsFromNode(dbNode)
	customFieldTemplate := neo4jentity.CustomFieldTemplateEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		Order:     utils.GetInt64PropOrNil(props, "order"),
		Required:  utils.GetBoolPropOrNil(props, "required"),
		Type:      utils.GetStringPropOrEmpty(props, "type"),
		Length:    utils.GetInt64PropOrNil(props, "length"),
		Min:       utils.GetInt64PropOrNil(props, "min"),
		Max:       utils.GetInt64PropOrNil(props, "max"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
	return &customFieldTemplate
}
