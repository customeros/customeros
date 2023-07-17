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
	rg.POST("/calcom", func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Printf("unable to read body: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}
		hSignature := ctx.Request.Header.Get("x-cal-signature-256")
		cSignature := util.Hmac(body, []byte(conf.CalCom.CalComWebhookSecret))
		if hSignature != *cSignature {
			log.Printf("Signature mismatch " + hSignature + " vs " + *cSignature)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"result": "unauthorized",
			})
			return
		}

		var request model.CalcomRequest
		if err = json.Unmarshal(body, &request); err != nil {
			log.Printf("unable to parse json: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		if request.TriggerEvent == "BOOKING_CREATED" {
			log.Printf("BOOKING_CREATED Trigger Event: %s", request.TriggerEvent)

			var createdBy []cosModel.MeetingParticipantInput
			userId, err := cosService.GetUserByEmail(&request.Payload.Organizer.Email)
			if err != nil {
				log.Printf("unable to get userId by email: %v", err.Error())
			} else {
				log.Printf("createdBy: %s %s", *userId, request.Payload.Organizer.Email)
				createdBy = []cosModel.MeetingParticipantInput{{UserID: userId}}
			}
			var attendedBy []cosModel.MeetingParticipantInput
			for _, attendee := range request.Payload.Attendees {
				contactId, err := cosService.GetContactByEmail(&request.Payload.Organizer.Email, &attendee.Email)
				if err != nil {
					log.Printf("unable to find contact with email %s: %v", attendee.Email, err.Error())
				} else {
					log.Printf("attendedBy: %s %s", *contactId, attendee.Email)
					attendedBy = append(attendedBy, cosModel.MeetingParticipantInput{ContactID: contactId})
				}
			}
			appSource := "calcom"
			noteInput := cosModel.NoteInput{HTML: request.Payload.AdditionalNotes, AppSource: &appSource}
			externalSystem := cosModel.ExternalSystemReferenceInput{
				ExternalID:     request.Payload.Uid,
				Type:           "CALCOM",
				ExternalURL:    &request.Payload.Metadata.VideoCallUrl,
				ExternalSource: &appSource,
			}
			meetingOptions := []s.MeetingOption{
				s.WithMeetingName(&request.Payload.Title),
				s.WithMeetingAppSource(&appSource),
				s.WithMeetingStartedAt(&request.Payload.StartTime),
				s.WithMeetingEndedAt(&request.Payload.EndTime),
				s.WithMeetingAttendedBy(attendedBy),
				s.WithMeetingCreatedBy(createdBy),
				s.WithMeetingUsername(&request.Payload.Organizer.Email),
				s.WithMeetingNote(&noteInput),
				s.WithExternalSystem(&externalSystem),
			}
			meeting, err := cosService.CreateMeeting(meetingOptions...)
			if err != nil {
				log.Printf("unable to create meeting: %v", err.Error())
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"result": fmt.Sprintf("Invalid input %s", err.Error()),
				})
				return
			} else {
				log.Printf("meeting created with id: %s", *meeting)
				ctx.JSON(http.StatusOK, gin.H{
					"result": fmt.Sprintf("meeting created with id: %s", *meeting),
				})
				return
			}
		} else {
			format := "Unhandled Trigger Event: %s"
			log.Printf(format, request.TriggerEvent)
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"result": fmt.Sprintf(format, request.TriggerEvent),
			})
			return
		}
	})
}
