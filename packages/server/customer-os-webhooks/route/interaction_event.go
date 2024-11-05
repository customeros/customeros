package route

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	"github.com/opentracing/opentracing-go"
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

//pending - contacts in flow that are not in the other stages
//completed - contacts that have received the email
//goal achieved - contacts that have received the sign-up email (Welcome to Embedd - Product Tips)

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

		appKey, err := services.CommonServices.PostgresRepositories.AppKeyRepository.FindByKey(ctx, string(security.CUSTOMER_OS_WEBHOOKS), apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		if appKey == nil {
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

		n, err := services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantByNameIgnoreCase(ctx, tenantByName)
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

		tenant := mapper.MapDbNodeToTenantEntity(n)
		tenantByName = tenant.Name

		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant: tenantByName,
		})

		span.LogFields(tracingLog.Bool("tenant.found", true))
		span.LogFields(tracingLog.String("tenant.name", tenantByName))
		span.SetTag(tracing.SpanTagTenant, tenantByName)

		htmlData := strings.ReplaceAll(postmarkEmailWebhookData.HtmlBody, "&amp;", "&")
		textData := strings.ReplaceAll(postmarkEmailWebhookData.TextBody, "&amp;", "&")
		emailExclusion := services.CommonServices.Cache.GetEmailExclusion(tenantByName)
		for _, exclusion := range emailExclusion {
			if exclusion.ExcludeSubject != nil {
				if strings.Contains(postmarkEmailWebhookData.Subject, *exclusion.ExcludeSubject) {
					span.LogFields(tracingLog.String("reason", "excluded by subject"))
					return
				}
			}
			if exclusion.ExcludeBody != nil {
				if strings.Index(htmlData, *exclusion.ExcludeBody) >= 0 {
					span.LogFields(tracingLog.String("reason", "excluded by html body"))
					return
				}
				if strings.Index(textData, *exclusion.ExcludeBody) >= 0 {
					span.LogFields(tracingLog.String("reason", "excluded by text body"))
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

		//TODO remove this hack
		err = processEmailForFlows(ctx, services, tenantByName, postmarkEmailWebhookData.FromFull.Email, participants, postmarkEmailWebhookData.Subject, postmarkEmailWebhookData)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncInteractionEvent) error processing email for flows: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err = processMailstackReply(ctx, services, tenantByName, postmarkEmailWebhookData, cfg.Slack.NotifyFlowGoalAchieved)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncInteractionEvent) error processing email for flows: %s", err.Error())
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
		if b.Email != "" && b.Email != "bcc@"+strings.ToLower(tenant)+".customeros.ai" {
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

//Deprecated

// if the sender is a user in the system, it means that this is outbound communication
// we mark the contacts that received this email as COMPLETED in the flows that they are in
// this is a hack for now as we should identify the flow that the contact is in and mark the contact as COMPLETED only in that specific flow
func processEmailForFlows(ctx context.Context, services *service.Services, tenant, fromEmailAddress string, participantsEmailAddresses []string, emailSubject string, input model.PostmarkEmailWebhookData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.processEmailForFlows")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	senderUser, err := services.CommonServices.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, fromEmailAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	outbound := false
	if senderUser != nil {
		outbound = true
	}

	if outbound {
		for _, p := range participantsEmailAddresses {
			contactsWithEmailNodes, err := services.CommonServices.Neo4jRepositories.ContactReadRepository.GetContactsWithEmail(ctx, tenant, p)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			for _, contactNode := range contactsWithEmailNodes {
				contactEntity := mapper.MapDbNodeToContactEntity(contactNode)

				flowsWithContact, err := services.CommonServices.FlowService.FlowsGetListWithContact(ctx, []string{contactEntity.Id})
				if err != nil {
					tracing.TraceErr(span, err)
					return err
				}

				for _, flow := range *flowsWithContact {
					flowContact, err := services.CommonServices.FlowService.FlowParticipantByEntity(ctx, flow.Id, contactEntity.Id, commonModel.CONTACT)
					if err != nil {
						tracing.TraceErr(span, err)
						return err
					}

					if flowContact == nil || flowContact.Status == neo4jentity.FlowParticipantStatusCompleted || flowContact.Status == neo4jentity.FlowParticipantStatusGoalAchieved {
						continue
					}

					if emailSubject == "Welcome to Embedd - Product Tips" {
						flowContact.Status = neo4jentity.FlowParticipantStatusGoalAchieved
					} else {
						flowContact.Status = neo4jentity.FlowParticipantStatusCompleted
					}

					_, err = services.CommonServices.Neo4jRepositories.FlowParticipantWriteRepository.Merge(ctx, nil, flowContact)
					if err != nil {
						tracing.TraceErr(span, err)
						return err
					}
				}

			}
		}
	}

	return nil
}

// check if the email is a reply to an email sent by mailstack
// if it is, mark the flow participant as GOAL_ACHIEVED
func processMailstackReply(ctx context.Context, services *service.Services, tenant string, input model.PostmarkEmailWebhookData, slackChannelUrl string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.processMailstackReply")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	inReplyTo, err := getInReplyTo(input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// if this email is a reply to a mailstack email
	// and the sender of the email is the same as the recipient of the mailstack email
	// mark the flow participant as GOAL_ACHIEVED
	if inReplyTo != "" {
		mailstackEmail, err := services.CommonServices.PostgresRepositories.EmailMessageRepository.GetByProviderMessageId(ctx, tenant, inReplyTo)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if mailstackEmail != nil && strings.Contains(mailstackEmail.ToString, input.FromFull.Email) && mailstackEmail.ProducerType == commonModel.NodeLabelFlowActionExecution {

			//TODO check that it isn't an automatic reply
			//header: Subject: Automatic reply: First restaurant killed by DoorDash reviews?

			flowActionExecution, err := services.CommonServices.FlowExecutionService.GetFlowActionExecutionById(ctx, mailstackEmail.ProducerId)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			flowParticipant, err := services.CommonServices.FlowService.FlowParticipantByEntity(ctx, flowActionExecution.FlowId, flowActionExecution.EntityId, flowActionExecution.EntityType)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			err = services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonModel.NodeLabelFlowParticipant, flowParticipant.Id, "status", string(neo4jentity.FlowParticipantStatusGoalAchieved))
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			primaryEmailForParticipant, err := services.CommonServices.EmailService.GetPrimaryEmailForEntityId(ctx, flowParticipant.EntityType, flowParticipant.EntityId)
			if err != nil {
				tracing.TraceErr(span, err)
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
					tracing.TraceErr(span, err)
				}
			}

			err = services.CommonServices.RabbitMQService.PublishEvent(ctx, flowActionExecution.FlowId, commonModel.FLOW, dto.FlowParticipantGoalAchieved{
				ParticipantId:   flowParticipant.EntityId,
				ParticipantType: flowParticipant.EntityType,
			})
			if err != nil {
				tracing.TraceErr(span, err)
			}

			services.CommonServices.RabbitMQService.PublishEventCompleted(ctx, tenant, flowParticipant.Id, commonModel.FLOW_PARTICIPANT, utils.NewEventCompletedDetails().WithUpdate())
		}

	}

	return nil
}
