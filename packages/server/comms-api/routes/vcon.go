package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes/ContactHub"
	s "github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

func vConPartyToEventParticipantInputArr(from []model.VConParty) []model.InteractionEventParticipantInput {
	var to = []model.InteractionEventParticipantInput{}
	for _, a := range from {
		if a.Mailto != nil {
			participantInput := model.InteractionEventParticipantInput{
				Email: a.Mailto,
			}
			to = append(to, participantInput)
		} else if a.Tel != nil {
			participantInput := model.InteractionEventParticipantInput{
				PhoneNumber: a.Tel,
			}
			to = append(to, participantInput)
		}
	}
	return to
}

func vConPartyToSessionParticipantInputArr(from []model.VConParty) []model.InteractionSessionParticipantInput {
	var to []model.InteractionSessionParticipantInput
	for _, a := range from {
		if a.Mailto != nil {
			participantInput := model.InteractionSessionParticipantInput{
				Email: a.Mailto,
			}
			to = append(to, participantInput)
		} else if a.Tel != nil {
			participantInput := model.InteractionSessionParticipantInput{
				PhoneNumber: a.Tel,
			}
			to = append(to, participantInput)
		} else if a.ContactId != nil {
			participantInput := model.InteractionSessionParticipantInput{
				ContactID: a.ContactId,
			}
			to = append(to, participantInput)
		} else if a.UserId != nil {
			participantInput := model.InteractionSessionParticipantInput{
				UserID: a.UserId,
			}
			to = append(to, participantInput)
		}
	}
	return to
}

func getDestination(parties []model.VConParty, dialog *model.VConDialog) *model.VConParty {
	if len(parties) == 0 {
		return nil
	}

	if len(parties) == 1 {
		return &parties[0]
	}

	if len(dialog.Parties) == 0 {
		return &parties[1]
	}
	if dialog.Parties[0] == 0 {
		return &parties[1]
	}

	return &parties[0]
}

func getInitator(parties []model.VConParty, dialog *model.VConDialog) *model.VConParty {
	if len(parties) == 0 {
		return nil
	}

	if len(parties) == 1 {
		return &parties[0]
	}

	if len(dialog.Parties) == 0 {
		return &parties[0]
	}

	if dialog.Parties[0] == 0 {
		return &parties[0]
	}
	return &parties[1]
}

func vConGetOrCreateSession(threadId string, name string, user string, attendants []model.InteractionSessionParticipantInput, cosService s.CustomerOSService) (*string, error) {
	var err error

	sessionId, err := cosService.GetInteractionSession(&threadId, nil, &user)
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
			s.WithSessionUsername(&user),
			s.WithSessionAttendedBy(attendants),
			s.WithSessionType(&sessionType),
		}
		sessionId, err = cosService.CreateInteractionSession(sessionOpts...)

		if err != nil {
			se, _ := status.FromError(err)
			log.Printf("failed creating interaction session: status=%s message=%s", se.Code(), se.Message())
			return nil, fmt.Errorf("vConGetOrCreateSession: failed creating interaction session: %v", err)
		}
		log.Printf("interaction session created: %s", *sessionId)
	}

	return sessionId, nil
}

func getUser(ctx *gin.Context, req *model.VCon) string {

	usernameHeader := ctx.GetHeader("X-Openline-USERNAME")

	if usernameHeader != "" {
		return usernameHeader
	}

	for _, p := range req.Parties {
		if p.Mailto != nil {
			return *p.Mailto
		}
	}
	return ""
}

func getContactWithIndex(req *model.VCon) (string, int) {
	for i, p := range req.Parties {
		if p.Tel != nil {
			return *p.Tel, i
		}
	}
	return "", 0
}

type VConEvent struct {
	Parties  []model.VConParty   `json:"parties,omitempty"`
	Dialog   *model.VConDialog   `json:"dialog,omitempty"`
	Analysis *model.VConAnalysis `json:"analysis,omitempty"`
}

func submitAnalysis(sessionId *string, req model.VCon, cosService s.CustomerOSService, ctx *gin.Context) ([]string, error) {
	user := getUser(ctx, &req)

	var ids []string
	for _, analysis := range req.Analysis {
		analysisType := string(analysis.Type)
		appSource := "COMMS_API"
		analysisOpts := []s.AnalysisOption{
			s.WithAnalysisUsername(&user),
			s.WithAnalysisType(&analysisType),
			s.WithAnalysisContent(&analysis.Body),
			s.WithAnalysisContentType(&analysis.MimeType),
			s.WithAnalysisAppSource(&appSource),
		}
		if req.Type != nil && *req.Type == model.MEETING {
			analysisOpts = append(analysisOpts, s.WithAnalysisDescribes(&model.AnalysisDescriptionInput{MeetingId: sessionId}))
		} else {
			analysisOpts = append(analysisOpts, s.WithAnalysisDescribes(&model.AnalysisDescriptionInput{InteractionSessionId: sessionId}))
		}

		response, err := cosService.CreateAnalysis(analysisOpts...)
		if err != nil {
			return nil, fmt.Errorf("submitDialog: failed creating interaction event: %v", err)
		}
		ids = append(ids, *response)
	}
	return ids, nil
}

func submitAttachmentsToEvent(eventId *string, req model.VCon, cosService s.CustomerOSService, ctx *gin.Context) ([]string, error) {
	user := getUser(ctx, &req)

	var ids []string
	for _, attachment := range req.Attachments {

		if attachment.MimeType == "application/x-openline-file-store-id" {
			if len(req.Dialog) > 0 {
				log.Printf("submitAttachments: adding attachment to interaction event: %s eventId: %s", attachment.Body, eventId)
				response, err := cosService.AddAttachmentToInteractionEvent(*eventId, attachment.Body, nil, &user)
				if err != nil {
					return nil, fmt.Errorf("submitAttachments: failed failed to link attachment to interaction event: %v", err)
				}
				ids = append(ids, *response)
			}

		} else {
			return nil, fmt.Errorf("submitAttachments: unsupported attachment type: %s", attachment.MimeType)
		}
	}
	return ids, nil
}

func submitDialog(sessionId *string, req model.VCon, cosService s.CustomerOSService, ctx *gin.Context) ([]string, error) {
	user := getUser(ctx, &req)

	var ids []string
	for _, d := range req.Dialog {
		initator := getInitator(req.Parties, &d)
		if initator == nil {
			return nil, fmt.Errorf("submitDialog: unable to determine initator")
		}
		destination := getDestination(req.Parties, &d)
		if destination == nil {
			return nil, fmt.Errorf("submitDialog: unable to determine destination")
		}

		channel := "VOICE"
		appSource := "COMMS_API"
		eventOpts := []s.EventOption{
			s.WithUsername(&user),
			s.WithChannel(&channel),
			s.WithContent(&d.Body),
			s.WithContentType(&d.MimeType),
			s.WithSentBy(vConPartyToEventParticipantInputArr([]model.VConParty{*initator})),
			s.WithSentTo(vConPartyToEventParticipantInputArr([]model.VConParty{*destination})),
			s.WithAppSource(&appSource),
		}

		if req.Type != nil && *req.Type == model.MEETING {
			eventOpts = append(eventOpts, s.WithMeetingId(sessionId))
		} else {
			eventOpts = append(eventOpts, s.WithSessionId(sessionId))
		}
		response, err := cosService.CreateInteractionEvent(eventOpts...)
		if err != nil {
			return nil, fmt.Errorf("submitDialog: failed creating interaction event: %v", err)
		}
		ids = append(ids, response.InteractionEventCreate.Id)
	}
	return ids, nil
}

func addVconRoutes(conf *c.Config, rg *gin.RouterGroup, cosService s.CustomerOSService, hub *ContactHub.ContactHub) {
	rg.POST("/vcon", func(ctx *gin.Context) {
		var req model.VCon
		if err := ctx.BindJSON(&req); err != nil {
			log.Printf("unable to parse json: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		if conf.VCon.ApiKey != ctx.GetHeader("X-Openline-VCon-Api-Key") {
			ctx.JSON(http.StatusForbidden, gin.H{"result": "Invalid API Key"})
			return
		}
		threadId := req.UUID
		if req.Appended != nil {
			threadId = req.Appended.UUID
		}

		contact, index := getContactWithIndex(&req)
		subject := ""
		if index == 0 {
			subject = fmt.Sprintf("Incoming call from %s", contact)
		} else {
			subject = fmt.Sprintf("Outgoing call to %s", contact)
		}
		var sessionId *string = nil
		var err error = nil

		if req.Type != nil && *req.Type == model.MEETING {
			if req.Appended != nil {
				sessionId = &req.Appended.UUID
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("Must specify a meeting id for type meeting"),
				})
				return
			}
		} else {
			sessionId, err = vConGetOrCreateSession(threadId, subject, getUser(ctx, &req), vConPartyToSessionParticipantInputArr(req.Parties), cosService)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("Unable to create InteractionSession! reasion: %v", err),
				})
				return
			}
		}

		var ids []string
		if req.Analysis != nil && len(req.Analysis) > 0 {
			newIds, err := submitAnalysis(sessionId, req, cosService, ctx)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("Unable to submit analysis! reasion: %v", err),
				})
				return
			}
			ids = append(ids, newIds...)
		}

		if req.Dialog != nil && len(req.Dialog) > 0 {
			newIds, err := submitDialog(sessionId, req, cosService, ctx)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("Unable to submit dialog! reason: %v", err),
				})
				return
			}
			if req.Attachments != nil && len(req.Attachments) > 0 {
				for _, newId := range newIds {
					_, err = submitAttachmentsToEvent(&newId, req, cosService, ctx)
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{
							"result": fmt.Sprintf("Unable to submit attachments! reason: %v", err),
						})
						return
					}
				}
			}

			ids = append(ids, newIds...)
		}

		log.Printf("message item created with ids: %v", ids)

		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprintf("message item created with ids: %v", ids),
		})
	})
}
