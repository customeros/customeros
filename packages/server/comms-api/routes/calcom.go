package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	s "github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/util"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"io"
	"log"
	"net/http"
)

func AddCalComRoutes(conf *c.Config, rg *gin.RouterGroup, cosService s.CustomerOSService) {
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
		//log.Printf("body: %s", body)
		hSignature := ctx.Request.Header.Get("x-cal-signature-256")
		cSignature := util.Hmac(body, []byte(conf.CalCom.CalComWebhookSecret))
		if false {
			log.Printf("Signature mismatch " + hSignature + " vs " + *cSignature)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"result": "unauthorized",
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
		if triggerEvent.TriggerEvent == "BOOKING_CREATED" {
			log.Printf("BOOKING_CREATED Trigger Event: %s", triggerEvent.TriggerEvent)

			request := model.BookingCreatedRequest{}
			if err = json.Unmarshal(body, &request); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}

			var createdBy []*cosModel.MeetingParticipantInput
			userId, err := cosService.GetUserByEmail(&request.Payload.Organizer.Email)
			if err != nil {
				log.Printf("unable to get userId by email: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("no user for meeting creation to parse json: %v", err.Error()),
				})
				return
			} else {
				log.Printf("createdBy: %s %s", *userId, request.Payload.Organizer.Email)
				createdBy = []*cosModel.MeetingParticipantInput{{UserID: userId}}
			}
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

			noteInput := cosModel.NoteInput{HTML: request.Payload.AdditionalNotes, AppSource: &appSource}
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
				AppSource:      appSource,
			}
			meeting, err := cosService.CreateMeeting(input, &request.Payload.Organizer.Email)
			if err != nil {
				log.Printf("unable to create meeting: %v", err.Error())
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"result": fmt.Sprintf("Invalid input %s", err.Error()),
				})
				return
			} else {
				log.Printf("meeting created externalId %s internalId: %s", externalSystem.ExternalID, *meeting)
				ctx.JSON(http.StatusOK, gin.H{
					"result": fmt.Sprintf("externalId %s internalId: %s", externalSystem.ExternalID, *meeting),
				})
				return
			}
		} else if triggerEvent.TriggerEvent == "BOOKING_RESCHEDULED" {
			request := model.BookingRescheduleRequest{}
			if err = json.Unmarshal(body, &request); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			log.Printf("BOOKING_RESCHEDULED Trigger Event: %s", request.TriggerEvent)
			exMeeting, err := cosService.ExternalMeeting("calcom", request.Payload.RescheduleUid, &request.Payload.Organizer.Email)
			if err != nil {
				log.Printf("unable to find external meeting: %v", err.Error())
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"result": fmt.Sprintf("Invalid input %s", err.Error()),
				})
				return
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
					AppSource:      appSource,
					ExternalSystem: &externalSystem,
				}
				meeting, err := cosService.UpdateMeeting(exMeeting.ID, input, &request.Payload.Organizer.Email)
				if err != nil {
					log.Printf("unable to update meeting: %v", err.Error())
					ctx.JSON(http.StatusUnprocessableEntity, gin.H{
						"result": fmt.Sprintf("Invalid input %s", err.Error()),
					})
					return
				} else {
					log.Printf("calcom meeting updated: externalId %s internalId: %s", externalSystem.ExternalID, *meeting)
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

					ctx.JSON(http.StatusOK, gin.H{
						"result": fmt.Sprintf("calcom meeting updated: externalId %s internalId: %s", externalSystem.ExternalID, *meeting),
					})
					return
				}
			}
		} else if triggerEvent.TriggerEvent == "BOOKING_CANCELLED" {
			request := model.BookingCancelRequest{}
			if err = json.Unmarshal(body, &request); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			log.Printf("BOOKING_CANCELLED Trigger Event: %s", request.TriggerEvent)
			exMeeting, err := cosService.ExternalMeeting("calcom", request.Payload.Uid, &request.Payload.Organizer.Email)
			if err != nil {
				log.Printf("unable to find external meetingId: %v", err.Error())
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"result": fmt.Sprintf("Invalid input %s", err.Error()),
				})
				return
			} else {

				canceled := cosModel.MeetingStatusCanceled
				input := cosModel.MeetingUpdateInput{
					Status: &canceled,
				}
				meeting, err := cosService.UpdateMeeting(exMeeting.ID, input, &request.Payload.Organizer.Email)
				if err != nil {
					log.Printf("unable to update meeting: %v", err.Error())
					ctx.JSON(http.StatusUnprocessableEntity, gin.H{
						"result": fmt.Sprintf("Invalid input %s", err.Error()),
					})
					return
				} else {
					log.Printf("calcom meeting canceled: externalId %s internalId: %s", request.Payload.Uid, *meeting)
					ctx.JSON(http.StatusOK, gin.H{
						"result": fmt.Sprintf("calcom meeting updated: externalId %s internalId: %s", request.Payload.Uid, *meeting),
					})
					return
				}
			}
		} else {
			format := "Unhandled Trigger Event: %s"
			log.Printf(format, triggerEvent.TriggerEvent)
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"result": fmt.Sprintf(format, triggerEvent.TriggerEvent),
			})
			return
		}
	})
}
