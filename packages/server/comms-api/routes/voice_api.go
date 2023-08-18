package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes/ContactHub"
	s "github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net/http"
	"time"
)

func callEventPartyToEventParticipantInput(party *model.CallEventParty) cosModel.InteractionEventParticipantInput {
	var participantInput cosModel.InteractionEventParticipantInput
	if party.Mailto != nil {
		participantInput = cosModel.InteractionEventParticipantInput{
			Email: party.Mailto,
		}
	} else if party.Tel != nil {
		participantInput = cosModel.InteractionEventParticipantInput{
			PhoneNumber: party.Tel,
		}
	}

	return participantInput
}

func callEventPartyToSessionParticipantInput(party *model.CallEventParty) cosModel.InteractionSessionParticipantInput {
	var participantInput cosModel.InteractionSessionParticipantInput
	if party.Mailto != nil {
		participantInput = cosModel.InteractionSessionParticipantInput{
			Email: party.Mailto,
		}
	} else if party.Tel != nil {
		participantInput = cosModel.InteractionSessionParticipantInput{
			PhoneNumber: party.Tel,
		}
	}

	return participantInput
}

func getForwardingInfoFromCallEventParty(party *model.CallEventParty) *string {
	if party.Mailto != nil {
		if party.Tel != nil {
			result := "tel:" + *party.Tel
			return &result
		}
		if party.Sip != nil {
			result := "sip:" + *party.Sip
			return &result
		}
	}
	return nil
}

func getForwardingInfoFromCallEventParties(from *model.CallEventParty, to *model.CallEventParty) *string {
	result := getForwardingInfoFromCallEventParty(from)
	if result != nil {
		return result
	}

	result = getForwardingInfoFromCallEventParty(to)
	if result != nil {
		return result
	}
	return nil
}

func callEventGetOrCreateSession(threadId string, name string, tenant string, attendants []cosModel.InteractionSessionParticipantInput, cosService s.CustomerOSService) (*string, error) {
	var err error

	sessionId, err := cosService.GetInteractionSession(&threadId, &tenant, nil)
	if err != nil {
		se, _ := status.FromError(err)
		log.Printf("failed retriving interaction session: status=%s message=%s", se.Code(), se.Message())
	} else {
		return sessionId, nil
	}

	if sessionId == nil {
		sessionChannel := "VOICE"
		sessionAppSource := "COMMS_API"
		sessionStatus := "ACTIVE"
		sessionType := "CALL"
		sessionOpts := []s.SessionOption{
			s.WithSessionIdentifier(&threadId),
			s.WithSessionChannel(&sessionChannel),
			s.WithSessionName(&name),
			s.WithSessionAppSource(&sessionAppSource),
			s.WithSessionStatus(&sessionStatus),
			s.WithSessionTenant(&tenant),
			s.WithSessionAttendedBy(attendants),
			s.WithSessionType(&sessionType),
		}
		sessionId, err = cosService.CreateInteractionSession(sessionOpts...)

		if err != nil {
			se, _ := status.FromError(err)
			log.Printf("failed creating interaction session: status=%s message=%s", se.Code(), se.Message())
			return nil, fmt.Errorf("callEventGetOrCreateSession: failed creating interaction session: %v", err)
		}
		log.Printf("interaction session created: %s", *sessionId)
	}

	return sessionId, nil
}

func getCallEventContactWithIndex(req *model.CallEvent) (string, int) {
	if req.From != nil && req.From.Tel != nil {
		return *req.From.Tel, 0
	} else if req.To != nil && req.To.Tel != nil {
		return *req.To.Tel, 1
	}
	return "", 0
}

type callProgressEventInfo struct {
	sessionId *string
	tenant    *string
	sentBy    []cosModel.InteractionEventParticipantInput
	sentTo    []cosModel.InteractionEventParticipantInput
	eventType string
	eventData string
	eventTime time.Time
}

type OpenlineCallProgressData struct {
	Version      string     `json:"version,default=1.0"`
	StartTime    *time.Time `json:"start_time,omitempty"`
	AnsweredTime *time.Time `json:"answered_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	Duration     *int64     `json:"duration,omitempty"`
	SentByType   *string    `json:"sent_by_type,omitempty"`
	SentToType   *string    `json:"sent_to_type,omitempty"`
	ForwardedTo  *string    `json:"forwarded_to,omitempty"`
}

func submitCallProgressEvent(event callProgressEventInfo, cosService s.CustomerOSService) (string, error) {

	channel := "VOICE"
	appSource := "COMMS_API"
	mimeTime := "application/x-openline-call-progress"
	utcTime := event.eventTime.UTC()
	eventOpts := []s.EventOption{
		s.WithTenant(event.tenant),
		s.WithContent(&event.eventData),
		s.WithContentType(&mimeTime),
		s.WithSentBy(event.sentBy),
		s.WithSentTo(event.sentTo),
		s.WithAppSource(&appSource),
		s.WithCreatedAt(&utcTime),
		s.WithEventType(&event.eventType),
	}

	eventOpts = append(eventOpts, s.WithSessionId(event.sessionId))
	eventOpts = append(eventOpts, s.WithChannel(&channel))

	response, err := cosService.CreateInteractionEvent(eventOpts...)
	if err != nil {
		return "", fmt.Errorf("submitCallProgressEvent: failed creating interaction event: %v", err)
	}
	return response.InteractionEventCreate.Id, nil
}

func convertCallEventPartyTypeToSourceType(partyType model.CallEventPartyType) string {
	switch partyType {
	case model.CALL_EVENT_TYPE_PSTN:
		return "PSTN"
	case model.CALL_EVENT_TYPE_SIP:
		return "ESIM"
	case model.CALL_EVENT_TYPE_WEBTRC:
		return "WEBRTC"
	case model.CALL_EVENT_TYPE_VOICEMAIL:
		return "VOICEMAIL"
	default:
		return "UNKNOWN"
	}
}

func addVoiceApiRoutes(conf *c.Config, rg *gin.RouterGroup, hub *ContactHub.ContactHub, services *s.Services) {

	rg.POST("/recording", func(ctx *gin.Context) {
		multipartFileHeader, err := ctx.FormFile("audio")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to get http body: %v", err.Error()),
			})
			return
		}

		callSessionIdentifier, found := ctx.GetPostForm("correlationId")
		if !found {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("missing form field: correlationId"),
			})
			return
		}

		isActive, tenant := services.RedisService.GetKeyInfo(ctx, "tenantKey", ctx.Request.Header.Get("X-API-KEY"))

		if !isActive || tenant == nil {
			ctx.JSON(http.StatusForbidden, gin.H{"result": "Invalid API Key"})
			return
		}

		sessionId, err := services.CustomerOsService.GetInteractionSession(&callSessionIdentifier, tenant, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to get interaction session: %v", err.Error()),
			})
			return
		}

		if sessionId == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to get interaction session: %v", err.Error()),
			})
			return
		}

		fileObject, err := services.FileStoreApiService.UploadSingleFile(*tenant, multipartFileHeader)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to upload file: %v", err.Error()),
			})
			return
		}

		// link recording to interaction session
		_, err = services.CustomerOsService.AddAttachmentToInteractionSession(*sessionId, fileObject.Id, tenant, nil)
		ctx.JSON(http.StatusOK, fileObject)
	})

	rg.POST("/call_progress", func(ctx *gin.Context) {
		var req model.CallEvent
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to get http body: %v", err.Error()),
			})
			return
		}
		err = json.Unmarshal(bodyBytes, &req)
		if err != nil {
			log.Printf("unable to parse json: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		isActive, tenant := services.RedisService.GetKeyInfo(ctx, "tenantKey", ctx.Request.Header.Get("X-API-KEY"))

		if !isActive {
			ctx.JSON(http.StatusForbidden, gin.H{"result": "Invalid API Key"})
			return
		}
		threadId := req.CorrelationId

		contact, index := getCallEventContactWithIndex(&req)
		subject := ""
		if index == 0 {
			subject = fmt.Sprintf("Incoming call from %s", contact)
		} else {
			subject = fmt.Sprintf("Outgoing call to %s", contact)
		}

		var sessionParticipants []cosModel.InteractionSessionParticipantInput

		sessionParticipants = append(sessionParticipants, callEventPartyToSessionParticipantInput(req.From))
		sessionParticipants = append(sessionParticipants, callEventPartyToSessionParticipantInput(req.To))

		sessionId, err := callEventGetOrCreateSession(threadId, subject, *tenant, sessionParticipants, services.CustomerOsService)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("Unable to create InteractionSession! reasion: %v", err),
			})
			return
		}

		var ids []string

		switch req.Event {
		case "CALL_START":
			var callStartEvent model.CallEventStart
			if err := json.Unmarshal(bodyBytes, &callStartEvent); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			fromType := convertCallEventPartyTypeToSourceType(req.From.Type)
			toType := convertCallEventPartyTypeToSourceType(req.To.Type)
			eventData := OpenlineCallProgressData{
				Version:     "1.0",
				StartTime:   &callStartEvent.StartTime,
				SentByType:  &fromType,
				SentToType:  &toType,
				ForwardedTo: getForwardingInfoFromCallEventParties(req.From, req.To),
			}

			eventDataBytes, err := json.Marshal(eventData)
			if err != nil {
				log.Printf("unable to marshal json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to marshal json: %v", err.Error()),
				})
				return
			}

			eventInfo := callProgressEventInfo{
				sessionId: sessionId,
				tenant:    tenant,
				sentBy:    []cosModel.InteractionEventParticipantInput{callEventPartyToEventParticipantInput(req.From)},
				sentTo:    []cosModel.InteractionEventParticipantInput{callEventPartyToEventParticipantInput(req.To)},
				eventType: "CALL_START",
				eventData: string(eventDataBytes),
				eventTime: callStartEvent.StartTime,
			}
			id, err := submitCallProgressEvent(eventInfo, services.CustomerOsService)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("Unable to submit call progress event! reasion: %v", err),
				})
				return
			}
			ids = append(ids, id)
		case "CALL_ANSWERED":
			var callAnsweredEvent model.CallEventAnswered
			if err := json.Unmarshal(bodyBytes, &callAnsweredEvent); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			fromType := convertCallEventPartyTypeToSourceType(req.From.Type)
			toType := convertCallEventPartyTypeToSourceType(req.To.Type)
			eventData := OpenlineCallProgressData{
				Version:      "1.0",
				StartTime:    &callAnsweredEvent.StartTime,
				AnsweredTime: &callAnsweredEvent.AnsweredTime,
				SentByType:   &toType,
				SentToType:   &fromType,
				ForwardedTo:  getForwardingInfoFromCallEventParties(req.From, req.To),
			}
			eventDataBytes, err := json.Marshal(eventData)
			if err != nil {
				log.Printf("unable to marshal json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to marshal json: %v", err.Error()),
				})
				return
			}
			eventInfo := callProgressEventInfo{
				sessionId: sessionId,
				tenant:    tenant,
				sentBy:    []cosModel.InteractionEventParticipantInput{callEventPartyToEventParticipantInput(req.To)},
				sentTo:    []cosModel.InteractionEventParticipantInput{callEventPartyToEventParticipantInput(req.From)},
				eventType: "CALL_ANSWERED",
				eventData: string(eventDataBytes),
				eventTime: callAnsweredEvent.AnsweredTime,
			}
			id, err := submitCallProgressEvent(eventInfo, services.CustomerOsService)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("Unable to submit call progress event! reasion: %v", err),
				})
				return
			}
			ids = append(ids, id)
		case "CALL_END":
			var callEndEvent model.CallEventEnd
			if err := json.Unmarshal(bodyBytes, &callEndEvent); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			fromType := convertCallEventPartyTypeToSourceType(req.From.Type)
			toType := convertCallEventPartyTypeToSourceType(req.To.Type)
			eventData := OpenlineCallProgressData{
				Version:      "1.0",
				StartTime:    callEndEvent.StartTime,
				AnsweredTime: callEndEvent.AnsweredTime,
				EndTime:      &callEndEvent.EndTime,
				Duration:     &callEndEvent.Duration,
				SentByType:   &fromType,
				SentToType:   &toType,
				ForwardedTo:  getForwardingInfoFromCallEventParties(req.From, req.To),
			}
			if callEndEvent.FromCaller {
				eventData.SentByType = &fromType
				eventData.SentToType = &toType
			} else {
				eventData.SentByType = &toType
				eventData.SentToType = &fromType
			}
			eventDataBytes, err := json.Marshal(eventData)
			if err != nil {
				log.Printf("unable to marshal json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to marshal json: %v", err.Error()),
				})
				return
			}
			eventInfo := callProgressEventInfo{
				sessionId: sessionId,
				tenant:    tenant,
				eventType: "CALL_END",
				eventData: string(eventDataBytes),
				eventTime: callEndEvent.EndTime,
			}
			if callEndEvent.FromCaller {
				eventInfo.sentBy = []cosModel.InteractionEventParticipantInput{callEventPartyToEventParticipantInput(req.From)}
				eventInfo.sentTo = []cosModel.InteractionEventParticipantInput{callEventPartyToEventParticipantInput(req.To)}
			} else {
				eventInfo.sentBy = []cosModel.InteractionEventParticipantInput{callEventPartyToEventParticipantInput(req.To)}
				eventInfo.sentTo = []cosModel.InteractionEventParticipantInput{callEventPartyToEventParticipantInput(req.From)}
			}
			id, err := submitCallProgressEvent(eventInfo, services.CustomerOsService)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("Unable to submit call progress event! reasion: %v", err),
				})
				return
			}
			ids = append(ids, id)
		}

		log.Printf("message item created with ids: %v", ids)

		ctx.JSON(http.StatusOK, gin.H{
			"result": "success",
			"ids":    ids,
		})
	})
}
