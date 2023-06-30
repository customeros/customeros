package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToCalendar(entity *entity.CalendarEntity) *model.Calendar {
	calendar := model.Calendar{
		ID:            entity.Id,
		CalType:       model.CalendarType(entity.CalType),
		Primary:       entity.Primary,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
	if len(entity.Link) > 0 {
		calendar.Link = utils.StringPtr(entity.Link)
	}

	return &calendar
}

func MapEntitiesToCalendars(entities *entity.CalendarEntities) []*model.Calendar {
	var calendars []*model.Calendar
	for _, calendarEntity := range *entities {
		calendars = append(calendars, MapEntityToCalendar(&calendarEntity))
	}
	return calendars
}
