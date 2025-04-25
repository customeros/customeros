package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

// MapCalendarAvailabilityResultToDaySlotsModel converts from interface CalendarAvailabilityResult to GraphQL model
func MapCalendarAvailabilityResultToDaySlotsModel(result *interfaces.CalendarAvailabilityResult) []*model.DaySlot {
	if result == nil {
		return []*model.DaySlot{}
	}

	daySlots := make([]*model.DaySlot, len(result.Days))
	for i, day := range result.Days {
		daySlots[i] = MapDaySlotToModel(day)
	}
	return daySlots
}

// MapDaySlotToModel converts from interface DaySlot to GraphQL model
func MapDaySlotToModel(daySlot *interfaces.DaySlot) *model.DaySlot {
	if daySlot == nil {
		return nil
	}

	timeSlots := make([]*model.TimeSlot, len(daySlot.TimeSlots))
	for i, slot := range daySlot.TimeSlots {
		timeSlots[i] = MapTimeSlotToModel(slot)
	}

	return &model.DaySlot{
		Date:      daySlot.Date,
		TimeSlots: timeSlots,
	}
}

// MapTimeSlotToModel converts from interface TimeSlot to GraphQL model
func MapTimeSlotToModel(timeSlot interfaces.TimeSlot) *model.TimeSlot {
	return &model.TimeSlot{
		StartTime: timeSlot.StartTime,
		EndTime:   timeSlot.EndTime,
	}
}
