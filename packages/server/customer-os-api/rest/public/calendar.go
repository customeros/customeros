package public

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type TimeSlot struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type DaySlot struct {
	Date      time.Time  `json:"date"`
	TimeSlots []TimeSlot `json:"timeSlots"`
}

type calendarAvailabilityResponse struct {
	Days               []DaySlot `json:"days"`
	Location           string    `json:"location"`
	TenantName         string    `json:"tenantName"`
	TenantLogoURL      string    `json:"tenantLogoUrl"`
	DurationMins       int64     `json:"durationMins"`
	BookingTitle       string    `json:"bookingTitle"`
	BookingDescription string    `json:"bookingDescription"`
}

// mapToRestDaySlots converts service response to REST API format
func mapToRestDaySlots(result *interfaces.CalendarAvailabilityResult) []DaySlot {
	if result == nil || len(result.Days) == 0 {
		return []DaySlot{}
	}

	restDaySlots := make([]DaySlot, len(result.Days))
	for i, day := range result.Days {
		timeSlots := make([]TimeSlot, len(day.TimeSlots))
		for j, slot := range day.TimeSlots {
			timeSlots[j] = TimeSlot{
				StartTime: slot.StartTime,
				EndTime:   slot.EndTime,
			}
		}
		restDaySlots[i] = DaySlot{
			Date:      day.Date,
			TimeSlots: timeSlots,
		}
	}
	return restDaySlots
}

// GetCalendarAvailability handles the public endpoint for getting calendar availability
func GetCalendarAvailability(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetCalendarAvailability", c.Request.Header)
		defer span.Finish()

		// Get and validate calendarId
		calendarId := c.Query("calendarId")
		if calendarId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Calendar ID is required"})
			return
		}

		// Get and validate startTime
		startTimeStr := c.Query("startTime")
		if startTimeStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start time is required"})
			return
		}
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format. Expected RFC3339 format (e.g., 2024-03-20T10:00:00Z)"})
			return
		}

		// Get and validate endTime
		endTimeStr := c.Query("endTime")
		if endTimeStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "End time is required"})
			return
		}
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format. Expected RFC3339 format (e.g., 2024-03-20T11:00:00Z)"})
			return
		}

		// Get timezone (optional)
		timezone := c.Query("timezone")

		// Validate time range
		if startTime.After(endTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
			return
		}

		// Convert input times to UTC for processing
		startTimeUTC := startTime.UTC()
		endTimeUTC := endTime.UTC()

		// Get meeting booking event to determine tenant
		meetingBookingEvent, err := s.Repositories.PostgresRepositories.MeetingBookingEventRepository.GetByIdCrossTenant(c, calendarId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.Log.Error("Failed to get meeting booking event: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Calendar not found"})
			return
		}
		if meetingBookingEvent == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Calendar not found"})
			return
		}

		// Set tenant in context
		ctx = common.WithCustomContext(
			ctx,
			&common.CustomContext{
				Tenant: meetingBookingEvent.Tenant,
			},
		)

		// Get tenant settings for workspace name
		tenantSettings, err := s.CommonServices.TenantSettingsService.GetTenantSettings(ctx)
		if err != nil {
			s.Log.Error("Failed to get tenant settings: %v", err)
		}
		tenantName := meetingBookingEvent.Tenant
		if tenantSettings != nil && tenantSettings.WorkspaceName != "" {
			tenantName = tenantSettings.WorkspaceName
		}

		// Round up duration to nearest 5 minutes if needed
		durationMins := meetingBookingEvent.DurationMins
		if durationMins%5 != 0 {
			durationMins = ((durationMins / 5) + 1) * 5
		}

		// Get calendar availability
		availabilityResult, err := s.CommonServices.MeetingService.GetCalendarAvailability(ctx, calendarId, startTimeUTC, endTimeUTC, timezone)
		if err != nil {
			s.Log.Error("Failed to get calendar availability: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get calendar availability"})
			return
		}

		// Convert to REST format
		restDaySlots := mapToRestDaySlots(availabilityResult)

		response := &calendarAvailabilityResponse{
			Days:               restDaySlots,
			Location:           meetingBookingEvent.Location,
			TenantName:         tenantName,
			TenantLogoURL:      "", // TODO: Get from tenant service when available
			DurationMins:       durationMins,
			BookingTitle:       meetingBookingEvent.Title,
			BookingDescription: meetingBookingEvent.Description,
		}

		c.JSON(http.StatusOK, response)
	}
}
