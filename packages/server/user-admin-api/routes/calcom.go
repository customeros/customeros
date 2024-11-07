package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"time"
)

func AddCalComRoutes(rg *gin.RouterGroup, secretsRepo postgresRepository.PersonalIntegrationRepository) {
	var appSource = "calcom"
	rg.POST("/calcom", func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Printf("unable to read body: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		var triggerEvent struct {
			TriggerEvent string `json:"triggerEvent"`
		}

		if err = json.Unmarshal(body, &triggerEvent); err != nil {
			log.Printf("unable to parse json: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		calcom := string(commonService.CALCOM)
		hSignature := ctx.GetHeader(commonService.CalComHeader)

		switch triggerEvent.TriggerEvent {
		case "BOOKING_CREATED":
			log.Printf("BOOKING_CREATED Trigger Event: %s", triggerEvent.TriggerEvent)

			request := BookingCreatedRequest{}
			if err = json.Unmarshal(body, &request); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			logrus.Info("calcom request body: ", request)
			meetingId, err := bookingCreatedHandler(request, body, secretsRepo, hSignature, calcom, appSource)
			if err != nil {
				log.Printf("unable to create meeting: %v", err.Error())
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"result": fmt.Sprintf("Invalid input %s", err.Error()),
				})
				return
			} else {
				log.Printf("calcom meeting created: externalId %s internalId: %s", request.Payload.Uid, *meetingId)
				ctx.JSON(http.StatusOK, gin.H{
					"result": fmt.Sprintf("calcom meeting created: externalId %s internalId: %s", request.Payload.Uid, *meetingId),
				})
				return
			}
		case "BOOKING_RESCHEDULED":
			request := BookingRescheduleRequest{}
			if err = json.Unmarshal(body, &request); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			logrus.Info("calcom request body: ", request)
			meetingId, err := bookingRescheduledHandler(request, body, secretsRepo, hSignature, calcom, appSource)
			if err != nil {
				log.Printf("unable to update meeting: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("nable to update meeting: %v", err.Error()),
				})
				return
			} else {
				log.Printf("calcom meeting rescheduled: externalId %s internalId: %s", request.Payload.Uid, *meetingId)
				ctx.JSON(http.StatusOK, gin.H{
					"result": fmt.Sprintf("calcom meeting rescheduled: externalId %s internalId: %s", request.Payload.Uid, *meetingId),
				})
				return
			}
		case "BOOKING_CANCELLED":
			request := BookingCancelRequest{}
			if err = json.Unmarshal(body, &request); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			logrus.Info("calcom request body: ", request)
			handler, err := bookingCanceledHandler(request, body, secretsRepo, hSignature, calcom)
			if err != nil {
				log.Printf("unable to cancel meeting: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to cancel meeting: %v", err.Error()),
				})
				return
			} else {
				log.Printf("calcom meeting canceled: externalId %s internalId: %s", request.Payload.Uid, *handler)
				ctx.JSON(http.StatusOK, gin.H{
					"result": fmt.Sprintf("calcom meeting canceled: externalId %s internalId: %s", request.Payload.Uid, *handler),
				})
				return
			}
		default:
			format := "Unhandled Trigger Event: %s"
			log.Printf(format, triggerEvent.TriggerEvent)
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"result": fmt.Sprintf(format, triggerEvent.TriggerEvent),
			})
			return
		}
	})
}

func bookingCanceledHandler(request BookingCancelRequest, body []byte, secretsRepo postgresRepository.PersonalIntegrationRepository, hSignature, calcom string) (*string, error) {
	//_, tenant, err := getParticipantAndTenant(cosService, &request.Payload.Organizer.Email)
	//if err != nil {
	//	return nil, fmt.Errorf("can't identify participant or tenant by for email: %v", err.Error())
	//}
	//sigCheck := commonService.SignatureCheck(hSignature, body, secretsRepo, *tenant, request.Payload.Organizer.Email, calcom)
	//if sigCheck != nil {
	//	return nil, fmt.Errorf("signature check failed: %v", sigCheck.Error())
	//}
	//log.Printf("BOOKING_CANCELLED Trigger Event: %s", request.TriggerEvent)
	//exMeeting, err := cosService.ExternalMeeting("calcom", request.Payload.Uid, &request.Payload.Organizer.Email)
	//if err != nil {
	//	log.Printf("unable to find external meeting: %v", err.Error())
	//	return nil, fmt.Errorf("invalid input %v", err.Error())
	//} else {
	//	canceled := cosModel.MeetingStatusCanceled
	//	input := cosModel.MeetingUpdateInput{
	//		Status: &canceled,
	//	}
	//	meeting, err := cosService.UpdateMeeting(exMeeting.ID, input, &request.Payload.Organizer.Email)
	//	if err != nil {
	//		log.Printf("unable to cancel booking: %v", err.Error())
	//		return nil, fmt.Errorf("unable to cancel booking %v", err.Error())
	//	} else {
	//		log.Printf("calcom booking updated: externalId %s internalId: %s", request.Payload.Uid, *meeting)
	//		return meeting, nil
	//	}
	//}
	return nil, nil
}

func bookingRescheduledHandler(request BookingRescheduleRequest, body []byte, secretsRepo postgresRepository.PersonalIntegrationRepository, hSignature, calcom string, appSource string) (*string, error) {
	//_, tenant, err := getParticipantAndTenant(cosService, &request.Payload.Organizer.Email)
	//if err != nil {
	//	return nil, fmt.Errorf("can't identify participant or tenant for email: %v", err.Error())
	//}
	//
	//sigCheck := commonService.SignatureCheck(hSignature, body, secretsRepo, *tenant, request.Payload.Organizer.Email, calcom)
	//if sigCheck != nil {
	//	return nil, fmt.Errorf("signature check failed: %v", sigCheck.Error())
	//}
	//log.Printf("BOOKING_RESCHEDULED Trigger Event: %s", request.TriggerEvent)
	//exMeeting, err := cosService.ExternalMeeting("calcom", request.Payload.RescheduleUid, &request.Payload.Organizer.Email)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to find external booking %v", err.Error())
	//} else {
	//	externalSystem := cosModel.ExternalSystemReferenceInput{
	//		ExternalID:     request.Payload.Uid,
	//		Type:           "CALCOM",
	//		ExternalURL:    &request.Payload.Metadata.VideoCallUrl,
	//		ExternalSource: &appSource,
	//	}
	//	input := cosModel.MeetingUpdateInput{
	//		Name:           &request.Payload.Title,
	//		StartedAt:      &request.Payload.RescheduleStartTime,
	//		EndedAt:        &request.Payload.RescheduleEndTime,
	//		AppSource:      &appSource,
	//		ExternalSystem: &externalSystem,
	//	}
	//	meeting, err := cosService.UpdateMeeting(exMeeting.ID, input, &request.Payload.Organizer.Email)
	//	if err != nil {
	//		return nil, fmt.Errorf("unable to update cos meeting %v", err.Error())
	//	} else {
	//		log.Printf("calcom meeting rescheduled: externalId %s internalId: %s", externalSystem.ExternalID, *meeting)
	//		var participantsEmailsList []string
	//		participantsEmailsSet := make(map[string]struct{})
	//		for _, participant := range exMeeting.AttendedBy {
	//			for _, contactEmail := range participant.ContactParticipant.Emails {
	//				participantsEmailsSet[contactEmail.Email] = struct{}{}
	//				participantsEmailsList = append(participantsEmailsList, contactEmail.Email)
	//			}
	//		}
	//		var attendeesEmailsList []string
	//		attendeesEmailsSet := make(map[string]struct{})
	//		for _, attendee := range request.Payload.Attendees {
	//			attendeesEmailsSet[attendee.Email] = struct{}{}
	//			attendeesEmailsList = append(attendeesEmailsList, attendee.Email)
	//		}
	//		// Add attendees that are not in participants
	//		for _, email := range attendeesEmailsList {
	//			if _, found := participantsEmailsSet[email]; !found {
	//				contactId, err := cosService.GetContactByEmail(&request.Payload.Organizer.Email, &email)
	//				if err != nil {
	//					log.Printf("unable to find contact with email. Creating contact %s: %v", email, err.Error())
	//					contactId, err = cosService.CreateContact(&request.Payload.Organizer.Email, &email)
	//					if err != nil {
	//						log.Printf("Unable to create contact with email %s: %v", email, err.Error())
	//					} else {
	//						log.Printf("attendedBy: %s %s", *contactId, email)
	//						meetingId, err := cosService.MeetingLinkAttendedBy(exMeeting.ID, cosModel.MeetingParticipantInput{ContactID: contactId}, &request.Payload.Organizer.Email)
	//						if err != nil {
	//							log.Printf("unable to link new meeting participant: %v", err.Error())
	//						} else {
	//							log.Printf("contact participant %s added to meeting: %s", *contactId, *meetingId)
	//						}
	//					}
	//				} else {
	//					log.Printf("attendedBy: %s %s", *contactId, email)
	//					meetingId, err := cosService.MeetingLinkAttendedBy(exMeeting.ID, cosModel.MeetingParticipantInput{ContactID: contactId}, &request.Payload.Organizer.Email)
	//					if err != nil {
	//						log.Printf("unable to link new meeting participant: %v", err.Error())
	//					} else {
	//						log.Printf("contact participant %s added to meeting: %s", *contactId, *meetingId)
	//					}
	//				}
	//			}
	//		}
	//
	//		// Remove participants that are not in attendees
	//		for _, email := range participantsEmailsList {
	//			if _, found := attendeesEmailsSet[email]; !found {
	//				contactId, err := cosService.GetContactByEmail(&request.Payload.Organizer.Email, &email)
	//				if err == nil {
	//					log.Printf("unlinking attendedBy: %s %s", *contactId, email)
	//					meetingId, err := cosService.MeetingUnLinkAttendedBy(exMeeting.ID, cosModel.MeetingParticipantInput{ContactID: contactId}, &request.Payload.Organizer.Email)
	//					if err != nil {
	//						log.Printf("unable to un link meeting participant: %v", err.Error())
	//					} else {
	//						log.Printf("contact participant %s removed to meeting: %s", *contactId, *meetingId)
	//					}
	//				}
	//			}
	//		}
	//		return meeting, nil
	//	}
	//}
	return nil, nil
}

func bookingCreatedHandler(request BookingCreatedRequest, body []byte, secretsRepo postgresRepository.PersonalIntegrationRepository, hSignature, calcom, appSource string) (*string, error) {
	//participant, tenant, err := getParticipantAndTenant(cosService, &request.Payload.Organizer.Email)
	//
	//if err != nil {
	//	return nil, fmt.Errorf("meeting created handler error: %v", err.Error())
	//}
	//
	//sigCheck := commonService.SignatureCheck(hSignature, body, secretsRepo, *tenant, request.Payload.Organizer.Email, calcom)
	//if sigCheck != nil {
	//	return nil, fmt.Errorf("unable to check signature: %v", sigCheck.Error())
	//}
	//createdBy := []*cosModel.MeetingParticipantInput{participant}
	//
	//var attendedBy []*cosModel.MeetingParticipantInput
	//for _, attendee := range request.Payload.Attendees {
	//	contactId, err := cosService.GetContactByEmail(&request.Payload.Organizer.Email, &attendee.Email)
	//	if err != nil {
	//		log.Printf("unable to find contact with email. Creating contact %s: %v", attendee.Email, err.Error())
	//		contactId, err = cosService.CreateContact(&request.Payload.Organizer.Email, &attendee.Email)
	//		if err != nil {
	//			log.Printf("Unable to create contact with email %s: %v", attendee.Email, err.Error())
	//		} else {
	//			log.Printf("attendedBy: %s %s", *contactId, attendee.Email)
	//			attendedBy = append(attendedBy, &cosModel.MeetingParticipantInput{ContactID: contactId})
	//		}
	//
	//	} else {
	//		log.Printf("attendedBy: %s %s", *contactId, attendee.Email)
	//		attendedBy = append(attendedBy, &cosModel.MeetingParticipantInput{ContactID: contactId})
	//	}
	//}
	//
	//noteInput := cosModel.NoteInput{Content: utils.StringPtr(request.Payload.AdditionalNotes), AppSource: &appSource}
	//externalSystem := cosModel.ExternalSystemReferenceInput{
	//	ExternalID:     request.Payload.Uid,
	//	Type:           "CALCOM",
	//	ExternalURL:    &request.Payload.Metadata.VideoCallUrl,
	//	ExternalSource: &appSource,
	//}
	//input := cosModel.MeetingInput{
	//	Name:           &request.Payload.Title,
	//	CreatedBy:      createdBy,
	//	AttendedBy:     attendedBy,
	//	StartedAt:      &request.Payload.StartTime,
	//	EndedAt:        &request.Payload.EndTime,
	//	Note:           &noteInput,
	//	ExternalSystem: &externalSystem,
	//	AppSource:      &appSource,
	//}
	//meeting, err := cosService.CreateMeeting(input, &request.Payload.Organizer.Email)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to create meeting: %v", err.Error())
	//}
	//return meeting, nil
	return nil, nil
}

//func getParticipantAndTenant(userEmail *string) (*cosModel.MeetingParticipantInput, *string, error) {
//user, err := cosService.GetUserByEmail(userEmail)
//if err != nil {
//	return nil, nil, fmt.Errorf("no user for meeting creation to parse json: %v", err.Error())
//} else {
//	log.Printf("createdBy: %s %s", user.UserByEmail.ID, *userEmail)
//}
//tenant, err := cosService.GetTenant(userEmail)
//if err != nil {
//	return nil, nil, fmt.Errorf("unable to get tenant by for user: %v", err.Error())
//}
//return &cosModel.MeetingParticipantInput{UserID: &user.UserByEmail.ID}, &tenant.Tenant, nil
//return nil, nil, nil
//}

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
