package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	ms "github.com/openline-ai/openline-customer-os/packages/server/message-store-api/proto/generated"
	c "github.com/openline-ai/openline-oasis/packages/server/channels-api/config"
	"github.com/openline-ai/openline-oasis/packages/server/channels-api/model"
	"github.com/openline-ai/openline-oasis/packages/server/channels-api/util"
	o "github.com/openline-ai/openline-oasis/packages/server/oasis-api/proto/generated"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

func encodePartyToParticipantId(party *model.VConParty) *ms.ParticipantId {
	if party.Mailto != nil {
		return &ms.ParticipantId{
			Type:       ms.ParticipantIdType_MAILTO,
			Identifier: *party.Mailto,
		}
	} else if party.Tel != nil {
		return &ms.ParticipantId{
			Type:       ms.ParticipantIdType_TEL,
			Identifier: *party.Tel,
		}
	}
	return nil
}

func getInitator(req *model.VCon) *ms.ParticipantId {
	if len(req.Parties) == 0 {
		return nil
	}

	if len(req.Analysis) != 0 {
		return encodePartyToParticipantId(&req.Parties[0])
	}

	if len(req.Dialog) == 0 {
		return nil
	}
	if len(req.Dialog[0].Parties) == 0 {
		return encodePartyToParticipantId(&req.Parties[0])
	}
	if req.Dialog[0].Parties[0] > int64(len(req.Parties)) {
		return encodePartyToParticipantId(&req.Parties[0])
	}
	return encodePartyToParticipantId(&req.Parties[req.Dialog[0].Parties[0]])
}

func getDirection(req *model.VCon) ms.MessageDirection {
	initator := getInitator(req)
	if initator == nil {
		return ms.MessageDirection_INBOUND
	}

	if initator.Type == ms.ParticipantIdType_MAILTO {
		return ms.MessageDirection_OUTBOUND
	}
	return ms.MessageDirection_INBOUND
}

func getUser(req *model.VCon) string {

	for _, p := range req.Parties {
		if p.Mailto != nil {
			return *p.Mailto
		}
	}
	return ""
}

type VConEvent struct {
	Parties  []model.VConParty   `json:"parties,omitempty"`
	Dialog   *model.VConDialog   `json:"dialog,omitempty"`
	Analysis *model.VConAnalysis `json:"analysis,omitempty"`
}

func makeMessage(req *model.VCon) *VConEvent {
	res := &VConEvent{}
	if req.Dialog != nil && len(req.Dialog) > 0 {
		res.Dialog = &req.Dialog[0]
	}
	if req.Analysis != nil && len(req.Analysis) > 0 {
		res.Analysis = &req.Analysis[0]
	}
	res.Parties = req.Parties
	return res
}

func AddVconRoutes(conf *c.Config, df util.DialFactory, rg *gin.RouterGroup) {
	rg.POST("/vcon", func(c *gin.Context) {
		var req model.VCon
		if err := c.BindJSON(&req); err != nil {
			log.Printf("unable to parse json: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		if conf.VCon.ApiKey != c.GetHeader("X-Openline-VCon-Api-Key") {
			c.JSON(http.StatusForbidden, gin.H{"result": "Invalid API Key"})
			return
		}
		threadId := req.UUID
		if req.Appended != nil {
			threadId = req.Appended.UUID
		}
		participants := make([]*ms.ParticipantId, len(req.Parties))
		for i, p := range req.Parties {
			participants[i] = encodePartyToParticipantId(&p)
		}
		initator := getInitator(&req)
		if initator == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "malformed vCon document!",
			})
			return
		}
		contentObject := makeMessage(&req)
		content, err := json.Marshal(contentObject)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to marshal json: %v", err.Error()),
			})
			return
		}
		contentStr := string(content)
		message := &ms.InputMessage{
			Type:                    ms.MessageType_VOICE,
			Subtype:                 ms.MessageSubtype_MESSAGE,
			Content:                 &contentStr,
			Direction:               getDirection(&req),
			InitiatorIdentifier:     initator,
			ThreadId:                &threadId,
			ParticipantsIdentifiers: participants,
		}
		//Store the message in message store
		msConn := util.GetMessageStoreConnection(c, df)
		defer util.CloseMessageStoreConnection(msConn)
		msClient := ms.NewMessageStoreServiceClient(msConn)

		ctx := context.Background()
		ctx = metadata.AppendToOutgoingContext(ctx, service.ApiKeyHeader, conf.Service.MessageStoreApiKey)
		ctx = metadata.AppendToOutgoingContext(ctx, service.UsernameHeader, getUser(&req))

		savedMessage, err := msClient.SaveMessage(ctx, message)
		if err != nil {
			se, _ := status.FromError(err)
			log.Printf("failed creating message item: status=%s message=%s", se.Code(), se.Message())
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("failed creating message item: status=%s message=%s", se.Code(), se.Message()),
			})
			return
		}
		log.Printf("message item created with id: %s", savedMessage.GetConversationEventId())

		//Set up a connection to the oasis-api server.
		oasisConn := GetOasisClient(c, df)
		defer closeOasisConnection(oasisConn)
		oasisClient := o.NewOasisApiServiceClient(oasisConn)
		_, mEventErr := oasisClient.NewMessageEvent(ctx, &o.NewMessage{ConversationId: savedMessage.ConversationId, ConversationItemId: savedMessage.GetConversationEventId()})
		if mEventErr != nil {
			se, _ := status.FromError(mEventErr)
			log.Printf("failed new message event: status=%s message=%s", se.Code(), se.Message())
		}

		c.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprintf("message item created with id: %s", savedMessage.GetConversationEventId()),
		})
	})
}
