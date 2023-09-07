package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	s "github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
)

func AddCalComRoutes(rg *gin.RouterGroup, cosService s.CustomerOSService, secretsRepo repository.PersonalIntegrationRepository) {
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

			request := model.BookingCreatedRequest{}
			if err = json.Unmarshal(body, &request); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			logrus.Info("calcom request body: ", request)
			meetingId, err := bookingCreatedHandler(cosService, request, body, secretsRepo, hSignature, calcom, appSource)
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
			request := model.BookingRescheduleRequest{}
			if err = json.Unmarshal(body, &request); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			logrus.Info("calcom request body: ", request)
			meetingId, err := bookingRescheduledHandler(cosService, request, body, secretsRepo, hSignature, calcom, appSource)
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
			request := model.BookingCancelRequest{}
			if err = json.Unmarshal(body, &request); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			logrus.Info("calcom request body: ", request)
			handler, err := bookingCanceledHandler(cosService, request, body, secretsRepo, hSignature, calcom)
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

func bookingCanceledHandler(cosService s.CustomerOSService, request model.BookingCancelRequest, body []byte, secretsRepo repository.PersonalIntegrationRepository, hSignature, calcom string) (*string, error) {
	_, tenant, err := getParticipantAndTenant(cosService, &request.Payload.Organizer.Email)
	if err != nil {
		return nil, fmt.Errorf("can't identify participant or tenant by for email: %v", err.Error())
	}
	sigCheck := commonService.SignatureCheck(hSignature, body, secretsRepo, *tenant, request.Payload.Organizer.Email, calcom)
	if sigCheck != nil {
		return nil, fmt.Errorf("signature check failed: %v", sigCheck.Error())
	}
	log.Printf("BOOKING_CANCELLED Trigger Event: %s", request.TriggerEvent)
	exMeeting, err := cosService.ExternalMeeting("calcom", request.Payload.Uid, &request.Payload.Organizer.Email)
	if err != nil {
		log.Printf("unable to find external meeting: %v", err.Error())
		return nil, fmt.Errorf("invalid input %v", err.Error())
	} else {
		canceled := cosModel.MeetingStatusCanceled
		input := cosModel.MeetingUpdateInput{
			Status: &canceled,
		}
		meeting, err := cosService.UpdateMeeting(exMeeting.ID, input, &request.Payload.Organizer.Email)
		if err != nil {
			log.Printf("unable to cancel booking: %v", err.Error())
			return nil, fmt.Errorf("unable to cancel booking %v", err.Error())
		} else {
			log.Printf("calcom booking updated: externalId %s internalId: %s", request.Payload.Uid, *meeting)
			return meeting, nil
		}
	}
}

func bookingRescheduledHandler(cosService s.CustomerOSService, request model.BookingRescheduleRequest, body []byte, secretsRepo repository.PersonalIntegrationRepository, hSignature, calcom string, appSource string) (*string, error) {
	_, tenant, err := getParticipantAndTenant(cosService, &request.Payload.Organizer.Email)
	if err != nil {
		return nil, fmt.Errorf("can't identify participant or tenant for email: %v", err.Error())
	}

	sigCheck := commonService.SignatureCheck(hSignature, body, secretsRepo, *tenant, request.Payload.Organizer.Email, calcom)
	if sigCheck != nil {
		return nil, fmt.Errorf("signature check failed: %v", sigCheck.Error())
	}
	log.Printf("BOOKING_RESCHEDULED Trigger Event: %s", request.TriggerEvent)
	exMeeting, err := cosService.ExternalMeeting("calcom", request.Payload.RescheduleUid, &request.Payload.Organizer.Email)
	if err != nil {
		return nil, fmt.Errorf("unable to find external booking %v", err.Error())
	} else {
		externalSystem := cosModel.ExternalSystemReferenceInput{
			ExternalID:     request.Payload.Uid,
			Type:           "CALCOM",
			ExternalURL:    &request.Payload.Metadata.VideoCallUrl,
			ExternalSource: &appSource,
		}
		input := cosModel.MeetingUpdateInput{
			Name:           &request.Payload.Title,
			StartedAt:      &request.Payload.RescheduleStartTime,
			EndedAt:        &request.Payload.RescheduleEndTime,
			AppSource:      &appSource,
			ExternalSystem: &externalSystem,
		}
		meeting, err := cosService.UpdateMeeting(exMeeting.ID, input, &request.Payload.Organizer.Email)
		if err != nil {
			return nil, fmt.Errorf("unable to update cos meeting %v", err.Error())
		} else {
			log.Printf("calcom meeting rescheduled: externalId %s internalId: %s", externalSystem.ExternalID, *meeting)
			var participantsEmailsList []string
			participantsEmailsSet := make(map[string]struct{})
			for _, participant := range exMeeting.AttendedBy {
				for _, contactEmail := range participant.ContactParticipant.Emails {
					participantsEmailsSet[contactEmail.Email] = struct{}{}
					participantsEmailsList = append(participantsEmailsList, contactEmail.Email)
				}
			}
			var attendeesEmailsList []string
			attendeesEmailsSet := make(map[string]struct{})
			for _, attendee := range request.Payload.Attendees {
				attendeesEmailsSet[attendee.Email] = struct{}{}
				attendeesEmailsList = append(attendeesEmailsList, attendee.Email)
			}
			// Add attendees that are not in participants
			for _, email := range attendeesEmailsList {
				if _, found := participantsEmailsSet[email]; !found {
					contactId, err := cosService.GetContactByEmail(&request.Payload.Organizer.Email, &email)
					if err != nil {
						log.Printf("unable to find contact with email. Creating contact %s: %v", email, err.Error())
						contactId, err = cosService.CreateContact(&request.Payload.Organizer.Email, &email)
						if err != nil {
							log.Printf("Unable to create contact with email %s: %v", email, err.Error())
						} else {
							log.Printf("attendedBy: %s %s", *contactId, email)
							meetingId, err := cosService.MeetingLinkAttendedBy(exMeeting.ID, cosModel.MeetingParticipantInput{ContactID: contactId}, &request.Payload.Organizer.Email)
							if err != nil {
								log.Printf("unable to link new meeting participant: %v", err.Error())
							} else {
								log.Printf("contact participant %s added to meeting: %s", *contactId, *meetingId)
							}
						}
					} else {
						log.Printf("attendedBy: %s %s", *contactId, email)
						meetingId, err := cosService.MeetingLinkAttendedBy(exMeeting.ID, cosModel.MeetingParticipantInput{ContactID: contactId}, &request.Payload.Organizer.Email)
						if err != nil {
							log.Printf("unable to link new meeting participant: %v", err.Error())
						} else {
							log.Printf("contact participant %s added to meeting: %s", *contactId, *meetingId)
						}
					}
				}
			}

			// Remove participants that are not in attendees
			for _, email := range participantsEmailsList {
				if _, found := attendeesEmailsSet[email]; !found {
					contactId, err := cosService.GetContactByEmail(&request.Payload.Organizer.Email, &email)
					if err == nil {
						log.Printf("unlinking attendedBy: %s %s", *contactId, email)
						meetingId, err := cosService.MeetingUnLinkAttendedBy(exMeeting.ID, cosModel.MeetingParticipantInput{ContactID: contactId}, &request.Payload.Organizer.Email)
						if err != nil {
							log.Printf("unable to un link meeting participant: %v", err.Error())
						} else {
							log.Printf("contact participant %s removed to meeting: %s", *contactId, *meetingId)
						}
					}
				}
			}
			return meeting, nil
		}
	}
}

func bookingCreatedHandler(cosService s.CustomerOSService, request model.BookingCreatedRequest, body []byte, secretsRepo repository.PersonalIntegrationRepository, hSignature, calcom, appSource string) (*string, error) {
	participant, tenant, err := getParticipantAndTenant(cosService, &request.Payload.Organizer.Email)

	if err != nil {
		return nil, fmt.Errorf("meeting created handler error: %v", err.Error())
	}

	sigCheck := commonService.SignatureCheck(hSignature, body, secretsRepo, *tenant, request.Payload.Organizer.Email, calcom)
	if sigCheck != nil {
		return nil, fmt.Errorf("unable to check signature: %v", sigCheck.Error())
	}
	createdBy := []*cosModel.MeetingParticipantInput{participant}

	var attendedBy []*cosModel.MeetingParticipantInput
	for _, attendee := range request.Payload.Attendees {
		contactId, err := cosService.GetContactByEmail(&request.Payload.Organizer.Email, &attendee.Email)
		if err != nil {
			log.Printf("unable to find contact with email. Creating contact %s: %v", attendee.Email, err.Error())
			contactId, err = cosService.CreateContact(&request.Payload.Organizer.Email, &attendee.Email)
			if err != nil {
				log.Printf("Unable to create contact with email %s: %v", attendee.Email, err.Error())
			} else {
				log.Printf("attendedBy: %s %s", *contactId, attendee.Email)
				attendedBy = append(attendedBy, &cosModel.MeetingParticipantInput{ContactID: contactId})
			}

		} else {
			log.Printf("attendedBy: %s %s", *contactId, attendee.Email)
			attendedBy = append(attendedBy, &cosModel.MeetingParticipantInput{ContactID: contactId})
		}
	}

	noteInput := cosModel.NoteInput{Content: utils.StringPtr(request.Payload.AdditionalNotes), AppSource: &appSource}
	externalSystem := cosModel.ExternalSystemReferenceInput{
		ExternalID:     request.Payload.Uid,
		Type:           "CALCOM",
		ExternalURL:    &request.Payload.Metadata.VideoCallUrl,
		ExternalSource: &appSource,
	}
	input := cosModel.MeetingInput{
		Name:           &request.Payload.Title,
		CreatedBy:      createdBy,
		AttendedBy:     attendedBy,
		StartedAt:      &request.Payload.StartTime,
		EndedAt:        &request.Payload.EndTime,
		Note:           &noteInput,
		ExternalSystem: &externalSystem,
		AppSource:      &appSource,
	}
	meeting, err := cosService.CreateMeeting(input, &request.Payload.Organizer.Email)
	if err != nil {
		return nil, fmt.Errorf("unable to create meeting: %v", err.Error())
	}
	return meeting, nil
}

func getParticipantAndTenant(cosService s.CustomerOSService, userEmail *string) (*cosModel.MeetingParticipantInput, *string, error) {
	userId, err := cosService.GetUserByEmail(userEmail)
	if err != nil {
		return nil, nil, fmt.Errorf("no user for meeting creation to parse json: %v", err.Error())
	} else {
		log.Printf("createdBy: %s %s", *userId, *userEmail)
	}
	tenant, err := cosService.GetTenant(userEmail)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get tenant by for user: %v", err.Error())
	}
	return &cosModel.MeetingParticipantInput{UserID: userId}, &tenant.Tenant, nil
}
