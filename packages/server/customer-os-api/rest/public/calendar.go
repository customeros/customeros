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

type calendarAvailabilityRequest struct {
	MeetingBookingEventID string    `form:"meetingBookingEventId" binding:"required"`
	StartTime             time.Time `form:"startTime" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime               time.Time `form:"endTime" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	Timezone              string    `form:"timezone"`
}

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

		var req calendarAvailabilityRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate time range
		if req.StartTime.After(req.EndTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
			return
		}

		// Convert input times to UTC for processing
		startTimeUTC := req.StartTime.UTC()
		endTimeUTC := req.EndTime.UTC()

		// Get meeting booking event to determine tenant
		meetingBookingEvent, err := s.Repositories.PostgresRepositories.MeetingBookingEventRepository.GetByIdCrossTenant(c, req.MeetingBookingEventID)
		if err != nil {
			s.Log.Error("Failed to get meeting booking event: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get meeting booking event"})
			return
		}
		if meetingBookingEvent == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meeting booking event not found"})
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
		availabilityResult, err := s.CommonServices.MeetingService.GetCalendarAvailability(ctx, req.MeetingBookingEventID, startTimeUTC, endTimeUTC, req.Timezone)
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
