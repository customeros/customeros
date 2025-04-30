package public

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/gin-gonic/gin"
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
	Days []DaySlot `json:"days"`
}

type calendarDetailsResponse struct {
	Location                 string `json:"location"`
	TenantName               string `json:"tenantName"`
	TenantLogoURL            string `json:"tenantLogoUrl"`
	DurationMins             int64  `json:"durationMins"`
	BookingTitle             string `json:"bookingTitle"`
	BookingDescription       string `json:"bookingDescription"`
	BookingFormNameEnabled   bool   `json:"bookingFormNameEnabled"`
	BookingFormEmailEnabled  bool   `json:"bookingFormEmailEnabled"`
	BookingFormPhoneEnabled  bool   `json:"bookingFormPhoneEnabled"`
	BookingFormPhoneRequired bool   `json:"bookingFormPhoneRequired"`
}

type bookMeetingRequest struct {
	CalendarId string `json:"calendarId"`
	StartTime  string `json:"startTime"`
	Timezone   string `json:"timezone"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

type timezoneResponse struct {
	Value string `json:"value"`
	Label string `json:"label"`
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
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "GetCalendarAvailability")
		defer spans.Finish()

		// Get and validate calendarId
		calendarId := c.Query("calendarId")
		if calendarId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Calendar ID is required"})
			return
		}

		var startTime *time.Time
		var endTime *time.Time
		var err error

		// Get and validate startTime
		startTimeStr := c.Query("startTime")
		if startTimeStr != "" {
			startTime, err = utils.UnmarshalDateTime(startTimeStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if startTime == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
				return
			}
		}

		// Get and validate endTime
		endTimeStr := c.Query("endTime")
		if endTimeStr != "" {
			endTime, err = utils.UnmarshalDateTime(endTimeStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if endTime == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
				return
			}
		}

		year := c.Query("year")
		month := c.Query("month")

		// if startTime is not set, set it to first day of the year and month, and end time to first day of the next month
		if startTime == nil && year != "" && month != "" {
			yearInt, _ := strconv.Atoi(year)
			monthInt, _ := strconv.Atoi(month)
			startTime = utils.TimePtr(time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC))
			// Set end time to first day of next month
			nextMonth := time.Month(monthInt) + 1
			nextYear := yearInt
			if nextMonth > 12 {
				nextMonth = 1
				nextYear++
			}
			endTime = utils.TimePtr(time.Date(nextYear, nextMonth, 1, 0, 0, 0, 0, time.UTC))
		}

		if startTime == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing start time"})
			return
		}

		// Get timezone (optional)
		timezone := c.Query("timezone")

		// Validate time range
		if startTime.After(*endTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
			return
		}

		// Convert input times to UTC for processing
		startTimeUTC := startTime.UTC()
		endTimeUTC := endTime.UTC()

		// Get meeting booking event to determine tenant
		meetingBookingEvent, err := s.Repositories.PostgresRepositories.MeetingBookingEventRepository.GetByIdCrossTenant(ctx, calendarId)
		if err != nil {
			spans.TraceError(err)
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
		spans.TagTenant(meetingBookingEvent.Tenant)

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
			Days: restDaySlots,
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetCalendarDetails(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "GetCalendarDetails")
		defer spans.Finish()

		// Get and validate calendarId
		calendarId := c.Query("calendarId")
		if calendarId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Calendar ID is required"})
			return
		}

		// Get meeting booking event to determine tenant
		meetingBookingEvent, err := s.Repositories.PostgresRepositories.MeetingBookingEventRepository.GetByIdCrossTenant(ctx, calendarId)
		if err != nil {
			spans.TraceError(err)
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
		spans.TagTenant(meetingBookingEvent.Tenant)

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

		response := &calendarDetailsResponse{
			Location:                 meetingBookingEvent.Location,
			TenantName:               tenantName,
			TenantLogoURL:            "", // TODO: Get from tenant service when available
			DurationMins:             durationMins,
			BookingTitle:             meetingBookingEvent.Title,
			BookingDescription:       meetingBookingEvent.Description,
			BookingFormNameEnabled:   meetingBookingEvent.BookingFormNameEnabled,
			BookingFormEmailEnabled:  meetingBookingEvent.BookingFormEmailEnabled,
			BookingFormPhoneEnabled:  meetingBookingEvent.BookingFormPhoneEnabled,
			BookingFormPhoneRequired: meetingBookingEvent.BookingFormPhoneRequired,
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetTimezones() gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, _ := telemetry.StartRestSpan(c.Request.Context(), "GetTimezones")
		defer spans.Finish()

		timezones := make([]timezoneResponse, len(data.Timezones()))
		for i, timezone := range data.Timezones() {
			timezones[i] = timezoneResponse{
				Value: timezone,
				Label: strings.ReplaceAll(timezone, "_", " "),
			}
		}
		c.JSON(http.StatusOK, gin.H{"timezones": timezones})
	}
}

func BookMeeting(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "BookMeeting")
		defer spans.Finish()

		var request bookMeetingRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse start time
		startTime, err := utils.UnmarshalDateTime(request.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
			return
		}
		if startTime == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing startTime in request"})
			return
		}

		// Email is required
		if request.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing email in request"})
			return
		}

		// Calendar id is required
		if request.CalendarId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing calendarId in request"})
			return
		}

		// Get meeting booking event to determine tenant
		meetingBookingEvent, err := s.Repositories.PostgresRepositories.MeetingBookingEventRepository.GetByIdCrossTenant(ctx, request.CalendarId)
		if err != nil {
			spans.TraceError(err)
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
		spans.TagTenant(meetingBookingEvent.Tenant)

		// Create meeting
		bookMeetingResult, err := s.CommonServices.MeetingService.BookMeeting(ctx, meetingBookingEvent.ID, *startTime, request.Timezone, request.Name, request.Email, request.Phone)
		if err != nil {
			s.Log.Error("Failed to create meeting: %v", err)
			if err == coserrors.ErrSlotNotAvailable {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Requested time slot is not available"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meeting"})
			}
			return
		}

		if bookMeetingResult == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meeting"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Meeting booked successfully",
			"startTime": bookMeetingResult.StartTime,
			"endTime":   bookMeetingResult.EndTime,
			"hostEmail": bookMeetingResult.HostEmail,
			"hostName":  bookMeetingResult.HostName,
		})
	}
}
