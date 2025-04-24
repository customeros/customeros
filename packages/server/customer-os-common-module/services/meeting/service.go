package meeting

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type meetingService struct {
	log      logger.Logger
	nylas    interfaces.NylasService
	postgres *postgres_repository.Repositories
}

func NewMeetingService(log logger.Logger, nylasService interfaces.NylasService, postgresRepository *postgres_repository.Repositories) interfaces.MeetingService {
	return &meetingService{
		log:      log,
		nylas:    nylasService,
		postgres: postgresRepository,
	}
}

// GetCalendarAvailabilityForEmail implements interfaces.MeetingService.
//func (s *meetingService) GetCalendarAvailabilityForEmail(ctx context.Context, email string, startTime time.Time, endTime time.Time, duration int, timezone string) (*interfaces.CalendarAvailability, error) {
//	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingService.GetCalendarAvailabilityForEmail")
//	defer spans.Finish()
//	spans.LogObjectAsJson("request", map[string]interface{}{
//		"email":     email,
//		"startTime": startTime,
//		"endTime":   endTime,
//		"duration":  duration,
//		"timezone":  timezone,
//	})
//
//	// validate tenant
//	err := common.ValidateTenant(ctx)
//	if err != nil {
//		spans.TraceError(err)
//		return nil, err
//	}
//	tenant := common.GetTenantFromContext(ctx)
//
//	// Get OAuth token for the email
//	oauthToken, err := s.postgresRepository.OAuthTokenRepository.GetByEmail(ctx, tenant, email)
//	if err != nil {
//		spans.TraceError(err)
//		return nil, err
//	}
//	if oauthToken == nil {
//		return nil, fmt.Errorf("no oauth token found for email: %s", email)
//	}
//
//	// check if token has calendar read scope
//	if !oauthToken.HasCalendarReadScope() {
//		return nil, fmt.Errorf("oauth token does not have calendar read scope")
//	}
//
//	// Get calendars for the user
//	calendars, err := s.nylasService.ListCalendars(ctx, email)
//	if err != nil {
//		spans.TraceError(err)
//		return nil, fmt.Errorf("failed to list calendars: %v", err)
//	}
//
//	if len(calendars) == 0 {
//		return &interfaces.CalendarAvailability{
//			TimeSlots:      []interfaces.TimeSlot{},
//			TotalUsers:     0,
//			AvailableUsers: 0,
//		}, nil
//	}
//
//	// Get events for the primary calendar
//	// TODO alexb to test if first calendar is primary or select default calendar
//	primaryCalendar := calendars[0] // Assuming first calendar is primary
//	events, err := s.nylasService.ListEvents(ctx, primaryCalendar.ID, startTime, endTime)
//	if err != nil {
//		spans.TraceError(err)
//		return nil, fmt.Errorf("failed to list events: %v", err)
//	}
//
//	// Generate time slots
//	timeSlots := generateTimeSlots(startTime, endTime, duration, timezone)
//
//	// Check availability for each time slot
//	availableSlots := make([]interfaces.TimeSlot, 0)
//	totalUsers := 1 // Single user
//	availableUsers := 0
//
//	for _, slot := range timeSlots {
//		isAvailable := true
//		for _, event := range events {
//			if isTimeSlotOverlapping(slot.StartTime, slot.EndTime, event.StartTime, event.EndTime) {
//				isAvailable = false
//				break
//			}
//		}
//
//		if isAvailable {
//			availableUsers++
//		}
//
//		availableSlots = append(availableSlots, interfaces.TimeSlot{
//			StartTime:   slot.StartTime,
//			EndTime:     slot.EndTime,
//			IsAvailable: isAvailable,
//		})
//	}
//
//	return &interfaces.CalendarAvailability{
//		TimeSlots:      availableSlots,
//		TotalUsers:     totalUsers,
//		AvailableUsers: availableUsers,
//	}, nil
//}

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

	// Save or update availability
	saved, err := s.postgres.UserCalendarAvailabilityRepository.SaveOrUpdate(ctx, availability)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to save calendar availability: %v", err)
	}

	return saved, nil
}

func (s *meetingService) SetDefaultUserCalendarAvailability(ctx context.Context, email string) (*postgres_entity.UserCalendarAvailability, error) {
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

	defaultAvailability := &postgres_entity.UserCalendarAvailability{
		Tenant:   tenant,
		Email:    email,
		Timezone: "Etc/UTC",
		Monday: postgres_entity.DayAvailability{
			Enabled:   true,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Tuesday: postgres_entity.DayAvailability{
			Enabled:   true,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Wednesday: postgres_entity.DayAvailability{
			Enabled:   true,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Thursday: postgres_entity.DayAvailability{
			Enabled:   true,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Friday: postgres_entity.DayAvailability{
			Enabled:   true,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Saturday: postgres_entity.DayAvailability{
			Enabled:   false,
			StartHour: "08:30",
			EndHour:   "17:00",
		},
		Sunday: postgres_entity.DayAvailability{
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
