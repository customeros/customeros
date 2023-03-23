package routes

import (
	"fmt"
	"github.com/DusanKasan/parsemail"
	"github.com/gin-gonic/gin"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	s "github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"net/mail"
	"strings"
	//pbOasis "openline-ai/oasis-api/proto"
	//"strings"
)

type MailPostRequest struct {
	Sender     string `json:"sender"`
	RawMessage string `json:"rawMessage"`
	Subject    string `json:"subject"`
	ApiKey     string `json:"api-key"`
	Tenant     string `json:"X-Openline-TENANT"`
}

type EmailContent struct {
	MessageId string   `json:"messageId"`
	Html      string   `json:"html"`
	Subject   string   `json:"subject"`
	From      string   `json:"from"`
	To        []string `json:"to"`
	Cc        []string `json:"cc"`
	Bcc       []string `json:"bcc"`
	InReplyTo []string `json:"InReplyTo"`
	Reference []string `json:"Reference"`
}

func addMailRoutes(conf *c.Config, df util.DialFactory, rg *gin.RouterGroup, cosService *s.CustomerOSService) {
	mailGroup := rg.Group("/mail")
	mailGroup.POST("/fwd/", func(c *gin.Context) {
		var req MailPostRequest
		if err := c.BindJSON(&req); err != nil {
			log.Printf("unable to parse json: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		if conf.Mail.ApiKey != req.ApiKey {
			c.JSON(http.StatusForbidden, gin.H{"result": "Invalid API Key"})
			return
		}

		mailReader := strings.NewReader(req.RawMessage)
		email, err := parsemail.Parse(mailReader) // returns Email struct and error
		if err != nil {
			log.Printf("Unable to parse Email: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("Unable to parse Email: %v", err.Error()),
			})
			return
		}

		if len(email.From) != 1 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("Email has more than one From: %v", email.From),
			})
			return
		}

		//  fromAddress := email.From[0].Address
		//  emailContent := EmailContent{
		//	MessageId: ensureRfcId(email.MessageID),
		//	Subject:   email.Subject,
		//	Html:      firstNotEmpty(email.HTMLBody, email.TextBody),
		//	From:      fromAddress,
		//	To:        toStringArr(email.To),
		//	Cc:        toStringArr(email.Cc),
		//	Bcc:       toStringArr(email.Bcc),
		//	InReplyTo: ensureRfcIds(email.InReplyTo),
		//	Reference: ensureRfcIds(email.References),
		//}
		//jsonContent, err := json.Marshal(emailContent)
		//if err != nil {
		//	se, _ := status.FromError(err)
		//	c.JSON(http.StatusInternalServerError, gin.H{
		//		"result": fmt.Sprintf("failed creating message content: status=%s message=%s", se.Code(), se.Message()),
		//	})
		//	return
		//}
		//Contact the server and print out its response.
		//jsonContentString := string(jsonContent)
		refSize := len(email.References)
		threadId := ""
		if refSize > 0 {
			threadId = ensureRfcId(email.References[0])
		} else {
			threadId = ensureRfcId(email.MessageID)
		}
		ctx := context.Background()
		sessionId, err := cosService.GetInteractionSession(ctx, threadId, req.Tenant)
		if err != nil {
			se, _ := status.FromError(err)
			log.Printf("failed retriving interaction session: status=%s message=%s", se.Code(), se.Message())
			return
		}
		if sessionId == nil {
			sessionId, err = cosService.CreateInteractionSession(ctx,
				s.WithSessionIdentifier(threadId),
				s.WithSessionChannel("EMAIL"),
				s.WithSessionName(email.Subject),
				s.WithSessionAppSource("CHANNELS"),
				s.WithSessionStatus("ACTIVE"),
				s.WithSessionTenant(req.Tenant))
			if err != nil {
				se, _ := status.FromError(err)
				log.Printf("failed creating interaction session: status=%s message=%s", se.Code(), se.Message())
				return
			}
			log.Printf("interaction session created: %s", sessionId)
		}

		participantTypeTO, participantTypeCC := "TO", "CC"
		response, err := cosService.CreateInteractionEvent(ctx,
			s.WithTenant(req.Tenant),
			s.WithSessionId(*sessionId),
			s.WithChannel("EMAIL"),
			s.WithContent(firstNotEmpty(email.HTMLBody, email.TextBody)),
			s.WithContentType(email.ContentType),
			s.WithSentBy(toParticipantInputArr(email.From, nil)),
			s.WithSentTo(append(toParticipantInputArr(email.To, &participantTypeTO), toParticipantInputArr(email.Cc, &participantTypeCC)...)),
		)

		if err != nil {
			se, _ := status.FromError(err)
			log.Printf("failed creating interaction event: status=%s message=%s", se.Code(), se.Message())
			return
		}

		log.Printf("interaction event created with id: %s", (*response).InteractionEventCreate.Id)

		////Set up a connection to the oasis-api server.
		//oasisConn := GetOasisClient(c, df)
		//defer closeOasisConnection(oasisConn)
		//oasisClient := o.NewOasisApiServiceClient(oasisConn)
		//_, mEventErr := oasisClient.NewMessageEvent(ctx, &o.NewMessage{ConversationId: savedMessage.ConversationId, ConversationItemId: savedMessage.GetConversationEventId()})
		//if mEventErr != nil {
		//	se, _ := status.FromError(mEventErr)
		//	log.Printf("failed new message event: status=%s message=%s", se.Code(), se.Message())
		//}

		c.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprintf("interaction event created with id: %s", "aaa"),
		})
	})
}

func ensureRfcId(id string) string {
	if !strings.HasPrefix(id, "<") {
		id = fmt.Sprintf("<%s>", id)
	}
	return id
}

func ensureRfcIds(to []string) []string {
	var result []string
	for _, id := range to {
		result = append(result, ensureRfcId(id))
	}
	return result
}

func toStringArr(from []*mail.Address) []string {
	var to []string
	for _, a := range from {
		to = append(to, a.Address)
	}
	return to
}

func toParticipantInputArr(from []*mail.Address, participantType *string) []s.InteractionEventParticipantInput {
	var to []s.InteractionEventParticipantInput
	for _, a := range from {
		participantInput := s.InteractionEventParticipantInput{
			Email:           &a.Address,
			ParticipantType: participantType,
		}
		to = append(to, participantInput)
	}
	return to
}

func firstNotEmpty(input ...string) string {
	for _, item := range input {
		if item != "" {
			return item
		}
	}
	return ""
}
