package meeting

import (
	"context"
	"fmt"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type meetingService struct {
	log                logger.Logger
	nylasService       interfaces.NylasService
	postgresRepository *postgres_repository.Repositories
}

func NewMeetingService(log logger.Logger, nylasService interfaces.NylasService, postgresRepository *postgres_repository.Repositories) interfaces.MeetingService {
	return &meetingService{
		log:                log,
		nylasService:       nylasService,
		postgresRepository: postgresRepository,
	}
}

func generateTimeSlots(startTime, endTime time.Time, duration int, timezone string) []interfaces.TimeSlot {
	loc, _ := time.LoadLocation(timezone)
	startTime = startTime.In(loc)
	endTime = endTime.In(loc)

	var slots []interfaces.TimeSlot
	currentTime := startTime

	for currentTime.Before(endTime) {
		slotEnd := currentTime.Add(time.Duration(duration) * time.Minute)
		if slotEnd.After(endTime) {
			break
		}

		slots = append(slots, interfaces.TimeSlot{
			StartTime:   currentTime,
			EndTime:     slotEnd,
			IsAvailable: true,
		})

		currentTime = slotEnd
	}

	return slots
}

func isTimeSlotOverlapping(slotStart, slotEnd, eventStart, eventEnd time.Time) bool {
	return !(slotEnd.Before(eventStart) || slotStart.After(eventEnd))
}

// GetCalendarAvailabilityForEmail implements interfaces.MeetingService.
func (s *meetingService) GetCalendarAvailabilityForEmail(ctx context.Context, email string, startTime time.Time, endTime time.Time, duration int, timezone string) (*interfaces.CalendarAvailability, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingService.GetCalendarAvailabilityForEmail")
	defer spans.Finish()
	spans.LogObjectAsJson("request", map[string]interface{}{
		"email":     email,
		"startTime": startTime,
		"endTime":   endTime,
		"duration":  duration,
		"timezone":  timezone,
	})

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	// Get OAuth token for the email
	oauthToken, err := s.postgresRepository.OAuthTokenRepository.GetByEmail(ctx, tenant, email)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if oauthToken == nil {
		return nil, fmt.Errorf("no oauth token found for email: %s", email)
	}

	// check if token has calendar read scope
	if !oauthToken.HasCalendarReadScope() {
		return nil, fmt.Errorf("oauth token does not have calendar read scope")
	}

	// Get calendars for the user
	calendars, err := s.nylasService.ListCalendars(ctx, email)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to list calendars: %v", err)
	}

	if len(calendars) == 0 {
		return &interfaces.CalendarAvailability{
			TimeSlots:      []interfaces.TimeSlot{},
			TotalUsers:     0,
			AvailableUsers: 0,
		}, nil
	}

	// Get events for the primary calendar
	// TODO alexb to test if first calendar is primary or select default calendar
	primaryCalendar := calendars[0] // Assuming first calendar is primary
	events, err := s.nylasService.ListEvents(ctx, primaryCalendar.ID, startTime, endTime)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to list events: %v", err)
	}

	// Generate time slots
	timeSlots := generateTimeSlots(startTime, endTime, duration, timezone)

	// Check availability for each time slot
	availableSlots := make([]interfaces.TimeSlot, 0)
	totalUsers := 1 // Single user
	availableUsers := 0

	for _, slot := range timeSlots {
		isAvailable := true
		for _, event := range events {
			if isTimeSlotOverlapping(slot.StartTime, slot.EndTime, event.StartTime, event.EndTime) {
				isAvailable = false
				break
			}
		}

		if isAvailable {
			availableUsers++
		}

		availableSlots = append(availableSlots, interfaces.TimeSlot{
			StartTime:   slot.StartTime,
			EndTime:     slot.EndTime,
			IsAvailable: isAvailable,
		})
	}

	return &interfaces.CalendarAvailability{
		TimeSlots:      availableSlots,
		TotalUsers:     totalUsers,
		AvailableUsers: availableUsers,
	}, nil
}

// GetCalendarAvailabilityForTenant implements interfaces.MeetingService.
func (s *meetingService) GetCalendarAvailabilityForTenant(ctx context.Context, startTime time.Time, endTime time.Time, duration int, timezone string) (*interfaces.CalendarAvailability, error) {
	// TODO implement
	return nil, nil
}

func (s *meetingService) GetUserCalendarAvailability(ctx context.Context, email string) (*postgres_entity.UserCalendarAvailability, error) {
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

	// Get user's calendar availability
	availability, err := s.postgresRepository.UserCalendarAvailabilityRepository.GetByTenantAndEmail(ctx, tenant, email)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to get calendar availability: %v", err)
	}

	return availability, nil
}

// SaveCalendarAvailableHours implements interfaces.MeetingService
func (s *meetingService) SaveUserCalendarAvailability(ctx context.Context, availability *postgres_entity.UserCalendarAvailability) (*postgres_entity.UserCalendarAvailability, error) {
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

	// Validate time formats
	if err := validateAvailabilityTimes(availability); err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// Save or update availability
	saved, err := s.postgresRepository.UserCalendarAvailabilityRepository.SaveOrUpdate(ctx, availability)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to save calendar availability: %v", err)
	}

	return saved, nil
}

// validateAvailabilityTimes validates the time formats in the availability settings
func validateAvailabilityTimes(availability *postgres_entity.UserCalendarAvailability) error {
	days := []postgres_entity.DayAvailability{
		availability.Monday,
		availability.Tuesday,
		availability.Wednesday,
		availability.Thursday,
		availability.Friday,
		availability.Saturday,
		availability.Sunday,
	}

	for _, day := range days {
		if !day.Enabled {
			continue
		}

		// Validate time format (HH:MM)
		if _, err := time.Parse("15:04", day.StartHour); err != nil {
			return fmt.Errorf("invalid start hour format: %s", day.StartHour)
		}
		if _, err := time.Parse("15:04", day.EndHour); err != nil {
			return fmt.Errorf("invalid end hour format: %s", day.EndHour)
		}

		// Validate start time is before end time
		start, _ := time.Parse("15:04", day.StartHour)
		end, _ := time.Parse("15:04", day.EndHour)
		if !start.Before(end) {
			return fmt.Errorf("start time must be before end time: %s - %s", day.StartHour, day.EndHour)
		}
	}

	// Validate timezone
	if _, err := time.LoadLocation(availability.Timezone); err != nil {
		return fmt.Errorf("invalid timezone: %s", availability.Timezone)
	}

	return nil
}
