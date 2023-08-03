package routes

import (
	"errors"
	"fmt"
	"github.com/DusanKasan/parsemail"
	"github.com/gin-gonic/gin"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes/ContactHub"
	s "github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/tracing"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"log"
	"net/http"
	"strings"
)

func addMailRoutes(conf *c.Config, rg *gin.RouterGroup, mailService s.MailService, hub *ContactHub.ContactHub) {
	rg.POST("/mail/send", func(c *gin.Context) {
		span, _ := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mail/send", c.Request.Header)
		defer span.Finish()

		var request model.MailReplyRequest

		if err := c.BindJSON(&request); err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		if conf.Mail.ApiKey != c.GetHeader("X-Openline-Mail-Api-Key") {
			errorMsg := "invalid mail API Key!"
			tracing.TraceErr(span, errors.New(errorMsg))
			log.Printf(errorMsg)
			c.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
			return
		}

		username := c.GetHeader("X-Openline-USERNAME")
		if username == "" {
			errMsg := "username header not found"
			tracing.TraceErr(span, errors.New(errMsg))
			c.JSON(http.StatusBadRequest, gin.H{"msg": errMsg})
			return
		}

		replyMail, err := mailService.SendMail(&request, &username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		mail, err := mailService.SaveMail(replyMail, nil, &request.Username)
		if err != nil {
			tracing.TraceErr(span, err)
			errorMsg := fmt.Sprintf("unable to save email: %v", err.Error())
			log.Printf(errorMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
			return
		}

		span.LogFields(tracingLog.String("result - interactionEventId", (*mail).InteractionEventCreate.Id))
		c.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprintf("interaction event created with id: %s", (*mail).InteractionEventCreate.Id),
		})

	})

	rg.POST("/mail/fwd/", func(c *gin.Context) {
		span, _ := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mail/fwd/", c.Request.Header)
		defer span.Finish()

		var req model.MailFwdRequest
		if err := c.BindJSON(&req); err != nil {
			tracing.TraceErr(span, err)
			errorMsg := fmt.Sprintf("unable to parse json: %v", err.Error())
			log.Printf(errorMsg)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errorMsg,
			})
			return
		}

		if conf.Mail.ApiKey != req.ApiKey {
			errorMsg := "invalid mail API Key!"
			tracing.TraceErr(span, errors.New(errorMsg))
			log.Printf(errorMsg)
			c.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
			return
		}

		if err := validateMailPostRequest(req); err != nil {
			tracing.TraceErr(span, err)
			log.Printf("Invalid request: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		email, err := parsemail.Parse(strings.NewReader(req.RawMessage))
		if err != nil {
			tracing.TraceErr(span, err)
			errorMsg := fmt.Sprintf("unable to parse email: %v", err.Error())
			log.Printf(errorMsg)
			c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
			return
		}

		saveResponse, err := mailService.SaveMail(&email, &req.Tenant, nil)
		if err != nil {
			tracing.TraceErr(span, err)
			errorMsg := fmt.Sprintf("unable to save forwarded email: %v", err.Error())
			log.Printf(errorMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
			return
		}

		for _, participant := range saveResponse.InteractionEventCreate.SentTo {
			contacts := participant.EmailParticipant.Contacts
			for _, contact := range contacts {
				log.Printf("Broadcasting to participant %s", contact.Id)
				contactItem := ContactHub.ContactEvent{
					ContactId:        contact.Id,
					InteractionEvent: saveResponse.InteractionEventCreate,
				}

				hub.Broadcast <- contactItem
			}
		}

		span.LogFields(tracingLog.String("result - interactionEventId", (*saveResponse).InteractionEventCreate.Id))

		c.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprintf("interaction event created with id: %s", (*saveResponse).InteractionEventCreate.Id),
		})
	})
}

func validateMailPostRequest(req model.MailFwdRequest) error {
	// Add validation checks for other fields in the MailPostRequest struct
	if req.Tenant == "" {
		return errors.New("missing tenant field")
	}
	if req.RawMessage == "" {
		return errors.New("missing raw message field")
	}
	return nil
}
