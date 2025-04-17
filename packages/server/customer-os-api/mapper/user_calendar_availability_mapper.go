package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

const defaultTimezone = "Etc/UTC"

func MapEntityToModel(entity *postgresEntity.UserCalendarAvailability) *model.UserCalendarAvailability {
	if entity == nil {
		return nil
	}

	timezone := entity.Timezone
	if timezone == "" {
		timezone = defaultTimezone
	}

	return &model.UserCalendarAvailability{
		ID:        entity.ID,
		Email:     entity.Email,
		Timezone:  timezone,
		Monday:    mapDayAvailabilityToModel(entity.Monday),
		Tuesday:   mapDayAvailabilityToModel(entity.Tuesday),
		Wednesday: mapDayAvailabilityToModel(entity.Wednesday),
		Thursday:  mapDayAvailabilityToModel(entity.Thursday),
		Friday:    mapDayAvailabilityToModel(entity.Friday),
		Saturday:  mapDayAvailabilityToModel(entity.Saturday),
		Sunday:    mapDayAvailabilityToModel(entity.Sunday),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func mapDayAvailabilityToModel(day postgresEntity.DayAvailability) *model.DayAvailability {
	return &model.DayAvailability{
		Enabled:   day.Enabled,
		StartHour: day.StartHour,
		EndHour:   day.EndHour,
	}
}

func MapInputToEntity(input *model.UserCalendarAvailabilityInput) *postgresEntity.UserCalendarAvailability {
	if input == nil {
		return nil
	}

	timezone := input.Timezone
	if timezone == "" {
		timezone = defaultTimezone
	}

	return &postgresEntity.UserCalendarAvailability{
		Email:     input.Email,
		Timezone:  timezone,
		Monday:    mapInputToDayAvailability(input.Monday),
		Tuesday:   mapInputToDayAvailability(input.Tuesday),
		Wednesday: mapInputToDayAvailability(input.Wednesday),
		Thursday:  mapInputToDayAvailability(input.Thursday),
		Friday:    mapInputToDayAvailability(input.Friday),
		Saturday:  mapInputToDayAvailability(input.Saturday),
		Sunday:    mapInputToDayAvailability(input.Sunday),
	}
}

func mapInputToDayAvailability(input *model.DayAvailabilityInput) postgresEntity.DayAvailability {
	if input == nil {
		return postgresEntity.DayAvailability{
			Enabled: false,
		}
	}

	return postgresEntity.DayAvailability{
		Enabled:   input.Enabled,
		StartHour: input.StartHour,
		EndHour:   input.EndHour,
	}
}
