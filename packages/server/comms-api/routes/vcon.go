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
	var to []model.InteractionEventParticipantInput
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

func vConGetOrCreateSession(threadId string, name string, user string, attendants []model.InteractionSessionParticipantInput, cosService s.CustomerOSService) (string, error) {
	var err error

	sessionId, err := cosService.GetInteractionSession(&threadId, nil, &user)
	if err != nil {
		se, _ := status.FromError(err)
		log.Printf("failed retriving interaction session: status=%s message=%s", se.Code(), se.Message())
	} else {
		return *sessionId, nil
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
			return "", fmt.Errorf("vConGetOrCreateSession: failed creating interaction session: %v", err)
		}
		log.Printf("interaction session created: %s", *sessionId)
	}

	return *sessionId, nil
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

func submitAnalysis(sessionId string, req model.VCon, cosService s.CustomerOSService, ctx *gin.Context) ([]string, error) {
	user := getUser(ctx, &req)

	var ids []string
	for _, a := range req.Analysis {
		analysisType := string(a.Type)
		appSource := "COMMS_API"
		analysisOpts := []s.AnalysisOption{
			s.WithAnalysisUsername(&user),
			s.WithAnalysisType(&analysisType),
			s.WithAnalysisContent(&a.Body),
			s.WithAnalysisContentType(&a.MimeType),
			s.WithAnalysisDescribes(&model.AnalysisDescriptionInput{InteractionSessionId: &sessionId}),
			s.WithAnalysisAppSource(&appSource),
		}
		response, err := cosService.CreateAnalysis(analysisOpts...)
		if err != nil {
			return nil, fmt.Errorf("submitDialog: failed creating interaction event: %v", err)
		}
		ids = append(ids, *response)
	}
	return ids, nil
}

func submitAttachments(sessionId string, req model.VCon, cosService s.CustomerOSService, ctx *gin.Context) ([]string, error) {
	user := getUser(ctx, &req)

	var ids []string
	for _, attachment := range req.Attachments {

		if attachment.MimeType == "application/x-openline-file-store-id" {
			log.Printf("submitAttachments: adding attachment to interaction session: %s sessionId: %s", attachment.Body, sessionId)
			response, err := cosService.AddAttachmentToInteractionSession(sessionId, attachment.Body, nil, &user)
			if err != nil {
				return nil, fmt.Errorf("submitAttachments: failed failed to link attachment to interaction event: %v", err)
			}
			ids = append(ids, *response)
		} else {
			return nil, fmt.Errorf("submitAttachments: unsupported attachment type: %s", attachment.MimeType)
		}
	}
	return ids, nil
}

func submitDialog(sessionId string, req model.VCon, cosService s.CustomerOSService, ctx *gin.Context) ([]string, error) {
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
			s.WithSessionId(&sessionId),
			s.WithChannel(&channel),
			s.WithContent(&d.Body),
			s.WithContentType(&d.MimeType),
			s.WithSentBy(vConPartyToEventParticipantInputArr([]model.VConParty{*initator})),
			s.WithSentTo(vConPartyToEventParticipantInputArr([]model.VConParty{*destination})),
			s.WithAppSource(&appSource),
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

		sessionId, err := vConGetOrCreateSession(threadId, subject, getUser(ctx, &req), vConPartyToSessionParticipantInputArr(req.Parties), cosService)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("Unable to create InteractionSession! reasion: %v", err),
			})
			return
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
					"result": fmt.Sprintf("Unable to submit dialog! reasion: %v", err),
				})
				return
			}
			ids = append(ids, newIds...)
		}

		if req.Attachments != nil && len(req.Attachments) > 0 {
			_, err = submitAttachments(sessionId, req, cosService, ctx)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("Unable to submit attachments! reasion: %v", err),
				})
				return
			}
		}

		log.Printf("message item created with ids: %v", ids)

		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprintf("message item created with ids: %v", ids),
		})
	})
}
