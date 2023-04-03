package routes

import (
	"errors"
	"fmt"
	"github.com/DusanKasan/parsemail"
	"github.com/gin-gonic/gin"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	s "github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"log"
	"net/http"
	"strings"
)

func addMailRoutes(conf *c.Config, rg *gin.RouterGroup, mailService s.MailService) {
	rg.POST("/mail/send", func(c *gin.Context) {
		var request model.MailReplyRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		if conf.Mail.ApiKey != c.GetHeader("X-Openline-Mail-Api-Key") {
			errorMsg := "invalid mail API Key!"
			log.Printf(errorMsg)
			c.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
			return
		}

		username := c.GetHeader("X-Openline-USERNAME")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "username header not found"})
			return
		}

		identityId := c.GetHeader("X-Openline-IDENTITY-ID")
		if identityId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "identity header not found"})
			return
		}
		replyMail, err := mailService.SendMail(&request, &username, &identityId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		mail, err := mailService.SaveMail(replyMail, nil, &request.Username)
		if err != nil {
			errorMsg := fmt.Sprintf("unable to save email: %v", err.Error())
			log.Printf(errorMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprintf("interaction event created with id: %s", (*mail).InteractionEventCreate.Id),
		})

	})

	rg.POST("/mail/fwd/", func(c *gin.Context) {
		var req model.MailFwdRequest
		if err := c.BindJSON(&req); err != nil {
			errorMsg := fmt.Sprintf("unable to parse json: %v", err.Error())
			log.Printf(errorMsg)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errorMsg,
			})
			return
		}

		if conf.Mail.ApiKey != req.ApiKey {
			errorMsg := "invalid mail API Key!"
			log.Printf(errorMsg)
			c.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
			return
		}

		if err := validateMailPostRequest(req); err != nil {
			log.Printf("Invalid request: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		email, err := parsemail.Parse(strings.NewReader(req.RawMessage))
		if err != nil {
			errorMsg := fmt.Sprintf("unable to parse email: %v", err.Error())
			log.Printf(errorMsg)
			c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
			return
		}

		mail, err := mailService.SaveMail(&email, &req.Tenant, nil)
		if err != nil {
			errorMsg := fmt.Sprintf("unable to save forwarded email: %v", err.Error())
			log.Printf(errorMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprintf("interaction event created with id: %s", (*mail).InteractionEventCreate.Id),
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
