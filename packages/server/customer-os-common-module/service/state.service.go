package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/neo4j/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type StateService interface {
	GetStatesByCountryId(ctx context.Context, countryId string) ([]*entity.StateEntity, error)
	MapDbNodeToStateEntity(node dbtype.Node) *entity.StateEntity
}

type stateService struct {
	repositories *repository.Repositories
}

func NewStateService(repository *repository.Repositories) StateService {
	return &stateService{
		repositories: repository,
	}
}

func (s *stateService) GetStatesByCountryId(ctx context.Context, countryId string) ([]*entity.StateEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StateService.GetStatesByCountryId")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "service")
	span.LogFields(log.String("countryId", countryId))

	nodes, err := s.repositories.StateRepository.GetStatesByCountryId(ctx, countryId)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.StateEntity, len(nodes))
	for i, stateNode := range nodes {
		result[i] = s.MapDbNodeToStateEntity(*stateNode)
	}

	return result, nil
}

func (s *stateService) MapDbNodeToStateEntity(node dbtype.Node) *entity.StateEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.StateEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		Code:      utils.GetStringPropOrEmpty(props, "code"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
	return &result
}
