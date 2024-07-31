package route

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/config"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
)

func AddInteractionEventRoutes(ctx context.Context, route *gin.Engine, services *service.Services, cfg *config.Config, log logger.Logger, cache *commoncaches.Cache) {
	route.POST("/sync/interaction-events",
		handler.TracingEnhancer(ctx, "/sync/interaction-events"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_WEBHOOKS, security.WithCache(cache)),
		syncInteractionEventsHandler(services, log))
	route.POST("/sync/interaction-event",
		handler.TracingEnhancer(ctx, "/sync/interaction-event"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_WEBHOOKS, security.WithCache(cache)),
		syncInteractionEventHandler(services, log))
	route.POST("/sync/postmark-interaction-event",
		handler.TracingEnhancer(ctx, "/sync/postmark-interaction-event"),
		syncPostmarkInteractionEventHandler(services, cfg, log))
}

func syncInteractionEventsHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncInteractionEvents", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeMessages)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncInteractionEvents) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var interactionEvents []model.InteractionEventData
		if err = json.Unmarshal(requestBody, &interactionEvents); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncInteractionEvents) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		if len(interactionEvents) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing interactionEvents in request"})
			return
		}

		// Context timeout, allocate per interactionEvent
		timeout := time.Duration(len(interactionEvents)) * common.Min1Duration
		if timeout > constants.RequestMaxTimeout {
			timeout = constants.RequestMaxTimeout
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		syncResult, err := services.InteractionEventService.SyncInteractionEvents(ctx, interactionEvents)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncInteractionEvents) error in sync interactionEvents: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing log entries"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
}

func syncInteractionEventHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncInteractionEvent", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeMessages)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncInteractionEvent) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var interactionEvent model.InteractionEventData
		if err = json.Unmarshal(requestBody, &interactionEvent); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncInteractionEvents) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		// Context timeout, allocate per interactionEvent
		ctx, cancel := context.WithTimeout(ctx, common.Min1Duration)
		defer cancel()

		syncResult, err := services.InteractionEventService.SyncInteractionEvents(ctx, []model.InteractionEventData{interactionEvent})
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncInteractionEvent) error in sync interactionEvent: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing log entry"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
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

		if username == "" {
			span.LogFields(tracingLog.Bool("mailbox.found", false))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		span.LogFields(tracingLog.Bool("mailbox.found", true))
		span.LogFields(tracingLog.String("mailbox.username", username))

		emailExists, err := services.CommonServices.PostgresRepositories.RawEmailRepository.EmailExistsByMessageId(externalSystem, tenantByName, username, postmarkEmailWebhookData.MessageID)
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

			err = services.CommonServices.PostgresRepositories.RawEmailRepository.Store(externalSystem, tenantByName, username, emailRawData.ProviderMessageId, postmarkEmailWebhookData.MessageID, string(jsonContent), emailRawData.Sent, entity.REAL_TIME)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

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

	messageId := ""
	threadId := ""
	//search in headers for Message-ID
	for k, v := range headers {
		if k == "Message-Id" {
			messageId = v
		}
		if k == "X-Session-ID" {
			threadId = v
		}
	}

	if messageId == "" {
		messageId = pmData.MessageID
	}

	return entity.EmailRawData{
		ProviderMessageId: pmData.MessageID,
		MessageId:         messageId,
		Sent:              *sentTime,
		Subject:           pmData.Subject,
		From:              from,
		To:                strings.Join(to, ", "),
		Cc:                strings.Join(cc, ", "),
		Bcc:               strings.Join(bcc, ", "),
		Html:              pmData.HtmlBody,
		Text:              pmData.TextBody,
		ThreadId:          threadId,
		InReplyTo:         pmData.ReplyTo,
		Reference:         pmData.OriginalRecipient,
		Headers:           headers,
	}, nil
}
