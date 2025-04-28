package integrations

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/gin-gonic/gin"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
)

const EXTERNAL_SYSTEM = "mailstack"

func (h *IntegrationHandler) PostmarkInboundEmail(c *gin.Context) {
	spans, _ := telemetry.StartRestSpan(c.Request.Context(), "IntegrationHandler.PostmarkInboundEmail")
	defer spans.Finish()

	// Validate Postmark User-Agent
	if c.Request.UserAgent() == "" || !strings.EqualFold(c.Request.UserAgent(), "Postmark") {
		spans.TraceError(fmt.Errorf("invalid user agent %s", c.Request.UserAgent()))
		h.responseHandler.HandleError(c, http.StatusForbidden, nil)
		return
	}

	// Parse email data
	emailData, err := h.parseInboundEmail(c)
	if err != nil {
		spans.LogObjectAsJson("body", c.Request.Body)
		spans.TraceError(err)
		h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
		return
	}

	// Process email asynchronously
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				err := fmt.Errorf("panic recovered in email processing: %v\n%s", r, stack)
				spans.TraceError(err)
			}
		}()

		if err := h.processInboundEmail(c, &emailData); err != nil {
			spans.TraceError(errors.Wrap(err, "failed to process inbound email from Postmark"))
		}
	}()
}

func (h *IntegrationHandler) parseInboundEmail(c *gin.Context) (PostmarkInboundEmailData, error) {
	var emailData PostmarkInboundEmailData
	err := c.BindJSON(&emailData)
	if err != nil {
		return emailData, err
	}

	return emailData, nil
}

func (h *IntegrationHandler) processInboundEmail(c *gin.Context, emailData *PostmarkInboundEmailData) error {
	spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "IntegrationHandler.processInboundEmail")
	defer spans.Finish()

	tenant, err := h.getTenant(c, emailData)
	if err != nil {
		return err
	}

	// update context with tenant, pass this where tenant is needed
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    tenant,
		AppSource: constants.AppSourceCustomerOsApiRest,
	})

	participants := emailData.AllParticipantEmails()
	username, err := h.getUsername(ctx, participants)
	if err != nil || username == "" {
		spans.LogFields(tracingLog.Bool("mailbox.found", false))
		spans.TraceError(err)
		return err
	}

	messageId := emailData.GetHeaderValue("Message-Id")
	emailExistsInDb, err := h.services.Repositories.PostgresRepositories.IngestEmailMessageRepository.EmailExistsByMessageId(
		ctx, EXTERNAL_SYSTEM, tenant, username, messageId,
	)
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("unable to determine if email exists in db for messageId %s: %v", messageId, err)
	}

	if emailExistsInDb {
		return nil
	}

	//TODO implement when we switch to this instead of webhooks app

	//dbEmailEntity := emailData.ToRawDbObject()
	//jsonEmailEntity, err := json.Marshal(dbEmailEntity)
	//if err != nil {
	//	return fmt.Errorf("Unable to produce JSON email object for db: %v", err)
	//}

	//ingestEmailMessage := postgres_entity.IngestEmailMessage{
	//	Tenant:   tenant,
	//	Username: username,
	//	Provider: enum.SourceMailstack.String(),
	//	State:    postgres_entity.IngestEmailMessageStatePending,
	//
	//	Subject:     emailRawData.Subject,
	//	TextContent: emailRawData.Text,
	//	HtmlContent: emailRawData.Html,
	//
	//	SentAt: emailRawData.Sent,
	//
	//	From: emailRawData.From,
	//	To:   emailRawData.To,
	//	Cc:   emailRawData.Cc,
	//	Bcc:  emailRawData.Bcc,
	//
	//	ProviderMessageId:  emailRawData.ProviderMessageId,
	//	ProviderThreadId:   emailRawData.ThreadId,
	//	ProviderInReplyTo:  emailRawData.InReplyTo,
	//	ProviderReferences: emailRawData.Reference,
	//
	//	Headers: string(headersString),
	//}
	//
	//dbErr := h.services.Repositories.PostgresRepositories.IngestEmailMessageRepository.Store(
	//	ctx, EXTERNAL_SYSTEM, tenant, username, dbEmailEntity.ProviderMessageId, messageId, string(jsonEmailEntity), dbEmailEntity.Sent, postgres_entity.REAL_TIME,
	//)
	//if dbErr != nil {
	//	span.LogFields(tracingLog.Object("raw_email", jsonEmailEntity))
	//	tracing.TraceErr(span, err)
	//	return fmt.Errorf("Error writing email to db: %v", err)
	//}

	// Check to see if email is a reply to a flow.  If so, mark as complete.
	// This should be handled in the email processor common service, not here.
	// Same with goal achieved.
	// Move slack notifications to processing service
	// Handle attachments

	return nil
}

func (h *IntegrationHandler) getTenant(c *gin.Context, emailData *PostmarkInboundEmailData) (string, error) {
	spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "IntegrationHandler.getTenant")
	defer spans.Finish()

	nameFromBcc := emailData.TenantFromBcc()

	n, err := h.services.Repositories.Neo4jRepositories.TenantReadRepository.GetTenantByNameIgnoreCase(ctx, nameFromBcc)
	if err != nil {
		spans.LogFields(tracingLog.Bool("tenant.found", false))
		spans.TraceError(err)
		return "", fmt.Errorf("cannot identify tenant %s: %v", nameFromBcc, err)
	}

	if n == nil {
		spans.LogKV("tenant.found", false)
		return "", fmt.Errorf("no valid tenant %s: %v", nameFromBcc, err)
	}

	tenant := mapper.MapDbNodeToTenantEntity(n)
	spans.LogKV("tenant.found", true)
	spans.LogKV("tenant.name", tenant.Name)

	return tenant.Name, nil
}

func (h *IntegrationHandler) getUsername(ctx context.Context, EmailParticipants []string) (string, error) {
	spans, _ := telemetry.StartRestSpan(ctx, "IntegrationHandler.getUsername")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	for _, p := range EmailParticipants {
		userByEmail, err := h.services.Repositories.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(
			ctx, tenant, p)
		if err != nil {
			return "", fmt.Errorf("error getting username from email %s: %v", p, err)
		}
		if userByEmail != nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("unable to find user amongst email participants")
}
