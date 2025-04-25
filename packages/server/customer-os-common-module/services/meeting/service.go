package meeting

import (
	"context"
	"fmt"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type meetingService struct {
	log      logger.Logger
	nylas    interfaces.NylasService
	postgres *postgresRepository.Repositories
}

func NewMeetingService(log logger.Logger, nylasService interfaces.NylasService, postgresRepository *postgresRepository.Repositories) interfaces.MeetingService {
	return &meetingService{
		log:      log,
		nylas:    nylasService,
		postgres: postgresRepository,
	}
}

func (s *meetingService) GetUserCalendarAvailability(ctx context.Context, email string) (*postgresEntity.UserCalendarAvailability, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingService.GetUserCalendarAvailability")
	defer spans.Finish()
	spans.LogKV("email", email)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	if email == "" {
		err = fmt.Errorf("missing email")
		spans.TraceError(err)
		return nil, err
	}

	// Get user's calendar availability
	availability, err := s.postgres.UserCalendarAvailabilityRepository.GetByTenantAndEmail(ctx, tenant, email)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to get calendar availability: %v", err)
	}
	if availability != nil {
		return availability, nil
	}

	defaultAvailability, err := s.SetDefaultUserCalendarAvailability(ctx, email)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to set default calendar availability: %v", err)
	}
	return defaultAvailability, nil
}

// SaveCalendarAvailableHours implements interfaces.MeetingService
func (s *meetingService) SaveUserCalendarAvailability(ctx context.Context, availability *postgresEntity.UserCalendarAvailability) (*postgresEntity.UserCalendarAvailability, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingService.SaveCalendarAvailableHours")
	defer spans.Finish()
	spans.LogKV("email", availability.Email)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	// Set tenant from context
	availability.Tenant = tenant

	// Save or update availability
	saved, err := s.postgres.UserCalendarAvailabilityRepository.SaveOrUpdate(ctx, availability)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to save calendar availability: %v", err)
	}

	return saved, nil
}

func (s *meetingService) SetDefaultUserCalendarAvailability(ctx context.Context, email string) (*postgresEntity.UserCalendarAvailability, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingService.SetDefaultUserCalendarAvailability")
	defer spans.Finish()
	spans.LogKV("email", email)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	// Get user's calendar availability
	availability, err := s.postgres.UserCalendarAvailabilityRepository.GetByTenantAndEmail(ctx, tenant, email)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to get calendar availability: %v", err)
	}
	if availability != nil {
		s.log.Info("calendar availability already exists for email: %s", email)
		return availability, nil
	}

	defaultAvailability := &postgresEntity.UserCalendarAvailability{
		Tenant:   tenant,
		Email:    email,
		Timezone: postgresEntity.DefaultTimezone,
		Monday: postgresEntity.DayAvailability{
			Enabled:   true,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Tuesday: postgresEntity.DayAvailability{
			Enabled:   true,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Wednesday: postgresEntity.DayAvailability{
			Enabled:   true,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Thursday: postgresEntity.DayAvailability{
			Enabled:   true,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Friday: postgresEntity.DayAvailability{
			Enabled:   true,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Saturday: postgresEntity.DayAvailability{
			Enabled:   false,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Sunday: postgresEntity.DayAvailability{
			Enabled:   false,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
	}
	// Save default availability
	defaultAvailability, err = s.postgres.UserCalendarAvailabilityRepository.SaveOrUpdate(ctx, defaultAvailability)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to save default calendar availability: %v", err)
	}

	return defaultAvailability, nil
}

// generateTimeSlots generates 15-minute time slots between start and end time
func generateTimeSlots(startTime, endTime time.Time) []*interfaces.TimeSlot {
	slots := []*interfaces.TimeSlot{}
	current := startTime

	for current.Before(endTime) {
		slotEnd := current.Add(15 * time.Minute)
		if slotEnd.After(endTime) {
			slotEnd = endTime
		}

		slots = append(slots, &interfaces.TimeSlot{
			StartTime:   current,
			EndTime:     slotEnd,
			IsAvailable: true,
		})

		current = slotEnd
	}

	return slots
}

// convertTimeSlotsToTimezone converts all time slots to the specified timezone
func convertTimeSlotsToTimezone(slots []*interfaces.TimeSlot, timezone string) ([]*interfaces.TimeSlot, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone %s: %v", timezone, err)
	}

	convertedSlots := make([]*interfaces.TimeSlot, len(slots))
	for i, slot := range slots {
		convertedSlots[i] = &interfaces.TimeSlot{
			StartTime:   slot.StartTime.In(location),
			EndTime:     slot.EndTime.In(location),
			IsAvailable: slot.IsAvailable,
		}
	}

	return convertedSlots, nil
}

// convertTimeSlotsToDaySlots converts time slots to day slots
func convertTimeSlotsToDaySlots(slots []*interfaces.TimeSlot) []*interfaces.DaySlot {
	// Group time slots by date
	slotsByDate := make(map[string][]*interfaces.TimeSlot)
	for _, slot := range slots {
		// Get date without time component
		date := time.Date(slot.StartTime.Year(), slot.StartTime.Month(), slot.StartTime.Day(), 0, 0, 0, 0, slot.StartTime.Location())
		dateStr := date.Format("2006-01-02")
		slotsByDate[dateStr] = append(slotsByDate[dateStr], slot)
	}

	// Convert to day slots
	daySlots := make([]*interfaces.DaySlot, 0, len(slotsByDate))
	for dateStr, timeSlots := range slotsByDate {
		// Parse date back to time.Time
		date, _ := time.Parse("2006-01-02", dateStr)

		// Check if any time slot is available
		isAvailable := false
		for _, slot := range timeSlots {
			if slot.IsAvailable {
				isAvailable = true
				break
			}
		}

		// Convert []*TimeSlot to []TimeSlot
		convertedTimeSlots := make([]interfaces.TimeSlot, len(timeSlots))
		for i, slot := range timeSlots {
			convertedTimeSlots[i] = *slot
		}

		daySlot := &interfaces.DaySlot{
			Date:        date,
			TimeSlots:   convertedTimeSlots,
			IsAvailable: isAvailable,
		}
		daySlots = append(daySlots, daySlot)
	}

	return daySlots
}

// isTimeInDayAvailability checks if a given time is within the day's availability window
func isTimeInDayAvailability(t time.Time, dayAvailability postgresEntity.DayAvailability) bool {
	if !dayAvailability.Enabled {
		return false
	}

	// Parse start and end hours
	startTime, _ := time.Parse("15:04", dayAvailability.StartHour)
	endTime, _ := time.Parse("15:04", dayAvailability.EndHour)

	// Get just the time part of the input time
	timeOnly := time.Date(0, 1, 1, t.Hour(), t.Minute(), 0, 0, t.Location())

	// Compare times
	return !timeOnly.Before(startTime) && !timeOnly.After(endTime)
}

// isSlotInUserAvailability checks if a time slot is within user's calendar availability
func isSlotInUserAvailability(slot *interfaces.TimeSlot, availability *postgresEntity.UserCalendarAvailability) bool {
	if availability == nil {
		return true // If no availability defined, consider it available
	}

	// Convert slot times to user's timezone
	userLoc, err := time.LoadLocation(availability.Timezone)
	if err != nil {
		return true // If timezone invalid, consider it available
	}

	slotStart := slot.StartTime.In(userLoc)
	slotEnd := slot.EndTime.In(userLoc)

	// Get day of week and corresponding availability
	var dayAvailability postgresEntity.DayAvailability
	switch slotStart.Weekday() {
	case time.Monday:
		dayAvailability = availability.Monday
	case time.Tuesday:
		dayAvailability = availability.Tuesday
	case time.Wednesday:
		dayAvailability = availability.Wednesday
	case time.Thursday:
		dayAvailability = availability.Thursday
	case time.Friday:
		dayAvailability = availability.Friday
	case time.Saturday:
		dayAvailability = availability.Saturday
	case time.Sunday:
		dayAvailability = availability.Sunday
	}

	// Check if both start and end times are within availability
	return isTimeInDayAvailability(slotStart, dayAvailability) && isTimeInDayAvailability(slotEnd, dayAvailability)
}

// GetCalendarAvailability implements interfaces.MeetingService
func (s *meetingService) GetCalendarAvailability(ctx context.Context, meetingBookingEventID string, startTime time.Time, endTime time.Time, timezone string) (*interfaces.CalendarAvailabilityResult, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingService.GetCalendarAvailability")
	defer spans.Finish()
	spans.LogObjectAsJson("request", map[string]interface{}{
		"meetingBookingEventID": meetingBookingEventID,
		"startTime":             startTime,
		"endTime":               endTime,
		"timezone":              timezone,
	})
	if timezone == "" {
		timezone = postgresEntity.DefaultTimezone
	}

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	// Get meeting booking event
	meetingBookingEvent, err := s.postgres.MeetingBookingEventRepository.GetById(ctx, tenant, meetingBookingEventID)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to get meeting booking event: %v", err)
	}
	if meetingBookingEvent == nil {
		return nil, fmt.Errorf("meeting booking event not found")
	}

	meetingDurationMins := meetingBookingEvent.DurationMins
	if meetingDurationMins%15 != 0 {
		meetingDurationMins = ((meetingDurationMins / 15) + 1) * 15
	}

	// 1. Get default calendar for each participant
	participantsData := make(map[string]*interfaces.NylasCalendar)
	for _, email := range meetingBookingEvent.AllowedParticipants {
		// Get default calendar for the participant
		calendar, err := s.nylas.GetDefaultCalendar(ctx, email)
		if err != nil {
			spans.TraceError(err)
			continue
		}
		if calendar == nil {
			s.log.Warn("No default calendar found for participant: %s", email)
			continue
		}

		participantsData[email] = calendar
	}

	spans.LogObjectAsJson("participants_with_calendars", participantsData)

	// 2a. Generate initial time slots for each participant
	participantTimeSlots := make(map[string][]*interfaces.TimeSlot)
	for email := range participantsData {
		// Generate 15-minute time slots for the participant
		slots := generateTimeSlots(startTime, endTime)
		participantTimeSlots[email] = slots
	}

	// 2b. Apply calendar availability restrictions for each participant
	for email, slots := range participantTimeSlots {
		// Get user's calendar availability
		userAvailability, err := s.postgres.UserCalendarAvailabilityRepository.GetByTenantAndEmail(ctx, tenant, email)
		if err != nil {
			s.log.Warn("Failed to get calendar availability for user %s: %v", email, err)
			continue
		}

		// If user has defined calendar availability, apply the restrictions
		if userAvailability != nil {
			// Check each slot against user's calendar availability
			for _, slot := range slots {
				if !isSlotInUserAvailability(slot, userAvailability) {
					slot.IsAvailable = false
				}
			}
		}
	}

	// 3. Convert all time slots to requested timezone
	convertedParticipantTimeSlots := make(map[string][]*interfaces.TimeSlot)
	for email, slots := range participantTimeSlots {
		convertedSlots, err := convertTimeSlotsToTimezone(slots, timezone)
		if err != nil {
			spans.TraceError(err)
			return nil, fmt.Errorf("failed to convert time slots to timezone %s: %v", timezone, err)
		}
		convertedParticipantTimeSlots[email] = convertedSlots
	}

	// 4. Group time slots into day slots
	var daySlots []*interfaces.DaySlot
	for _, slots := range convertedParticipantTimeSlots {
		daySlots = convertTimeSlotsToDaySlots(slots)
		// Only process first participant
		break
	}

	return &interfaces.CalendarAvailabilityResult{
		Days: daySlots,
	}, nil
}
