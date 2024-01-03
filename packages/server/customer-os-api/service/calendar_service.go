package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CalendarService interface {
	GetAllForUsers(ctx context.Context, userIds []string) (*entity.CalendarEntities, error)
}

type calendarService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewCalendarService(log logger.Logger, repositories *repository.Repositories, services *Services) CalendarService {
	return &calendarService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *calendarService) getDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *calendarService) GetAllForUsers(ctx context.Context, userIds []string) (*entity.CalendarEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CalendarService.GetAllForUsers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("userIds", userIds))

	calendars, err := s.repositories.CalendarRepository.GetAllForUsers(ctx, common.GetTenantFromContext(ctx), userIds)
	if err != nil {
		return nil, err
	}
	calendarEntities := entity.CalendarEntities{}
	for _, v := range calendars {
		calendarEntity := s.mapDbNodeToCalendarEntity(*v.Node)
		calendarEntity.DataloaderKey = v.LinkedNodeId
		calendarEntities = append(calendarEntities, *calendarEntity)
	}
	return &calendarEntities, nil
}

func (s *calendarService) mapDbNodeToCalendarEntity(node dbtype.Node) *entity.CalendarEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.CalendarEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CalType:       utils.GetStringPropOrEmpty(props, "calType"),
		Link:          utils.GetStringPropOrEmpty(props, "link"),
		Primary:       utils.GetBoolPropOrFalse(props, "primary"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
	return &result
}
