package route

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func AddInteractionEventRoutes(ctx context.Context, route *gin.Engine, services *service.Services, cfg *config.Config, log logger.Logger, cache *commoncaches.Cache) {
	route.POST("/sync/postmark-interaction-event",
		tracing.TracingEnhancer(ctx, "/sync/postmark-interaction-event"),
		syncPostmarkInteractionEventHandler(services, cfg, log))
}

func syncPostmarkInteractionEventHandler(services *service.Services, cfg *config.Config, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "syncPostmarkInteractionEventHandler", c.Request.Header)
		defer span.Finish()

		//check API key as param
		apiKey := c.Query(security.ApiKeyHeader)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		key := services.CommonServices.PostgresRepositories.AppKeyRepository.FindByKey(ctx, string(security.CUSTOMER_OS_WEBHOOKS), apiKey)
		if key.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		if key.Result == nil {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		body := c.Request.Body
		requestBody, err := io.ReadAll(body)
		if err != nil {
			tracing.LogObjectAsJson(span, "body", body)
			tracing.TraceErr(span, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var postmarkEmailWebhookData model.PostmarkEmailWebhookData
		if err = json.Unmarshal(requestBody, &postmarkEmailWebhookData); err != nil {
			tracing.LogObjectAsJson(span, "requestBody", requestBody)
			tracing.TraceErr(span, err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		tracing.LogObjectAsJson(span, "webhookData", postmarkEmailWebhookData)

		pattern := `@([^.]+)\.`
		tenantNamePattern, err := regexp.Compile(pattern)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		tenantByName := ""
		for _, email := range postmarkEmailWebhookData.BccFull {
			matches := tenantNamePattern.FindStringSubmatch(email.Email)
			if len(matches) < 2 {
				continue
			}
			tenantByName = matches[1]
			break
		}

		if tenantByName == "" {
			span.LogFields(tracingLog.Bool("tenant.found", false))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		n, err := services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantByName(ctx, tenantByName)
		if err != nil {
			span.LogFields(tracingLog.Bool("tenant.found", false))
			tracing.TraceErr(span, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if n == nil {
			span.LogFields(tracingLog.Bool("tenant.found", false))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		span.LogFields(tracingLog.Bool("tenant.found", true))
		span.LogFields(tracingLog.String("tenant.name", tenantByName))
		span.SetTag(tracing.SpanTagTenant, tenantByName)

		externalSystem := "mailstack"

		participants := make([]string, 0)
		participants = append(participants, postmarkEmailWebhookData.FromFull.Email)
		for _, to := range postmarkEmailWebhookData.ToFull {
			participants = append(participants, to.Email)
		}
		if postmarkEmailWebhookData.CcFull != nil {
			for _, cc := range postmarkEmailWebhookData.CcFull {
				participants = append(participants, cc.Email)
			}
		}
		if postmarkEmailWebhookData.BccFull != nil {
			for _, bcc := range postmarkEmailWebhookData.BccFull {
				if bcc.Email != "" && bcc.Email != "bcc@"+tenantByName+".customeros.ai" {
					participants = append(participants, bcc.Email)
				}
			}
		}

		//identify mailbox
		username := ""
		for _, p := range participants {
			userByEmail, err := services.CommonServices.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenantByName, p)
			if err != nil {
				tracing.TraceErr(span, err)
				log.Errorf("(SyncInteractionEvent) error getting user by email: %s", err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			if userByEmail != nil {
				username = p
				break
			}
		}

		messageId, err := getMessageId(postmarkEmailWebhookData)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if username == "" {
			span.LogFields(tracingLog.Bool("mailbox.found", false))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		span.LogFields(tracingLog.Bool("mailbox.found", true))
		span.LogFields(tracingLog.String("mailbox.username", username))

		emailExists, err := services.CommonServices.PostgresRepositories.RawEmailRepository.EmailExistsByMessageId(ctx, externalSystem, tenantByName, username, messageId)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncInteractionEvent) error checking email exists: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if !emailExists {
			emailRawData, err := mapPostmarkToEmailRawData(tenantByName, postmarkEmailWebhookData)

			jsonContent, err := JSONMarshal(emailRawData)
			if err != nil {
				span.LogFields(tracingLog.Object("emailRawData", emailRawData))
				tracing.TraceErr(span, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			err = services.CommonServices.PostgresRepositories.RawEmailRepository.Store(ctx, externalSystem, tenantByName, username, emailRawData.ProviderMessageId, messageId, string(jsonContent), emailRawData.Sent, entity.REAL_TIME)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			if cfg.Slack.NotifyPostmarkEmail != "" {
				slackMessageText := "*From:* " + postmarkEmailWebhookData.FromFull.Email + " - " + postmarkEmailWebhookData.FromFull.Name + "\n"
				for _, t := range postmarkEmailWebhookData.ToFull {
					slackMessageText += "*To:* " + t.Email + " - " + t.Name + "\n"
				}
				for _, t := range postmarkEmailWebhookData.CcFull {
					slackMessageText += "*CC:* " + t.Email + " - " + t.Name + "\n"
				}
				for _, t := range postmarkEmailWebhookData.BccFull {
					slackMessageText += "*BCC:* " + t.Email + " - " + t.Name + "\n"
				}
				slackMessageText += "*Subject:* " + postmarkEmailWebhookData.Subject + "\n"
				slackMessageText += "*Body:* " + postmarkEmailWebhookData.HtmlBody

				utils.SendSlackMessage(ctx, cfg.Slack.NotifyPostmarkEmail, slackMessageText)
			}
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func getMessageId(pmData model.PostmarkEmailWebhookData) (string, error) {
	messageId := ""
	for _, header := range pmData.Headers {
		if header.Name == "Message-ID" || header.Name == "Message-Id" || strings.ToLower(header.Name) == "message-id" {
			messageId = header.Value
		}
	}

	if messageId == "" {
		return "", errors.New("Message-ID not found in headers")
	}
	return messageId, nil
}

func getReferences(pmData model.PostmarkEmailWebhookData) (string, error) {
	references := ""
	for _, header := range pmData.Headers {
		if header.Name == "References" {
			references = header.Value
		}
	}

	return references, nil
}

func getInReplyTo(pmData model.PostmarkEmailWebhookData) (string, error) {
	inReplyTo := ""
	for _, header := range pmData.Headers {
		if header.Name == "In-Reply-To" {
			inReplyTo = header.Value
		}
	}
	return inReplyTo, nil
}

func mapPostmarkToEmailRawData(tenant string, pmData model.PostmarkEmailWebhookData) (entity.EmailRawData, error) {
	// Parse the Date field to time.Time
	sentTime, err := utils.UnmarshalDateTime(pmData.Date)
	if err != nil {
		return entity.EmailRawData{}, err
	}

	// Map headers from slice to map
	headers := make(map[string]string)
	for _, header := range pmData.Headers {
		headers[header.Name] = header.Value
	}

	from := "<" + pmData.FromFull.Email + ">"
	to := make([]string, 0)
	for _, t := range pmData.ToFull {
		if t.Email != "" {
			to = append(to, "<"+t.Email+">")
		}
	}

	cc := make([]string, 0)
	for _, c := range pmData.CcFull {
		if c.Email != "" {
			cc = append(cc, "<"+c.Email+">")
		}
	}

	bcc := make([]string, 0)
	for _, b := range pmData.BccFull {
		if b.Email != "" && b.Email != "bcc@"+tenant+".customeros.ai" {
			bcc = append(bcc, b.Email)
		}
	}

	messageId, err := getMessageId(pmData)
	if err != nil {
		return entity.EmailRawData{}, err
	}

	references, err := getReferences(pmData)
	if err != nil {
		return entity.EmailRawData{}, err
	}

	inReplyTo, err := getInReplyTo(pmData)
	if err != nil {
		return entity.EmailRawData{}, err
	}

	return entity.EmailRawData{
		ProviderMessageId: messageId,
		MessageId:         messageId,
		Sent:              *sentTime,
		Subject:           pmData.Subject,
		From:              from,
		To:                strings.Join(to, ", "),
		Cc:                strings.Join(cc, ", "),
		Bcc:               strings.Join(bcc, ", "),
		Html:              pmData.HtmlBody,
		Text:              pmData.TextBody,
		ThreadId:          "",
		InReplyTo:         inReplyTo,
		Reference:         references,
		Headers:           headers,
	}, nil
}
