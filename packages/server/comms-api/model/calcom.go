package model

import "time"

type BookingCreatedRequest struct {
	TriggerEvent string    `json:"triggerEvent"`
	CreatedAt    time.Time `json:"createdAt"`
	Payload      struct {
		Type            string `json:"type"`
		Title           string `json:"title"`
		Description     string `json:"description"`
		AdditionalNotes string `json:"additionalNotes"`
		CustomInputs    struct {
		} `json:"customInputs"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
		Organizer struct {
			Id       int    `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
			TimeZone string `json:"timeZone"`
			Language struct {
				Locale string `json:"locale"`
			} `json:"language"`
			TimeFormat string `json:"timeFormat"`
		} `json:"organizer"`
		Responses struct {
			Name struct {
				Label string `json:"label"`
				Value string `json:"value"`
			} `json:"name"`
			Email struct {
				Label string `json:"label"`
				Value string `json:"value"`
			} `json:"email"`
			Location struct {
				Label string `json:"label"`
				Value struct {
					OptionValue string `json:"optionValue"`
					Value       string `json:"value"`
				} `json:"value"`
			} `json:"location"`
			Notes struct {
				Label string `json:"label"`
			} `json:"notes"`
			Guests struct {
				Label string `json:"label"`
			} `json:"guests"`
			RescheduleReason struct {
				Label string `json:"label"`
			} `json:"rescheduleReason"`
		} `json:"responses"`
		UserFieldsResponses struct {
		} `json:"userFieldsResponses"`
		Attendees []struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			TimeZone string `json:"timeZone"`
			Language struct {
				Locale string `json:"locale"`
			} `json:"language"`
		} `json:"attendees"`
		Location            string `json:"location"`
		DestinationCalendar struct {
			Id           int         `json:"id"`
			Integration  string      `json:"integration"`
			ExternalId   string      `json:"externalId"`
			UserId       int         `json:"userId"`
			EventTypeId  interface{} `json:"eventTypeId"`
			CredentialId interface{} `json:"credentialId"`
		} `json:"destinationCalendar"`
		HideCalendarNotes    bool        `json:"hideCalendarNotes"`
		RequiresConfirmation interface{} `json:"requiresConfirmation"`
		EventTypeId          int         `json:"eventTypeId"`
		SeatsShowAttendees   bool        `json:"seatsShowAttendees"`
		SeatsPerTimeSlot     interface{} `json:"seatsPerTimeSlot"`
		Uid                  string      `json:"uid"`
		ConferenceData       struct {
			CreateRequest struct {
				RequestId string `json:"requestId"`
			} `json:"createRequest"`
		} `json:"conferenceData"`
		VideoCallData struct {
			Type     string `json:"type"`
			Id       string `json:"id"`
			Password string `json:"password"`
			Url      string `json:"url"`
		} `json:"videoCallData"`
		AppsStatus []struct {
			AppName  string        `json:"appName"`
			Type     string        `json:"type"`
			Success  int           `json:"success"`
			Failures int           `json:"failures"`
			Errors   []interface{} `json:"errors"`
		} `json:"appsStatus"`
		EventTitle       string      `json:"eventTitle"`
		EventDescription interface{} `json:"eventDescription"`
		Price            int         `json:"price"`
		Currency         string      `json:"currency"`
		Length           int         `json:"length"`
		BookingId        int         `json:"bookingId"`
		Metadata         struct {
			VideoCallUrl string `json:"videoCallUrl"`
		} `json:"metadata"`
		Status string `json:"status"`
	} `json:"payload"`
}

type CreateMeetingResponse struct {
	MeetingCreate struct {
		Id         string    `json:"id"`
		Name       string    `json:"name"`
		Source     string    `json:"source"`
		StartedAt  time.Time `json:"startedAt"`
		EndedAt    time.Time `json:"endedAt"`
		AttendedBy []struct {
			Typename           string `json:"__typename"`
			ContactParticipant struct {
				Id        string `json:"id"`
				FirstName string `json:"firstName"`
			} `json:"contactParticipant"`
		} `json:"attendedBy"`
		CreatedBy []struct {
			Typename        string `json:"__typename"`
			UserParticipant struct {
				Id        string `json:"id"`
				FirstName string `json:"firstName"`
			} `json:"userParticipant"`
		} `json:"createdBy"`
		Note []struct {
			Id            string    `json:"id"`
			Html          string    `json:"html"`
			CreatedAt     time.Time `json:"createdAt"`
			UpdatedAt     time.Time `json:"updatedAt"`
			AppSource     string    `json:"appSource"`
			SourceOfTruth string    `json:"sourceOfTruth"`
		} `json:"note"`
		CreatedAt     time.Time `json:"createdAt"`
		UpdatedAt     time.Time `json:"updatedAt"`
		AppSource     string    `json:"appSource"`
		SourceOfTruth string    `json:"sourceOfTruth"`
	} `json:"meeting_Create"`
}

type BookingRescheduleRequest struct {
	TriggerEvent string    `json:"triggerEvent"`
	CreatedAt    time.Time `json:"createdAt"`
	Payload      struct {
		BookerUrl       string `json:"bookerUrl"`
		Type            string `json:"type"`
		Title           string `json:"title"`
		Description     string `json:"description"`
		AdditionalNotes string `json:"additionalNotes"`
		CustomInputs    struct {
		} `json:"customInputs"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
		Organizer struct {
			Id       int    `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
			TimeZone string `json:"timeZone"`
			Language struct {
				Locale string `json:"locale"`
			} `json:"language"`
			TimeFormat string `json:"timeFormat"`
		} `json:"organizer"`
		Responses struct {
			Name struct {
				Label string `json:"label"`
				Value string `json:"value"`
			} `json:"name"`
			Email struct {
				Label string `json:"label"`
				Value string `json:"value"`
			} `json:"email"`
			Location struct {
				Label string `json:"label"`
				Value struct {
					OptionValue string `json:"optionValue"`
					Value       string `json:"value"`
				} `json:"value"`
			} `json:"location"`
			Notes struct {
				Label string `json:"label"`
				Value string `json:"value"`
			} `json:"notes"`
			Guests struct {
				Label string `json:"label"`
			} `json:"guests"`
			RescheduleReason struct {
				Label string `json:"label"`
				Value string `json:"value"`
			} `json:"rescheduleReason"`
		} `json:"responses"`
		UserFieldsResponses struct {
		} `json:"userFieldsResponses"`
		Attendees []struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			TimeZone string `json:"timeZone"`
			Language struct {
				Locale string `json:"locale"`
			} `json:"language"`
		} `json:"attendees"`
		Location            string `json:"location"`
		DestinationCalendar struct {
			Id           int         `json:"id"`
			Integration  string      `json:"integration"`
			ExternalId   string      `json:"externalId"`
			UserId       int         `json:"userId"`
			EventTypeId  interface{} `json:"eventTypeId"`
			CredentialId interface{} `json:"credentialId"`
		} `json:"destinationCalendar"`
		HideCalendarNotes    bool        `json:"hideCalendarNotes"`
		RequiresConfirmation interface{} `json:"requiresConfirmation"`
		EventTypeId          int         `json:"eventTypeId"`
		SeatsShowAttendees   bool        `json:"seatsShowAttendees"`
		SeatsPerTimeSlot     interface{} `json:"seatsPerTimeSlot"`
		Uid                  string      `json:"uid"`
		VideoCallData        struct {
			Type     string `json:"type"`
			Id       string `json:"id"`
			Password string `json:"password"`
			Url      string `json:"url"`
		} `json:"videoCallData"`
		ConferenceData struct {
			CreateRequest struct {
				RequestId string `json:"requestId"`
			} `json:"createRequest"`
		} `json:"conferenceData"`
		AppsStatus []struct {
			AppName  string        `json:"appName"`
			Type     string        `json:"type"`
			Success  int           `json:"success"`
			Failures int           `json:"failures"`
			Errors   []interface{} `json:"errors"`
		} `json:"appsStatus"`
		EventTitle          string      `json:"eventTitle"`
		EventDescription    interface{} `json:"eventDescription"`
		Price               int         `json:"price"`
		Currency            string      `json:"currency"`
		Length              int         `json:"length"`
		BookingId           int         `json:"bookingId"`
		RescheduleUid       string      `json:"rescheduleUid"`
		RescheduleStartTime time.Time   `json:"rescheduleStartTime"`
		RescheduleEndTime   time.Time   `json:"rescheduleEndTime"`
		Metadata            struct {
			VideoCallUrl string `json:"videoCallUrl"`
		} `json:"metadata"`
		Status string `json:"status"`
	} `json:"payload"`
}

type BookingCancelRequest struct {
	TriggerEvent string `json:"triggerEvent"`
	Payload      struct {
		Organizer struct {
			Email string `json:"email"`
		} `json:"organizer"`
		Uid string `json:"uid"`
	} `json:"payload"`
}

type ExternalMeetingsResponse struct {
	ExternalMeetings struct {
		Content []struct {
			ID string `json:"id"`
		} `json:"content"`
		TotalElements int `json:"totalElements"`
		TotalPages    int `json:"totalPages"`
	} `json:"externalMeetings"`
}
