package route

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

	commoncaches "github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-webhooks/config"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/model"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/service"
)

func AddInteractionEventRoutes(ctx context.Context, route *gin.Engine, services *service.Services, cfg *config.Config, log logger.Logger, cache *commoncaches.Cache) {
	route.POST("/sync/postmark-interaction-event",
		RestTracingEnhancer(ctx, "/sync/postmark-interaction-event"),
		syncPostmarkInteractionEventHandler(services, cfg, log))
}

// pending - contacts in flow that are not in the other stages
// completed - contacts that have received the email
// goal achieved - contacts that have received the sign-up email (Welcome to Embedd - Product Tips)
func syncPostmarkInteractionEventHandler(services *service.Services, cfg *config.Config, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "syncPostmarkInteractionEventHandler", c.Request.Header)
		defer spans.Finish()

		// check API key as param
		apiKey := c.Query(security.ApiKeyHeader)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		if cfg.App.AppKey != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		body := c.Request.Body
		requestBody, err := io.ReadAll(body)
		if err != nil {
			spans.LogObjectAsJson("body", body)
			spans.TraceError(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var postmarkEmailWebhookData model.PostmarkEmailWebhookData
		if err = json.Unmarshal(requestBody, &postmarkEmailWebhookData); err != nil {
			spans.LogObjectAsJson("requestBody", requestBody)
			spans.TraceError(err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		spans.LogObjectAsJson("webhookData", postmarkEmailWebhookData)

		pattern := `@([^.]+)\.`
		tenantNamePattern, err := regexp.Compile(pattern)
		if err != nil {
			spans.TraceError(err)
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
			spans.LogKV("tenant.found", false)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		n, err := services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantByNameIgnoreCase(ctx, tenantByName)
		if err != nil {
			spans.LogKV("tenant.found", false)
			spans.TraceError(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if n == nil {
			spans.LogKV("tenant.found", false)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		tenant := mapper.MapDbNodeToTenantEntity(n)
		tenantByName = tenant.Name

		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant: tenantByName,
		})

		spans.LogKV("tenant.found", true)
		spans.LogKV("tenant.name", tenantByName)
		spans.TagTenant(tenantByName)

		htmlData := strings.ReplaceAll(postmarkEmailWebhookData.HtmlBody, "&amp;", "&")
		textData := strings.ReplaceAll(postmarkEmailWebhookData.TextBody, "&amp;", "&")
		emailExclusion := services.CommonServices.Cache.GetEmailExclusion(tenantByName)
		for _, exclusion := range emailExclusion {
			if exclusion.ExcludeSubject != nil {
				if strings.Contains(postmarkEmailWebhookData.Subject, *exclusion.ExcludeSubject) {
					spans.LogKV("reason", "excluded by subject")
					return
				}
			}
			if exclusion.ExcludeBody != nil {
				if strings.Index(htmlData, *exclusion.ExcludeBody) >= 0 {
					spans.LogKV("reason", "excluded by html body")
					return
				}
				if strings.Index(textData, *exclusion.ExcludeBody) >= 0 {
					spans.LogKV("reason", "excluded by text body")
					return
				}
			}
		}

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
				if bcc.Email != "" && bcc.Email != "bcc@"+strings.ToLower(tenantByName)+".customeros.ai" {
					participants = append(participants, bcc.Email)
				}
			}
		}

		// identify mailbox
		username := ""
		for _, p := range participants {
			mailboxRecord, err := services.CommonServices.MailstackService.GetByMailbox(ctx, tenantByName, p)
			if err != nil {
				spans.TraceError(err)
				log.Errorf("(SyncInteractionEvent) error getting mailbox: %s", err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			if mailboxRecord != nil {
				username = p
				break
			}
		}

		messageId, err := getMessageId(postmarkEmailWebhookData)
		if err != nil {
			spans.TraceError(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if username == "" {
			spans.LogKV("mailbox.found", false)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		spans.LogKV("mailbox.found", true)
		spans.LogKV("mailbox.username", username)

		emailExists, err := services.CommonServices.PostgresRepositories.IngestEmailMessageRepository.EmailExistsByMessageId(ctx, tenantByName, username, externalSystem, messageId)
		if err != nil {
			spans.TraceError(err)
			log.Errorf("(SyncInteractionEvent) error checking email exists: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if !emailExists {
			emailRawData, err := mapPostmarkToEmailRawData(tenantByName, postmarkEmailWebhookData)

			headersString, err := JSONMarshal(emailRawData.Headers)
			if err != nil {
				spans.TraceError(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			ingestEmailMessage := postgres_entity.IngestEmailMessage{
				Tenant:   tenantByName,
				Username: username,
				Provider: enum.SourceMailstack.String(),
				State:    postgres_entity.IngestEmailMessageStatePending,

				Subject:     emailRawData.Subject,
				TextContent: emailRawData.Text,
				HtmlContent: emailRawData.Html,

				SentAt: emailRawData.Sent,

				From: emailRawData.From,
				To:   emailRawData.To,
				Cc:   emailRawData.Cc,
				Bcc:  emailRawData.Bcc,

				MessageId:          emailRawData.MessageId,
				ProviderMessageId:  emailRawData.ProviderMessageId,
				ProviderThreadId:   emailRawData.ThreadId,
				ProviderInReplyTo:  emailRawData.InReplyTo,
				ProviderReferences: emailRawData.Reference,

				Headers: string(headersString),
			}

			err = services.CommonServices.PostgresRepositories.IngestEmailMessageRepository.Store(ctx, tenantByName, username, externalSystem, messageId, &ingestEmailMessage)
			if err != nil {
				spans.TraceError(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			storedRawEmail, err := services.CommonServices.PostgresRepositories.IngestEmailMessageRepository.GetByMessageId(ctx, externalSystem, tenantByName, username, messageId)
			if err != nil {
				spans.TraceError(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			loadedEmail, err := services.CommonServices.MailService.LoadIngestEmailMessage(ctx, storedRawEmail)
			if err != nil {
				spans.TraceError(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			processEmailCheck := services.CommonServices.MailService.ProcessEmailCheck(ctx, tenantByName, &loadedEmail)

			if processEmailCheck.ProcessEmail ||
				processEmailCheck.SkipReason == "BULK | FROM NON-PRIMARY DOMAIN" { // allow personal emails to be processed
				err = processMailstackReply(ctx, services, tenantByName, postmarkEmailWebhookData, cfg.Common.External.SlackConfig.NotifyFlowGoalAchieved)
				if err != nil {
					spans.TraceError(err)
					log.Errorf("(SyncInteractionEvent) error processing email for flows: %s", err.Error())
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
					return
				}

				if cfg.Common.External.SlackConfig.NotifyPostmarkEmail != "" {
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

					utils.SendSlackMessage(ctx, cfg.Common.External.SlackConfig.NotifyPostmarkEmail, slackMessageText)
				}
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

func mapPostmarkToEmailRawData(tenant string, pmData model.PostmarkEmailWebhookData) (postgres_entity.EmailRawData, error) {
	// Parse the Date field to time.Time
	sentTime, err := utils.UnmarshalDateTime(pmData.Date)
	if err != nil {
		return postgres_entity.EmailRawData{}, err
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
		if b.Email != "" && b.Email != "bcc@"+strings.ToLower(tenant)+".customeros.ai" {
			bcc = append(bcc, b.Email)
		}
	}

	messageId, err := getMessageId(pmData)
	if err != nil {
		return postgres_entity.EmailRawData{}, err
	}

	references, err := getReferences(pmData)
	if err != nil {
		return postgres_entity.EmailRawData{}, err
	}

	inReplyTo, err := getInReplyTo(pmData)
	if err != nil {
		return postgres_entity.EmailRawData{}, err
	}

	return postgres_entity.EmailRawData{
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

// check if the email is a reply to an email sent by mailstack
// if it is, mark the flow participant as GOAL_ACHIEVED
func processMailstackReply(ctx context.Context, services *service.Services, tenant string, input model.PostmarkEmailWebhookData, slackChannelUrl string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.processMailstackReply")
	defer spans.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	inReplyTo, err := getInReplyTo(input)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	// if this email is a reply to a mailstack email
	// and the sender of the email is the same as the recipient of the mailstack email
	// mark the flow participant as GOAL_ACHIEVED
	if inReplyTo != "" {
		mailstackEmail, err := services.CommonServices.PostgresRepositories.EmailMessageRepository.GetByProviderMessageId(ctx, tenant, inReplyTo)
		if err != nil {
			spans.TraceError(err)
			return err
		}

		if mailstackEmail != nil && strings.Contains(mailstackEmail.ToString, input.FromFull.Email) && mailstackEmail.ProducerType == commonModel.NodeLabelFlowActionExecution {

			flowActionExecution, err := services.CommonServices.FlowExecutionService.GetFlowActionExecutionById(ctx, mailstackEmail.ProducerId)
			if err != nil {
				spans.TraceError(err)
				return err
			}

			flowParticipant, err := services.CommonServices.FlowService.FlowParticipantByEntity(ctx, flowActionExecution.FlowId, flowActionExecution.EntityId, flowActionExecution.EntityType)
			if err != nil {
				spans.TraceError(err)
				return err
			}

			err = services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonModel.NodeLabelFlowParticipant, flowParticipant.Id, "status", string(neo4jentity.FlowParticipantStatusGoalAchieved))
			if err != nil {
				spans.TraceError(err)
				return err
			}

			primaryEmailForParticipant, err := services.CommonServices.EmailService.GetPrimaryEmailForEntityId(ctx, flowParticipant.EntityType, flowParticipant.EntityId)
			if err != nil {
				spans.TraceError(err)
				return err
			}

			if primaryEmailForParticipant == nil {
				primaryEmailForParticipant = &neo4jentity.EmailEntity{
					RawEmail: "primary email missing",
				}
			}

			if slackChannelUrl != "" {
				slackMessageText := "*Tenant:* " + tenant + "\n"
				slackMessageText += "*Goal achieved for:* " + primaryEmailForParticipant.RawEmail + "\n"

				err := utils.SendSlackMessage(ctx, slackChannelUrl, slackMessageText)
				if err != nil {
					spans.TraceError(err)
				}
			}

			err = services.CommonServices.Events.Publisher.PublishFanoutEvent(ctx, flowActionExecution.FlowId, commonModel.FLOW, dto.FlowParticipantGoalAchieved{
				ParticipantId:   flowParticipant.EntityId,
				ParticipantType: flowParticipant.EntityType,
			})
			if err != nil {
				spans.TraceError(err)
			}

			services.CommonServices.Events.Publisher.PublishNotification(ctx, tenant, flowParticipant.Id, commonModel.FLOW_PARTICIPANT, utils.NewEventCompletedDetails().WithUpdate())
		}

	}

	return nil
}
