package model

import (
	"time"
)

type BookingCreatedRequest struct {
	TriggerEvent string `json:"triggerEvent"`
	Payload      struct {
		Title           string    `json:"title"`
		AdditionalNotes string    `json:"additionalNotes"`
		StartTime       time.Time `json:"startTime"`
		EndTime         time.Time `json:"endTime"`
		Organizer       struct {
			Email string `json:"email"`
		} `json:"organizer"`
		Attendees []struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			TimeZone string `json:"timeZone"`
			Language struct {
				Locale string `json:"locale"`
			} `json:"language"`
		} `json:"attendees"`
		Uid            string `json:"uid"`
		ConferenceData struct {
			CreateRequest struct {
				RequestId string `json:"requestId"`
			} `json:"createRequest"`
		} `json:"conferenceData"`
		Metadata struct {
			VideoCallUrl string `json:"videoCallUrl"`
		} `json:"metadata"`
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
	TriggerEvent string `json:"triggerEvent"`
	Payload      struct {
		Title               string    `json:"title"`
		RescheduleUid       string    `json:"rescheduleUid"`
		RescheduleStartTime time.Time `json:"rescheduleStartTime"`
		RescheduleEndTime   time.Time `json:"rescheduleEndTime"`
		Organizer           struct {
			Email string `json:"email"`
		} `json:"organizer"`
		Attendees []struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			TimeZone string `json:"timeZone"`
			Language struct {
				Locale string `json:"locale"`
			} `json:"language"`
		} `json:"attendees"`
		Uid      string `json:"uid"`
		Metadata struct {
			VideoCallUrl string `json:"videoCallUrl"`
		} `json:"metadata"`
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
type Email struct {
	Email string `json:"email"`
}

type ContactParticipant struct {
	ID     string   `json:"id"`
	Emails []*Email `json:"emails"`
}

type AttendedBy struct {
	ContactParticipant *ContactParticipant `json:"contactParticipant"`
}

type Note struct {
	HTML string `json:"html"`
	ID   string `json:"id"`
}

type ExternalMeeting struct {
	AttendedBy []*AttendedBy `json:"attendedBy"`
	Note       []*Note       `json:"note"`
	ID         string        `json:"id"`
}

type ExternalMeetings struct {
	Content       []*ExternalMeeting `json:"content"`
	TotalElements int64              `json:"totalElements"`
	TotalPages    int64              `json:"totalPages"`
}

type Response struct {
	ExternalMeetings *ExternalMeetings `json:"externalMeetings"`
}
