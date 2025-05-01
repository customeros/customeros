package meeting

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type meetingService struct {
	log         logger.Logger
	nylas       interfaces.NylasService
	userService interfaces.UserService
	postgres    *postgresRepository.Repositories
}

func NewMeetingService(log logger.Logger, nylasService interfaces.NylasService, userService interfaces.UserService, postgresRepository *postgresRepository.Repositories) interfaces.MeetingService {
	return &meetingService{
		log:         log,
		nylas:       nylasService,
		userService: userService,
		postgres:    postgresRepository,
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

	defaultAvailability, err := s.SetDefaultUserCalendarAvailability(ctx, email, "")
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

func (s *meetingService) SetDefaultUserCalendarAvailability(ctx context.Context, email, timezone string) (*postgresEntity.UserCalendarAvailability, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingService.SetDefaultUserCalendarAvailability")
	defer spans.Finish()
	spans.LogKV("email", email, "timezone", timezone)

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

	if timezone == "" {
		timezone = postgresEntity.DefaultTimezone
	}

	defaultAvailability := &postgresEntity.UserCalendarAvailability{
		Tenant:   tenant,
		Email:    email,
		Timezone: timezone,
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

// generateTimeSlots generates 5-minute time slots between start and end time
func generateTimeSlots(startTime, endTime time.Time) []*interfaces.TimeSlot {
	var slots []*interfaces.TimeSlot

	// Round start time down to nearest 5 minutes
	startMins := startTime.Minute()
	roundedStartMins := (startMins / 5) * 5
	roundedStart := startTime.Add(-time.Duration(startMins-roundedStartMins) * time.Minute)

	// Round end time up to nearest 5 minutes
	endMins := endTime.Minute()
	roundedEndMins := ((endMins + 4) / 5) * 5
	roundedEnd := endTime.Add(time.Duration(roundedEndMins-endMins) * time.Minute)

	current := roundedStart

	for current.Before(roundedEnd) {
		slotEnd := current.Add(5 * time.Minute)
		if slotEnd.After(roundedEnd) {
			slotEnd = roundedEnd
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
	convertedSlots := make([]*interfaces.TimeSlot, len(slots))
	for i, slot := range slots {
		// Convert times to target timezone using the utility function
		startInZone, err := utils.GetTimeInTimeZone(slot.StartTime, timezone)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone %s: %v", timezone, err)
		}
		endInZone, err := utils.GetTimeInTimeZone(slot.EndTime, timezone)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone %s: %v", timezone, err)
		}

		convertedSlots[i] = &interfaces.TimeSlot{
			StartTime:   startInZone,
			EndTime:     endInZone,
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

	// Sort day slots by date in ascending order
	sort.Slice(daySlots, func(i, j int) bool {
		return daySlots[i].Date.Before(daySlots[j].Date)
	})

	return daySlots
}

// isTimeInDayAvailability checks if a given time is within the day's availability window
func isTimeInDayAvailability(t time.Time, dayAvailability postgresEntity.DayAvailability) bool {
	if !dayAvailability.Enabled {
		return false
	}

	// Parse start and end hours in local time
	startTime, _ := time.Parse("15:04", dayAvailability.StartHour)
	endTime, _ := time.Parse("15:04", dayAvailability.EndHour)

	// Create reference time points in the user's timezone for today
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	startTimeToday := today.Add(time.Duration(startTime.Hour())*time.Hour + time.Duration(startTime.Minute())*time.Minute)
	endTimeToday := today.Add(time.Duration(endTime.Hour())*time.Hour + time.Duration(endTime.Minute())*time.Minute)

	// Compare the actual time against the window
	return !t.Before(startTimeToday) && !t.After(endTimeToday)
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

// getValidStartMinutes returns the valid minutes past the hour when a meeting can start
func getValidStartMinutes(durationMins int) []int {
	switch durationMins {
	case 5:
		return []int{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55}
	case 10:
		return []int{0, 10, 20, 30, 40, 50}
	case 15:
		return []int{0, 15, 30, 45}
	case 20:
		return []int{0, 20, 40}
	case 25, 30:
		return []int{0, 30}
	case 45:
		return []int{0}
	case 60, 90, 120:
		return []int{0}
	default:
		// For non-standard durations, round up to nearest 5 and start at hour
		return []int{0}
	}
}

// filterValidMeetingStartSlots filters 5-minute slots to only keep valid meeting start times
func filterValidMeetingStartSlots(slots []*interfaces.TimeSlot, durationMins int32) []*interfaces.TimeSlot {
	if len(slots) == 0 {
		return slots
	}

	validStartMinutes := getValidStartMinutes(int(durationMins))
	filteredSlots := make([]*interfaces.TimeSlot, 0)
	durationInMinutes := int(durationMins)

	for i, slot := range slots {
		if !slot.IsAvailable {
			continue
		}

		// Check if this slot's minute is a valid start time
		slotMinute := slot.StartTime.Minute()
		isValidStartTime := false
		for _, validMin := range validStartMinutes {
			if slotMinute == validMin {
				isValidStartTime = true
				break
			}
		}

		if !isValidStartTime {
			slot.IsAvailable = false
			continue
		}

		// Check if we have enough consecutive available slots for the meeting duration
		hasEnoughTime := true
		slotsNeeded := (durationInMinutes + 4) / 5 // Round up to nearest 5 minutes

		for j := 0; j < slotsNeeded && i+j < len(slots); j++ {
			if !slots[i+j].IsAvailable {
				hasEnoughTime = false
				break
			}
		}

		if !hasEnoughTime {
			slot.IsAvailable = false
			continue
		}

		// Mark subsequent slots as unavailable
		for j := 1; j < slotsNeeded && i+j < len(slots); j++ {
			slots[i+j].IsAvailable = false
		}

		filteredSlots = append(filteredSlots, slot)
	}

	return filteredSlots
}

// unifyTimeSlots combines time slots from all participants based on the assignment method
func unifyTimeSlots(participantSlots map[string][]*interfaces.TimeSlot) []*interfaces.TimeSlot {
	if len(participantSlots) == 0 {
		return []*interfaces.TimeSlot{}
	}

	// Create a map to track unique time slots
	slotMap := make(map[string]*interfaces.TimeSlot)

	// Process slots from each participant
	for _, slots := range participantSlots {
		for _, slot := range slots {
			// Use time as key (format that includes timezone to handle DST correctly)
			key := slot.StartTime.Format(time.RFC3339)
			if _, exists := slotMap[key]; !exists {
				// Create a new slot copy
				slotMap[key] = &interfaces.TimeSlot{
					StartTime:   slot.StartTime,
					EndTime:     slot.EndTime,
					IsAvailable: slot.IsAvailable,
				}
			}
		}
	}

	// Convert map to slice
	unified := make([]*interfaces.TimeSlot, 0, len(slotMap))
	for _, slot := range slotMap {
		unified = append(unified, slot)
	}

	return unified
}

// sortTimeSlots sorts time slots by start time in ascending order
func sortTimeSlots(slots []*interfaces.TimeSlot) []*interfaces.TimeSlot {
	// Create a copy to avoid modifying the original slice
	sorted := make([]*interfaces.TimeSlot, len(slots))
	copy(sorted, slots)

	// Sort by start time
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].StartTime.Before(sorted[j].StartTime)
	})

	return sorted
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

	// Ensure times are in UTC
	if startTime.Location() != time.UTC {
		startTime = startTime.UTC()
	}
	if endTime.Location() != time.UTC {
		endTime = endTime.UTC()
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

	filteredParticipantTimeSlots, err := s.prepareAvailableParticipantTimeSlots(ctx, meetingBookingEvent, startTime, endTime)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// 4c. Convert all time slots to requested timezone
	convertedParticipantTimeSlots := make(map[string][]*interfaces.TimeSlot)
	for email, slots := range filteredParticipantTimeSlots {
		convertedSlots, err := convertTimeSlotsToTimezone(slots, timezone)
		if err != nil {
			spans.TraceError(err)
			return nil, fmt.Errorf("failed to convert time slots to timezone %s: %v", timezone, err)
		}
		convertedParticipantTimeSlots[email] = convertedSlots
	}

	// 4d. Group slots
	unifiedSlots := make([]*interfaces.TimeSlot, 0)
	switch meetingBookingEvent.AssignmentMethod {
	case enum.MeetingBookingAssignmentMethodRoundRobinMaxAvailability:
		// For round-robin, we want slots that are available for any participant
		unifiedSlots = unifyTimeSlots(convertedParticipantTimeSlots)
	default:
		// Default behavior: same as round-robin
		unifiedSlots = unifyTimeSlots(convertedParticipantTimeSlots)
	}

	// Sort slots by start time
	unifiedSlots = sortTimeSlots(unifiedSlots)

	// Group time slots into day slots
	var daySlots = convertTimeSlotsToDaySlots(unifiedSlots)

	return &interfaces.CalendarAvailabilityResult{
		Days: daySlots,
	}, nil
}

// getParticipantTimeSlots returns time slots for all participants in the given time range
func (s *meetingService) prepareAvailableParticipantTimeSlots(ctx context.Context, meetingBookingEvent *postgresEntity.MeetingBookingEvent, startTime time.Time, endTime time.Time) (map[string][]*interfaces.TimeSlot, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingService.prepareAvailableParticipantTimeSlots")
	defer spans.Finish()
	spans.LogObjectAsJson("request", map[string]interface{}{
		"startTime": startTime,
		"endTime":   endTime,
	})
	// Ensure times are in UTC
	if startTime.Location() != time.UTC {
		startTime = startTime.UTC()
	}
	if endTime.Location() != time.UTC {
		endTime = endTime.UTC()
	}

	meetingDurationMins := int(meetingBookingEvent.DurationMins)
	if meetingDurationMins%5 != 0 {
		meetingDurationMins = ((meetingDurationMins / 5) + 1) * 5
	}

	// 1. Get default calendar for each participant
	participantsData := make(map[string]*interfaces.NylasCalendar)
	for _, email := range meetingBookingEvent.ParticipantEmails {
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
		// Generate 5-minute time slots for the participant
		slots := generateTimeSlots(startTime, endTime)
		participantTimeSlots[email] = slots
	}

	// 2b. Apply calendar availability restrictions for each participant
	for email, slots := range participantTimeSlots {
		// Get user's calendar availability
		userAvailability, err := s.postgres.UserCalendarAvailabilityRepository.GetByTenantAndEmail(ctx, common.GetTenantFromContext(ctx), email)
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

	// 2c. Apply booking option restrictions for minimum notice and maximum advance if enabled
	bufferBefore := 0
	bufferAfter := 0
	if meetingBookingEvent.BookOptionEnabled {
		bufferBefore = int(meetingBookingEvent.BookOptionBufferBetweenMeetingsMins)
		bufferAfter = int(meetingBookingEvent.BookOptionBufferBetweenMeetingsMins)

		now := utils.Now()
		minNotice := meetingBookingEvent.BookOptionMinNoticeMins
		if minNotice < 0 {
			minNotice = 0
		}
		maxAdvance := meetingBookingEvent.BookOptionDaysInAdvance
		if maxAdvance < 0 {
			maxAdvance = 0
		}
		minNoticeTime := now.Add(time.Duration(minNotice) * time.Minute)
		maxAdvanceTime := now.AddDate(0, 0, int(maxAdvance))

		// Apply restrictions to all participants' slots
		for _, slots := range participantTimeSlots {
			for _, slot := range slots {
				// Mark slot as unavailable if it's before minimum notice time
				if slot.StartTime.Before(minNoticeTime) {
					slot.IsAvailable = false
					continue
				}

				// Mark slot as unavailable if it's after maximum advance time
				if maxAdvance > 0 && slot.StartTime.After(maxAdvanceTime) {
					slot.IsAvailable = false
				}
			}
		}
	}

	// 2d. Get calendar availability for each participant
	for email, calendar := range participantsData {
		// Get calendar availability from Nylas
		calendarAvailability, err := s.nylas.GetCalendarAvailability(ctx, email, calendar.ID, startTime, endTime, bufferBefore, bufferAfter)
		if err != nil {
			spans.TraceError(err)
			s.log.Warn("Failed to get calendar availability for participant %s: %v", email, err)
			// If error occurs, mark all slots as unavailable for this participant
			for _, slot := range participantTimeSlots[email] {
				slot.IsAvailable = false
			}
			continue
		}

		// Create a map of available times from Nylas
		availableTimes := make(map[int64]bool)
		for _, slot := range calendarAvailability.Data.TimeSlots {
			// Mark all times between start and end as available
			for t := slot.StartTime; t < slot.EndTime; t += 300 { // 300 seconds = 5 minutes
				availableTimes[t] = true
			}
		}

		// Update participant's time slots based on Nylas availability
		// Only mark as unavailable if the slot was previously available
		for _, slot := range participantTimeSlots[email] {
			if slot.IsAvailable { // Only check Nylas availability if slot is currently available
				slotStart := slot.StartTime.UTC().Unix()
				slotEnd := slot.EndTime.UTC().Unix()

				// Check if all 5-minute intervals are available in Nylas
				for t := slotStart; t < slotEnd; t += 300 {
					if !availableTimes[t] {
						slot.IsAvailable = false
						break
					}
				}
			}
		}
	}

	// 3. Filter slots to valid meeting start times
	for email, slots := range participantTimeSlots {
		participantTimeSlots[email] = filterValidMeetingStartSlots(slots, int32(meetingDurationMins))
	}

	// 4a. Remove time slots that are not available
	filteredParticipantTimeSlots := make(map[string][]*interfaces.TimeSlot)
	for email, slots := range participantTimeSlots {
		availableSlots := make([]*interfaces.TimeSlot, 0)
		for _, slot := range slots {
			if slot.IsAvailable {
				availableSlots = append(availableSlots, slot)
			}
		}
		filteredParticipantTimeSlots[email] = availableSlots
	}

	// 4b. Set end time for each slot based on meeting duration
	for _, participant5MinSlots := range filteredParticipantTimeSlots {
		for _, slot := range participant5MinSlots {
			slot.EndTime = slot.StartTime.Add(time.Duration(meetingDurationMins) * time.Minute)
		}
	}

	return filteredParticipantTimeSlots, nil
}

// GetAvailableCalendarParticipantEmailsForTimeRange returns all participants that are available in the given time range
func (s *meetingService) GetAvailableCalendarParticipantEmailsForTimeRange(ctx context.Context, meetingBookingEventID string, startTime, endTime time.Time) ([]string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingService.GetAvailableCalendarParticipantEmailsForTimeRange")
	defer spans.Finish()
	spans.LogObjectAsJson("request", map[string]interface{}{
		"meetingBookingEventID": meetingBookingEventID,
		"startTime":             startTime,
		"endTime":               endTime,
	})

	// Ensure times are in UTC
	if startTime.Location() != time.UTC {
		startTime = startTime.UTC()
	}
	if endTime.Location() != time.UTC {
		endTime = endTime.UTC()
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

	// Get participant time slots using common functionality
	participantTimeSlots, err := s.prepareAvailableParticipantTimeSlots(ctx, meetingBookingEvent, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// Find participants with at least one available slot
	availableParticipants := make([]string, 0)
	for email, slots := range participantTimeSlots {
		// Check if any slot is available
		for _, slot := range slots {
			if slot.IsAvailable {
				availableParticipants = append(availableParticipants, email)
				break
			}
		}
	}

	return availableParticipants, nil
}

// BookMeeting implements interfaces.MeetingService
func (s *meetingService) BookMeeting(ctx context.Context, meetingBookingEventID string, startTime time.Time, timezone, clientName, clientEmail, clientPhone string) (*interfaces.BookMeetingResult, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingService.BookMeeting")
	defer spans.Finish()
	spans.LogObjectAsJson("request", map[string]interface{}{
		"meetingBookingEventID": meetingBookingEventID,
		"startTime":             startTime,
		"timezone":              timezone,
		"clientName":            clientName,
		"clientEmail":           clientEmail,
		"clientPhone":           clientPhone,
	})

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

	if timezone == "" {
		timezone = postgresEntity.DefaultTimezone
	}

	// Ensure times are in UTC
	startTimeInUTC := startTime
	if startTime.Location() != time.UTC {
		startTimeInUTC = startTime.UTC()
	}

	durationMins := int(meetingBookingEvent.DurationMins)
	if durationMins%5 != 0 {
		durationMins = ((durationMins / 5) + 1) * 5
	}
	// Calculate end time based on meeting duration
	endTimeInUTC := startTimeInUTC.Add(time.Duration(durationMins) * time.Minute)

	// Get available participant emails
	availableParticipantEmails, err := s.GetAvailableCalendarParticipantEmailsForTimeRange(ctx, meetingBookingEventID, startTimeInUTC, endTimeInUTC)
	if err != nil {
		return nil, fmt.Errorf("failed to get available participant emails: %v", err)
	}
	if len(availableParticipantEmails) == 0 {
		return nil, coserrors.ErrSlotNotAvailable
	}

	// pick participant
	hostEmail := ""
	switch meetingBookingEvent.AssignmentMethod {
	case enum.MeetingBookingAssignmentMethodRoundRobinMaxAvailability:
		// pick participant with most availability
		hostEmail = pickRandomParticipant(availableParticipantEmails)
	default:
		// pick first participant
		hostEmail = pickRandomParticipant(availableParticipantEmails)
	}

	// get host user name
	hostName := ""
	hostUsers, err := s.userService.GetUsersByEmailAddresses(ctx, []string{hostEmail})
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to get host user: %v", err)
	}
	if len(*hostUsers) > 0 {
		hostName = (*hostUsers)[0].FullName()
	}

	calendar, err := s.nylas.GetDefaultCalendar(ctx, hostEmail)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to get default calendar: %v", err)
	}

	// create meeting event
	createdEvent, nylasResponse, err := s.nylas.CreateEvent(ctx, interfaces.NylasCreateEventRequest{
		Busy:        true,
		Title:       meetingBookingEvent.Title,
		Description: meetingBookingEvent.Description,
		Location:    meetingBookingEvent.Location,
		Participants: []interfaces.NylasEventParticipant{
			{
				Email: hostEmail,
			},
			{
				Email:       clientEmail,
				Name:        clientName,
				PhoneNumber: clientPhone,
			},
		},
		When: struct {
			StartTime     int64  `json:"start_time"`
			EndTime       int64  `json:"end_time"`
			StartTimezone string `json:"start_timezone"`
			EndTimezone   string `json:"end_timezone"`
		}{
			StartTime:     startTimeInUTC.Unix(),
			EndTime:       endTimeInUTC.Unix(),
			StartTimezone: "Etc/UTC",
			EndTimezone:   "Etc/UTC",
		},
	}, hostEmail, calendar.ID, meetingBookingEvent.EmailNotificationEnabled)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to create meeting event: %v", err)
	}
	if createdEvent == nil {
		spans.TraceError(fmt.Errorf("failed to create meeting event, createdEvent is nil"))
		return nil, fmt.Errorf("failed to create meeting event")
	}

	// store booked event
	err = s.postgres.MeetingBookedEventRepository.Create(ctx, &postgresEntity.MeetingBookedEvent{
		MeetingBookingEventID: meetingBookingEventID,
		Tenant:                tenant,
		HostEmail:             hostEmail,
		ClientEmail:           clientEmail,
		ClientName:            clientName,
		ClientPhone:           clientPhone,
		StartTime:             startTimeInUTC,
		EndTime:               endTimeInUTC,
		DurationMins:          int64(durationMins),
		NylasResponse:         nylasResponse,
		Canceled:              false,
	})
	if err != nil {
		spans.TraceError(err)
	}

	startTimeInTimeZone, err := utils.GetTimeInTimeZone(startTimeInUTC, timezone)
	if err != nil {
		spans.TraceError(err)
	}
	endTimeInTimeZone, err := utils.GetTimeInTimeZone(endTimeInUTC, timezone)
	if err != nil {
		spans.TraceError(err)
	}
	result := interfaces.BookMeetingResult{
		StartTime: startTimeInTimeZone,
		EndTime:   endTimeInTimeZone,
		HostEmail: hostEmail,
		HostName:  hostName,
	}

	spans.LogObjectAsJson("result", result)
	return &result, nil
}

func pickRandomParticipant(availableParticipantEmails []string) string {
	return availableParticipantEmails[rand.Intn(len(availableParticipantEmails))]
}
